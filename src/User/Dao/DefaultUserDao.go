package UserDao

import (
	"errors"
	"strconv"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"

	"DAOHelper"
	nf_datastore "NavitasFitness/datastore"
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
	userCollectionParentKey = datastore.NameKey(USER_KIND, USER_PARENT_STRING_ID, nil)
)

func (u *DefaultUserDAO) StringToKey(ctx context.Context, key string) *datastore.Key {
	intId, _ := strconv.ParseInt(key, 10, 64)
	return datastore.IDKey(USER_KIND, intId, userCollectionParentKey)
}

func (u *DefaultUserDAO) GetByEmail(ctx context.Context, email string) (*UserDTO, error) {
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery(USER_KIND).
		Ancestor(userCollectionParentKey).
		Filter("Email=", email).
		Limit(1)

	userDtoList := make([]UserDTO, 0, 1)

	keys, err := dsClient.GetAll(ctx, query, &userDtoList)
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
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery(USER_KIND).
		Ancestor(userCollectionParentKey).
		Filter("AccessId=", accessId).
		Limit(1)

	userDtoList := make([]UserDTO, 0, 1)

	keys, err := dsClient.GetAll(ctx, query, &userDtoList)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, UserNotFoundError
	}

	userDtoList[0].Key = keys[0]

	return &userDtoList[0], nil
}

func (u *DefaultUserDAO) GetAll(ctx context.Context) ([]*datastore.Key, UserList, error) {
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, nil, err
	}

	query := datastore.NewQuery(USER_KIND).
		Ancestor(userCollectionParentKey)

	users := new(UserList)
	keys, err := dsClient.GetAll(ctx, query, users)
	if err != nil {
		return keys, *users, err
	}

	return keys, *users, nil
}

func (u *DefaultUserDAO) GetByKeys(ctx context.Context, keys []*datastore.Key) (UserList, error) {
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, err
	}

	users := make(UserList, len(keys))
	err = dsClient.GetMulti(ctx, keys, users)

	for i, user := range users {
		if user != nil {
			user.Key = keys[i]
		}
	}

	return users, err
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

	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return err
	}
	if user, _ := u.GetByEmail(ctx, user.Email); user != nil {
		if user.Verified {
			return UniqueConstraint_email
		} else if err := dsClient.Delete(ctx, user.Key); err != nil {
			return err
		}
	}

	if user, _ := u.GetByAccessId(ctx, user.AccessId); user != nil {
		if user.Verified {
			return UniqueConstraint_accessId
		} else if err := dsClient.Delete(ctx, user.Key); err != nil {
			return err
		}
	}

	if err := user.UpdatePasswordHash(user.Password); err != nil {
		return err
	}

	key := datastore.IncompleteKey(USER_KIND, userCollectionParentKey)
	newKey, err := dsClient.Put(ctx, key, user)
	if err != nil {
		return err
	}

	user.Key = newKey

	return nil
}

func (u *DefaultUserDAO) DeleteUsers(ctx context.Context, ids []*datastore.Key) error {
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return err
	}
	return dsClient.DeleteMulti(ctx, ids)
}

func (u *DefaultUserDAO) GetByKey(ctx context.Context, key *datastore.Key) (*UserDTO, error) {
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, err
	}
	var user = &UserDTO{}
	if err := dsClient.Get(ctx, key, user); err != nil {
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

	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return err
	}

	key, err := dsClient.Put(ctx, user.Key, user)
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
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, err
	}

	user := &UserDTO{}

	err = dsClient.Get(ctx, userKey, user)
	if err == datastore.ErrNoSuchEntity {
		err = nil
	} else if err != nil {
		return nil, err
	}

	if user.CurrentSessionUUID != uuid {
		return nil, invalidSessionError
	}

	user.Key = userKey
	return user, nil
}
