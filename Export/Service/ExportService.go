package ExportService

import (
	"IPN/Transaction"
	"User/Dao"
	"appengine"
	"github.com/gorilla/mux"
	"github.com/tealeg/xlsx"
	"net/http"
)

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/export"

	router.
		Methods("GET").
		Path(path + "/xlsx").
		Name("export").
		HandlerFunc(exportXsltHandler)

}

func getTransactionList(ctx appengine.Context) ([]UserDao.UserDTO, error) {

	userKeys, users, err := UserDao.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	usersWithActiveSubscription := make([]UserDao.UserDTO, 0, len(userKeys))

	for i, userKey := range userKeys {
		userHasActiveSubscription, err := TransactionDao.UserHasActiveSubscription(ctx, userKey)
		if err != nil {
			return nil, err
		}

		if userHasActiveSubscription {
			usersWithActiveSubscription = append(usersWithActiveSubscription, users[i])
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

func addRow(sheet *xlsx.Sheet, headers ...string) {
	row := sheet.AddRow()

	for _, header := range headers {
		cell := row.AddCell()
		cell.Value = header
	}
}

func exportXsltHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)

	users, err := getTransactionList(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addRow(sheet, "email")

	for _, user := range users {
		addRow(sheet, user.Email)
	}

	httpHeader := w.Header()
	configureHeaderForFileDownload(&httpHeader, "ActiveSubscriptions.xlsx")

	err = file.Write(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
