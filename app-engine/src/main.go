package main
import (
	"net/http"
	"fmt"
	"time"

	"appengine"
	"appengine/datastore"
	"encoding/json"
	"strconv"
)

type BlogEntry struct {
	Author  string
	Content string `datastore:",noindex"`
	Date    time.Time
	Id      int64 `datastore:"-"`
}

const KIND = "BlogEntry"
const PARENT_STRING_ID = "default_blogentry"

func blogEntryParentKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, KIND, PARENT_STRING_ID, 0, nil)
}

func writeJSON(w http.ResponseWriter, data interface{}) ([]byte, error) {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(data)

	w.Write(js)

	return js, err
}


func init() {
	http.HandleFunc("/rest/blogEntry", blogEntry)
	http.HandleFunc("/rest/", root)
}

func root(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Location", "static/index.html")
	//	w.WriteHeader(http.StatusFound)
	fmt.Fprint(w, "Hello, world!")
}

func blogEntry(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	switch r.Method {
	case "GET":
		blogEntryGet(c, w, r)
	case "POST":
		blogEntryPost(c, w, r)
	case "PUT":
		blogEntryPut(c, w, r)
	case "DELETE":
		blogEntryDelete(c, w, r)
	}
}

func blogEntryGet(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	q := datastore.NewQuery(KIND).Ancestor(blogEntryParentKey(c)).Order("Date").Limit(10)

	blogEntries := make([]BlogEntry, 0, 10)

	keys, err := q.GetAll(c, &blogEntries);
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, key := range keys {
		blogEntries[i].Id = key.IntID()
	}

	if _, err := writeJSON(w, blogEntries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func blogEntryPost(c appengine.Context, w http.ResponseWriter, r *http.Request) {

}

func blogEntryPut(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	//	if u := user.Current(c); u != nil {
	//		b.Author = u.String()
	//	}

	b := BlogEntry{
		Author: "skk",
		Content: "test",
		Date: time.Now(),
	}

	key := datastore.NewIncompleteKey(c, KIND, blogEntryParentKey(c))

	key, err := datastore.Put(c, key, &b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b.Id = key.IntID()

	if _, err := writeJSON(w, b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func blogEntryDelete(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	intIdStr := r.URL.Query().Get("id")
	intId, _ := strconv.Atoi(intIdStr)
	key := datastore.NewKey(c, KIND, "", int64(intId), blogEntryParentKey(c))

	if err := datastore.Delete(c, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}