package UserService

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"appengine"

	"github.com/gorilla/mux"

	"AccessIdValidator"
	"AppEngineHelper"
	"Auth"
	"DAOHelper"
	"IPN/Transaction"
	"User/Dao"
)

var (
	userDao        = UserDao.GetInstance()
	transactionDao = TransactionDao.GetInstance()
)

const accessIdKey = "accessId"

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/user"

	router.
		Methods("GET").
		Path(path).
		Name("Get User Current User Info").
		HandlerFunc(AsUser(getUserFromSessionHandler))

	router.
		Methods("GET").
		Path(path + "/transactions").
		Name("Get Latest Transactions").
		HandlerFunc(AsUser(getUserTransactionsHandler))

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

func AsAdmin(f func(http.ResponseWriter, *http.Request, *UserDao.UserDTO)) func(http.ResponseWriter, *http.Request) {
	return AsUser(func(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) {
		if !user.IsAdmin {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		f(w, r, user)
	})
}

func AsUser(f func(http.ResponseWriter, *http.Request, *UserDao.UserDTO)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)

		user, err := getUserFromSession(ctx, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		f(w, r, user)
	}
}

func getUserFromSession(ctx appengine.Context, r *http.Request) (*UserDao.UserDTO, error) {
	uuid, err := AuthService.GetSessionUUID(r)

	if uuid == "" && err == nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return userDao.GetUserFromSessionUUID(ctx, uuid)
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

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	user, err := createUser(ctx, r.Body)

	if err == nil {
		_, err = AppEngineHelper.WriteJSON(w, user)
	}

	DAOHelper.ReportError(ctx, w, err)
}

func createUser(ctx appengine.Context, respBody io.ReadCloser) (*UserDao.UserDTO, error) {
	user := &UserDao.UserDTO{}

	decoder := json.NewDecoder(respBody)
	if err := decoder.Decode(user); err != nil {
		return nil, err
	}

	if err := user.ValidateUser(ctx); err != nil {
		return nil, err
	}

	if err := userDao.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func getUserTransactionsHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) {
	ctx := appengine.NewContext(r)

	transactions, err := transactionDao.GetTransactionsByUser(ctx, user.GetDataStoreKey(ctx))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	txnClientDtoList := make([]*TransactionMsgClientDTO, len(transactions))

	for i, txn := range transactions {
		txnClientDtoList[i] = newTransactionMsgClientDTO(txn)
	}

	if _, err := AppEngineHelper.WriteJSON(w, txnClientDtoList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type TransactionMsgClientDTO struct {
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	PaymentDate time.Time `json:"paymentDate"`
	Status      string    `json:"status"`
	IsActive    bool      `json:"isActive"`
	// IpnMessages           []string  `json:"ipnMessages"` // History of IpnMessages
}

func newTransactionMsgClientDTO(source *TransactionDao.TransactionMsgDTO) *TransactionMsgClientDTO {

	txClient := TransactionMsgClientDTO{
		Amount:      source.GetAmount(),
		Currency:    source.GetCurrency(),
		PaymentDate: source.GetPaymentDate(),
		Status:      source.GetPaymentStatus(),
		IsActive:    source.IsActive(),
		// IpnMessages:           source.GetIpnMessages(),
	}

	return &txClient
}
