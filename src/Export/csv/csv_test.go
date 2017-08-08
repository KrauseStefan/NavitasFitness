package csv

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
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

func mockUserRetriever(keys []*datastore.Key, users []UserDao.UserDTO, err error) *UserDaoTestHelper.UserRetrieverMock {
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

func createMessages(dates []time.Time) []*TransactionDao.TransactionMsgDTO {
	return createMessagesWithEmail(dates, "gpmac_1231902686_biz@paypal.com")
}

func createMessagesWithEmail(dates []time.Time, email string) []*TransactionDao.TransactionMsgDTO {
	const layout = "15:04:05 Jan 02, 2006 MST"
	messages := make([]*TransactionDao.TransactionMsgDTO, 0, 5)
	for _, date := range dates {
		dateStr := date.In(utc).Format(layout)
		dateIpnMsg := TransactionDao.FIELD_PAYMENT_DATE + "=" + dateStr
		receiverEmail := TransactionDao.FIELD_RECEIVER_EMAIL + "=" + email
		dateTxnMsg := TransactionDao.NewTransactionMsgDTOFromIpn(dateIpnMsg + "&" + receiverEmail)
		messages = append(messages, dateTxnMsg)
	}

	return messages
}

func createUsers(accessIds []string) ([]*datastore.Key, []UserDao.UserDTO) {
	keys := make([]*datastore.Key, 0, 6)
	users := make([]UserDao.UserDTO, 0, 5)
	for _, accessId := range accessIds {
		users = append(users, UserDao.UserDTO{AccessId: accessId})
		keys = append(keys, &datastore.Key{})
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
	//mockAccessIdValidator()
	mockTransactionRetriever(nil, nil)

	assert(t, createCsvFile(ctx, buffer)).Equals(nil)
	doc, _ := ioutil.ReadAll(buffer)

	assert(t, doc).Equals([]byte{})
	//assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, txnDaoMock.CallCount).Equals(0)
	assert(t, userDaoMock.LatestCallCtxArg).Equals(ctx)
}

func TestShouldReturnPassErrorFromUserDaoThrough(t *testing.T) {
	err := errors.New("Some User Error")
	mockUserRetriever(nil, nil, err)
	mockAccessIdValidator()
	mockTransactionRetriever(nil, nil)

	assert(t, createCsvFile(ctx, &bytes.Buffer{})).Equals(err)

	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, txnDaoMock.CallCount).Equals(0)
}

func TestShouldCreateEmptyCsvIfNoUserHasAnyTxn(t *testing.T) {
	buffer := &bytes.Buffer{}
	keys, users := createUsers([]string{"AccessId1"})
	mockUserRetriever(keys, users, nil)
	mockAccessIdValidator()
	mockTransactionRetriever(nil, nil)

	assert(t, createCsvFile(ctx, buffer)).Equals(nil)
	csv, _ := ioutil.ReadAll(buffer)

	assert(t, csv).Equals([]byte{})
	assert(t, accessIdValidatorMock.CallCount).Equals(1)
	assert(t, txnDaoMock.CallCount).Equals(1)
}

func TestShouldReturnPassErrorFromTxnDaoThrough(t *testing.T) {
	buffer := &bytes.Buffer{}
	err := errors.New("Some User Error")
	keys, users := createUsers([]string{"AccessId1"})
	mockUserRetriever(keys, users, nil)
	mockAccessIdValidator()
	txnDaoMock := mockTransactionRetriever(nil, err)

	assert(t, createCsvFile(ctx, buffer)).Equals(err)

	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, accessIdValidatorMock.CallCount).Equals(1)
	assert(t, txnDaoMock.CallCount).Equals(1)
}

