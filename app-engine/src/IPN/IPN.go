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
)

const (
	localIpn					= "http://localhost:8081/cgi-bin/webscr" //(Will behave like a live IPN)
	PaypalIpn					= "https://www.paypal.com/cgi-bin/webscr" //(for live IPNs)
	PaypalIpnSandBox	= "https://www.sandbox.paypal.com/cgi-bin/webscr" // (for Sandbox IPNs)

	IpnQueueName					= "paypalIpn"

	basePath 					= "/rest/paypal"
	ipnUrl 						= basePath + "/ipn"
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

func processIPN(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)

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

func ipnDoResponseTask(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	extraData := []byte("cmd=_notify-validate&")

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		ctx.Infof("error: " + err.Error())
		return
	}

	//ctx.Infof("content: " + string(content))

	client := urlfetch.Client(ctx)
	resp, err := client.Post(localIpn, FromEncodedContentType, bytes.NewBuffer(append(extraData, content...)))
	defer resp.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	//ctx.Infof("resp body: " + string(respBody))

	var t TransActionDao.TransactionMsgDTO
	t.IpnMessage = string(content)
	t.StatusResp = string(respBody)


	if err := TransActionDao.PersistIpnMessage(ctx, &t); err != nil {
		ctx.Errorf(err.Error())
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
