package main

import (
	"net/http"
	"src/Services"
)

func init() {
	http.HandleFunc("/rest/blogEntry", Services.HandleBlogEntryRequest)
	http.HandleFunc("/rest/user", Services.HandleUserServiceRequest)
//	http.HandleFunc("/rest/", root)
}

//func root(w http.ResponseWriter, r *http.Request) {
//	//	w.Header().Set("Location", "static/index.html")
//	//	w.WriteHeader(http.StatusFound)
//	fmt.Fprint(w, "Hello, world!")
//}
