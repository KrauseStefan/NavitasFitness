package xlsx

import (
	"appengine"
	"appengine/datastore"

	"User/Dao"
)

type UserRetrieverMock struct {
	keys  []*datastore.Key
	users []UserDao.UserDTO
	err   error

	CallCount        int
	LatestCallCtxArg appengine.Context
}

func mockUserRetriever(keys []*datastore.Key, users []UserDao.UserDTO, err error) *UserRetrieverMock {
	mock := &UserRetrieverMock{
		keys:  keys,
		users: users,
		err:   err,

		CallCount:        0,
		LatestCallCtxArg: nil,
	}
	userDAO = mock
	return mock
}

func (mock *UserRetrieverMock) GetAll(ctx appengine.Context) ([]*datastore.Key, []UserDao.UserDTO, error) {
	mock.CallCount++
	mock.LatestCallCtxArg = ctx
	return mock.keys, mock.users, mock.err
}
