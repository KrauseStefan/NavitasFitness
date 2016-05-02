package UserService

import (
	"net/http"

	"appengine"

	"github.com/gorilla/mux"

	"src/Common"
	"encoding/json"
	"src/Auth"

	"src/User/Dao"
	"src/IPN/Transaction"
)

const emailParam = "email"

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/user"

	router.
		Methods("GET").
		Path(path).
		Name("GetUserCurrentUserInfo").
		HandlerFunc(getUserFromSessionHandler)

	router.
		Methods("GET").
		Path(path + "/transactions").
		Name("GetLatestTransactions").
		HandlerFunc(getUserTransactionsHandler)

	//router.
	//	Methods("GET").
	//	Path(path + "/{" + emailParam + "}").
	//	Name("GetUserCurrentUserInfo").
	//	HandlerFunc(userGetByEmail)

	router.
		Methods("POST").
		Path(path).
		Name("CreateUserInfo").
		HandlerFunc(userPost)

}

func getUserFromSession(ctx appengine.Context, r *http.Request) (*UserDao.UserDTO, error){
	uuid, err := AuthService.GetSessionUUID(r)

	if err != nil {
		return nil, err
	}

	return UserDao.GetUserFromSessionUUID(ctx, uuid)
}

func getUserFromSessionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	user, err := getUserFromSession(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if _, err := Common.WriteJSON(w, user); err != nil {
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

func getUserTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	user, err := getUserFromSession(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	transactions, err := TransActionDao.GetTransactionsByUser(ctx, user.GetDataStoreKey(ctx))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := Common.WriteJSON(w, transactions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//func userGetByEmail(w http.ResponseWriter, r *http.Request) {
//
//	vars := mux.Vars(r)
//	email := vars[emailParam]
//
//	ctx := appengine.NewContext(r)
//
//	userDto, err := UserDao.GetUserByEmail(ctx, email)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	if _, err := Common.WriteJSON(w, userDto); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}
