package UserDao

import (
	"errors"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

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
	userHasNoIdError  = errors.New("Cannot update user without ID")

	invalidSessionError   = errors.New("Invalid user session")
	passwordCanNotBeEmpty = errors.New("Can not set update password when password is blank")
)

var (
	userCollectionParentKey = AppEngineHelper.CollectionParentKeyGetFnGenerator(USER_KIND, USER_PARENT_STRING_ID, 0)
	userIntIDToKeyInt64     = AppEngineHelper.IntIDToKeyInt64(USER_KIND, userCollectionParentKey)
)

func (u *DefaultUserDAO) StringToKey(ctx context.Context, key string) *datastore.Key {
	return userIntIDToKeyInt64(ctx, key)
}

func (u *DefaultUserDAO) GetByEmail(ctx context.Context, email string) (*UserDTO, error) {
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

func (u *DefaultUserDAO) GetByAccessId(ctx context.Context, accessId string) (*UserDTO, error) {
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

func (u *DefaultUserDAO) GetAll(ctx context.Context) ([]*datastore.Key, []UserDTO, error) {
	query := datastore.NewQuery(USER_KIND).
		Ancestor(userCollectionParentKey(ctx))

	users := new([]UserDTO)
	keys, err := query.GetAll(ctx, users)
	if err != nil {
		return nil, nil, err
	}

	return keys, *users, nil
}

// This function tries its best to ensure no user is created with duplicated accessId or email
func (u *DefaultUserDAO) Create(ctx context.Context, user *UserDTO, keyHint *datastore.Key) error {

	if existingUser, err := u.GetByKey(ctx, keyHint); existingUser != nil && err == nil {
		if existingUser.Email == user.Email && existingUser.Verified {
			return UniqueConstraint_email
		} else if existingUser.AccessId == user.AccessId && existingUser.Verified {
			return UniqueConstraint_accessId
		}

		// At this point it is concluded that the existing user found by the key hit is not equivalent
	}

	if user, _ := u.GetByEmail(ctx, user.Email); user != nil {
		if user.Verified {
			return UniqueConstraint_email
		} else if err := datastore.Delete(ctx, user.Key); err != nil {
			return err
		}
	}

	if user, _ := u.GetByAccessId(ctx, user.AccessId); user != nil {
		if user.Verified {
			return UniqueConstraint_accessId
		} else if err := datastore.Delete(ctx, user.Key); err != nil {
			return err
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

func (u *DefaultUserDAO) GetByKey(ctx context.Context, key *datastore.Key) (*UserDTO, error) {
	var user = &UserDTO{}

	if err := datastore.Get(ctx, key, user); err != nil {
		return nil, err
	}

	user.Key = key
	return user, nil
}

func (u *DefaultUserDAO) SaveUser(ctx context.Context, user *UserDTO) error {
	if !user.hasKey() {
		return userHasNoIdError
	}

	if user.Password != "" {
		if err := user.UpdatePasswordHash(user.Password); err != nil {
			return err
		}
	}

	key, err := datastore.Put(ctx, user.Key, user)
	if err != nil {
		return err
	}
	user.Key = key
	return nil
}

func (u *DefaultUserDAO) SetSessionUUID(ctx context.Context, user *UserDTO, uuid string) error {

	user.CurrentSessionUUID = uuid

	return u.SaveUser(ctx, user)
}

func (u *DefaultUserDAO) GetUserFromSessionUUID(ctx context.Context, userKey *datastore.Key, uuid string) (*UserDTO, error) {
	user := &UserDTO{}

	if err := datastore.Get(ctx, userKey, user); err != nil {
		return nil, err
	}

	if user.CurrentSessionUUID != uuid {
		return nil, invalidSessionError
	}

	user.Key = userKey
	return user, nil
}
