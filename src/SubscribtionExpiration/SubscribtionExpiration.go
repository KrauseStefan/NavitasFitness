package subscriptionExpiration

import (
	"context"
	"fmt"
	"google.golang.org/appengine/log"
	"net/http"

	"IPN/Transaction"
	"User/Dao"
	"User/Service"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/mail"
)

var userDao = UserDao.GetInstance()

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/subscriptionExpiration"

	router.
		Methods("GET").
		Path(path + "/dryRun").
		Name("subscriptionExpiration dryRun").
		HandlerFunc(UserService.AsAdmin(dryRun))

	router.
		Methods("GET").
		Path(path + "/send").
		Name("subscriptionExpiration send").
		HandlerFunc(UserService.AsAdmin(send))

}

func dryRun(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
	ctx := appengine.NewContext(r)

	return dryRunImpl(ctx)
}

func send(w http.ResponseWriter, r *http.Request, callingUser *UserDao.UserDTO) (interface{}, error) {
	ctx := appengine.NewContext(r)

	users, err := dryRunImpl(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		sendEmail(ctx, user, callingUser.Email)
	}

	return nil, nil
}

func dryRunImpl(ctx context.Context) ([]*UserDao.UserDTO, error) {
	txnKeys, err := TransactionDao.GetTransactionsAboutToExpire(ctx)
	if err != nil {
		return nil, err
	}

	userKeys := make([]*datastore.Key, 0, len(txnKeys))

	for _, key := range txnKeys {
		userKeys = append(userKeys, key.Parent())
	}

	return userDao.GetByKeys(ctx, userKeys)
}

var subscriptionExpiredEmailBodyTbl = `
This is a reminder that your membership to Navitas fitnes is about to expire

Your membership to Navitas fitness will expire within 7 days, after that date you will no longer have access.

If you which to extend your membership you may follow the below directions:
<ol>
	<li>Login at <a href="https://navitas-fitness-aarhus.appspot.com/">navitas-fitness-aarhus.appspot.com</a> on our website with your AccessId (%s) and password</li>
	<li>Click on the "Payment Status" tab</li>
	<li>Perform a new payment following the instructions on the page</li>
</ol>

Kind regards
Navitas-Fitness
`

func sendEmail(ctx context.Context, user *UserDao.UserDTO, sendTo string) error {
	log.Infof(ctx, "Subscription expiration mail sent to '%s'", user.Email)

	msg := &mail.Message{
		Sender:   "noreply - Navitass Fitness <navitas-fitness-aarhus@appspot.gserviceaccount.com>",
		To:       []string{sendTo},
		Subject:  "Your membership to Navitas-Fitness is about to expire - " + user.Email,
		HTMLBody: fmt.Sprintf(subscriptionExpiredEmailBodyTbl, user.AccessId),
	}

	return mail.Send(ctx, msg)
}
