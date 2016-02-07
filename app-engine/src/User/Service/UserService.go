package UserService

import (
	"net/http"

	"appengine"

	"github.com/gorilla/mux"

	"src/Common"
	"encoding/json"
	"src/Auth"

	"src/User/Dao"
)

const emailParam = "email"

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/user"

	router.
		Methods("GET").
		Path(path).
		Name("GetUserCurrentUserInfo").
		HandlerFunc(getUserFromSession)

	router.
		Methods("GET").
		Path(path + "/{" + emailParam + "}").
		Name("GetUserCurrentUserInfo").
		HandlerFunc(userGetByEmail)

	router.
		Methods("POST").
		Path(path).
		Name("CreateUserInfo").
		HandlerFunc(userPost)

}

func getUserFromSession(w http.ResponseWriter, r *http.Request) {
	uuid, err := AuthService.GetSessionUUID(r)

	if uuid == "" {
		http.Error(w, "", http.StatusNoContent)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
//		http.Error(w, "", http.StatusNoContent)
		return
	}

	ctx := appengine.NewContext(r)

	user, err := UserDao.GetUserFromSessionUUID(ctx, uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if _, err := Common.WriteJSON(w, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func userGetByEmail(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	email := vars[emailParam]

	ctx := appengine.NewContext(r)

	userDto, err := UserDao.GetUserByEmail(ctx, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := Common.WriteJSON(w, userDto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func userPost(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	user := &UserDao.UserDTO{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := UserDao.CreateUser(ctx, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := Common.WriteJSON(w, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
