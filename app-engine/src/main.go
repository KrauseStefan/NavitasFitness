package main

import (
	"net/http"
	"src/Services/Blog"
	"src/Services/User"
)

func init() {
	http.HandleFunc("/rest/blogEntry", BlogEntry.HandleBlogEntryRequest)
	http.HandleFunc("/rest/user", User.HandleUserServiceRequest)
//	http.HandleFunc("/rest/", root)
}

//func root(w http.ResponseWriter, r *http.Request) {
//	//	w.Header().Set("Location", "static/index.html")
//	//	w.WriteHeader(http.StatusFound)
//	fmt.Fprint(w, "Hello, world!")
//}
