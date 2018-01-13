package UserRest

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"github.com/gorilla/mux"

	"AccessIdValidator"
	"AppEngineHelper"
	"Auth"
	"DAOHelper"
	"User/Dao"
	"User/Service"
)

const emailKey = "email"
const userKey = "userKey"

var accessIdValidator = AccessIdValidator.GetInstance()

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
		HandlerFunc(createUserHandler)

	router.
		Methods("GET").
		Path(path + "/verify").
		Name("VerifyEmailCallback").
		HandlerFunc(verifyUserRequestHandler)

	router.
		Methods("POST").
		Path(path + "/resetPassword/{" + emailKey + "}").
		Name("ResetPassword").
		HandlerFunc(requestResetUserPasswordHandler)

	router.
		Methods("POST").
		Path(path + "/changePassword").
		Name("ChangePassword").
		HandlerFunc(resetUserPasswordHandler)

	router.
		Methods("GET").
		Path(path + "/all").
		Name("Retrieve all users").
		HandlerFunc(UserService.AsAdmin(getAllUsersHandler))

}

func getAllUsersHandler(w http.ResponseWriter, r *http.Request, _ *UserDao.UserDTO) {
	type UserAndKeys struct {
		Keys []string          `json:"keys"`
		User []UserDao.UserDTO `json:"users"`
	}

	ctx := appengine.NewContext(r)

	keys, users, err := UserService.GetAllUsers(ctx)

	data := &UserAndKeys{
		keys,
		users,
	}

	if err == nil {
		_, err = AppEngineHelper.WriteJSON(w, data)
	}
}

type UserSessionDto struct {
	User          *UserDao.UserDTO `json:"user"`
	IsAdmin       bool             `json:"isAdmin"`
	ValidAccessId bool             `json:"validAccessId"`
}

func getUserFromSessionHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) {
	ctx := appengine.NewContext(r)
	if user == nil {
		log.Infof(ctx, "User is not logged in")
		return
	}
	isValid, err := accessIdValidator.ValidateAccessIdPrimary(ctx, []byte(user.AccessId))

	if !isValid && err == nil {
		isValid, err = accessIdValidator.ValidateAccessIdSecondary(ctx, []byte(user.AccessId))
	}

	us := UserSessionDto{
		User:          user,
		IsAdmin:       user.IsAdmin,
		ValidAccessId: isValid,
	}

	if err == nil {
		_, err = AppEngineHelper.WriteJSON(w, us)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getCurrentUserTransactionsHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) {
	ctx := appengine.NewContext(r)

	txnClientDtoList, err := UserService.GetUserTransactions(ctx, user.Key)

	if err == nil {
		_, err = AppEngineHelper.WriteJSON(w, txnClientDtoList)
	}

	DAOHelper.ReportError(ctx, w, err)
}

func getUserTransactionsHandler(w http.ResponseWriter, r *http.Request, _ *UserDao.UserDTO) {
	ctx := appengine.NewContext(r)

	userKeyStr := mux.Vars(r)[userKey]
	userKey, err := datastore.DecodeKey(userKeyStr)
	if err == nil {
		txnClientDtoList, err := UserService.GetUserTransactions(ctx, userKey)

		if err == nil {
			_, err = AppEngineHelper.WriteJSON(w, txnClientDtoList)
		}
	}

	DAOHelper.ReportError(ctx, w, err)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	user, err := createUserHandlerInternal(w, r)

	if err == nil {
		_, err = AppEngineHelper.WriteJSON(w, user)
	}

	DAOHelper.ReportError(ctx, w, err)
}

func createUserHandlerInternal(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := appengine.NewContext(r)

	sessionData, err := Auth.GetSessionData(r)
	if err != nil {
		return nil, err
	}

	user, err := UserService.CreateUser(ctx, r, sessionData)
	if err != nil {
		return nil, err
	}

	if err := Auth.UpdateSessionDataUserKey(r, w, user.Key); err != nil {
		return nil, err
	}

	return user, nil
}

func verifyUserRequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if err := r.ParseForm(); err != nil {
		DAOHelper.ReportError(ctx, w, err)
		return
	}

	key := r.Form.Get("code")
	if err := UserService.MarkUserVerified(ctx, key); err != nil {
		http.Redirect(w, r, "/?Verified=false", http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, "/?Verified=true", http.StatusTemporaryRedirect)
	}
}

func requestResetUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	email := mux.Vars(r)[emailKey]

	if err := UserService.RequestResetUserPassword(ctx, email); err != nil {
		DAOHelper.ReportError(ctx, w, err)
	}
}

func resetUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if err := UserService.ResetUserPassword(ctx, r.Body); err != nil {
		DAOHelper.ReportError(ctx, w, err)
	}
}
