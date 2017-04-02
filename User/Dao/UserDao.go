package UserDao

import (
	"appengine"
	"appengine/datastore"

	"AppEngineHelper"
)

type UserCreator interface {
	Create(ctx appengine.Context, user *UserDTO) error
}

type UserRetriever interface {
	GetByEmail(ctx appengine.Context, email string) (*UserDTO, error)
	GetByAccessId(ctx appengine.Context, accessId string) (*UserDTO, error)
	GetAll(ctx appengine.Context) ([]*datastore.Key, []UserDTO, error)
}

type UserSessionSetter interface {
	SetSessionUUID(ctx appengine.Context, user *UserDTO, uuid string) error
}

type UserSessionRetriever interface {
	GetUserFromSessionUUID(ctx appengine.Context, uuid string) (*UserDTO, error)
}

type UserDAO interface {
	AppEngineHelper.AppEngineDaoBase
	UserCreator
	UserRetriever
	UserSessionRetriever
	UserSessionSetter
}
