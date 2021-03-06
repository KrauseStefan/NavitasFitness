package TransactionDao

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type TransactionIpnManipulator interface {
	UpdateIpnMessage(ctx context.Context, ipnTxn *TransactionMsgDTO) error

	PersistNewIpnMessage(ctx context.Context, ipnTxn *TransactionMsgDTO, userKey *datastore.Key) error
}

type TransactionRetriever interface {
	GetTransaction(ctx context.Context, txnId string) (*TransactionMsgDTO, error)

	GetTransactionsByUser(ctx context.Context, parentUserKey *datastore.Key) (TransactionList, error)

	GetCurrentTransactionsAfter(ctx context.Context, date time.Time) (TransactionList, error)
}

type TransactionDao interface {
	TransactionIpnManipulator
	TransactionRetriever
}
