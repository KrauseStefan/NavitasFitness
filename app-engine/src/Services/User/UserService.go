package UserService

import (
	"net/http"

	"appengine"
	"appengine/datastore"

	"encoding/json"

	"github.com/gorilla/mux"

	"src/Services/Common"
)

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
	userName := "name"
	ctx := appengine.NewContext(r)

	userDto, err := GetUserByUserName(ctx, userName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := Common.WriteJSON(w, userDto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func userPut(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var user UserDTO
	var key *datastore.Key

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if (user.HasId()) {
		CreateUser(ctx, &user)
	} else {
		UpdateUser(ctx, &user)
	}

	key, err = datastore.Put(ctx, key, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := Common.WriteJSON(w, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}