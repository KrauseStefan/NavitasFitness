package ExportService

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/gorilla/mux"
	"github.com/tealeg/xlsx"

	"Dropbox"
	"IPN/Transaction"
	"User/Dao"
	"User/Service"
)

const (
	xlsxDateFormat = "02.01.2006"
	csvDateFormat  = "02-01-2006"
)

var (
	userDao_GetAllUsers                        = UserDao.GetAllUsers
	transactionDao_GetCurrentTransactionsAfter = func(ctx appengine.Context, userKey *datastore.Key, date time.Time) (time.Time, time.Time, error) {
		activeSubscriptions, err := TransactionDao.GetCurrentTransactionsAfter(ctx, userKey, date)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		if len(activeSubscriptions) >= 1 {
			firstTxn, lastTxn := getExtrema(activeSubscriptions)

			return firstTxn.GetPaymentDate(), lastTxn.GetPaymentDate(), nil
		}

		return time.Time{}, time.Time{}, nil
	}
)

type UserTxnTuple struct {
	user      UserDao.UserDTO
	firstDate time.Time
	lastDate  time.Time
}

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/export"

	router.
		Methods("GET").
		Path(path + "/xlsx").
		Name("export").
		HandlerFunc(UserService.AsAdmin(exportXlsxHandler))

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

func getActiveTransactionList(ctx appengine.Context) ([]UserTxnTuple, error) {

	userKeys, users, err := userDao_GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	usersWithActiveSubscription := make([]UserTxnTuple, 0, len(userKeys))

	for i, userKey := range userKeys {
		firstDate, lastDate, err := transactionDao_GetCurrentTransactionsAfter(ctx, userKey, time.Now().AddDate(0, -6, 0))
		if err != nil {
			return nil, err
		}

		if !firstDate.IsZero() && !lastDate.IsZero() {

			tuple := UserTxnTuple{
				user:      users[i],
				firstDate: firstDate,
				lastDate:  lastDate,
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

func addXlsxRow(sheet *xlsx.Sheet, values ...string) {
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

	addXlsxRow(
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
		addXlsxRow(
			sheet,
			user.user.AccessId,
			user.firstDate.Format(xlsxDateFormat),
			user.user.AccessId,
			user.firstDate.Format(xlsxDateFormat),
			user.lastDate.AddDate(0, 6, 0).Format(xlsxDateFormat),
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

func createCsvFile(ctx appengine.Context, w io.Writer) error {
	userTxnTuple, err := getActiveTransactionList(ctx)
	if err != nil {
		return err
	}

	bomPrefix := []byte{0xef, 0xbb, 0xbf}
	windowsNewline := []byte{0x0D, 0x0A}
	comma := []byte{','}
	w.Write(bomPrefix)

	//N0774,27-06-2016,03-01-2017
	//AAMS-asa,27-06-2016,03-01-2017
	//201505600,27-06-2016,03-01-2017

	if len(userTxnTuple) > 0 {
		user := userTxnTuple[0]
		ctx.Infof("%s, %s, %s", user.user.AccessId, user.firstDate.String(), user.lastDate.String())
		w.Write([]byte(user.user.AccessId))
		w.Write(comma)
		w.Write([]byte(user.firstDate.Format(csvDateFormat)))
		w.Write(comma)
		w.Write([]byte(user.lastDate.AddDate(0, 6, 0).Format(csvDateFormat)))
	}

	for _, user := range userTxnTuple[1:] {
		ctx.Infof("%s, %s, %s", user.user.AccessId, user.firstDate.String(), user.lastDate.String())
		w.Write([]byte(windowsNewline))
		w.Write([]byte(user.user.AccessId))
		w.Write(comma)
		w.Write([]byte(user.firstDate.Format(csvDateFormat)))
		w.Write(comma)
		w.Write([]byte(user.lastDate.AddDate(0, 6, 0).Format(csvDateFormat)))
	}

	return nil
}

func exportCsvHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) {
	ctx := appengine.NewContext(r)
	fileName := "ActiveSubscriptions.csv"

	var buffer bytes.Buffer

	err := createCsvFile(ctx, &buffer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = Dropbox.UploadDoc(ctx, fileName, &buffer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
