package UserService
import (
	"strconv"

	"appengine/datastore"
	"appengine"

	"src/Services/Common"
	"errors"
)


var (
	userHasIdError = errors.New("Cannot create new user with id")
	userAlreadyExistsError = errors.New("Cannot update an already existing user")
	userNotFoundError = errors.New("User does not exist in datastore")
)

type UserDTO struct {
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	password  string `datastore:"noindex"`
	NavitasId string `json:"navitasId"`
	Id        string `json:"id",datastore:"-"`
}

func (user UserDTO) HasId() bool {
	return len(user.Id) == 0
}

const USER_KIND = "User"
const PARENT_STRING_ID = "default_user"

var userCollectionParentKey = Common.CollectionParentKeyGetFnGenerator(USER_KIND, PARENT_STRING_ID, 0)
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

func GetUserByUserName(ctx appengine.Context, userName string) (*UserDTO, error) {
	q := datastore.NewQuery(USER_KIND).
	Ancestor(userCollectionParentKey(ctx)).
	Filter("UserName=", userName).
	Limit(1)

	userDtoList := make([]UserDTO, 0, 1)

	keys, err := q.GetAll(ctx, &userDtoList)
	if err != nil {
		return nil, err
	}

	if (len(keys) == 0) {
		return nil, userNotFoundError
	}

	userDtoList[0].Id = strconv.FormatInt(keys[0].IntID(), 10)

	return &userDtoList[0], nil
}

func CreateUser(ctx appengine.Context, user *UserDTO) error {

	if user.HasId() {
		return userHasIdError
	}

	if _, err := GetUserById(ctx, user.Id); err != userNotFoundError {
		if (err != nil) {
			return err
		}
		return userAlreadyExistsError
	}

	putUser(ctx, user)

	return nil
}

//TODO: figure out if this should be possible
func UpdateUser(ctx appengine.Context, user *UserDTO) error {

	if _, err := GetUserById(ctx, user.Id); err == nil {
		return err
	}

	putUser(ctx, user)

	return nil
}

func putUser(ctx appengine.Context, user *UserDTO) (*	datastore.Key, error) {
	var (
		err error
		key *datastore.Key
	)

	if (user.HasId()) {
		key = datastore.NewIncompleteKey(ctx, USER_KIND, userCollectionParentKey(ctx))
	} else {
		key = userIntIDToKeyInt64(ctx, user.Id)
	}

	if key, err = datastore.Put(ctx, key, &user); err != nil {
		return nil, err
	}

	user.Id = strconv.FormatInt(key.IntID(), 10)

	return key, nil
}
