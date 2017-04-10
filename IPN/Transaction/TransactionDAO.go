package TransactionDao

import (
	"time"

	"appengine"
	"appengine/datastore"
)

type TransactionIpnManipulator interface {
	UpdateIpnMessage(ctx appengine.Context, ipnTxn *TransactionMsgDTO) error

	PersistNewIpnMessage(ctx appengine.Context, ipnTxn *TransactionMsgDTO, userKey *datastore.Key) error
}

type TransactionRetriever interface {
	GetTransaction(ctx appengine.Context, txnId string) (*TransactionMsgDTO, error)

	GetTransactionsByUser(ctx appengine.Context, parentUserKey *datastore.Key) ([]*TransactionMsgDTO, error)

	GetCurrentTransactionsAfter(ctx appengine.Context, userKey *datastore.Key, date time.Time) ([]*TransactionMsgDTO, error)
}

type TransactionDao interface {
	TransactionIpnManipulator
	TransactionRetriever
}
