package ExportService

import (
	"net/http"

	"appengine"

	"github.com/gorilla/mux"
	"github.com/tealeg/xlsx"

	"IPN/Transaction"
	"User/Dao"
	"User/Service"
)

var (
	userDao_GetAllUsers                      = UserDao.GetAllUsers
	transactionDao_UserHasActiveSubscription = TransactionDao.UserHasActiveSubscription
)

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/export"

	router.
		Methods("GET").
		Path(path + "/xlsx").
		Name("export").
		HandlerFunc(UserService.AsAdmin(exportXlsxHandler))

}

func getTransactionList(ctx appengine.Context) ([]UserDao.UserDTO, error) {

	userKeys, users, err := userDao_GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	usersWithActiveSubscription := make([]UserDao.UserDTO, 0, len(userKeys))

	for i, userKey := range userKeys {
		userHasActiveSubscription, err := transactionDao_UserHasActiveSubscription(ctx, userKey)
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

	for _, user := range users {
		addRow(
			sheet,
			user.NavitasId,
			"TODO",
			user.NavitasId,
			"TODO",
			"TODO",
			"24 Timers",
			user.Email,
		)
	}

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

	file, err := exportXslt(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = file.Write(w)
}
