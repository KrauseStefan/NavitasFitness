package UserService

import (
	"encoding/json"
	"net/http"
	"time"

	"appengine"

	"github.com/gorilla/mux"

	"AppEngineHelper"
	"Auth"
	"IPN/Transaction"
	"User/Dao"
)

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

	router.
		Methods("POST").
		Path(path).
		Name("CreateUserInfo").
		HandlerFunc(userPost)

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

	return UserDao.GetUserFromSessionUUID(ctx, uuid)
}

func getUserFromSessionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	user, err := getUserFromSession(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if _, err := AppEngineHelper.WriteJSON(w, user); err != nil {
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
		ctx.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := UserDao.CreateUser(ctx, user); err != nil {
		if err == UserDao.EmailAlreadyExistsError {
			http.Error(w, err.Error(), http.StatusConflict)
			ctx.Infof(err.Error())
			return
		}
		ctx.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := AppEngineHelper.WriteJSON(w, user); err != nil {
		ctx.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getUserTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	user, err := getUserFromSession(ctx, r)
	if err != nil {
		// User is not authorised sending back an empty response
		return
	}

	transactions, err := TransactionDao.GetTransactionsByUser(ctx, user.GetDataStoreKey(ctx))
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
