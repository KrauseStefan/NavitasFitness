package UserDao

import (
	"appengine"
	"appengine/datastore"
)

type IUserDAO interface {
	StringToKey(ctx appengine.Context, key string) *datastore.Key
	GetUserByEmail(ctx appengine.Context, email string) (*UserDTO, error)
	GetUserByAccessId(ctx appengine.Context, accessId string) (*UserDTO, error)
	GetAllUsers(ctx appengine.Context) ([]*datastore.Key, []UserDTO, error)
	CreateUser(ctx appengine.Context, user *UserDTO) error
	SetSessionUUID(ctx appengine.Context, user *UserDTO, uuid string) error
	GetUserFromSessionUUID(ctx appengine.Context, uuid string) (*UserDTO, error)
}
