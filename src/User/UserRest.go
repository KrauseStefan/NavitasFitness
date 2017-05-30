package UserRest

import (
	"net/http"

	"google.golang.org/appengine"

	"github.com/gorilla/mux"

	"AccessIdValidator"
	"AppEngineHelper"
	"DAOHelper"
	"User/Dao"
	"User/Service"
)

const accessIdKey = "accessId"
const emailKey = "email"

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
		HandlerFunc(UserService.AsUser(getUserTransactionsHandler))

	router.
		Methods("POST").
		Path(path).
		Name("Create User Info").
		HandlerFunc(createUserHandler)

	router.
		Methods("GET").
		Path(path + "/validate_id/{" + accessIdKey + "}").
		Name("Validate Access Id").
		HandlerFunc(validateAccessId)

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

}

type UserSessionDto struct {
	User          *UserDao.UserDTO `json:"user"`
	IsAdmin       bool             `json:"isAdmin"`
	ValidAccessId bool             `json:"validAccessId"`
}

func getUserFromSessionHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) {
	ctx := appengine.NewContext(r)
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

func getUserTransactionsHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) {
	ctx := appengine.NewContext(r)

	txnClientDtoList, err := UserService.GetUserTransactions(ctx, user)

	if err == nil {
		_, err = AppEngineHelper.WriteJSON(w, txnClientDtoList)
	}

	DAOHelper.ReportError(ctx, w, err)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	user, err := UserService.CreateUser(ctx, r.Body)

	if err == nil {
		_, err = AppEngineHelper.WriteJSON(w, user)
	}

	DAOHelper.ReportError(ctx, w, err)
}

func validateAccessId(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	accessId_bytes := []byte(mux.Vars(r)[accessIdKey])

	isValid, err := accessIdValidator.ValidateAccessIdPrimary(ctx, accessId_bytes)
	if err != nil {
		DAOHelper.ReportError(ctx, w, err)
	}

	if isValid {
		w.Write(accessId_bytes)
	}
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
