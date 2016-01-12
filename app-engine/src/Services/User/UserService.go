package UserService

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

type UserEntry struct {
	email     string
	password  string `datastore:",noindex"`
	navitasId time.Time
	Id        string `datastore:"-"`
}

const USER_KIND = "User"
const PARENT_STRING_ID = "default_user"

var userCollectionParentKey = Common.CollectionParentKeyGetFnGenerator(USER_KIND, PARENT_STRING_ID, 0)
var userIntIDToKeyInt64 = Common.IntIDToKeyInt64(USER_KIND, userCollectionParentKey)

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/user"

	router.
		Methods("GET").
		Path(path).
		Name("GetUserInfo").
		HandlerFunc(userGet)

	router.
		Methods("PUT").
		Path(path).
		Name("PersistUserInfo").
		HandlerFunc(userPut)

}

func userGet(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery(USER_KIND).Ancestor(userCollectionParentKey(c)).Order("Date").Limit(10)

	userEntries := make([]UserEntry, 0, 10)

	keys, err := q.GetAll(c, &userEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, key := range keys {
		userEntries[i].Id = strconv.FormatInt(key.IntID(), 10)
	}

	if _, err := Common.WriteJSON(w, userEntries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func userPut(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var b UserEntry
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
		key = datastore.NewIncompleteKey(c, USER_KIND, userCollectionParentKey(c))
	} else {
		key = userIntIDToKeyInt64(c, b.Id)
	}

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