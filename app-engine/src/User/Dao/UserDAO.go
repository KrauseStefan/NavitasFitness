package UserDao

import (
	"strconv"

	"appengine"
	"appengine/datastore"

	"errors"
	"src/Common"
	"time"
)

var (
	userHasIdError					= errors.New("Cannot create new user, key must be nil")
	userHasNoIDError				= errors.New("Cannot create new user, key must be defined")
	userAlreadyExistsError	= errors.New("Cannot update an already existing user")
	userNotFoundError				= errors.New("User does not exist in datastore")
	invalidSessionError			= errors.New("Invalid user session")
)

const USER_KIND = "User"
const USER_PARENT_STRING_ID = "default_user"

var userCollectionParentKey = Common.CollectionParentKeyGetFnGenerator(USER_KIND, USER_PARENT_STRING_ID, 0)
var userIntIDToKeyInt64 = Common.IntIDToKeyInt64(USER_KIND, userCollectionParentKey)

type UserDTO struct {
	Key									string 		`json:"key",datastore:"-"`
	Email								string 		`json:"email"`
	Password						string 		`json:"password",datastore:",noindex"`
	NavitasId						string 		`json:"navitasId"`
	CreatedDate					time.Time	`json:"createdDate"`
	LastLogin						time.Time	`json:"lastLogin"`
	CurrentSessionUUID	string 		`json:"currentSessionKey"`
	IsAdmin							string		`json:"isAdmin,omitempty"`
}

func (user UserDTO) hasKey() bool {
	return len(user.Key) > 0
}

func (user UserDTO) getDataStoreKey(ctx appengine.Context) *datastore.Key {
	return userIntIDToKeyInt64(ctx, user.Key)
}

func(user UserDTO) setKey(key *datastore.Key) UserDTO {
	user.Key = strconv.FormatInt(key.IntID(), 10)
	return user
}

func GetUserByKey(ctx appengine.Context, key *datastore.Key) (*UserDTO, error) {
	var user UserDTO

	if err := datastore.Get(ctx, key, &user); err != nil {
		return nil, err
	}

	user.setKey(key)
	return &user, nil
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

func CreateUser(ctx appengine.Context, user *UserDTO) error {

	if user.hasKey() {
		return userHasIdError
	}

	key := datastore.NewIncompleteKey(ctx, USER_KIND, userCollectionParentKey(ctx))
	newKey, err := datastore.Put(ctx, key, user)
	if  err != nil {
		return err
	}

	user.Key = strconv.FormatInt(newKey.IntID(), 10)

	return nil
}

func saveUser(ctx appengine.Context, user *UserDTO) error {
	if !user.hasKey() {
		return userHasIdError
	}

	key, err := datastore.Put(ctx, user.getDataStoreKey(ctx), user)

	if err == nil {
		user.setKey(key)
	}

	return err
}

func SetSessionUUID(ctx appengine.Context, user *UserDTO, uuid string) error {

	user.CurrentSessionUUID = uuid

	return saveUser(ctx, user)
}

func GetUserFromSessionUUID(ctx appengine.Context, uuid string) (*UserDTO, error){

	users := make([]UserDTO, 0, 2)

	keys, err := datastore.NewQuery(USER_KIND).
		Ancestor(userCollectionParentKey(ctx)).
		Filter("CurrentSessionUUID =", uuid).
		Limit(2).
		GetAll(ctx, &users)

	if(err != nil) {
		return nil, err
	} else if len(keys) != 1 {
		return nil, errors.New(invalidSessionError.Error() + " - uuid: " + uuid)
	}

	return &users[0], nil
}
