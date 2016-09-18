package UserDao

import (
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"

	"AppEngineHelper"
)

const (
	USER_KIND             = "User"
	USER_PARENT_STRING_ID = "default_user"
)

const (
	PW_SALT_BYTES = 32
)

var (
	userHasIdError         = errors.New("Cannot create new user, key must be nil")
	userHasNoIdError       = errors.New("Cannot create new user, key must be defined")
	userAlreadyExistsError = errors.New("Cannot update an already existing user")
	userNotFoundError      = errors.New("User does not exist in datastore")
	invalidSessionError    = errors.New("Invalid user session")
	passwordCanNotBeEmpty  = errors.New("Can not set update password when password is blank")
)

var (
	userCollectionParentKey = AppEngineHelper.CollectionParentKeyGetFnGenerator(USER_KIND, USER_PARENT_STRING_ID, 0)
	userIntIDToKeyInt64     = AppEngineHelper.IntIDToKeyInt64(USER_KIND, userCollectionParentKey)
)

type UserDTO struct {
	Key                string    `json:"key",datastore:"-"`
	Email              string    `json:"email"`
	Password           string    `json:"password,omitempty",datastore:",noindex"`
	PasswordHash       []byte    `json:"-",datastore:",noindex"`
	PasswordSalt       []byte    `json:"-",datastore:",noindex"`
	NavitasId          string    `json:"navitasId"`
	CreatedDate        time.Time `json:"createdDate"`
	CurrentSessionUUID string    `json:"currentSessionKey"`
	IsAdmin            bool      `json:"isAdmin,omitempty"`
}

func (user *UserDTO) hasKey() bool {
	return len(user.Key) > 0
}

func (user *UserDTO) GetDataStoreKey(ctx appengine.Context) *datastore.Key {
	return StringToKey(ctx, user.Key)
}

func (user *UserDTO) setKey(key *datastore.Key) *UserDTO {
	user.Key = strconv.FormatInt(key.IntID(), 10)
	return user
}

func StringToKey(ctx appengine.Context, key string) *datastore.Key {
	return userIntIDToKeyInt64(ctx, key)
}

func (user *UserDTO) getPasswordWithSalt(password []byte) []byte {
	return append(user.PasswordSalt, password...)
}

func (user *UserDTO) UpdatePasswordHash(password []byte) error {
	if password == nil && user.Password != "" {
		password = []byte(user.Password)
	}
	if password == nil {
		return passwordCanNotBeEmpty
	}
	// https://crackstation.net/hashing-security.htm
	user.PasswordSalt = make([]byte, PW_SALT_BYTES)
	rand.Read(user.PasswordSalt)

	passwordHash, err := bcrypt.GenerateFromPassword(user.getPasswordWithSalt(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash
	user.Password = ""

	return nil
}

func (user *UserDTO) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(user.PasswordHash, user.getPasswordWithSalt([]byte(password))) == nil
}

func GetUserByEmail(ctx appengine.Context, email string) (*UserDTO, error) {
	q := datastore.NewQuery(USER_KIND).
		Ancestor(userCollectionParentKey(ctx)).
		Filter("Email=", email).
		Limit(1)

	userDtoList := make([]UserDTO, 0, 1)

	keys, err := q.GetAll(ctx, &userDtoList)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, userNotFoundError
	}

	userDtoList[0].Key = strconv.FormatInt(keys[0].IntID(), 10)

	return &userDtoList[0], nil
}

func GetAllUsers(ctx appengine.Context) ([]*datastore.Key, []UserDTO, error) {
	query := datastore.NewQuery(USER_KIND).
		Ancestor(userCollectionParentKey(ctx))

	count, err := query.Count(ctx)
	if err != nil {
		return nil, nil, err
	}
	if count <= 0 {
		return nil, nil, nil
	}

	users := make([]UserDTO, 0, count)
	keys, err := query.GetAll(ctx, &users)
	if err != nil {
		ctx.Criticalf("error in txn 2")
		return nil, nil, err
	}

	return keys, users, nil
}

func CreateUser(ctx appengine.Context, user *UserDTO) error {

	if user.hasKey() {
		return userHasIdError
	}

	if user, _ := GetUserByEmail(ctx, user.Email); user != nil {
		return userAlreadyExistsError
	}

	if err := user.UpdatePasswordHash(nil); err != nil {
		return err
	}

	key := datastore.NewIncompleteKey(ctx, USER_KIND, userCollectionParentKey(ctx))
	newKey, err := datastore.Put(ctx, key, user)
	if err != nil {
		return err
	}

	user.Key = strconv.FormatInt(newKey.IntID(), 10)

	return nil
}

func saveUser(ctx appengine.Context, user *UserDTO) error {
	if !user.hasKey() {
		return userHasIdError
	}

	//Only updates if password field has been set
	user.UpdatePasswordHash(nil)

	key, err := datastore.Put(ctx, user.GetDataStoreKey(ctx), user)

	if err == nil {
		user.setKey(key)
	}

	return err
}

func SetSessionUUID(ctx appengine.Context, user *UserDTO, uuid string) error {

	user.CurrentSessionUUID = uuid

	return saveUser(ctx, user)
}

func GetUserFromSessionUUID(ctx appengine.Context, uuid string) (*UserDTO, error) {

	users := make([]UserDTO, 0, 2)

	keys, err := datastore.NewQuery(USER_KIND).
		Ancestor(userCollectionParentKey(ctx)).
		Filter("CurrentSessionUUID =", uuid).
		Limit(2).
		GetAll(ctx, &users)

	if err != nil {
		return nil, err
	} else if len(keys) != 1 {
		return nil, errors.New(invalidSessionError.Error() + " - uuid: " + uuid)
	}

	return &users[0], nil
}
