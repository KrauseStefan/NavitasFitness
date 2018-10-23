package UserDao

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type UserCreator interface {
	Create(ctx context.Context, user *UserDTO, keyHint *datastore.Key) error
}

type SingleUserRetriever interface {
	GetByEmail(ctx context.Context, email string) (*UserDTO, error)
	GetByAccessId(ctx context.Context, accessId string) (*UserDTO, error)
}

type UsersRetriever interface {
	GetAll(ctx context.Context) ([]*datastore.Key, []*UserDTO, error)
	GetByKey(ctx context.Context, key *datastore.Key) (*UserDTO, error)
	GetByKeys(ctx context.Context, keys []*datastore.Key) ([]*UserDTO, error)
}

type UserModifier interface {
	SaveUser(ctx context.Context, user *UserDTO) error
	SetSessionUUID(ctx context.Context, user *UserDTO, uuid string) error
	DeleteUsers(ctx context.Context, ids []*datastore.Key) error
}

type UserSessionRetriever interface {
	GetUserFromSessionUUID(ctx context.Context, userKey *datastore.Key, uuid string) (*UserDTO, error)
}

type UserDAO interface {
	UserCreator
	SingleUserRetriever
	UsersRetriever
	UserSessionRetriever
	UserModifier
}
