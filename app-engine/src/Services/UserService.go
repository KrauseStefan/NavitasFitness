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

type UserEntry struct {
	email     string
	password  string `datastore:",noindex"`
	navitasId time.Time
	Id        string `datastore:"-"`
}

const USER_KIND = "User"
const PARENT_STRING_ID = "default_user"

func userEntryParentKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, USER_KIND, PARENT_STRING_ID, 0, nil)
}

func userIntIDToKeyInt64(c appengine.Context, id string) *datastore.Key {
	intId, _ := strconv.ParseInt(id, 10, 64)
	return datastore.NewKey(c, USER_KIND, "", intId, userEntryParentKey(c))
}

func HandleUserServiceRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	switch r.Method {
	case "GET":
		userGet(c, w, r)
	case "POST":
		userPost(c, w, r)
	case "PUT":
		userPut(c, w, r)
	case "DELETE":
		userDelete(c, w, r)
	}
}

func userGet(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	q := datastore.NewQuery(USER_KIND).Ancestor(userEntryParentKey(c)).Order("Date").Limit(10)

	userEntries := make([]UserEntry, 0, 10)

	keys, err := q.GetAll(c, &userEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, key := range keys {
		userEntries[i].Id = strconv.FormatInt(key.IntID(), 10)
	}

	if _, err := writeJSON(w, userEntries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func userPost(c appengine.Context, w http.ResponseWriter, r *http.Request) {

}

func userPut(c appengine.Context, w http.ResponseWriter, r *http.Request) {
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
		key = datastore.NewIncompleteKey(c, USER_KIND, userEntryParentKey(c))
	} else {
		key = userIntIDToKeyInt64(c, b.Id)
	}

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

func userDelete(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	intIdStr := r.URL.Query().Get("id")
	key := userIntIDToKeyInt64(c, intIdStr)

	if err := datastore.Delete(c, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}