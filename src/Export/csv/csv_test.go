package csv

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"google.golang.org/appengine/datastore"

	"IPN/Transaction"
	"User/Dao"

	"AccessIdValidator/AccessIdValidatorTestHelper"
	"IPN/Transaction/TransactionDaoTestHelper"
	"TestHelper"
	"User/Dao/UserDaoTestHelper"
)

var (
	userDaoMock           *UserDaoTestHelper.UserRetrieverMock
	accessIdValidatorMock *AccessIdValidatorTestHelper.AccessIdValidatorMock
	txnDaoMock            *TransactionDaoTestHelper.TransactionRetrieverMock
)

var (
	windowsNewline = []byte{0x0D, 0x0A}

	assert = TestHelper.Assert
	utc, _ = time.LoadLocation("UTC")

	ctx = TestHelper.GetContext()
	now = time.Now()

	validIds = []string{
		"AccessId1",
		"AccessId2",
		"AccessId3",
		"AccessId4",
	}
)

func mockUserRetriever(keys []*datastore.Key, users []*UserDao.UserDTO, err error) *UserDaoTestHelper.UserRetrieverMock {
	userDaoMock = UserDaoTestHelper.NewUserRetrieverMock(keys, users, err)
	userDAO = userDaoMock
	return userDaoMock
}

func mockAccessIdValidator() *AccessIdValidatorTestHelper.AccessIdValidatorMock {
	accessIdValidatorMock = AccessIdValidatorTestHelper.NewAccessIdValidatorMock(validIds, nil)
	accessIdValidator = accessIdValidatorMock
	return accessIdValidatorMock
}

func mockTransactionRetriever(messages []*TransactionDao.TransactionMsgDTO, err error) *TransactionDaoTestHelper.TransactionRetrieverMock {
	txnDaoMock = TransactionDaoTestHelper.NewTransactionRetrieverMock(messages, err)
	transactionDao = txnDaoMock
	return txnDaoMock
}

func createMessages(dates []time.Time, userKeys []*datastore.Key) []*TransactionDao.TransactionMsgDTO {
	return createMessagesWithEmail(dates, userKeys, "gpmac_1231902686_biz@paypal.com")
}

func createMessagesWithEmail(dates []time.Time, userKeys []*datastore.Key, email string) []*TransactionDao.TransactionMsgDTO {
	const layout = "15:04:05 Jan 02, 2006 MST"
	messages := make([]*TransactionDao.TransactionMsgDTO, 0, 5)
	for i, date := range dates {
		txnKey := datastore.NewKey(ctx, "txn", "", int64(i)+1, userKeys[i])
		dateStr := date.In(utc).Format(layout)
		dateIpnMsg := TransactionDao.FIELD_PAYMENT_DATE + "=" + dateStr
		receiverEmail := TransactionDao.FIELD_RECEIVER_EMAIL + "=" + email
		dateTxnMsg := TransactionDao.NewTransactionMsgDTOFromIpnWithKey(dateIpnMsg+"&"+receiverEmail, txnKey)
		messages = append(messages, dateTxnMsg)
	}

	return messages
}

func createUsers(accessIds []string) ([]*datastore.Key, []*UserDao.UserDTO) {
	keys := make([]*datastore.Key, 0, 6)
	users := make([]*UserDao.UserDTO, 0, 5)
	for i, accessId := range accessIds {
		userKey := datastore.NewKey(ctx, "user", "", int64(i)+1, nil)
		users = append(users, &UserDao.UserDTO{AccessId: accessId, Key: userKey})
		keys = append(keys, userKey)
	}
	return keys, users
}

func convertDates(dates []time.Time) [][]string {
	dateStrs := make([][]string, 0, 8)
	for _, date := range dates {
		startTimeBStr := date.Format(csvDateFormat)
		endTimeBStr := date.AddDate(0, 6, 0).Format(csvDateFormat)
		dateStrs = append(dateStrs, []string{startTimeBStr, endTimeBStr})
	}
	return dateStrs
}

func TestShouldBeAbleToCreateAnEmptyCsvFileWithBom(t *testing.T) {
	buffer := &bytes.Buffer{}
	mockUserRetriever(nil, nil, nil)
	mockTransactionRetriever(nil, nil)

	assert(t, createCsvFile(ctx, buffer, nil)).Equals(nil)
	doc, _ := ioutil.ReadAll(buffer)

	txnDate := txnDaoMock.CallArgs[0].Date
	assert(t, txnDate.Before(now.AddDate(0, -6, 1)))
	assert(t, txnDate.After(now.AddDate(0, -6, -1)))
	assert(t, userDaoMock.CallCount).Equals(0)
	assert(t, doc).Equals([]byte{})
}

