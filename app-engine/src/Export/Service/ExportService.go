package ExportService

import (
	"appengine"
	"github.com/gorilla/mux"
	"github.com/tealeg/xlsx"
	"net/http"
	"src/IPN/Transaction"
	"src/User/Dao"
)

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/export"

	router.
		Methods("GET").
		Path(path + "/ActiveSubscriptions.xlsx").
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

	header := sheet.AddRow()
	nameCell := header.AddCell()
	nameCell.Value = "email"

	for _, user := range users {
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.Value = user.Email
	}

	w.Header().Add("Content-Disposition", "attachment")
	w.Header().Add("filename", "name_of_excel_file.xls")
	w.Header().Add("Content-Type", "application/vnd.ms-excel")

	err = file.Write(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
