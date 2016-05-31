package ExportService

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/tealeg/xlsx"
)

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/export"

	router.
	Methods("GET").
	Path(path + "/xlsx").
	Name("export").
	HandlerFunc(exportXslt)

}

func exportXslt(w http.ResponseWriter, r *http.Request) {

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	row := sheet.AddRow()
	cell := row.AddCell()
	cell.Value = "I am a cell!"

	err = file.Write(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}