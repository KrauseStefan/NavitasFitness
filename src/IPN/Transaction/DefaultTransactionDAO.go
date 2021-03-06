package TransactionDao

import (
	"errors"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type DefaultTransactionDao struct{}

var defaultTransactionDao = DefaultTransactionDao{}

func GetInstance() TransactionDao {
	return &defaultTransactionDao
}

var (
	TxnDuplicateTxnMsg = errors.New("Doublicate message recived, this is likely not a programming error")
)

func (t *DefaultTransactionDao) UpdateIpnMessage(ctx context.Context, ipnTxn *TransactionMsgDTO) error {

	key := ipnTxn.GetKey()

	// Make sure indexed fields are updated
	ipnTxn.parseMessage()

	if _, err := datastore.Put(ctx, key, ipnTxn.dsDto); err != nil {
		return err
	}

	return nil
}

func (t *DefaultTransactionDao) PersistNewIpnMessage(ctx context.Context, ipnTxn *TransactionMsgDTO, userKey *datastore.Key) error {

	var newKey *datastore.Key

	if ipnTxn.hasKey() {
		return errors.New("ipnTxn has already been persisted, use update function ínstead")
	}

	if userKey == nil {
		newKey = datastore.NewIncompleteKey(ctx, TXN_KIND, txnCollectionParentKey(ctx))
	} else {
		newKey = datastore.NewIncompleteKey(ctx, TXN_KIND, userKey)
	}

	//Make sure indexed fields are updated
	ipnTxn.parseMessage()
	ipnTxn.key = newKey
	if _, err := datastore.Put(ctx, newKey, ipnTxn.dsDto); err != nil {
		return err
	}

	return nil
}

func (t *DefaultTransactionDao) GetTransaction(ctx context.Context, txnId string) (*TransactionMsgDTO, error) {
	q := datastore.NewQuery(TXN_KIND).
		Filter("TxnId=", txnId).
		Limit(1)

	var txnDtoList []*transactionMsgDsDTO
	keys, err := q.GetAll(ctx, &txnDtoList)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, nil
	}

	return NewTransactionMsgDTOFromDs(txnDtoList[0], keys[0]), nil
}

func (t *DefaultTransactionDao) GetTransactionsByUser(ctx context.Context, parentUserKey *datastore.Key) (TransactionList, error) {

	q := datastore.NewQuery(TXN_KIND).
		Ancestor(parentUserKey).
		Order("PaymentDate")

	var txnDsDtoList []*transactionMsgDsDTO
	keys, err := q.GetAll(ctx, &txnDsDtoList)
	if err != nil {
		return nil, err
	}

	return NewTransactionMsgDTOList(txnDsDtoList, keys), nil
}

func (t *DefaultTransactionDao) GetCurrentTransactionsAfter(ctx context.Context, date time.Time) (TransactionList, error) {
	q := datastore.NewQuery(TXN_KIND).
		Filter("PaymentDate>=", date)

	var txnDsDtoList []*transactionMsgDsDTO
	keys, err := q.GetAll(ctx, &txnDsDtoList)
	if err != nil {
		return nil, err
	}

	return NewTransactionMsgDTOList(txnDsDtoList, keys), nil
}

func GetTransactionsPayedBetween(ctx context.Context, start time.Time, end time.Time) (TransactionList, error) {
	q := datastore.NewQuery(TXN_KIND).
		Filter("PaymentDate>=", start).
		Filter("PaymentDate<=", end)

	var txns []*transactionMsgDsDTO
	keys, err := q.GetAll(ctx, &txns)
	if err != nil {
		return nil, err
	}

	return NewTransactionMsgDTOList(txns, keys), err
}

func SetExpirationWarningGiven(ctx context.Context, txns TransactionList, value bool) error {
	for _, txn := range txns {
		txn.dsDto.ExpirationWarningGiven = value
	}

	_, err := putMulti(ctx, txns)
	return err
}

func putMulti(ctx context.Context, txns TransactionList) ([]*datastore.Key, error) {
	txnKeys, dsTxns := txns.getDatastoreKeyAndDtos()
	return datastore.PutMulti(ctx, txnKeys, dsTxns)
}
