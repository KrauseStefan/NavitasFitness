package subscriptionExpiration

import (
	"AppEngineHelper"
	"DAOHelper"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	TransactionDao "IPN/Transaction"
	UserDao "User/Dao"
	UserService "User/Service"
	"constants"
	log "logger"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/mail"
)

var userDao = UserDao.GetInstance()

const ExpirationWarningOffsetDays = 7

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

	users, _, usersErr := getAboutToExpireTxnsWithUsers(ctx)
	usersJson, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}

	strs := []string{string(usersJson)}

	if usersErr != nil {
		httpError := DAOHelper.DefaultHttpError{InnerError: usersErr}
		strs = append(strs, httpError.Error())
	}

	return strings.Join(strs, "\n"), nil
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
		bcc := []string{"stefan.krausekjaer@gmail.com"}
		if callingUser != nil && callingUser.Email != "" {
			bcc = []string{callingUser.Email}
		}

		err := sendEmail(ctx, user, bcc)
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

func getAboutToExpireTxnsWithUsers(ctx context.Context) (UserDao.UserList, TransactionDao.TransactionList, error) {
	paymentExpiratinDate := time.Now().AddDate(0, -constants.SubscriptionDurationInMonth, 0)
	paymentWarningDate := paymentExpiratinDate.AddDate(0, 0, ExpirationWarningOffsetDays)

	format := "02-01-06T15:04:05-07:00"
	log.Infof(ctx, "PaymentDate>=DATETIME('%s')", paymentExpiratinDate.Format(format))
	log.Infof(ctx, "PaymentDate<=DATETIME('%s')", paymentWarningDate.Format(format))

	txns, err := TransactionDao.GetTransactionsPayedBetween(ctx, paymentExpiratinDate, paymentWarningDate)
	if err != nil {
		return nil, nil, err
	}
	log.Infof(ctx, "txns found %d", len(txns))

	aboutToExpireTxn := txns.
		Filter(func(txn *TransactionDao.TransactionMsgDTO) bool { return !txn.ExpirationWarningGiven() })
	log.Infof(ctx, "aboutToExpireTxn found %d", len(aboutToExpireTxn))

	usersResp, err := userDao.GetByKeys(ctx, aboutToExpireTxn.GetUserKeys())

	err = AppEngineHelper.ToMultiError(err).
		Filter(func(err error, i int) bool {
			if err != nil {
				if err == datastore.ErrNoSuchEntity {
					txn := aboutToExpireTxn[i]
					log.Warningf(ctx, "Ignoring Error from email: '%s', id: '%s', error: '%s'", txn.GetPayerEmail(), txn.GetTxnId(), err.Error())
				} else {
					return true
				}
			}
			return false
		}).
		ToError()

	users := usersResp.Filter(func(user *UserDao.UserDTO) bool {
		return user != nil
	})

	return users, aboutToExpireTxn, err
}

var subscriptionExpiredEmailBodyTbl = `
This is a reminder that your membership to Navitas fitnes is about to expire
Your membership to Navitas fitness will expire within %d days, after that date you will no longer have access.

If you which to extend your membership you may follow the below directions:
<ol>
	<li>Login at <a href="https://navitas-fitness-aarhus.appspot.com/">navitas-fitness-aarhus.appspot.com</a> on our website with your AccessId (%s) and password</li>
	<li>Click on the "Payment Status" tab</li>
	<li>Perform a new payment following the instructions on the page</li>
</ol>

Kind regards
Navitas-Fitness
`

func sendEmail(ctx context.Context, user *UserDao.UserDTO, bcc []string) error {
	msg := &mail.Message{
		Sender:   "noreply - Navitass Fitness <navitas-fitness-aarhus@appspot.gserviceaccount.com>",
		To:       []string{user.Email},
		Bcc:      bcc,
		Subject:  "Your membership to Navitas-Fitness is about to expire",
		HTMLBody: fmt.Sprintf(subscriptionExpiredEmailBodyTbl, ExpirationWarningOffsetDays, user.AccessId),
	}

	err := mail.Send(ctx, msg)
	if err == nil {
		log.Infof(ctx, "Subscription expiration mail sent to '%s'", user.Email)
	}
	return err
}
