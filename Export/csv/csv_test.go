package csv

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"appengine/datastore"

	"IPN/Transaction"
	"User/Dao"

	"IPN/Transaction/TransactionDaoTestHelper"
	"TestHelper"
	"User/Dao/UserDaoTestHelper"
)

var (
	bom            = []byte{0xef, 0xbb, 0xbf}
	windowsNewline = []byte{0x0D, 0x0A}

	assert = TestHelper.Assert
	utc, _ = time.LoadLocation("UTC")

	ctx = &TestHelper.ContextMock{}
	now = time.Now()
)

func mockUserRetriever(keys []*datastore.Key, users []UserDao.UserDTO, err error) *UserDaoTestHelper.UserRetrieverMock {
	mock := UserDaoTestHelper.NewUserRetrieverMock(keys, users, err)
	userDAO = mock
	return mock
}

func mockTransactionRetriever(messages []*TransactionDao.TransactionMsgDTO, err error) *TransactionDaoTestHelper.TransactionRetrieverMock {
	mock := TransactionDaoTestHelper.NewTransactionRetrieverMock(messages, err)
	transactionDao = mock
	return mock
}

func createMessages(dates []time.Time) []*TransactionDao.TransactionMsgDTO {
	const layout = "15:04:05 Jan 02, 2006 MST"
	messages := make([]*TransactionDao.TransactionMsgDTO, 0, 5)
	for _, date := range dates {
		dateStr := date.In(utc).Format(layout)
		dateIpnMsg := TransactionDao.FIELD_PAYMENT_DATE + "=" + dateStr
		dateTxnMsg := TransactionDao.NewTransactionMsgDTOFromIpn(dateIpnMsg)
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
	userDaoMock := mockUserRetriever(nil, nil, nil)
	txnDaoMock := mockTransactionRetriever(nil, nil)

	assert(t, createCsvFile(ctx, buffer)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)

	assert(t, csvBytes).Equals(bom)
	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, txnDaoMock.CallCount).Equals(0)
	assert(t, userDaoMock.LatestCallCtxArg).Equals(ctx)
}

func TestShouldReturnPassErrorFromUserDaoThrough(t *testing.T) {
	err := errors.New("Some User Error")
	userDaoMock := mockUserRetriever(nil, nil, err)
	txnDaoMock := mockTransactionRetriever(nil, nil)

	assert(t, createCsvFile(ctx, &bytes.Buffer{})).Equals(err)

	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, txnDaoMock.CallCount).Equals(0)
}

func TestShouldCreateEmptyCsvIfNoUserHasAnyTxn(t *testing.T) {
	buffer := &bytes.Buffer{}
	keys, users := createUsers([]string{""})
	mockUserRetriever(keys, users, nil)
	txnDaoMock := mockTransactionRetriever(nil, nil)

	assert(t, createCsvFile(ctx, buffer)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)

	assert(t, csvBytes).Equals(bom)
	assert(t, txnDaoMock.CallCount).Equals(1)
}

func TestShouldReturnPassErrorFromTxnDaoThrough(t *testing.T) {
	buffer := &bytes.Buffer{}
	err := errors.New("Some User Error")
	userDaoMock := mockUserRetriever([]*datastore.Key{{}}, []UserDao.UserDTO{{}}, nil)
	txnDaoMock := mockTransactionRetriever(nil, err)

	assert(t, createCsvFile(ctx, buffer)).Equals(err)

	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, txnDaoMock.CallCount).Equals(1)
}

