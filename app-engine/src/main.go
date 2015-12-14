package main

import (
	"fmt"
	"net/http"
	"time"

	"encoding/json"
	"appengine"
	"appengine/datastore"
	"strconv"
)

type BlogEntry struct {
	Author  string
	Content string `datastore:",noindex"`
	Date    time.Time
	Id     string `datastore:"-"`
}

const KIND = "BlogEntry"
const PARENT_STRING_ID = "default_blogentry"

func blogEntryParentKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, KIND, PARENT_STRING_ID, 0, nil)
}

func intIDToKeyInt64(c appengine.Context, id string) *datastore.Key {
	intId, _ := strconv.ParseInt(id, 10, 64)
	return datastore.NewKey(c, KIND, "", intId, blogEntryParentKey(c))
}

func writeJSON(w http.ResponseWriter, data interface{}) ([]byte, error) {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(data)

	w.Write(js)

	return js, err
}

func init() {
	http.HandleFunc("/rest/blogEntry", blogEntry)
//	http.HandleFunc("/rest/", root)
}

//func root(w http.ResponseWriter, r *http.Request) {
//	//	w.Header().Set("Location", "static/index.html")
//	//	w.WriteHeader(http.StatusFound)
//	fmt.Fprint(w, "Hello, world!")
//}

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

	keys, err := q.GetAll(c, &blogEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, key := range keys {
		blogEntries[i].Id = strconv.FormatInt(key.IntID(), 10)
	}

	if _, err := writeJSON(w, blogEntries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func blogEntryPost(c appengine.Context, w http.ResponseWriter, r *http.Request) {

}

func blogEntryPut(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	var b BlogEntry
	var key *datastore.Key

	//	if u := user.Current(c); u != nil {
	//		b.Author = u.String()
	//	}

	fmt.Print("content: ", r.Body)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if (len(b.Id) == 0) {
		key = datastore.NewIncompleteKey(c, KIND, blogEntryParentKey(c))
	} else {
		key = intIDToKeyInt64(c, b.Id)
	}
	b.Date = time.Now()
	b.Author = "skk"

	key, err = datastore.Put(c, key, &b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := writeJSON(w, b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func blogEntryDelete(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	intIdStr := r.URL.Query().Get("id")
	key := intIDToKeyInt64(c, intIdStr)

	if err := datastore.Delete(c, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
