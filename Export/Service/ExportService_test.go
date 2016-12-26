package ExportService

import (
	"net/http"
	"testing"

	"appengine"
	"appengine/datastore"

	"TestHelper"
	"User/Dao"
	"errors"
	"github.com/tealeg/xlsx"
)

var assert = TestHelper.Assert

func mockoutGetAllUsers(keys []*datastore.Key, users []UserDao.UserDTO, err error) *TestHelper.Spy {
	spy := new(TestHelper.Spy)
	userDao_GetAllUsers = func(ctx appengine.Context) ([]*datastore.Key, []UserDao.UserDTO, error) {
		spy.RegisterCall()
		spy.RegisterArg1(ctx)
		return keys, users, err
	}
	return spy
}

func mockoutUserHasActiveSubscription(isActive []bool, err []error) *TestHelper.Spy {
	spy := new(TestHelper.Spy)
	transactionDao_GetCurrentTransactionsAfter = func(ctx appengine.Context, userKey *datastore.Key) (bool, error) {
		spy.RegisterCall()
		spy.RegisterArg1(ctx)
		spy.RegisterArg1(userKey)
		_isActive := isActive[0]
		_err := err[0]
		if len(isActive) > 1 {
			isActive = isActive[1:]
		}
		if len(err) > 1 {
			err = err[1:]
		}
		return _isActive, _err
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
		&datastore.Key{},
		&datastore.Key{},
	}
	users := []UserDao.UserDTO{
		UserDao.UserDTO{Email: "NO_Subscription"},
		UserDao.UserDTO{Email: "hasSubscription"},
	}

	spy := mockoutGetAllUsers(keys, users, nil)
	spyHasActiveSub := mockoutUserHasActiveSubscription([]bool{false, true, true}, []error{nil, nil, nil})

	usersWithActiveSubscription, err := getActiveTransactionList(ctx)

	assert(t, spy.CallCount()).Equals(1)
	assert(t, spy.GetLatestArg1()).Equals(ctx)
	assert(t, len(usersWithActiveSubscription)).Equals(1)
	assert(t, usersWithActiveSubscription[0]).Equals(users[1])
	assert(t, err).Equals(nil)

	assert(t, spyHasActiveSub.CallCount()).Equals(2)

}

func TestShouldPassOnErrorsFromDataStore_GetAllUsers(t *testing.T) {
	ctx := &TestHelper.ContextMock{}
	testError := errors.New("test error")

	getAllUsersSpy := mockoutGetAllUsers(nil, nil, testError)

	usersWithActiveSubscription, err := getActiveTransactionList(ctx)

	assert(t, getAllUsersSpy.CallCount()).Equals(1)
	assert(t, getAllUsersSpy.GetLatestArg1()).Equals(ctx)
	assert(t, usersWithActiveSubscription).Equals(nil)
	assert(t, err).Equals(testError)
}

func TestShouldPassOnErrorsFromDataStore_HasSubscription(t *testing.T) {
	ctx := &TestHelper.ContextMock{}
	testError := errors.New("test error")

	keys := []*datastore.Key{&datastore.Key{}}
	users := []UserDao.UserDTO{UserDao.UserDTO{}}

	mockoutGetAllUsers(keys, users, nil)
	mockoutUserHasActiveSubscription([]bool{false}, []error{testError})

	usersWithActiveSubscription, err := getActiveTransactionList(ctx)

	assert(t, usersWithActiveSubscription).Equals(nil)
	assert(t, err).Equals(testError)
}

func TestShouldAddHeaderRowBasedOnPassedArgumentsToAddRow(t *testing.T) {
	col1 := "test1"
	col2 := "test2"
	sheet := xlsx.Sheet{}

	addRow(&sheet, col1, col2)

	assert(t, len(sheet.Rows)).Equals(1)
	assert(t, sheet.Rows[0].Cells[0].Value).Equals(col1)
	assert(t, sheet.Rows[0].Cells[1].Value).Equals(col2)
}

func TestShouldCreateXlsxSheetWithAllUserHavingActiveSubscription(t *testing.T) {
	ctx := &TestHelper.ContextMock{}

	keys := []*datastore.Key{&datastore.Key{}}
	users := []UserDao.UserDTO{UserDao.UserDTO{Email: "testMail"}}

	mockoutGetAllUsers(keys, users, nil)
	mockoutUserHasActiveSubscription([]bool{true}, []error{nil})

	file, error := createXlsxFile(ctx)

	headerCells := file.Sheets[0].Rows[0].Cells
	firstRowCells := file.Sheets[0].Rows[1].Cells

	assert(t, len(headerCells)).Equals(7)
	assert(t, len(firstRowCells)).Equals(7)

	assert(t, firstRowCells[0].Value).Equals(users[0].NavitasId)
	//assert(t, firstRowCells[1].Value).Equals("TODO")
	assert(t, firstRowCells[2].Value).Equals(users[0].NavitasId)
	//assert(t, firstRowCells[3].Value).Equals("TODO")
	//assert(t, firstRowCells[4].Value).Equals("TODO")
	assert(t, firstRowCells[5].Value).Equals("24 Timers")
	assert(t, firstRowCells[6].Value).Equals(users[0].Email)
	assert(t, error).Equals(nil)
}

func TestExportXlsxHandler(t *testing.T) {
	ctx := new(TestHelper.ContextMock)
	testError := errors.New("test error")

	mockoutGetAllUsers(nil, nil, testError)

	file, err := createXlsxFile(ctx)

	assert(t, file).Equals(nil)
	assert(t, err).Equals(testError)
}