func TestShouldBeAbleToCreateCsvWithOneEntry(t *testing.T) {
	buffer := &bytes.Buffer{}
	keys, users := createUsers([]string{"1234"})
	mockUserRetriever(keys, users, nil)
	txnDaoMock := mockTransactionRetriever(createMessages([]time.Time{now}), nil)

	assert(t, createCsvFile(ctx, buffer)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(bytes.TrimPrefix(csvBytes, bom))

	startTimeStr := now.Format(csvDateFormat)
	endTimeStr := now.AddDate(0, 6, 0).Format(csvDateFormat)
	assert(t, csvString).Equals(fmt.Sprintf("%s,%s,%s", users[0].AccessId, startTimeStr, endTimeStr))
	assert(t, txnDaoMock.CallCount).Equals(1)
}

func TestShouldBeAbleToCreateCsvWithTwoEntries(t *testing.T) {
	buffer := &bytes.Buffer{}
	keys, users := createUsers([]string{"A", "B"})
	userDaoMock := mockUserRetriever(keys, users, nil)

	plusFive := now.AddDate(0, 0, 5)

	txnDaoMock := mockTransactionRetriever(createMessages([]time.Time{now}), nil).
		AddReturn(createMessages([]time.Time{plusFive}), nil)

	assert(t, createCsvFile(ctx, buffer)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(bytes.TrimPrefix(csvBytes, bom))

	dateStrs := convertDates([]time.Time{now, plusFive})

	assert(t, csvString).Equals(
		fmt.Sprintf("%s,%s,%s%s", users[0].AccessId, dateStrs[0][0], dateStrs[0][1], windowsNewline) +
			fmt.Sprintf("%s,%s,%s", users[1].AccessId, dateStrs[1][0], dateStrs[1][1]))
	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, txnDaoMock.CallCount).Equals(2)
}

func TestShouldBeAbleToCreateCsvWithTreeEntries(t *testing.T) {
	buffer := &bytes.Buffer{}
	keys, users := createUsers([]string{"A", "B", "C"})
	userDaoMock := mockUserRetriever(keys, users, nil)

	plusFive := now.AddDate(0, 0, 5)
	plusFour := now.AddDate(0, 0, 6)

	txnDaoMock := mockTransactionRetriever(createMessages([]time.Time{now}), nil).
		AddReturn(createMessages([]time.Time{plusFive}), nil).
		AddReturn(createMessages([]time.Time{plusFour}), nil)

	assert(t, createCsvFile(ctx, buffer)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(bytes.TrimPrefix(csvBytes, bom))

	dateStrs := convertDates([]time.Time{now, plusFive, plusFour})

	assert(t, csvString).Equals(
		fmt.Sprintf("%s,%s,%s%s", users[0].AccessId, dateStrs[0][0], dateStrs[0][1], windowsNewline) +
			fmt.Sprintf("%s,%s,%s%s", users[1].AccessId, dateStrs[1][0], dateStrs[1][1], windowsNewline) +
			fmt.Sprintf("%s,%s,%s", users[2].AccessId, dateStrs[2][0], dateStrs[2][1]))
	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, txnDaoMock.CallCount).Equals(3)
}

func TestShouldBeAbleToCreateCsvWithMultipleTxnEntries(t *testing.T) {
	buffer := &bytes.Buffer{}
	keys, users := createUsers([]string{"A", "B"})
	userDaoMock := mockUserRetriever(keys, users, nil)

	plusA := now.AddDate(0, 0, 5)
	plusB := now.AddDate(0, 6, 0)
	plusC := now.AddDate(0, 9, 6)
	extremeHigh := now.AddDate(1, 0, 6)
	extremeLow := now.AddDate(0, 0, -3)

	txnDaoMock := mockTransactionRetriever(createMessages([]time.Time{now}), nil).
		AddReturn(createMessages([]time.Time{plusA, plusB, plusC, extremeHigh, extremeLow}), nil)

	assert(t, createCsvFile(ctx, buffer)).Equals(nil)
	csvBytes, _ := ioutil.ReadAll(buffer)
	csvString := string(bytes.TrimPrefix(csvBytes, bom))

	dateStrs := convertDates([]time.Time{now, plusA, plusB, plusC, extremeHigh, extremeLow})

	assert(t, csvString).Equals(
		fmt.Sprintf("%s,%s,%s%s", users[0].AccessId, dateStrs[0][0], dateStrs[0][1], windowsNewline) +
			fmt.Sprintf("%s,%s,%s", users[1].AccessId, dateStrs[5][0], dateStrs[4][1]))
	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, txnDaoMock.CallCount).Equals(2)
}
