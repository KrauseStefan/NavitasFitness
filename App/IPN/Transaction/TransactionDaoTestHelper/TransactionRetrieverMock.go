package TransactionDaoTestHelper

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"IPN/Transaction"
)

type CallArgs struct {
	Ctx   context.Context
	Key   *datastore.Key
	TxnId string
	Date  time.Time
}

type ReturnValues struct {
	messages []*TransactionDao.TransactionMsgDTO
	err      error
}

type TransactionRetrieverMock struct {
	returnValues []ReturnValues

	CallCount int
	CallArgs  []CallArgs
}

func NewTransactionRetrieverMock(messages []*TransactionDao.TransactionMsgDTO, err error) *TransactionRetrieverMock {
	mock := &TransactionRetrieverMock{}
	mock.AddReturn(messages, err)
	return mock
}

func (mock *TransactionRetrieverMock) AddReturn(messages []*TransactionDao.TransactionMsgDTO, err error) *TransactionRetrieverMock {
	mock.returnValues = append(mock.returnValues, ReturnValues{
		messages: messages,
		err:      err,
	})
	return mock
}

func (mock *TransactionRetrieverMock) GetTransaction(ctx context.Context, txnId string) (*TransactionDao.TransactionMsgDTO, error) {
	rtnValues := mock.returnValues[mock.CallCount]
	mock.CallCount++
	mock.CallArgs = append(mock.CallArgs, CallArgs{
		Ctx:   ctx,
		TxnId: txnId,
	})
	return rtnValues.messages[0], rtnValues.err
}

func (mock *TransactionRetrieverMock) GetTransactionsByUser(ctx context.Context, parentUserKey *datastore.Key) ([]*TransactionDao.TransactionMsgDTO, error) {
	rtnValues := mock.returnValues[mock.CallCount]
	mock.CallCount++
	mock.CallArgs = append(mock.CallArgs, CallArgs{
		Ctx: ctx,
		Key: parentUserKey,
	})
	return rtnValues.messages, rtnValues.err
}

func (mock *TransactionRetrieverMock) GetCurrentTransactionsAfter(ctx context.Context, userKey *datastore.Key, date time.Time) ([]*TransactionDao.TransactionMsgDTO, error) {
	rtnValues := mock.returnValues[mock.CallCount]
	mock.CallCount++
	mock.CallArgs = append(mock.CallArgs, CallArgs{
		Ctx:  ctx,
		Key:  userKey,
		Date: date,
	})
	return rtnValues.messages, rtnValues.err
}
