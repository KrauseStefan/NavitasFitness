package main
import (
	"net/http"
	"fmt"
	"time"

	"appengine"
	"appengine/datastore"
	"encoding/json"
)

type BlogEntry struct {
	Author string
	Content string
	Date time.Time
}

func blogEntryKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "BlogEntry", "default_guestbook", 0, nil)
}

func writeJSON(w http.ResponseWriter, data interface{}) ([]byte, error){
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(data)

	w.Write(js)

	return js, err
}


func init() {
	http.HandleFunc("/rest/blogEntry", blogEntryGet)
	http.HandleFunc("/rest/", root)
}

func root(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Location", "static/index.html")
//	w.WriteHeader(http.StatusFound)
	fmt.Fprint(w, "Hello, world!")
}

func blogEntryGet(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	q:= datastore.NewQuery("BlogEntry").Ancestor(blogEntryKey(c)).Order("Date").Limit(10)

	blogEntries := make([]BlogEntry, 0, 10)

	if _, err := q.GetAll(c, &blogEntries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	if _, err := writeJSON(w, blogEntries); err != nil {
		http.Error(w, err.Error(),  http.StatusInternalServerError)
		return
	}
}