func TestShouldReturnPassErrorFromUserDaoThrough(t *testing.T) {
	err := errors.New("Some User Error")
	mockUserRetriever(nil, nil, err)
	mockAccessIdValidator()
	mockTransactionRetriever(nil, nil)

	assert(t, createCsvFile(ctx, &bytes.Buffer{}, nil)).Equals(err)
}

func TestShouldReturnPassErrorFromTxnDaoThrough(t *testing.T) {
	buffer := &bytes.Buffer{}
	err := errors.New("Some User Error")
	mockUserRetriever(nil, nil, nil)
	mockAccessIdValidator()
	mockTransactionRetriever(nil, err)

	assert(t, createCsvFile(ctx, buffer, nil)).Equals(err)
}

func TestShouldNotBeAbleToCreateCsvWhenTxnHasWrongEmail(t *testing.T) {
	buffer := &bytes.Buffer{}
	userKeys, users := createUsers([]string{"AccessId1"})
	mockUserRetriever(userKeys, users, nil)
	mockAccessIdValidator()
	txnDaoMock := mockTransactionRetriever(createMessagesWithEmail([]time.Time{now}, userKeys, "bad@email.com"), nil)

	assert(t, createCsvFile(ctx, buffer, nil)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(csvBytes)

	assert(t, csvString).Equals("")
	assert(t, userDaoMock.CallCount).Equals(0)
	assert(t, txnDaoMock.CallCount).Equals(1)
}

func TestShouldBeAbleToCreateCsvWithOneEntry(t *testing.T) {
	buffer := &bytes.Buffer{}
	userKeys, users := createUsers([]string{"AccessId1"})
	mockUserRetriever(userKeys, users, nil)
	mockAccessIdValidator()
	txnDaoMock := mockTransactionRetriever(createMessages([]time.Time{now}, userKeys), nil)

	assert(t, createCsvFile(ctx, buffer, nil)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(csvBytes)

	startTimeStr := now.Format(csvDateFormat)
	endTimeStr := now.AddDate(0, 6, 0).Format(csvDateFormat)
	assert(t, csvString).Equals(fmt.Sprintf("%s,%s,%s%s", users[0].AccessId, startTimeStr, endTimeStr, windowsNewline))
	assert(t, accessIdValidatorMock.CallCount).Equals(1)
	assert(t, txnDaoMock.CallCount).Equals(1)
}

func TestShouldBeAbleToCreateCsvWithTwoEntries(t *testing.T) {
	buffer := &bytes.Buffer{}
	userKeys, users := createUsers([]string{"AccessId1", "AccessId2"})
	userDaoMock := mockUserRetriever(userKeys, users, nil)
	mockAccessIdValidator()

	plusFive := now.AddDate(0, 0, 5)

	mockTransactionRetriever(createMessages([]time.Time{now, plusFive}, userKeys), nil)

	assert(t, createCsvFile(ctx, buffer, nil)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(csvBytes)

	dateStrs := convertDates([]time.Time{now, plusFive})

	assert(t, strings.Contains(csvString, fmt.Sprintf("%s,%s,%s%s", users[0].AccessId, dateStrs[0][0], dateStrs[0][1], windowsNewline)))
	assert(t, strings.Contains(csvString, fmt.Sprintf("%s,%s,%s%s", users[1].AccessId, dateStrs[1][0], dateStrs[1][1], windowsNewline)))
	assert(t, strings.Count(csvString, "\n")).Equals(2)

	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, accessIdValidatorMock.CallCount).Equals(2)
}

func TestShouldNotIncludeUsersWithInvalidAccessIds(t *testing.T) {
	buffer := &bytes.Buffer{}
	userKeys, users := createUsers([]string{"AccessId1", "InvalidId"})
	userDaoMock := mockUserRetriever(userKeys, users, nil)
	mockAccessIdValidator()

	plusFive := now.AddDate(0, 0, 5)

	mockTransactionRetriever(createMessages([]time.Time{now}, userKeys[:1]), nil).
		AddReturn(createMessages([]time.Time{plusFive}, userKeys[1:2]), nil)

	assert(t, createCsvFile(ctx, buffer, nil)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(csvBytes)

	dateStrs := convertDates([]time.Time{now})

	assert(t, csvString).Equals(
		fmt.Sprintf("%s,%s,%s%s", users[0].AccessId, dateStrs[0][0], dateStrs[0][1], windowsNewline))
	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, accessIdValidatorMock.CallCount).Equals(2)
}

func TestShouldBeAbleToCreateCsvWithTreeEntries(t *testing.T) {
	buffer := &bytes.Buffer{}
	userKeys, users := createUsers([]string{"AccessId1", "AccessId2", "AccessId3"})
	userDaoMock := mockUserRetriever(userKeys, users, nil)
	mockAccessIdValidator()

	dates := []time.Time{
		now,
		now.AddDate(0, 0, 5),
		now.AddDate(0, 0, 6),
	}

	mockTransactionRetriever(createMessages(dates, userKeys), nil)

	assert(t, createCsvFile(ctx, buffer, nil)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(csvBytes)

	dateStrs := convertDates(dates)

	assert(t, strings.Contains(csvString, fmt.Sprintf("%s,%s,%s%s", users[0].AccessId, dateStrs[0][0], dateStrs[0][1], windowsNewline)))
	assert(t, strings.Contains(csvString, fmt.Sprintf("%s,%s,%s%s", users[1].AccessId, dateStrs[1][0], dateStrs[1][1], windowsNewline)))
	assert(t, strings.Contains(csvString, fmt.Sprintf("%s,%s,%s%s", users[2].AccessId, dateStrs[2][0], dateStrs[2][1], windowsNewline)))
	assert(t, strings.Count(csvString, "\n")).Equals(3)

	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, accessIdValidatorMock.CallCount).Equals(3)
}

func TestShouldBeAbleToCreateCsvWithMultipleTxnEntries(t *testing.T) {
	buffer := &bytes.Buffer{}
	userKeys, users := createUsers([]string{"AccessId1", "AccessId2", "AccessId3"})
	userDaoMock := mockUserRetriever(userKeys, users, nil)
	mockAccessIdValidator()

	dates := []time.Time{
		now,                   // User 1
		now.AddDate(0, 0, 5),  // User 2 A
		now.AddDate(0, 6, 0),  // User 2 B
		now.AddDate(0, 9, 6),  // User 2 C
		now.AddDate(1, 0, 6),  // User 2 extremeHigh
		now.AddDate(0, 0, -3), // User 2 extremeLow
	}

	txnUserKeys := []*datastore.Key{
		userKeys[0],
		userKeys[1], userKeys[1], userKeys[1], userKeys[1], userKeys[1],
	}

	mockTransactionRetriever(createMessages(dates, txnUserKeys), nil)
	newTxn := createMessages([]time.Time{now}, []*datastore.Key{userKeys[2]})[0]

	assert(t, createCsvFile(ctx, buffer, newTxn)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(csvBytes)

	dateStrs := convertDates(dates)

	assert(t, strings.Contains(csvString, fmt.Sprintf("%s,%s,%s%s", users[0].AccessId, dateStrs[0][0], dateStrs[0][1], windowsNewline)))
	assert(t, strings.Contains(csvString, fmt.Sprintf("%s,%s,%s%s", users[1].AccessId, dateStrs[5][0], dateStrs[4][1], windowsNewline)))
	assert(t, strings.Contains(csvString, fmt.Sprintf("%s,%s,%s%s", users[2].AccessId, dateStrs[0][0], dateStrs[0][1], windowsNewline)))
	assert(t, strings.Count(csvString, "\n")).Equals(3)

	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, accessIdValidatorMock.CallCount).Equals(3)
}

func TestShouldNotCreateCsvWithDoublicatedAccessIdEntries(t *testing.T) {
	buffer := &bytes.Buffer{}
	userKeys, users := createUsers([]string{"AccessId1"})
	userDaoMock := mockUserRetriever(userKeys, users, nil)
	mockAccessIdValidator()

	dates := []time.Time{
		now,
		now.AddDate(0, 0, -5),
		now.AddDate(0, -6, -0),
		now.AddDate(0, -9, -6),
		now.AddDate(-1, 0, -6),
		now.AddDate(0, 0, 3),
	}

	txnUserKeys := make([]*datastore.Key, len(dates))
	for i, _ := range dates {
		txnUserKeys[i] = userKeys[0]
	}

	mockTransactionRetriever(createMessages(dates, txnUserKeys), nil)
	newTxn := createMessages([]time.Time{now}, []*datastore.Key{userKeys[0]})[0]

	assert(t, createCsvFile(ctx, buffer, newTxn)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(csvBytes)

	dateStrs := convertDates(dates)

	assert(t, csvString).Equals(
		fmt.Sprintf("%s,%s,%s%s", users[0].AccessId, dateStrs[4][0], dateStrs[5][1], windowsNewline))
	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, accessIdValidatorMock.CallCount).Equals(1)
}
