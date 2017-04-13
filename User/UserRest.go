package UserRest

import (
	"net/http"

	"appengine"

	"github.com/gorilla/mux"

	"AccessIdValidator"
	"AppEngineHelper"
	"DAOHelper"
	"User/Dao"
	"User/Service"
)

const accessIdKey = "accessId"

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
}

type UserSessionDto struct {
	User    *UserDao.UserDTO `json:"user"`
	IsAdmin bool             `json:"isAdmin"`
}

func getUserFromSessionHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) {
	us := UserSessionDto{user, user.IsAdmin}

	if _, err := AppEngineHelper.WriteJSON(w, us); err != nil {
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

	isValid, err := AccessIdValidator.ValidateAccessId(ctx, accessId_bytes)
	if err != nil {
		ctx.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
