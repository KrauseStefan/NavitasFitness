package TransactionDao

import (
	"errors"
	"time"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"

	nf_datastore "NavitasFitness/datastore"
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
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return err
	}

	key := ipnTxn.GetKey()

	// Make sure indexed fields are updated
	ipnTxn.parseMessage()

	if _, err := dsClient.Put(ctx, key, ipnTxn.dsDto); err != nil {
		return err
	}

	return nil
}

func (t *DefaultTransactionDao) PersistNewIpnMessage(ctx context.Context, ipnTxn *TransactionMsgDTO, userKey *datastore.Key) error {
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return err
	}

	var newKey *datastore.Key

	if ipnTxn.hasKey() {
		return errors.New("ipnTxn has already been persisted, use update function ínstead")
	}

	if userKey == nil {
		newKey = datastore.IncompleteKey(TXN_KIND, txnCollectionParentKey)
	} else {
		newKey = datastore.IncompleteKey(TXN_KIND, userKey)
	}

	//Make sure indexed fields are updated
	ipnTxn.parseMessage()
	ipnTxn.key = newKey
	if _, err := dsClient.Put(ctx, newKey, ipnTxn.dsDto); err != nil {
		return err
	}

	return nil
}

func (t *DefaultTransactionDao) GetTransaction(ctx context.Context, txnId string) (*TransactionMsgDTO, error) {
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery(TXN_KIND).
		Filter("TxnId=", txnId).
		Limit(1)

	var txnDtoList []*transactionMsgDsDTO
	keys, err := dsClient.GetAll(ctx, query, &txnDtoList)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, nil
	}

	return NewTransactionMsgDTOFromDs(txnDtoList[0], keys[0]), nil
}

func (t *DefaultTransactionDao) GetTransactionsByUser(ctx context.Context, parentUserKey *datastore.Key) (TransactionList, error) {
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery(TXN_KIND).
		Ancestor(parentUserKey).
		Order("PaymentDate")

	var txnDsDtoList []*transactionMsgDsDTO
	keys, err := dsClient.GetAll(ctx, query, &txnDsDtoList)
	if err != nil {
		return nil, err
	}

	return NewTransactionMsgDTOList(txnDsDtoList, keys), nil
}

func (t *DefaultTransactionDao) GetCurrentTransactionsAfter(ctx context.Context, date time.Time) (TransactionList, error) {
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery(TXN_KIND).
		Filter("PaymentDate>=", date)

	var txnDsDtoList []*transactionMsgDsDTO
	keys, err := dsClient.GetAll(ctx, query, &txnDsDtoList)
	if err != nil {
		return nil, err
	}

	return NewTransactionMsgDTOList(txnDsDtoList, keys), nil
}

func GetTransactionsPayedBetween(ctx context.Context, start time.Time, end time.Time) (TransactionList, error) {
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery(TXN_KIND).
		Filter("PaymentDate>=", start).
		Filter("PaymentDate<=", end)

	var txns []*transactionMsgDsDTO
	keys, err := dsClient.GetAll(ctx, query, &txns)
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
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, err
	}

	txnKeys, dsTxns := txns.getDatastoreKeyAndDtos()
	return dsClient.PutMulti(ctx, txnKeys, dsTxns)
}
