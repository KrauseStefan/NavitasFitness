package xlsx

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"google.golang.org/appengine/datastore"

	"github.com/tealeg/xlsx"

	"IPN/Transaction"
	"User/Dao"

	"IPN/Transaction/TransactionDaoTestHelper"
	"TestHelper"
	"User/Dao/UserDaoTestHelper"
)

var (
	assert = TestHelper.Assert
	utc, _ = time.LoadLocation("UTC")
	ctx    = TestHelper.GetContext()
)

func mockTransactionRetriever(messages []*TransactionDao.TransactionMsgDTO, err error) *TransactionDaoTestHelper.TransactionRetrieverMock {
	mock := TransactionDaoTestHelper.NewTransactionRetrieverMock(messages, err)
	transactionDao = mock
	return mock
}

func mockUserRetriever(keys []*datastore.Key, users []UserDao.UserDTO, err error) *UserDaoTestHelper.UserRetrieverMock {
	mock := UserDaoTestHelper.NewUserRetrieverMock(keys, users, err)
	userDAO = mock
	return mock
}

func createMessage(date time.Time) []*TransactionDao.TransactionMsgDTO {
	const layout = "15:04:05 Jan 02, 2006 MST"
	dateStr := date.In(utc).Format(layout)
	dateIpnMsg := TransactionDao.FIELD_PAYMENT_DATE + "=" + dateStr
	dateTxnMsg := TransactionDao.NewTransactionMsgDTOFromIpn(dateIpnMsg)
	return []*TransactionDao.TransactionMsgDTO{dateTxnMsg, dateTxnMsg}
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
	keys := []*datastore.Key{{}, {}}
	users := []UserDao.UserDTO{
		{Email: "NO_Subscription"},
		{Email: "hasSubscription"},
	}

	invalidMessages := createMessage(time.Time{})
	nowMessages := createMessage(time.Now())

	userDaoMock := mockUserRetriever(keys, users, nil)
	transactionRetrieverMock := mockTransactionRetriever(invalidMessages, nil).
		AddReturn(nowMessages, nil)

	userTxnTuple, err := getActiveTransactionList(ctx)

	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, userDaoMock.LatestCallCtxArg).Equals(ctx)
	assert(t, len(userTxnTuple)).Equals(1)
	assert(t, userTxnTuple[0].user).Equals(users[1])
	assert(t, err).Equals(nil)

	assert(t, transactionRetrieverMock.CallCount).Equals(2)
}

func TestShouldPassOnErrorsFromDataStore_GetAllUsers(t *testing.T) {
	testError := errors.New("test error")

	userDaoMock := mockUserRetriever(nil, nil, testError)

	userTxnTuple, err := getActiveTransactionList(ctx)

	assert(t, userDaoMock.CallCount).Equals(1)
	assert(t, userDaoMock.LatestCallCtxArg).Equals(ctx)
	assert(t, userTxnTuple).Equals(nil)
	assert(t, err).Equals(testError)
}

func TestShouldPassOnErrorsFromDataStore_HasSubscription(t *testing.T) {
	testError := errors.New("test error")

	keys := []*datastore.Key{{}}
	users := []UserDao.UserDTO{{}}
	now := time.Now()

	mockTransactionRetriever(createMessage(now), testError)
	mockUserRetriever(keys, users, nil)

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
	keys := []*datastore.Key{{}}
	users := []UserDao.UserDTO{{Email: "testMail"}}
	now := time.Now()

	mockTransactionRetriever(createMessage(now), nil)
	mockUserRetriever(keys, users, nil)

	file, err := createXlsxFile(ctx)

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
	assert(t, err).Equals(nil)
}

func TestExportXlsxHandler(t *testing.T) {
	testError := errors.New("test error")

	mockUserRetriever(nil, nil, testError)

	file, err := createXlsxFile(ctx)

	assert(t, file).Equals(nil)
	assert(t, err).Equals(testError)
}
