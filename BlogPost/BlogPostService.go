package BlogPostService

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"appengine"
	"appengine/datastore"

	"AppEngineHelper"
	"User/Dao"
	"User/Service"
)

type BlogEntry struct {
	Key          string    `json:"key",datastore:"-"`
	Author       string    `json:"author"`
	LastEditedBy string    `json:"lastEditedBy"`
	Content      string    `json:"content",datastore:",noindex"`
	Date         time.Time `json:"date"`
}

func (blogPost BlogEntry) hasId() bool {
	return len(blogPost.Key) > 0
}

const BLOG_KIND = "BlogEntry"
const BLOG_PARENT_STRING_ID = "default_blogentry"

var blogCollectionParentKey = AppEngineHelper.CollectionParentKeyGetFnGenerator(BLOG_KIND, BLOG_PARENT_STRING_ID, 0)
var blogIntIDToKeyInt64 = AppEngineHelper.IntIDToKeyInt64(BLOG_KIND, blogCollectionParentKey)

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
		HandlerFunc(UserService.AsAdmin(blogEntryPut))

	router.
		Methods("DELETE").
		Path(path).
		Name("DeleteBlogPost").
		HandlerFunc(UserService.AsAdmin(blogEntryDelete))

}

func blogEntryGet(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	q := datastore.NewQuery(BLOG_KIND).Ancestor(blogCollectionParentKey(ctx)).Order("Date")

	blogEntries := make([]BlogEntry, 0, 10)

	keys, err := q.GetAll(ctx, &blogEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, key := range keys {
		blogEntries[i].Key = strconv.FormatInt(key.IntID(), 10)
	}

	if _, err := AppEngineHelper.WriteJSON(w, blogEntries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func blogEntryPut(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) {
	ctx := appengine.NewContext(r)
	var blog BlogEntry
	var key *datastore.Key

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&blog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if blog.hasId() {
		key = blogIntIDToKeyInt64(ctx, blog.Key)
	} else {
		key = datastore.NewIncompleteKey(ctx, BLOG_KIND, blogCollectionParentKey(ctx))
		blog.Author = user.Email
	}
	blog.Date = time.Now()
	blog.LastEditedBy = user.Email

	key, err = datastore.Put(ctx, key, &blog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := AppEngineHelper.WriteJSON(w, blog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func blogEntryDelete(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) {
	c := appengine.NewContext(r)

	intIdStr := r.URL.Query().Get("id")
	key := blogIntIDToKeyInt64(c, intIdStr)

	if err := datastore.Delete(c, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
