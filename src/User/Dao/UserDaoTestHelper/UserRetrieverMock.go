package UserDaoTestHelper

import (
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"

	UserDao "github.com/KrauseStefan/NavitasFitness/User/Dao"
)

type CallArgs struct {
	Ctx  context.Context
	Keys []*datastore.Key
}

type ReturnValues struct {
	keys     []*datastore.Key
	userDtos UserDao.UserList
	err      error
}

type UserRetrieverMock struct {
	CallCount int

	returnValues []ReturnValues
	CallArgs     []CallArgs
}

func NewUserRetrieverMock(keys []*datastore.Key, users UserDao.UserList, err error) *UserRetrieverMock {
	mock := &UserRetrieverMock{}
	mock.AddReturn(keys, users, err)

	return mock
}

func (mock *UserRetrieverMock) AddReturn(keys []*datastore.Key, userDtos UserDao.UserList, err error) *UserRetrieverMock {
	mock.returnValues = append(mock.returnValues, ReturnValues{
		keys:     keys,
		userDtos: userDtos,
		err:      err,
	})
	return mock
}

func (mock *UserRetrieverMock) GetAll(ctx context.Context) ([]*datastore.Key, UserDao.UserList, error) {
	rtnValues := mock.returnValues[mock.CallCount]
	mock.CallCount++
	mock.CallArgs = append(mock.CallArgs, CallArgs{
		Ctx: ctx,
	})
	return rtnValues.keys, rtnValues.userDtos, rtnValues.err
}

func (mock *UserRetrieverMock) GetByKey(ctx context.Context, key *datastore.Key) (*UserDao.UserDTO, error) {
	panic("not implemented")
}

func (mock *UserRetrieverMock) GetByKeys(ctx context.Context, keys []*datastore.Key) (UserDao.UserList, error) {
	rtnValues := mock.returnValues[mock.CallCount]
	mock.CallCount++
	mock.CallArgs = append(mock.CallArgs, CallArgs{
		Ctx:  ctx,
		Keys: keys,
	})
	return rtnValues.userDtos, rtnValues.err
}
