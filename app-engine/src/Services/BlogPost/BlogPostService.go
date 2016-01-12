package BlogPostService

import (
	"fmt"

	"net/http"

	"appengine"
	"appengine/datastore"

	"strconv"
	"encoding/json"
	"time"
	"src/Services/Common"
	 "github.com/gorilla/mux"
)

type BlogEntry struct {
	Author  string
	Content string `datastore:",noindex"`
	Date    time.Time
	Id      string `datastore:"-"`
}

const BLOG_KIND = "BlogEntry"
const BLOG_PARENT_STRING_ID = "default_blogentry"

var blogCollectionParentKey = Common.CollectionParentKeyGetFnGenerator(BLOG_KIND, BLOG_PARENT_STRING_ID, 0)
var blogIntIDToKeyInt64 = Common.IntIDToKeyInt64(BLOG_KIND, blogCollectionParentKey)

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/blogEntry"

 	router.
		Methods("GET").
		Path(path).
		Name("GetAllBlogPosts").
		HandlerFunc(blogEntryGet)

	router.
		Methods("PUT").
		Path(path).
		Name("PersistBlogPost").
		HandlerFunc(blogEntryPut)

	router.
		Methods("DELETE").
		Path(path).
		Name("DeleteBlogPost").
		HandlerFunc(blogEntryDelete)

}

func blogEntryGet(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery(BLOG_KIND).Ancestor(blogCollectionParentKey(c)).Order("Date")

	blogEntries := make([]BlogEntry, 0, 10)

	keys, err := q.GetAll(c, &blogEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, key := range keys {
		blogEntries[i].Id = strconv.FormatInt(key.IntID(), 10)
	}

	if _, err := Common.WriteJSON(w, blogEntries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func blogEntryPut(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
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
		key = datastore.NewIncompleteKey(c, BLOG_KIND, blogCollectionParentKey(c))
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

	if _, err := Common.WriteJSON(w, b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func blogEntryDelete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	intIdStr := r.URL.Query().Get("id")
	key := blogIntIDToKeyInt64(c, intIdStr)

	if err := datastore.Delete(c, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}