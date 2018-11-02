package UserRest

import (
	"encoding/json"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"github.com/gorilla/mux"

	"AccessIdValidator"
	"AppEngineHelper"
	"Auth"
	"User/Dao"
	"User/Service"
)

const emailKey = "email"
const userKey = "userKey"

var accessIdValidator = AccessIdValidator.GetInstance()

type UserSessionDto struct {
	User          *UserDao.UserDTO `json:"user"`
	IsAdmin       bool             `json:"isAdmin"`
	ValidAccessId bool             `json:"validAccessId"`
}

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/user"

	router.
		Methods("GET").
		Path(path).
		Name("Get User Current User Info").
		HandlerFunc(UserService.AsUser(getUserFromSessionHandler))

	router.
		Methods("GET").
		Path(path + "/transactions").
		Name("Get Latest Transactions").
		HandlerFunc(UserService.AsUser(getCurrentUserTransactionsHandler))

	router.
		Methods("GET").
		Path(path + "/transactions/{" + userKey + "}").
		Name("Get Latest Transactions").
		HandlerFunc(UserService.AsAdmin(getUserTransactionsHandler))

	router.
		Methods("POST").
		Path(path).
		Name("Create User Info").
		HandlerFunc(AppEngineHelper.HandlerW(createUserHandler))

	router.
		Methods("GET").
		Path(path + "/verify").
		Name("VerifyEmailCallback").
		HandlerFunc(AppEngineHelper.HandlerW(verifyUserRequestHandler))

	router.
		Methods("POST").
		Path(path + "/resetPassword/{" + emailKey + "}").
		Name("ResetPassword").
		HandlerFunc(AppEngineHelper.HandlerW(requestResetUserPasswordHandler))

	router.
		Methods("POST").
		Path(path + "/changePassword").
		Name("ChangePassword").
		HandlerFunc(AppEngineHelper.HandlerW(resetUserPasswordHandler))

	router.
		Methods("GET").
		Path(path + "/all").
		Name("Retrieve all users").
		HandlerFunc(UserService.AsAdmin(getAllUsersHandler))
}

func getAllUsersHandler(w http.ResponseWriter, r *http.Request, _ *UserDao.UserDTO) (interface{}, error) {
	type UserAndKeys struct {
		Keys []string           `json:"keys"`
		User []*UserDao.UserDTO `json:"users"`
	}

	ctx := appengine.NewContext(r)

	keys, users, err := UserService.GetAllUsers(ctx)

	data := &UserAndKeys{
		keys,
		users,
	}

	return data, err
}

func getUserFromSessionHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
	ctx := appengine.NewContext(r)
	if user == nil {
		log.Debugf(ctx, "User is not logged in")
		return nil, nil
	}

	if err := accessIdValidator.EnsureUpdatedIds(ctx); err != nil {
		return nil, err
	}

	isValid, err := accessIdValidator.ValidateAccessId(ctx, []byte(user.AccessId))
	if err != nil {
		return nil, err
	}

	userSessionDto := &UserSessionDto{
		User:          user,
		IsAdmin:       user.IsAdmin,
		ValidAccessId: isValid,
	}

	return userSessionDto, err
}

func getCurrentUserTransactionsHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
	ctx := appengine.NewContext(r)

	return UserService.GetUserTransactions(ctx, user.Key)
}

func getUserTransactionsHandler(w http.ResponseWriter, r *http.Request, _ *UserDao.UserDTO) (interface{}, error) {
	ctx := appengine.NewContext(r)

	userKeyStr := mux.Vars(r)[userKey]
	userKey, err := datastore.DecodeKey(userKeyStr)
	if err == nil {
		return nil, err
	}
	return UserService.GetUserTransactions(ctx, userKey)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := appengine.NewContext(r)
	user := &UserDao.UserDTO{}

	sessionData, err := Auth.GetSessionData(r)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(user); err != nil {
		return nil, err
	}

	createdUser, err := UserService.CreateUser(ctx, user, sessionData)
	if err != nil {
		return nil, err
	}

	if err := Auth.UpdateSessionDataUserKey(r, w, createdUser.Key); err != nil {
		return nil, err
	}

	return createdUser, nil
}

func verifyUserRequestHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := appengine.NewContext(r)

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	key := r.Form.Get("code")
	if err := UserService.MarkUserVerified(ctx, key); err != nil {
		http.Redirect(w, r, "/?Verified=false", http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, "/?Verified=true", http.StatusTemporaryRedirect)
	}
	return nil, nil
}

func requestResetUserPasswordHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := appengine.NewContext(r)
	email := mux.Vars(r)[emailKey]

	err := UserService.RequestResetUserPassword(ctx, email)
	return nil, err
}

func resetUserPasswordHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := appengine.NewContext(r)

	err := UserService.ResetUserPassword(ctx, r.Body)
	return nil, err
}
