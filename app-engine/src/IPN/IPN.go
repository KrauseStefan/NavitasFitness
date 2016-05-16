package IPN

import (
	"github.com/gorilla/mux"
	"net/http"
	"bytes"
	"io/ioutil"

	"appengine"
	"appengine/urlfetch"
	"appengine/taskqueue"
	"src/IPN/Transaction"
	"src/User/Dao"
	"errors"
	"fmt"
	"encoding/json"
)

const (
	localIpn = "http://localhost:8081/cgi-bin/webscr" //(Will behave like a live IPN)
	PaypalIpn = "https://www.paypal.com/cgi-bin/webscr" //(for live IPNs)
	PaypalIpnSandBox = "https://www.sandbox.paypal.com/cgi-bin/webscr" // (for Sandbox IPNs)

	IpnQueueName = "paypalIpn"

	basePath = "/rest/paypal"
	ipnUrl = basePath + "/ipn"
	ipnRespondTaskUrl = basePath + "/ipnDoRespone"

	FromEncodedContentType = "application/x-www-form-urlencoded"

	ReceiverEmail = "stefan.krausekjaer@gmail.com" //TODO Verify that you are the intended recipient of the IPN message. To do this, check the email address in the message. This check prevents another merchant from accidentally or intentionally using your listener.
)

func IntegrateRoutes(router *mux.Router) {

	router.
	Methods("POST").
	Path(ipnUrl).
	Name("ipn notification").
	HandlerFunc(processIPN)

	router.
	Methods("POST").
	Path(ipnRespondTaskUrl).
	Name("ipn notification responder task").
	HandlerFunc(ipnDoResponseTaskHandler)

}

//Receives the IPN message from PayPal and sends a empty response back code 200
//The received package is parsed to the task queue, where the appropriate response is made
//This is need in order to close the original request before responding
func processIPN(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task := newTask(ipnRespondTaskUrl, content)
	if _, err := taskqueue.Add(ctx, task, IpnQueueName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Persist transaction details in unconfirmed state (processing?)
}

func sendVerificationMassageToPaypal(ctx appengine.Context, content string, testIpnField string) (string, error) {

	paypalIpnUrl := PaypalIpn
	if testIpnField != "" {
		if testIpnField == "1" {
			paypalIpnUrl = PaypalIpnSandBox
		} else {
			paypalIpnUrl = localIpn
		}
	}

	ctx.Infof("Sending msg to: " + paypalIpnUrl)
	extraData := []byte("cmd=_notify-validate&")
	client := urlfetch.Client(ctx)
	resp, err := client.Post(paypalIpnUrl, FromEncodedContentType, bytes.NewBuffer(append(extraData, content...)))
	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	resp.Body.Close()
	return string(respBody), err
}

// Send message for verification with cmd=_notify-validate prepended
// Verify message //discard if message cannot be verified
// Lookup Transaction to update (or create new)
// Lookup user, if it exists verify that it is the same as a found transaction
func ipnDoResponseTaskHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if err := ipnDoResponseTask(ctx, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		ctx.Errorf(err.Error())
	}
}

func ipnDoResponseTask(ctx appengine.Context, r *http.Request) error {

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	transaction := TransactionDao.NewTransactionMsgDTO()
	transaction.AddNewIpnMessage(string(content))
	testIpnField := transaction.GetField(TransactionDao.FIELD_TEST_IPN)
	email := transaction.GetField(TransactionDao.FIELD_CUSTOM) //The custom field should contain the email
	respBody, err := sendVerificationMassageToPaypal(ctx, string(content), testIpnField)
	if err != nil {
		return err
	}
	if string(respBody) != "VERIFIED" {
		return errors.New("Message was not verified by paypal, either a fake message or program error, msg: " + string(respBody))
	}

	//message is now verified and should be persisted

	ctx.Debugf(fmt.Sprintf("%s: %q", TransactionDao.FIELD_PAYMENT_STATUS, transaction.GetField(TransactionDao.FIELD_PAYMENT_STATUS)))

	savedTransaction, err := TransactionDao.GetTransaction(ctx, transaction.GetField(TransactionDao.FIELD_TXN_ID))
	if err != nil {
		return err
	}
	if savedTransaction != nil {
		if transaction.GetPaymentStatus() == savedTransaction.GetPaymentStatus() {
			//Verify that the IPN is not a duplicate. To do this, save the transaction ID and last payment status in each IPN message in a database and verify that the current IPN's values for these fields are not already in this database.
			//Duplicate txnMsg
			//Persist anyway?, with status duplicate?

			return TransactionDao.TxnDuplicateTxnMsg
		}
		savedTransaction.AddNewIpnMessage(string(content))
		savedTransaction.StatusResp = string(respBody)

		ctx.Debugf("IpnData: %q", respBody)

		js, _ := json.Marshal(savedTransaction)
		ctx.Debugf("IpnSaved: %q", js)


		if err := TransactionDao.UpdateIpnMessage(ctx, savedTransaction); err != nil {
			return err
		}
	} else {
		ctx.Infof(fmt.Sprintf("TxnId not found: %q", transaction.GetField(TransactionDao.FIELD_TXN_ID)))

		if email == "" {
			// TODO: handle bad request in a way other then discarding, save without a parent user
			//http.Error(w, "No email received", http.StatusBadRequest)
			ctx.Infof(fmt.Sprintf("TxnId nor email found in msg: %q", string(respBody)))
			return errors.New("Neither transaction ID nor email could be used to lookup user")
		}

		transaction.StatusResp = string(respBody)

		ctx.Infof(fmt.Sprintf("Recived transaction from: %q", email))
		user, err := UserDao.GetUserByEmail(ctx, email)
		if err != nil {
			return err
		}
		if user == nil && savedTransaction == nil {
			//http.Error(w, "User does not exist", http.StatusBadRequest)
			return errors.New("User does not exist")
		}

		if user != nil {
			ctx.Infof(fmt.Sprintf("User key: %q", user.Key))
		}

		ctx.Debugf("IpnData: %q", respBody)

		js, _ := json.Marshal(transaction)
		ctx.Debugf("IpnSaved: %q", js)

		if err := TransactionDao.PersistNewIpnMessage(ctx, transaction, user.Key); err != nil {
			return err
		}
	}

	return nil
	//if(string(respBody) == "VERIFIED") {
	//	// register ipn message
	//	// Verify amount is correct
	//} else if(string(respBody) == "INVALID") {
	//	// Log Invalid payment
	//} else {
	//	// Log Severe error (We are properly being hacked at this point) How to make sure this could never happen
	//}
}

func newTask(path string, data []byte) *taskqueue.Task {
	h := make(http.Header)
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	return &taskqueue.Task{
		Path:    path,
		Payload: data,
		Header:  h,
		Method:  "POST",
	}
}
