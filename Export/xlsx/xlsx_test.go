package xlsx

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"appengine"
	"appengine/datastore"
	"github.com/tealeg/xlsx"

	"TestHelper"
	"User/Dao"
)

var assert = TestHelper.Assert

func mockUserHasActiveSubscription(firstDate []time.Time, lastDate []time.Time, err []error) *TestHelper.Spy {
	spy := new(TestHelper.Spy)
	transactionDao_GetCurrentTransactionsAfter = func(ctx appengine.Context, userKey *datastore.Key, date time.Time) (time.Time, time.Time, error) {
		spy.RegisterCall()
		spy.RegisterArg3(ctx, userKey, date)
		_firstDate := firstDate[0]
		_lastDate := lastDate[0]
		_err := err[0]
		if len(firstDate) > 1 {
			firstDate = firstDate[1:]
		}
		if len(lastDate) > 1 {
			lastDate = lastDate[1:]
		}
		if len(err) > 1 {
			err = err[1:]
		}
		return _firstDate, _lastDate, _err
	}
	return spy
}

func TestShouldConfigureHeaderForDownload(t *testing.T) {
	header := make(http.Header)
	configureHeaderForFileDownload(&header, "test.file")

	assert(t, header.Get("Content-Disposition") == "attachment; filename=test.file")
	assert(t, header.Get("Content-Type") == "application/vnd.ms-excel")
}

func TestShouldConfigureHeaderForNoCache(t *testing.T) {
	header := make(http.Header)
	configureHeaderForFileDownload(&header, "test.file")

	assert(t, header.Get("Cache-Control") == "no-cache, no-store, must-revalidate")
	assert(t, header.Get("Pragma") == "no-cache")
	assert(t, header.Get("Expires") == "0")
}

func TestShouldGetTransactionsFromDataStore(t *testing.T) {
	ctx := &TestHelper.ContextMock{}

	keys := []*datastore.Key{
		{},
		{},
	}
	users := []UserDao.UserDTO{
		{Email: "NO_Subscription"},
		{Email: "hasSubscription"},
	}

	now := time.Now()
	invalid := time.Time{}
	userDaoMock := mockUserRetriever(keys, users, nil)
	spyHasActiveSub := mockUserHasActiveSubscription([]time.Time{invalid, now, now}, []time.Time{invalid, now, now}, []error{nil, nil, nil})

	userTxnTuple, err := getActiveTransactionList(ctx)

	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, userDaoMock.LatestCallCtxArg).Equals(ctx)
	assert(t, len(userTxnTuple)).Equals(1)
	assert(t, userTxnTuple[0].user).Equals(users[1])
	assert(t, err).Equals(nil)

	assert(t, spyHasActiveSub.CallCount()).Equals(2)
}

func TestShouldPassOnErrorsFromDataStore_GetAllUsers(t *testing.T) {
	ctx := &TestHelper.ContextMock{}
	testError := errors.New("test error")

	userDaoMock := mockUserRetriever(nil, nil, testError)

	userTxnTuple, err := getActiveTransactionList(ctx)

	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, userDaoMock.LatestCallCtxArg).Equals(ctx)
	assert(t, userTxnTuple).Equals(nil)
	assert(t, err).Equals(testError)
}

func TestShouldPassOnErrorsFromDataStore_HasSubscription(t *testing.T) {
	ctx := &TestHelper.ContextMock{}
	testError := errors.New("test error")

	keys := []*datastore.Key{{}}
	users := []UserDao.UserDTO{{}}

	now := time.Now()

	mockUserRetriever(keys, users, nil)
	mockUserHasActiveSubscription([]time.Time{now}, []time.Time{now}, []error{testError})

	userTxnTuple, err := getActiveTransactionList(ctx)

	assert(t, userTxnTuple).Equals(nil)
	assert(t, err).Equals(testError)
}

func TestShouldAddHeaderRowBasedOnPassedArgumentsToAddRow(t *testing.T) {
	col1 := "test1"
	col2 := "test2"
	sheet := xlsx.Sheet{}

	addXlsxRow(&sheet, col1, col2)

	assert(t, len(sheet.Rows)).Equals(1)
	assert(t, sheet.Rows[0].Cells[0].Value).Equals(col1)
	assert(t, sheet.Rows[0].Cells[1].Value).Equals(col2)
}

func TestShouldCreateXlsxSheetWithAllUserHavingActiveSubscription(t *testing.T) {
	ctx := &TestHelper.ContextMock{}

	keys := []*datastore.Key{{}}
	users := []UserDao.UserDTO{{Email: "testMail"}}

	now := time.Now()

	mockUserRetriever(keys, users, nil)
	mockUserHasActiveSubscription([]time.Time{now}, []time.Time{now}, []error{nil})

	file, error := createXlsxFile(ctx)

	headerCells := file.Sheets[0].Rows[0].Cells
	firstRowCells := file.Sheets[0].Rows[1].Cells

	assert(t, len(headerCells)).Equals(7)
	assert(t, len(firstRowCells)).Equals(7)

	assert(t, firstRowCells[0].Value).Equals(users[0].AccessId)
	assert(t, firstRowCells[1].Value).Equals(now.Format(xlsxDateFormat))
	assert(t, firstRowCells[2].Value).Equals(users[0].AccessId)
	assert(t, firstRowCells[3].Value).Equals(now.Format(xlsxDateFormat))
	assert(t, firstRowCells[4].Value).Equals(now.AddDate(0, 6, 0).Format(xlsxDateFormat))
	assert(t, firstRowCells[5].Value).Equals("24 Timers")
	assert(t, firstRowCells[6].Value).Equals(users[0].Email)
	assert(t, error).Equals(nil)
}

func TestExportXlsxHandler(t *testing.T) {
	ctx := new(TestHelper.ContextMock)
	testError := errors.New("test error")

	mockUserRetriever(nil, nil, testError)

	file, err := createXlsxFile(ctx)

	assert(t, file).Equals(nil)
	assert(t, err).Equals(testError)
}
