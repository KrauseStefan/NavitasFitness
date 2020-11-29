package IPN

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"

	"Export/csv"
	TransactionDao "IPN/Transaction"
	UserDao "User/Dao"
	log "logger"
)

var (
	userDAO        = UserDao.GetInstance()
	transactionDao = TransactionDao.GetInstance()
)

const (
	localIpn         = "http://localhost:8082/cgi-bin/webscr"          // (Will behave like a live IPN)
	PaypalIpn        = "https://www.paypal.com/cgi-bin/webscr"         // (for live IPNs)
	PaypalIpnSandBox = "https://www.sandbox.paypal.com/cgi-bin/webscr" // (for Sandbox IPNs)

	IpnQueueName = "paypalIpn"

	basePath          = "/rest/paypal"
	ipnUrl            = basePath + "/ipn"
	ipnRespondTaskUrl = basePath + "/ipnDoResponse"

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

}

// Receives the IPN message from PayPal and sends a empty response back code 200
// The received package is proccessed after the connection is closed
// This is needed in order to close the original request before responding
func processIPN(w http.ResponseWriter, r *http.Request) {

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	server, _ := ctx.Value(http.ServerContextKey).(*http.Server)
	server.ConnState = func(n net.Conn, cs http.ConnState) {
		log.Debugf(ctx, "Callback State Change: %s", cs.String())
		if cs == http.StateClosed {
			bgCtx := context.Background()
			server.ConnState = nil
			if err := ipnDoResponseTask(bgCtx, content); err != nil {
				log.Errorf(ctx, "Error accepting IPN response: %s", err.Error())
			}
		}
	}
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
	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Post(paypalIpnUrl, FromEncodedContentType, bytes.NewBuffer(append(extraData, content...)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

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
func ipnDoResponseTask(ctx context.Context, content []byte) error {
	const expectedAmount = 300 // kr

	txn := TransactionDao.NewTransactionMsgDTOFromIpn(string(content))
	testIpnField := txn.GetField(TransactionDao.FIELD_TEST_IPN)
	email := txn.GetField(TransactionDao.FIELD_CUSTOM) //The custom field should contain the email

	if err := verifyMassageWithPaypal(ctx, string(content), testIpnField); err != nil {
		return err
	}

	//message is now verified and should be persisted

	log.Debugf(ctx, fmt.Sprintf("%s: %q", TransactionDao.FIELD_PAYMENT_STATUS, txn.GetField(TransactionDao.FIELD_PAYMENT_STATUS)))

	savedTransaction, err := transactionDao.GetTransaction(ctx, txn.GetField(TransactionDao.FIELD_TXN_ID))
	if err != nil {
		return err
	}
	if savedTransaction != nil {
		if txn.GetPaymentStatus() == savedTransaction.GetPaymentStatus() {
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
				log.Warningf(ctx, "The amount for the transaction was wrong, recived %f expected %d", savedTransaction.GetAmount(), expectedAmount)
			}
		}

		if err := transactionDao.UpdateIpnMessage(ctx, savedTransaction); err != nil {
			return err
		}
	} else {
		log.Infof(ctx, fmt.Sprintf("No previus transaction with ID: %q", txn.GetField(TransactionDao.FIELD_TXN_ID)))
		log.Infof(ctx, fmt.Sprintf("Recived transaction from: %q", email))

		user, err := userDAO.GetByEmail(ctx, email)
		if err != nil {
			return err
		}
		if user == nil && savedTransaction == nil {
			return errors.New("User does not exist")
		}

		var userKey *datastore.Key
		if user != nil {
			log.Debugf(ctx, fmt.Sprintf("User key: %q", user.Key.Encode()))
			userKey = user.Key
		} else {
			log.Errorf(ctx, "Recived paypal IPN message for unknown user")
		}

		js, _ := json.Marshal(txn)
		log.Debugf(ctx, "IpnSaved: %q", js)

		if txn.PaymentIsCompleted() {
			if txn.GetAmount() != expectedAmount {
				log.Warningf(ctx, "The amount for the transaction was wrong, recived %f expected %d", txn.GetAmount(), expectedAmount)
			}
		}

		if err := transactionDao.PersistNewIpnMessage(ctx, txn, userKey); err != nil {
			return err
		}
	}

	if err := csv.CreateAndUploadFile(ctx, txn); err != nil {
		return err
	}

	return nil
}