func TestShouldNotBeAbleToCreateCsvWhenTxnHasWrongEmail(t *testing.T) {
	buffer := &bytes.Buffer{}
	keys, users := createUsers([]string{"AccessId1"})
	mockUserRetriever(keys, users, nil)
	mockAccessIdValidator()
	txnDaoMock := mockTransactionRetriever(createMessagesWithEmail([]time.Time{now}, "bad@email.com"), nil)

	assert(t, createCsvFile(ctx, buffer)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(csvBytes)

	assert(t, csvString).Equals("")
	assert(t, accessIdValidatorMock.CallCount).Equals(1)
	assert(t, txnDaoMock.CallCount).Equals(1)
}

func TestShouldBeAbleToCreateCsvWithOneEntry(t *testing.T) {
	buffer := &bytes.Buffer{}
	keys, users := createUsers([]string{"AccessId1"})
	mockUserRetriever(keys, users, nil)
	mockAccessIdValidator()
	txnDaoMock := mockTransactionRetriever(createMessages([]time.Time{now}), nil)

	assert(t, createCsvFile(ctx, buffer)).Equals(nil)
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
	keys, users := createUsers([]string{"AccessId1", "AccessId2"})
	userDaoMock := mockUserRetriever(keys, users, nil)
	mockAccessIdValidator()

	plusFive := now.AddDate(0, 0, 5)

	txnDaoMock := mockTransactionRetriever(createMessages([]time.Time{now}), nil).
		AddReturn(createMessages([]time.Time{plusFive}), nil)

	assert(t, createCsvFile(ctx, buffer)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(csvBytes)

	dateStrs := convertDates([]time.Time{now, plusFive})

	assert(t, csvString).Equals(
		fmt.Sprintf("%s,%s,%s%s", users[0].AccessId, dateStrs[0][0], dateStrs[0][1], windowsNewline) +
			fmt.Sprintf("%s,%s,%s%s", users[1].AccessId, dateStrs[1][0], dateStrs[1][1], windowsNewline))
	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, accessIdValidatorMock.CallCount).Equals(2)
	assert(t, txnDaoMock.CallCount).Equals(2)
}

func TestShouldNotIncludeUsersWithInvalidAccessIds(t *testing.T) {
	buffer := &bytes.Buffer{}
	keys, users := createUsers([]string{"AccessId1", "InvalidId"})
	userDaoMock := mockUserRetriever(keys, users, nil)
	mockAccessIdValidator()

	txnDaoMock := mockTransactionRetriever(createMessages([]time.Time{now}), nil)

	assert(t, createCsvFile(ctx, buffer)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(csvBytes)

	dateStrs := convertDates([]time.Time{now})

	assert(t, csvString).Equals(
		fmt.Sprintf("%s,%s,%s%s", users[0].AccessId, dateStrs[0][0], dateStrs[0][1], windowsNewline))
	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, accessIdValidatorMock.CallCount).Equals(2)
	assert(t, txnDaoMock.CallCount).Equals(1)
}

func TestShouldBeAbleToCreateCsvWithTreeEntries(t *testing.T) {
	buffer := &bytes.Buffer{}
	keys, users := createUsers([]string{"AccessId1", "AccessId2", "AccessId3"})
	userDaoMock := mockUserRetriever(keys, users, nil)
	mockAccessIdValidator()

	plusFive := now.AddDate(0, 0, 5)
	plusFour := now.AddDate(0, 0, 6)

	txnDaoMock := mockTransactionRetriever(createMessages([]time.Time{now}), nil).
		AddReturn(createMessages([]time.Time{plusFive}), nil).
		AddReturn(createMessages([]time.Time{plusFour}), nil)

	assert(t, createCsvFile(ctx, buffer)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(csvBytes)

	dateStrs := convertDates([]time.Time{now, plusFive, plusFour})

	assert(t, csvString).Equals(
		fmt.Sprintf("%s,%s,%s%s", users[0].AccessId, dateStrs[0][0], dateStrs[0][1], windowsNewline) +
			fmt.Sprintf("%s,%s,%s%s", users[1].AccessId, dateStrs[1][0], dateStrs[1][1], windowsNewline) +
			fmt.Sprintf("%s,%s,%s%s", users[2].AccessId, dateStrs[2][0], dateStrs[2][1], windowsNewline))
	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, accessIdValidatorMock.CallCount).Equals(3)
	assert(t, txnDaoMock.CallCount).Equals(3)
}

func TestShouldBeAbleToCreateCsvWithMultipleTxnEntries(t *testing.T) {
	buffer := &bytes.Buffer{}
	keys, users := createUsers([]string{"AccessId1", "AccessId2"})
	userDaoMock := mockUserRetriever(keys, users, nil)
	mockAccessIdValidator()

	plusA := now.AddDate(0, 0, 5)
	plusB := now.AddDate(0, 6, 0)
	plusC := now.AddDate(0, 9, 6)
	extremeHigh := now.AddDate(1, 0, 6)
	extremeLow := now.AddDate(0, 0, -3)

	txnDaoMock := mockTransactionRetriever(createMessages([]time.Time{now}), nil).
		AddReturn(createMessages([]time.Time{plusA, plusB, plusC, extremeHigh, extremeLow}), nil)

	assert(t, createCsvFile(ctx, buffer)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(csvBytes)

	dateStrs := convertDates([]time.Time{now, plusA, plusB, plusC, extremeHigh, extremeLow})

	assert(t, csvString).Equals(
		fmt.Sprintf("%s,%s,%s%s", users[0].AccessId, dateStrs[0][0], dateStrs[0][1], windowsNewline) +
			fmt.Sprintf("%s,%s,%s%s", users[1].AccessId, dateStrs[5][0], dateStrs[4][1], windowsNewline))
	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, accessIdValidatorMock.CallCount).Equals(2)
	assert(t, txnDaoMock.CallCount).Equals(2)
}
