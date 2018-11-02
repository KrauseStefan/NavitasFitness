package csv

import (
	"bytes"
	"errors"
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

type UserSubscriptionInfo struct {
	userKey   *datastore.Key
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

type UserTransactionMap map[datastore.Key][]*TransactionDao.TransactionMsgDTO

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
		userKey := *txn.GetUser()
		userTxns := userTransactionMap[userKey]

		if userTxns == nil {
			userTxns = make([]*TransactionDao.TransactionMsgDTO, 0, 2)
		}

		userTransactionMap[userKey] = append(userTxns, txn)
	}
	return userTransactionMap
}

func getUsers(ctx context.Context, usersTxnMap UserTransactionMap) ([]*UserDao.UserDTO, []*TransactionDao.TransactionMsgDTO, error) {
	if len(usersTxnMap) == 0 {
		return nil, nil, nil
	}
	userKeys := make([]*datastore.Key, 0, len(usersTxnMap))
	for userKey, _ := range usersTxnMap {
		userKeys = append(userKeys, &userKey)
	}

	users, err := userDAO.GetByKeys(ctx, userKeys)
	if multiErr, ok := err.(appengine.MultiError); ok {
		txnsWithNoUser := make([]*TransactionDao.TransactionMsgDTO, 0, len(usersTxnMap))
		foundUsers := make([]*UserDao.UserDTO, 0, len(usersTxnMap))
		foundOtherError := false
		for i, e := range multiErr {
			if e == nil {
				foundUsers = append(foundUsers, users[i])
			} else if e == datastore.ErrNoSuchEntity {
				txnsWithNoUser = append(txnsWithNoUser, usersTxnMap[*userKeys[i]]...)
			} else {
				foundOtherError = true
				return nil, nil, err
			}
		}

		if foundOtherError {
			return foundUsers, txnsWithNoUser, err
		} else {
			return foundUsers, txnsWithNoUser, nil
		}
	}

	return users, nil, err
}

func partitionUsersByValidity(ctx context.Context, users []*UserDao.UserDTO) ([]*UserDao.UserDTO, []*UserDao.UserDTO, error) {
	if err := accessIdValidator.EnsureUpdatedIds(ctx); err != nil {
		return nil, nil, err
	}

	validUsers := make([]*UserDao.UserDTO, 0, len(users))
	invalidUsers := make([]*UserDao.UserDTO, 0, 10)
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

func mapToUserNames(users []*UserDao.UserDTO) []string {
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

func mapUsersToActivePeriod(ctx context.Context, validUsers []*UserDao.UserDTO, usersWithActiveSubscriptions UserTransactionMap) map[string]*UserSubscriptionInfo {
	usersWithPeroid := make(map[string]*UserSubscriptionInfo)

	for _, user := range validUsers {
		txnDates := mapTxnToDate(usersWithActiveSubscriptions[*user.Key])
		min, max := findMinMax(txnDates)

		newUserSubscriptionInfo := &UserSubscriptionInfo{
			userKey:   user.Key,
			startDate: min,
			endDate:   max.AddDate(0, subscriptionPeriodMonth, 0),
		}

		prevUserSubscriptionInfo := usersWithPeroid[user.AccessId]
		if prevUserSubscriptionInfo == nil {
			usersWithPeroid[user.AccessId] = newUserSubscriptionInfo
		} else {
			log.Errorf(ctx, "Doublicated accessId detected %s, key1: %s, key2: %s", user.AccessId, user.Key, prevUserSubscriptionInfo.userKey)
			if newUserSubscriptionInfo.endDate.After(prevUserSubscriptionInfo.endDate) {
				usersWithPeroid[user.AccessId] = newUserSubscriptionInfo
			}
		}
	}

	return usersWithPeroid
}

func getActiveTransactionList(ctx context.Context, newTxn *TransactionDao.TransactionMsgDTO) (map[string]*UserSubscriptionInfo, error) {
	txns, err := getAllActiveSubscriptionsTxns(ctx)
	if err != nil {
		return nil, err
	}
	if newTxn != nil {
		txns = append(txns, newTxn)
	}
	validTxns, _ := validateTransactions(ctx, txns)

	usersWithActiveSubscriptions := getLatestTxnByUser(validTxns)
	users, badTxns, err := getUsers(ctx, usersWithActiveSubscriptions)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		log.Errorf(ctx, "No users cannot generate csv")
	}
	if len(badTxns) > 0 {
		log.Errorf(ctx, "Txns without user: %d", len(badTxns))
		badTxnIds := make([]string, len(badTxns))
		for i, txn := range badTxns {
			badTxnIds[i] = txn.GetTxnId() + " - " + txn.GetPayerEmail()
		}

		log.Errorf(ctx, "Bad transactions id - email: %s", strings.Join(badTxnIds, ", "))
	}

	validUsers, invalidUsers, err := partitionUsersByValidity(ctx, users)
	if err != nil {
		return nil, err
	}

	activeUsersWithPariod := mapUsersToActivePeriod(ctx, validUsers, usersWithActiveSubscriptions)

	if len(invalidUsers) > 0 {
		usersWithActiveSubscriptionButInvalidIdsStr := strings.Join(mapToUserNames(invalidUsers), ", ")
		log.Infof(ctx, "%s has paid for access but ID is not valid, skipped in csv export", usersWithActiveSubscriptionButInvalidIdsStr)
	}

	return activeUsersWithPariod, nil
}

func createCsvFile(ctx context.Context, w io.Writer, newTxn *TransactionDao.TransactionMsgDTO) error {
	activeUsersWithPariod, err := getActiveTransactionList(ctx, newTxn)
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

	for accessId, userInfo := range activeUsersWithPariod {
		log.Infof(ctx, "%s, %s, %s", accessId, userInfo.startDate.String(), userInfo.endDate.String())
		w.Write([]byte(accessId))
		w.Write(comma)
		w.Write([]byte(userInfo.startDate.Format(csvDateFormat)))
		w.Write(comma)
		w.Write([]byte(userInfo.endDate.Format(csvDateFormat)))
		w.Write([]byte(windowsNewline))
	}

	return nil
}

func exportCsvHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
	ctx := appengine.NewContext(r)

	err := CreateAndUploadFile(ctx, nil)
	return nil, err
}

func CreateAndUploadFile(ctx context.Context, newTxn *TransactionDao.TransactionMsgDTO) error {
	var buffer bytes.Buffer

	tokens, err := Dropbox.GetAccessTokens(ctx)
	if err != nil {
		return err
	}

	if err := createCsvFile(ctx, &buffer, newTxn); err != nil {
		return errors.New("Error generating CSV file: " + err.Error())
	}

	for _, token := range tokens {
		if _, err := Dropbox.UploadDoc(ctx, token, AccessIdValidator.GetAccessListPath(ctx), &buffer); err != nil {
			return err
		}
	}

	return nil
}
