package UserService

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

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

func AsAdmin(f func(http.ResponseWriter, *http.Request, *UserDao.UserDTO)) func(http.ResponseWriter, *http.Request) {
	return AsUser(func(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) {
		if user == nil || !user.IsAdmin {
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

func getUserFromSession(ctx context.Context, r *http.Request) (*UserDao.UserDTO, error) {
	sessionData, err := Auth.GetSessionData(r)
	if err != nil {
		return nil, err
	}

	if !sessionData.HasLoginInfo() && err == nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return userDao.GetUserFromSessionUUID(ctx, sessionData.UserKey, sessionData.Uuid)
}

func GetAllUsers(ctx context.Context) ([]string, []UserDao.UserDTO, error) {
	keys, users, err := userDao.GetAll(ctx)
	keyStrings := make([]string, len(keys))

	for i, key := range keys {
		keyStrings[i] = key.Encode()
	}

	return keyStrings, users, err
}

func GetDuplicatedUsers(ctx context.Context) ([]string, []UserDao.UserDTO, error) {
	keys, users, err := userDao.GetAll(ctx)

	sort.Sort(UserDao.ByAccessId(users))

	prev := users[0]
	prevKey := keys[0]

	filteredUsers := make([]UserDao.UserDTO, 0, len(users))
	keyStrings := make([]string, 0, len(keys))

	for i, user := range users[1:] {
		key := keys[i+1]

		if user.IsEquivalent(&prev) {
			previousIndex := len(filteredUsers) - 1

			if previousIndex < 0 || !filteredUsers[previousIndex].IsEquivalent(&prev) {
				keyStrings = append(keyStrings, prevKey.Encode())
				filteredUsers = append(filteredUsers, prev)
			}

			keyStrings = append(keyStrings, key.Encode())
			filteredUsers = append(filteredUsers, user)
		}
		prev = user
	}

	return keyStrings, filteredUsers, err
}

// This function tries its best to validate and ensure no user is created with duplicated accessId or email
func CreateUser(ctx context.Context, r *http.Request, sessionData Auth.SessionData) (*UserDao.UserDTO, error) {
	user := &UserDao.UserDTO{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(user); err != nil {
		return nil, err
	}

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

func GetDuplicatedInactiveUsers(ctx context.Context, ids []string) ([]string, error) {
	createError := func(msg string) ([]string, error) {
		return nil, &DAOHelper.DefaultHttpError{
			InnerError: errors.New(msg),
			StatusCode: http.StatusBadRequest,
		}
	}

	idKeys, err := AppEngineHelper.StringIdsToDsKeys(ids)
	if err != nil {
		return nil, err
	}

	users, err := userDao.GetByKeys(ctx, idKeys)
	if err != nil {
		return nil, err
	}

	accessId := users[0].AccessId
	email := users[0].Email
	name := users[0].Name
	sessionId := users[0].CurrentSessionUUID
	// verified := users[0].Verified

	usersIdsToDelete := make([]string, 0, len(users)-1)

	transactions, err := transactionDao.GetTransactionsByUser(ctx, idKeys[0])
	if err != nil {
		return nil, err
	}

	if len(transactions) <= 0 && !users[0].Verified {
		log.Errorf(ctx, "Marking for deletion %v", users[0])
		usersIdsToDelete = append(usersIdsToDelete, idKeys[0].Encode())
	}

	for i, user := range users[1:] {
		key := idKeys[i+1]
		if accessId != user.AccessId || email != user.Email || name != user.Name {
			return createError("All the ides does not match in terms of AccessId, Email and Name, aborting")
		}

		if len(sessionId) == 0 && len(user.CurrentSessionUUID) != 0 {
			sessionId = user.CurrentSessionUUID
		} else if len(sessionId) != 0 && len(user.CurrentSessionUUID) != 0 {
			return createError("Multiple users has active user sessions, aborting")
		}

		userTransactions, err := transactionDao.GetTransactionsByUser(ctx, key)
		if err != nil {
			return nil, err
		}

		log.Errorf(ctx, "Transactions - prev: %d - current: %d", len(transactions), len(userTransactions))
		if len(transactions) != 0 && len(userTransactions) != 0 {
			return createError("Multiple users to be merged had transactions, aborting")
		} else {
			transactions = userTransactions
		}

		if len(transactions) <= 0 {
			log.Errorf(ctx, "Marking for deletion %v", user)
			usersIdsToDelete = append(usersIdsToDelete, key.Encode())
		}
	}

	return usersIdsToDelete, nil
}

func DeleteInactiveUsers(ctx context.Context, ids []string) error {
	idKeys, err := AppEngineHelper.StringIdsToDsKeys(ids)
	if err != nil {
		return err
	}

	for i, idKey := range idKeys {
		userTransactions, err := transactionDao.GetTransactionsByUser(ctx, idKey)
		if err != nil {
			return err
		}

		if len(userTransactions) > 0 {
			return errors.New("User has transactions, deletion aborted, userId: " + ids[i])
		}
	}

	return userDao.DeleteUsers(ctx, idKeys)
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
