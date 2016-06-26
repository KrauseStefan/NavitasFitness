package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"./Auth"
	"./BlogPost"
	"./Export/Service"
	"./IPN"
	"./User/Service"
)

// http://blog.golang.org/context
// http://blog.golang.org/go-videos-from-google-io-2012

func init() {

	router := mux.NewRouter().StrictSlash(true)
	BlogPostService.IntegrateRoutes(router)
	UserService.IntegrateRoutes(router)
	AuthService.IntegrateRoutes(router)
	IPN.IntegrateRoutes(router)
	ExportService.IntegrateRoutes(router)

	http.Handle("/", router)
	//	http.HandleFunc("/rest/", root)
}

//func root(w http.ResponseWriter, r *http.Request) {
//	//	w.Header().Set("Location", "static/index.html")
//	//	w.WriteHeader(http.StatusFound)
//	fmt.Fprint(w, "Hello, world!")
//}
