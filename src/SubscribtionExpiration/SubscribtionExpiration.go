package subscriptionExpiration

import (
	"context"
	"fmt"
	"google.golang.org/appengine/log"
	"net/http"
	"time"

	"IPN/Transaction"
	"User/Dao"
	"User/Service"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
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

	users, _, err := getAboutToExpireTxnsWithUsers(ctx)
	return users, err
}

func send(w http.ResponseWriter, r *http.Request, callingUser *UserDao.UserDTO) (interface{}, error) {
	ctx := appengine.NewContext(r)

	users, txns, err := getAboutToExpireTxnsWithUsers(ctx)
	if err != nil {
		return nil, err
	}

	errs := make([]error, 0)
	warnedTxns := make(TransactionDao.TransactionList, 0, len(txns))
	for i, user := range users {
		err := sendEmail(ctx, user, callingUser.Email)
		if err != nil {
			errs = append(errs, err)
		} else {
			warnedTxns = append(warnedTxns, txns[i])
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			log.Errorf(ctx, err.Error())
		}
		return nil, errs[0]
	}

	err = TransactionDao.SetExpirationWarningGiven(ctx, warnedTxns, true)
	return nil, err
}

func getAboutToExpireTxnsWithUsers(ctx context.Context) ([]*UserDao.UserDTO, TransactionDao.TransactionList, error) {
	subscriptionDurationInMonth := 6
	warningDeltaDays := 7

	paymentExpiratinDate := time.Now().AddDate(0, -subscriptionDurationInMonth, 0)
	paymentWarningStartDate := paymentExpiratinDate.AddDate(0, 0, -warningDeltaDays)

	log.Infof(ctx, "PaymentDate>=%s", paymentWarningStartDate)
	log.Infof(ctx, "PaymentDate<=%s", paymentExpiratinDate)

	txns, err := TransactionDao.GetTransactionsPayedBetween(ctx, paymentWarningStartDate, paymentExpiratinDate)
	if err != nil {
		return nil, nil, err
	}

	aboutToExpireTxn := txns.
		Filter(func(txn *TransactionDao.TransactionMsgDTO) bool { return txn.ExpirationWarningGiven() })

	users, err := userDao.GetByKeys(ctx, aboutToExpireTxn.GetUserKeys())
	return users, aboutToExpireTxn, err
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
