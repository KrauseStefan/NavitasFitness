package UserService

import (
	"net/http"

	"appengine"

	"github.com/gorilla/mux"

	"src/Services/Common"
	"encoding/json"
)

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/user"

	router.
		Methods("GET").
		Path(path).
		Name("GetUserInfo").
		HandlerFunc(userGet)

	router.
		Methods("POST").
		Path(path).
		Name("CreateUserInfo").
		HandlerFunc(userPost)

}

func userGet(w http.ResponseWriter, r *http.Request) {
	userName := "name"
	ctx := appengine.NewContext(r)

	userDto, err := GetUserByEmail(ctx, userName)
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
	user := &UserDTO{}

//	user.Email = "test"
//	user.Password = "test"
//	user.NavitasId = "test"

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := CreateUser(ctx, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := Common.WriteJSON(w, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
