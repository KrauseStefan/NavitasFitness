package csv

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"AccessIdValidator"
	"Dropbox"
	"IPN/Transaction"
	"User/Dao"
	"User/Service"
	"strings"
)

var (
	userDAO           UserDao.UsersRetriever              = UserDao.GetInstance()
	transactionDao    TransactionDao.TransactionRetriever = TransactionDao.GetInstance()
	accessIdValidator                                     = AccessIdValidator.GetInstance()
)

const (
	csvDateFormat = "02-01-2006"
)

func getFirstAndLastTxn(ctx context.Context, userKey *datastore.Key, date time.Time) (time.Time, time.Time, error) {
	activeSubscriptions, err := transactionDao.GetCurrentTransactionsAfter(ctx, userKey, date)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	validationEmails := strings.Split(AccessIdValidator.GetPaypalValidationEmail(ctx), ":")

	validActiveSubscriptions := make([]*TransactionDao.TransactionMsgDTO, 0, len(activeSubscriptions))
	for _, txn := range activeSubscriptions {
		email := txn.GetReceiverEmail()
		for _, validationEmail := range validationEmails {
			if email == validationEmail {
				validActiveSubscriptions = append(validActiveSubscriptions, txn)
				break
			}
		}
	}

	if len(validActiveSubscriptions) >= 1 {
		firstTxn, lastTxn := getExtrema(validActiveSubscriptions)

		return firstTxn.GetPaymentDate(), lastTxn.GetPaymentDate(), nil
	}

	return time.Time{}, time.Time{}, nil
}

type UserTxnTuple struct {
	user      UserDao.UserDTO
	firstDate time.Time
	lastDate  time.Time
}

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/export"

	router.
		Methods("GET").
		Path(path + "/csv").
		Name("export").
		HandlerFunc(UserService.AsAdmin(exportCsvHandler))
}

func getExtrema(txns []*TransactionDao.TransactionMsgDTO) (*TransactionDao.TransactionMsgDTO, *TransactionDao.TransactionMsgDTO) {
	firstTxn := txns[0]
	lastTxn := txns[0]

	for _, txn := range txns {
		if txn.GetPaymentDate().Before(firstTxn.GetPaymentDate()) {
			firstTxn = txn
		}

		if txn.GetPaymentDate().After(lastTxn.GetPaymentDate()) {
			lastTxn = txn
		}
	}

	return firstTxn, lastTxn
}

func getActiveTransactionList(ctx context.Context) ([]UserTxnTuple, error) {

	userKeys, users, err := userDAO.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	usersWithActiveSubscription := make([]UserTxnTuple, 0, len(userKeys))

	for i, userKey := range userKeys {
		user := users[i]
		isValid, err := accessIdValidator.ValidateAccessIdPrimary(ctx, []byte(user.AccessId))
		if err != nil || !isValid {
			log.Infof(ctx, "%s has paid for access but ID is not valid, skipped in csv export", user.AccessId)
			continue // Skip uses with invalid access ids they are not allowed access
		}

		firstDate, lastDate, err := getFirstAndLastTxn(ctx, userKey, time.Now().AddDate(0, -6, 0))
		if err != nil {
			return nil, err
		}

		if !firstDate.IsZero() && !lastDate.IsZero() {

			tuple := UserTxnTuple{
				user:      user,
				firstDate: firstDate,
				lastDate:  lastDate,
			}
			usersWithActiveSubscription = append(usersWithActiveSubscription, tuple)
		}
	}

	return usersWithActiveSubscription, nil
}

func createCsvFile(ctx context.Context, w io.Writer) error {
	userTxnTuple, err := getActiveTransactionList(ctx)
	if err != nil {
		return err
	}

	//bomPrefix := []byte{0xef, 0xbb, 0xbf}
	windowsNewline := []byte{0x0D, 0x0A}
	comma := []byte{','}
	//w.Write(bomPrefix)

	//N0774,27-06-2016,03-01-2017
	//AAMS-asa,27-06-2016,03-01-2017
	//201505600,27-06-2016,03-01-2017

	for _, user := range userTxnTuple {
		log.Infof(ctx, "%s, %s, %s", user.user.AccessId, user.firstDate.String(), user.lastDate.String())
		w.Write([]byte(user.user.AccessId))
		w.Write(comma)
		w.Write([]byte(user.firstDate.Format(csvDateFormat)))
		w.Write(comma)
		w.Write([]byte(user.lastDate.AddDate(0, 6, 0).Format(csvDateFormat)))
		w.Write([]byte(windowsNewline))
	}

	return nil
}

func exportCsvHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) {
	ctx := appengine.NewContext(r)

	if err := CreateAndUploadFile(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateAndUploadFile(ctx context.Context) error {
	var buffer bytes.Buffer

	tokens, err := Dropbox.GetAccessTokens(ctx)
	if err != nil {
		return err
	}

	if err := createCsvFile(ctx, &buffer); err != nil {
		return err
	}

	for _, token := range tokens {
		if _, err := Dropbox.UploadDoc(ctx, token, AccessIdValidator.GetAccessListPath(ctx), &buffer); err != nil {
			return err
		}
	}

	return nil
}
