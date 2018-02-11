package UserDaoTestHelper

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	".."
)

type UserRetrieverMock struct {
	keys  []*datastore.Key
	users []UserDao.UserDTO
	err   error

	CallCount        int
	LatestCallCtxArg context.Context
}

func NewUserRetrieverMock(keys []*datastore.Key, users []UserDao.UserDTO, err error) *UserRetrieverMock {
	mock := &UserRetrieverMock{
		keys:  keys,
		users: users,
		err:   err,
	}
	return mock
}

func (mock *UserRetrieverMock) GetAll(ctx context.Context) ([]*datastore.Key, []UserDao.UserDTO, error) {
	mock.CallCount++
	mock.LatestCallCtxArg = ctx
	return mock.keys, mock.users, mock.err
}

func (mock *UserRetrieverMock) GetByKey(ctx context.Context, key *datastore.Key) (*UserDao.UserDTO, error) {
	panic("not implemented")
}
