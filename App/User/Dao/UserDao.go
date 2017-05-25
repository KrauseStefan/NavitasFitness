package UserDao

import (
	"appengine"
	"appengine/datastore"
)

type UserCreator interface {
	Create(ctx appengine.Context, user *UserDTO) error
}

type SingleUserRetriever interface {
	GetByEmail(ctx appengine.Context, email string) (*UserDTO, error)
	GetByAccessId(ctx appengine.Context, accessId string) (*UserDTO, error)
}

type UsersRetriever interface {
	GetAll(ctx appengine.Context) ([]*datastore.Key, []UserDTO, error)
}

type UserModifier interface {
	SaveUser(ctx appengine.Context, user *UserDTO) error
	SetSessionUUID(ctx appengine.Context, user *UserDTO, uuid string) error
}

type UserSessionRetriever interface {
	GetUserFromSessionUUID(ctx appengine.Context, uuid string) (*UserDTO, error)
}

type UserDAO interface {
	UserCreator
	SingleUserRetriever
	UsersRetriever
	UserSessionRetriever
	UserModifier
}
