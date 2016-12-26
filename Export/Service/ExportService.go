package ExportService

import (
	"net/http"

	"appengine"

	"github.com/gorilla/mux"
	"github.com/tealeg/xlsx"

	"IPN/Transaction"
	"User/Dao"
	"User/Service"
	"time"
)

const xlsxDateFormat = "02.01.2006"

var (
	userDao_GetAllUsers                        = UserDao.GetAllUsers
	transactionDao_GetCurrentTransactionsAfter = TransactionDao.GetCurrentTransactionsAfter
)

type UserTxnTuple struct {
	user     UserDao.UserDTO
	firstTxn *TransactionDao.TransactionMsgDTO
	lastTxn  *TransactionDao.TransactionMsgDTO
}

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/export"

	router.
		Methods("GET").
		Path(path + "/xlsx").
		Name("export").
		HandlerFunc(UserService.AsAdmin(exportXlsxHandler))
}

func getExtrema(txns []*TransactionDao.TransactionMsgDTO) (*TransactionDao.TransactionMsgDTO, *TransactionDao.TransactionMsgDTO) {
	firstTxn := txns[0]
	lastTxn := txns[0]

	for _, txn := range txns {
		if txn.GetPaymentActivationDate().Before(firstTxn.GetPaymentActivationDate()) {
			firstTxn = txn
		}

		if txn.GetPaymentActivationDate().After(lastTxn.GetPaymentActivationDate()) {
			lastTxn = txn
		}
	}

	return firstTxn, lastTxn
}

func getActiveTransactionList(ctx appengine.Context) ([]UserTxnTuple, error) {

	userKeys, users, err := userDao_GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	usersWithActiveSubscription := make([]UserTxnTuple, 0, len(userKeys))

	for i, userKey := range userKeys {
		activeSubscriptions, err := transactionDao_GetCurrentTransactionsAfter(ctx, userKey, time.Now().AddDate(0, -6, 0))
		if err != nil {
			return nil, err
		}

		if len(activeSubscriptions) >= 1 {
			firstTxn, lastTxn := getExtrema(activeSubscriptions)

			tuple := UserTxnTuple{
				user:     users[i],
				firstTxn: firstTxn,
				lastTxn:  lastTxn,
			}
			usersWithActiveSubscription = append(usersWithActiveSubscription, tuple)
		}
	}

	return usersWithActiveSubscription, nil
}

func configureHeaderForFileDownload(header *http.Header, filename string) {
	header.Add("Content-Disposition", "attachment; filename="+filename)
	header.Add("Content-Type", "application/vnd.ms-excel")
	header.Add("Cache-Control", "no-cache, no-store, must-revalidate")
	header.Add("Pragma", "no-cache")
	header.Add("Expires", "0")
}

func addRow(sheet *xlsx.Sheet, values ...string) {
	row := sheet.AddRow()

	for _, value := range values {
		cell := row.AddCell()
		cell.Value = value
	}
}

func createXlsxFile(ctx appengine.Context) (*xlsx.File, error) {

	userTxnTuple, err := getActiveTransactionList(ctx)
	if err != nil {
		return nil, err
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return nil, err
	}

	addRow(
		sheet,
		"Medarbejder nr i ADK",
		"Aktiveringsdato",
		"Nr.",
		"Fra dato",
		"Til dato",
		"Tidsskema",
		"Bemærkninger",
	)

	for _, user := range userTxnTuple {
		addRow(
			sheet,
			user.user.NavitasId,
			user.firstTxn.GetPaymentActivationDate().Format(xlsxDateFormat),
			user.user.NavitasId,
			user.firstTxn.GetPaymentActivationDate().Format(xlsxDateFormat),
			user.lastTxn.GetPaymentActivationDate().AddDate(0, 6, 0).Format(xlsxDateFormat),
			"24 Timers",
			user.user.Email,
		)
	}

	// Example Data:
	// "Medarbejder nr i ADK" : "N0416"
	// "Aktiveringsdato" 			: "30.06.2015"
	// "Nr."									: "N0416"
	// "Fra dato"							: "30-06-2015"
	// "Til dato"							: "06-01-2016"
	// "Tidsskema"						: "24 Timers"
	// "Bemærkninger"					: ""

	return file, nil
}

func exportXlsxHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) {
	ctx := appengine.NewContext(r)

	httpHeader := w.Header()
	configureHeaderForFileDownload(&httpHeader, "ActiveSubscriptions.xlsx")

	file, err := createXlsxFile(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = file.Write(w)
}
