package UserService

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"appengine"

	"Auth"
	"IPN/Transaction"
	"User/Dao"
)

var (
	userDao        = UserDao.GetInstance()
	transactionDao = TransactionDao.GetInstance()
)

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

func CreateUser(ctx appengine.Context, respBody io.ReadCloser) (*UserDao.UserDTO, error) {
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

	if err := SendConfirmationMail(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserTransactions(ctx appengine.Context, user *UserDao.UserDTO) ([]*TransactionMsgClientDTO, error) {
	transactions, err := transactionDao.GetTransactionsByUser(ctx, user.Key)
	if err != nil {
		return nil, err
	}

	txnClientDtoList := make([]*TransactionMsgClientDTO, len(transactions))
	for i, txn := range transactions {
		txnClientDtoList[i] = newTransactionMsgClientDTO(txn)
	}

	return txnClientDtoList, nil
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

func MarkUserVerified(ctx appengine.Context, encodedKey string) error {

	return userDao.MarkUserVerified(ctx, encodedKey)

}
