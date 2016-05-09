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
	"net/url"
	"src/User/Dao"
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
	HandlerFunc(ipnDoResponseTask)

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

func sendVerificationMassageToPaypal(ctx appengine.Context, content string) (string, error) {

	extraData := []byte("cmd=_notify-validate&")
	client := urlfetch.Client(ctx)
	resp, err := client.Post(PaypalIpnSandBox, FromEncodedContentType, bytes.NewBuffer(append(extraData, content...)))
	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	resp.Body.Close()
	return string(respBody), err
}

// Send message for verification with cmd=_notify-validate prepended
// Verify message
// Lookup user
// Lookup Transaction?
func ipnDoResponseTask(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var t TransActionDao.TransactionMsgDTO

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		ctx.Errorf("error: " + err.Error())
		return
	}
	t.AddNewIpnMessage(string(content))

	respBody, err := sendVerificationMassageToPaypal(ctx, t.GetLatestIPNMessage())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		ctx.Errorf("error: " + err.Error())
		return
	}
	t.StatusResp = string(respBody)

	body, err := url.ParseQuery(t.GetLatestIPNMessage())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		ctx.Errorf("error: " + err.Error())
		return
	}

	email := body.Get(TransActionDao.FIELD_CUSTOM)
	if email == "" { // TODO: handle bad request in a way other then discarding
		http.Error(w, "No email recived", http.StatusBadRequest)
		ctx.Errorf("No email recived")
		return
	}

	ctx.Infof("Recived transaction from: " + email)
	user, err := UserDao.GetUserByEmail(ctx, email)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		ctx.Errorf(err.Error())
		return
	}
	if user == nil {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		ctx.Errorf("User does not exist")
		return
	}

	if (user != nil) {
		ctx.Infof("id: " + user.Key)

		if err := TransActionDao.PersistIpnMessage(ctx, &t, user.Key); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			ctx.Errorf(err.Error())
		}

	}

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
