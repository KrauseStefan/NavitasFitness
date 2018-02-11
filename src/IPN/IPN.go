package IPN

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
	"google.golang.org/appengine/urlfetch"

	"../Export/csv"
	"../User/Dao"
	"./Transaction"
)

var (
	userDAO        = UserDao.GetInstance()
	transactionDao = TransactionDao.GetInstance()
)

const (
	localIpn         = "http://localhost:8081/cgi-bin/webscr"          //(Will behave like a live IPN)
	PaypalIpn        = "https://www.paypal.com/cgi-bin/webscr"         //(for live IPNs)
	PaypalIpnSandBox = "https://www.sandbox.paypal.com/cgi-bin/webscr" // (for Sandbox IPNs)

	IpnQueueName = "paypalIpn"

	basePath          = "/rest/paypal"
	ipnUrl            = basePath + "/ipn"
	ipnRespondTaskUrl = basePath + "/ipnDoRespone"

	FromEncodedContentType = "application/x-www-form-urlencoded"

	ReceiverEmail = "stefan.krausekjaer@gmail.com" //TODO Verify that you are the intended recipient of the IPN message. To do this, check the email address in the message. This check prevents another merchant from accidentally or intentionally using your listener.
)

var (
	IpnMessageCouldNotBeValidated = errors.New("Ipn message was not vallidated by paypal")
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

func verifyMassageWithPaypal(ctx context.Context, content string, testIpnField string) error {

	paypalIpnUrl := PaypalIpn
	if testIpnField != "" {
		if testIpnField == "1" {
			paypalIpnUrl = PaypalIpnSandBox
		} else {
			paypalIpnUrl = localIpn
		}
	}

	log.Infof(ctx, "Sending msg to: "+paypalIpnUrl)
	extraData := []byte("cmd=_notify-validate&")
	client := urlfetch.Client(ctx)
	resp, err := client.Post(paypalIpnUrl, FromEncodedContentType, bytes.NewBuffer(append(extraData, content...)))
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	log.Debugf(ctx, "Ipn Validation response "+string(respBody))

	if err == nil && string(respBody) != "VERIFIED" {
		return IpnMessageCouldNotBeValidated
	}

	return err
}

// Send message for verification with cmd=_notify-validate prepended
// Verify message //discard if message cannot be verified
// Lookup Transaction to update (or create new)
// Lookup user, if it exists verify that it is the same as a found transaction
func ipnDoResponseTaskHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if err := ipnDoResponseTask(ctx, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}

}

func ipnDoResponseTask(ctx context.Context, r *http.Request) error {
	const expectedAmount = 300 // kr

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	transaction := TransactionDao.NewTransactionMsgDTOFromIpn(string(content))
	testIpnField := transaction.GetField(TransactionDao.FIELD_TEST_IPN)
	email := transaction.GetField(TransactionDao.FIELD_CUSTOM) //The custom field should contain the email

	if err := verifyMassageWithPaypal(ctx, string(content), testIpnField); err != nil {
		return err
	}

	//message is now verified and should be persisted

	log.Debugf(ctx, fmt.Sprintf("%s: %q", TransactionDao.FIELD_PAYMENT_STATUS, transaction.GetField(TransactionDao.FIELD_PAYMENT_STATUS)))

	savedTransaction, err := transactionDao.GetTransaction(ctx, transaction.GetField(TransactionDao.FIELD_TXN_ID))
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

		js, _ := json.Marshal(savedTransaction)
		log.Debugf(ctx, "IpnSaved: %q", js)

		if savedTransaction.PaymentIsCompleted() {
			if savedTransaction.GetAmount() < expectedAmount {
				log.Warningf(ctx, "The amount for the transaction was wrong, recived %f expected %f", savedTransaction.GetAmount(), expectedAmount)
			}
		}

		if err := transactionDao.UpdateIpnMessage(ctx, savedTransaction); err != nil {
			return err
		}
	} else {
		log.Infof(ctx, fmt.Sprintf("No previus transaction with ID: %q", transaction.GetField(TransactionDao.FIELD_TXN_ID)))
		log.Infof(ctx, fmt.Sprintf("Recived transaction from: %q", email))

		user, err := userDAO.GetByEmail(ctx, email)
		if err != nil {
			return err
		}
		if user == nil && savedTransaction == nil {
			return errors.New("User does not exist")
		}

		var userKey *datastore.Key = nil
		if user != nil {
			log.Debugf(ctx, fmt.Sprintf("User key: %q", user.Key.Encode()))
			userKey = user.Key
		} else {
			log.Errorf(ctx, "Recived paypal IPN message for unknown user")
		}

		js, _ := json.Marshal(transaction)
		log.Debugf(ctx, "IpnSaved: %q", js)

		if transaction.PaymentIsCompleted() {
			if transaction.GetAmount() != expectedAmount {
				log.Warningf(ctx, "The amount for the transaction was wrong, recived %f expected %f", transaction.GetAmount(), expectedAmount)
			}
		}

		if err := transactionDao.PersistNewIpnMessage(ctx, transaction, userKey); err != nil {
			return err
		}
	}

	if err := csv.CreateAndUploadFile(ctx); err != nil {
		return err
	}

	return nil
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
