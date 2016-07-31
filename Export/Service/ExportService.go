package ExportService

import (
	"net/http"

	"appengine"

	"github.com/gorilla/mux"
	"github.com/tealeg/xlsx"

	"NavitasFitness/IPN/Transaction"
	"NavitasFitness/User/Dao"
)

var (
userDao_GetAllUsers = UserDao.GetAllUsers
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

	userKeys, users, err := userDao_GetAllUsers(ctx)
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

	httpHeader := w.Header()
	configureHeaderForFileDownload(&httpHeader, "ActiveSubscriptions.xlsx")

	file, err := exportXslt(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = file.Write(w)
}

func exportXslt(ctx appengine.Context) (*xlsx.File, error) {

	users, err := getTransactionList(ctx)
	if err != nil {
		return nil, err
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return nil, err
	}

	addRow(sheet, "email")

	for _, user := range users {
		addRow(sheet, user.Email)
	}

	if err != nil {
		return nil, err
	}

	return file, nil
}
