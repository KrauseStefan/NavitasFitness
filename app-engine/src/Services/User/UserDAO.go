package UserService

import (
	"strconv"

	"appengine"
	"appengine/datastore"

	"errors"
	"src/Services/Common"
	"time"
)

var (
	userHasIdError         = errors.New("Cannot create new user, Id must be nil")
	userAlreadyExistsError = errors.New("Cannot update an already existing user")
	userNotFoundError      = errors.New("User does not exist in datastore")
)

type UserDTO struct {
	Email     	string 		`json:"email"`
	Password  	string 		`json:"password",datastore:",noindex"`
	NavitasId 	string 		`json:"navitasId"`
	CreatedDate	time.Time	`json:"createdDate"`
	Id        	string 		`json:"id",datastore:"-"`
}

func (user UserDTO) hasId() bool {
	return len(user.Id) > 0
}

const USER_KIND = "User"
const USER_PARENT_STRING_ID = "default_user"

var userCollectionParentKey = Common.CollectionParentKeyGetFnGenerator(USER_KIND, USER_PARENT_STRING_ID, 0)
var userIntIDToKeyInt64 = Common.IntIDToKeyInt64(USER_KIND, userCollectionParentKey)

func GetUserById(ctx appengine.Context, Id string) (*UserDTO, error) {
	var user UserDTO
	//	user := new(UserDTO)
	key := userIntIDToKeyInt64(ctx, Id)

	if err := datastore.Get(ctx, key, &user); err != nil {
		return nil, err
	}

	user.Id = strconv.FormatInt(key.IntID(), 10)
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

	userDtoList[0].Id = strconv.FormatInt(keys[0].IntID(), 10)

	return &userDtoList[0], nil
}

func CreateUser(ctx appengine.Context, user *UserDTO) error {

	if user.hasId() {
		return userHasIdError
	}

	key := datastore.NewIncompleteKey(ctx, USER_KIND, userCollectionParentKey(ctx))
	newKey, err := datastore.Put(ctx, key, user)
	if  err != nil {
		return err
	}

	user.Id = strconv.FormatInt(newKey.IntID(), 10)

	return nil
}

//TODO: figure out if this should be possible
//func UpdateUser(ctx appengine.Context, user *UserDTO) error {
//
//	if _, err := GetUserById(ctx, user.Id); err == nil {
//		return err
//	}
//
//	putUser(ctx, user)
//
//	return nil
//}
