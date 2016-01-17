package main

import (
	"net/http"
	"src/Services/BlogPost"
	"src/Services/User"
	"src/Services/Auth"

	"github.com/gorilla/mux"
)


// http://blog.golang.org/context
// http://blog.golang.org/go-videos-from-google-io-2012

func init() {

	router := mux.NewRouter().StrictSlash(true)
	BlogPostService.IntegrateRoutes(router)
	UserService.IntegrateRoutes(router)
	AuthService.IntegrateRoutes(router)
//	router.HandleFunc("/rest/user", UserService.HandleUserServiceRequest)

	http.Handle("/", router)
	//	http.HandleFunc("/rest/", root)
}

//func root(w http.ResponseWriter, r *http.Request) {
//	//	w.Header().Set("Location", "static/index.html")
//	//	w.WriteHeader(http.StatusFound)
//	fmt.Fprint(w, "Hello, world!")
//}
