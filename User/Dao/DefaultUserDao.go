package UserDao

import (
	"errors"

	"appengine"
	"appengine/datastore"

	"AppEngineHelper"
	"DAOHelper"
)

type DefaultUserDAO struct{}

var defaultUserDaoInstance = DefaultUserDAO{}

func GetInstance() UserDAO {
	return &defaultUserDaoInstance
}

const (
	USER_KIND             = "User"
	USER_PARENT_STRING_ID = "default_user"
)

const (
	PW_SALT_BYTES = 32
)

var (
	UniqueConstraint_email    = &DAOHelper.ConstraintError{Field: "email", Type: DAOHelper.UniqueConstraint}
	UniqueConstraint_accessId = &DAOHelper.ConstraintError{Field: "accessId", Type: DAOHelper.UniqueConstraint}
	Invalid_accessId          = &DAOHelper.ConstraintError{Field: "accessId", Type: DAOHelper.Invalid}

	UserNotFoundError = errors.New("User does not exist in datastore")
	userHasIdError    = errors.New("Cannot create new user, key must be nil")

	invalidSessionError   = errors.New("Invalid user session")
	passwordCanNotBeEmpty = errors.New("Can not set update password when password is blank")
)

var (
	userCollectionParentKey = AppEngineHelper.CollectionParentKeyGetFnGenerator(USER_KIND, USER_PARENT_STRING_ID, 0)
	userIntIDToKeyInt64     = AppEngineHelper.IntIDToKeyInt64(USER_KIND, userCollectionParentKey)
)

func (u *DefaultUserDAO) StringToKey(ctx appengine.Context, key string) *datastore.Key {
	return userIntIDToKeyInt64(ctx, key)
}

func (u *DefaultUserDAO) GetByEmail(ctx appengine.Context, email string) (*UserDTO, error) {
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
		return nil, UserNotFoundError
	}

	userDtoList[0].Key = keys[0]

	return &userDtoList[0], nil
}

func (u *DefaultUserDAO) GetByAccessId(ctx appengine.Context, accessId string) (*UserDTO, error) {
	q := datastore.NewQuery(USER_KIND).
		Ancestor(userCollectionParentKey(ctx)).
		Filter("AccessId=", accessId).
		Limit(1)

	userDtoList := make([]UserDTO, 0, 1)

	keys, err := q.GetAll(ctx, &userDtoList)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, UserNotFoundError
	}

	userDtoList[0].Key = keys[0]

	return &userDtoList[0], nil
}

func (u *DefaultUserDAO) GetAll(ctx appengine.Context) ([]*datastore.Key, []UserDTO, error) {
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

func (u *DefaultUserDAO) Create(ctx appengine.Context, user *UserDTO) error {

	if user.hasKey() {
		return userHasIdError
	}

	if user, _ := u.GetByEmail(ctx, user.Email); user != nil {
		if user.Verified {
			return UniqueConstraint_email
		} else {
			datastore.Delete(ctx, user.Key)
		}
	}

	if user, _ := u.GetByAccessId(ctx, user.AccessId); user != nil {
		if user.Verified {
			return UniqueConstraint_accessId
		} else {
			datastore.Delete(ctx, user.Key)
		}
	}

	if err := user.UpdatePasswordHash(user.Password); err != nil {
		return err
	}

	key := datastore.NewIncompleteKey(ctx, USER_KIND, userCollectionParentKey(ctx))
	newKey, err := datastore.Put(ctx, key, user)
	if err != nil {
		return err
	}

	user.Key = newKey

	return nil
}

func (u *DefaultUserDAO) SaveUser(ctx appengine.Context, user *UserDTO) error {
	if !user.hasKey() {
		return userHasIdError
	}

	if err := user.UpdatePasswordHash(user.Password); err != nil {
		return nil
	}

	key, err := datastore.Put(ctx, user.Key, user)
	if err == nil {
		user.Key = key
	}

	return err
}

func (u *DefaultUserDAO) SetSessionUUID(ctx appengine.Context, user *UserDTO, uuid string) error {

	user.CurrentSessionUUID = uuid

	return u.SaveUser(ctx, user)
}

func (u *DefaultUserDAO) GetUserFromSessionUUID(ctx appengine.Context, uuid string) (*UserDTO, error) {

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
