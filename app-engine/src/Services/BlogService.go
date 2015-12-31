package Services

import (
	"fmt"

	"net/http"

	"appengine"
	"appengine/datastore"

	"strconv"
	"encoding/json"
	"time"
)

type BlogEntry struct {
	Author  string
	Content string `datastore:",noindex"`
	Date    time.Time
	Id     string `datastore:"-"`
}

const BLOG_KIND = "BlogEntry"
const BLOG_PARENT_STRING_ID = "default_blogentry"

func blogIntIDToKeyInt64(c appengine.Context, id string) *datastore.Key {
	intId, _ := strconv.ParseInt(id, 10, 64)
	return datastore.NewKey(c, BLOG_KIND, "", intId, blogEntryParentKey(c))
}

func blogEntryParentKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, BLOG_KIND, BLOG_PARENT_STRING_ID, 0, nil)
}

func HandleBlogEntryRequest(w http.ResponseWriter, r *http.Request) {
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
	q := datastore.NewQuery(BLOG_KIND).Ancestor(blogEntryParentKey(c)).Order("Date").Limit(10)

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
		key = datastore.NewIncompleteKey(c, BLOG_KIND, blogEntryParentKey(c))
	} else {
		key = blogIntIDToKeyInt64(c, b.Id)
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
	key := blogIntIDToKeyInt64(c, intIdStr)

	if err := datastore.Delete(c, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}