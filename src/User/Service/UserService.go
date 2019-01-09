package UserService

import (
	"AppEngineHelper"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"Auth"
	"DAOHelper"
	"IPN/Transaction"
	"User/Dao"
)

var (
	userDao        = UserDao.GetInstance()
	transactionDao = TransactionDao.GetInstance()
)

type CustomHttpHandler = func(http.ResponseWriter, *http.Request, *UserDao.UserDTO) (interface{}, error)

func AsAdmin(f CustomHttpHandler) func(http.ResponseWriter, *http.Request) {
	return AsUser(func(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
		if user == nil || !user.IsAdmin {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return nil, nil
		}
		return f(w, r, user)
	})
}

func AsUser(f CustomHttpHandler) func(http.ResponseWriter, *http.Request) {
	return AppEngineHelper.HandlerW(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		ctx := appengine.NewContext(r)

		user, err := getUserFromSession(ctx, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return nil, nil
		}

		return f(w, r, user)
	})
}

func getUserFromSession(ctx context.Context, r *http.Request) (*UserDao.UserDTO, error) {
	sessionData, err := Auth.GetSessionData(r)

	if err == nil && sessionData.HasLoginInfo() {
		return userDao.GetUserFromSessionUUID(ctx, sessionData.UserKey, sessionData.Uuid)
	}

	return nil, err
}

func GetAllUsers(ctx context.Context) ([]string, []*UserDao.UserDTO, error) {
	keys, users, err := userDao.GetAll(ctx)
	keyStrings := make([]string, len(keys))

	for i, key := range keys {
		keyStrings[i] = key.Encode()
	}

	return keyStrings, users, err
}

// This function tries its best to validate and ensure no user is created with duplicated accessId or email
func CreateUser(ctx context.Context, user *UserDao.UserDTO, sessionData Auth.SessionData) (*UserDao.UserDTO, error) {

	if err := user.ValidateUser(ctx); err != nil {
		return nil, err
	}

	if err := userDao.Create(ctx, user, sessionData.UserKey); err != nil {
		return nil, err
	}

	if err := SendConfirmationMail(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserTransactions(ctx context.Context, userKey *datastore.Key) ([]*TransactionMsgClientDTO, error) {
	transactions, err := transactionDao.GetTransactionsByUser(ctx, userKey)
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

func MarkUserVerified(ctx context.Context, encodedKey string) error {
	key, err := datastore.DecodeKey(encodedKey)
	if err != nil {
		return err
	}

	userDto := &UserDao.UserDTO{}
	if err := datastore.Get(ctx, key, userDto); err != nil {
		return err
	}

	userDto.Verified = true

	if _, err := datastore.Put(ctx, key, userDto); err != nil {
		return err
	}

	log.Infof(ctx, "Verified email: %s for user with accessId: %s", userDto.Email, userDto.AccessId)

	return nil
}

func RandString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func RequestResetUserPassword(ctx context.Context, email string) error {
	rndStr := RandString(10)

	user, err := userDao.GetByEmail(ctx, email)
	if err == UserDao.UserNotFoundError {
		return &DAOHelper.DefaultHttpError{
			InnerError: err,
			StatusCode: http.StatusNotFound,
		}
	} else if err != nil {
		return err
	}

	user.PasswordResetTime = time.Now()
	user.PasswordResetSecret = rndStr

	if err := userDao.SaveUser(ctx, user); err != nil {
		return err
	}

	if err := SendPasswordResetMail(ctx, user, rndStr); err != nil {
		return err
	}

	return nil
}

type PasswordChangeDto struct {
	Key      string `json:"key"`
	Password string `json:"password"`
	Secret   string `json:"secret"`
}

var resetInputInvalidError = &DAOHelper.DefaultHttpError{
	StatusCode: http.StatusBadRequest,
	InnerError: errors.New("Invalid password reset token"),
}

func ResetUserPassword(ctx context.Context, respBody io.ReadCloser) error {
	dto := &PasswordChangeDto{}
	user := &UserDao.UserDTO{}
	maxAge := time.Now().Add(time.Minute * -30)

	decoder := json.NewDecoder(respBody)
	if err := decoder.Decode(dto); err != nil {
		return resetInputInvalidError
	}

	key, err := datastore.DecodeKey(dto.Key)
	if err != nil {
		return resetInputInvalidError
	}

	if err := datastore.Get(ctx, key, user); err != nil {
		return resetInputInvalidError
	}

	if user.PasswordResetSecret == "" || user.PasswordResetSecret != dto.Secret || !user.PasswordResetTime.After(maxAge) {
		log.Infof(ctx, "serects user: %q", user.PasswordResetSecret)
		log.Infof(ctx, "serects dto : %q", dto.Secret)
		log.Infof(ctx, "serects equals: %v", dto.Secret == user.PasswordResetSecret)
		log.Infof(ctx, "PasswordResetTime: %v, should be after Maxage: %v", user.PasswordResetTime.Format(time.Stamp), maxAge.Format(time.Stamp))
		log.Infof(ctx, "PasswordResetTime ok: %v", user.PasswordResetTime.After(maxAge))
		return resetInputInvalidError
	}

	if err := user.UpdatePasswordHash(dto.Password); err != nil {
		return err
	}

	user.PasswordResetTime = time.Time{}
	user.PasswordResetSecret = ""

	if _, err := datastore.Put(ctx, key, user); err != nil {
		return err
	}

	return nil
}
