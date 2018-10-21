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
	csvDateFormat           = "02-01-2006"
	subscriptionPeriodMonth = 6
)

type UserTxnTuple struct {
	user      UserDao.UserDTO
	startDate time.Time
	endDate   time.Time
}

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/export"

	router.
		Methods("GET").
		Path(path + "/csv").
		Name("export").
		HandlerFunc(UserService.AsAdmin(exportCsvHandler))
}

type UserTransactionMap map[*datastore.Key][]*TransactionDao.TransactionMsgDTO

func validateTransactions(ctx context.Context, txns []*TransactionDao.TransactionMsgDTO) ([]*TransactionDao.TransactionMsgDTO, []*TransactionDao.TransactionMsgDTO) {
	validationEmails := strings.Split(AccessIdValidator.GetPaypalValidationEmail(ctx), ":")
	validationEmailsLendth := len(validationEmails)
	validTxns := make([]*TransactionDao.TransactionMsgDTO, 0, len(txns))
	invalidTxns := make([]*TransactionDao.TransactionMsgDTO, 0, 5)
	for _, txn := range txns {
		email := txn.GetReceiverEmail()
		for i, validationEmail := range validationEmails {
			if email == validationEmail {
				validTxns = append(validTxns, txn)
				break
			} else if i == validationEmailsLendth-1 {
				log.Warningf(ctx, "Invalid Email used to validate Txn entry: %s, txnId: %s", txn.GetReceiverEmail(), txn.GetTxnId())
				invalidTxns = append(invalidTxns, txn)
			}
		}
	}

	return validTxns, invalidTxns
}

func getAllActiveSubscriptionsTxns(ctx context.Context) ([]*TransactionDao.TransactionMsgDTO, error) {
	expiredTransactionDate := time.Now().AddDate(0, -subscriptionPeriodMonth, 0)
	return transactionDao.GetCurrentTransactionsAfter(ctx, expiredTransactionDate)
}

func getLatestTxnByUser(txns []*TransactionDao.TransactionMsgDTO) UserTransactionMap {
	userTransactionMap := make(UserTransactionMap)

	for _, txn := range txns {
		userKey := txn.GetUser()
		userTxns := userTransactionMap[userKey]

		if userTxns == nil {
			userTxns = make([]*TransactionDao.TransactionMsgDTO, 0, 2)
		}

		userTransactionMap[userKey] = append(userTxns, txn)
	}
	return userTransactionMap
}

func getUsers(ctx context.Context, usersTxnMap UserTransactionMap) ([]UserDao.UserDTO, error) {
	if len(usersTxnMap) == 0 {
		return nil, nil
	}
	userKeys := make([]*datastore.Key, 0, len(usersTxnMap))
	for userKey, _ := range usersTxnMap {
		userKeys = append(userKeys, userKey)
	}

	return userDAO.GetByKeys(ctx, userKeys)
}

func partitionUsersByValidity(ctx context.Context, users []UserDao.UserDTO) ([]UserDao.UserDTO, []UserDao.UserDTO, error) {
	if err := accessIdValidator.EnsureUpdatedIds(ctx); err != nil {
		return nil, nil, err
	}

	validUsers := make([]UserDao.UserDTO, 0, len(users))
	invalidUsers := make([]UserDao.UserDTO, 0, 10)
	for _, user := range users {
		isValid, err := accessIdValidator.ValidateAccessId(ctx, []byte(user.AccessId))
		if err != nil {
			return nil, nil, err
		}

		if isValid {
			validUsers = append(validUsers, user)
		} else {
			invalidUsers = append(invalidUsers, user)
		}
	}

	return validUsers, invalidUsers, nil
}

func mapToUserNames(users []UserDao.UserDTO) []string {
	var userNames = make([]string, len(users))
	for i, user := range users {
		userNames[i] = user.Name
	}
	return userNames
}

func mapTxnToDate(txns []*TransactionDao.TransactionMsgDTO) []time.Time {
	dates := make([]time.Time, len(txns))
	for i, txn := range txns {
		dates[i] = txn.GetPaymentDate()
	}

	return dates
}

func findMinMax(dates []time.Time) (time.Time, time.Time) {
	min := dates[0]
	max := dates[0]
	for _, date := range dates {
		if date.Before(min) {
			min = date
		} else if date.After(max) {
			max = date
		}
	}
	return min, max
}

func mapUsersToActivePeriod(validUsers []UserDao.UserDTO, usersWithActiveSubscriptions UserTransactionMap) []UserTxnTuple {
	usersWithPeroid := make([]UserTxnTuple, len(validUsers))

	for i, user := range validUsers {
		txnDates := mapTxnToDate(usersWithActiveSubscriptions[user.Key])
		min, max := findMinMax(txnDates)

		usersWithPeroid[i] = UserTxnTuple{
			user:      user,
			startDate: min,
			endDate:   max.AddDate(0, subscriptionPeriodMonth, 0),
		}
	}

	return usersWithPeroid
}

func getActiveTransactionList(ctx context.Context) ([]UserTxnTuple, error) {
	txns, err := getAllActiveSubscriptionsTxns(ctx)
	if err != nil {
		return nil, err
	}
	validTxns, _ := validateTransactions(ctx, txns)

	usersWithActiveSubscriptions := getLatestTxnByUser(validTxns)

	users, err := getUsers(ctx, usersWithActiveSubscriptions)
	if err != nil {
		return nil, err
	}

	validUsers, invalidUsers, err := partitionUsersByValidity(ctx, users)
	if err != nil {
		return nil, err
	}

	activeUsersWithPariod := mapUsersToActivePeriod(validUsers, usersWithActiveSubscriptions)

	if len(invalidUsers) > 0 {
		usersWithActiveSubscriptionButInvalidIdsStr := strings.Join(mapToUserNames(invalidUsers), ", ")
		log.Infof(ctx, "%s has paid for access but ID is not valid, skipped in csv export", usersWithActiveSubscriptionButInvalidIdsStr)
	}

	return activeUsersWithPariod, nil
}

func createCsvFile(ctx context.Context, w io.Writer) error {
	activeUsersWithPariod, err := getActiveTransactionList(ctx)
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

	for _, user := range activeUsersWithPariod {
		log.Infof(ctx, "%s, %s, %s", user.user.AccessId, user.startDate.String(), user.endDate.String())
		w.Write([]byte(user.user.AccessId))
		w.Write(comma)
		w.Write([]byte(user.startDate.Format(csvDateFormat)))
		w.Write(comma)
		w.Write([]byte(user.endDate.Format(csvDateFormat)))
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
