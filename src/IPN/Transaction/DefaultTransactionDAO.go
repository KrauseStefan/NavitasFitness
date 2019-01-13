package TransactionDao

import (
	"errors"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
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

	key := ipnTxn.GetDataStoreKey(ctx)

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
	if _, err := datastore.Put(ctx, newKey, &ipnTxn.dsDto); err != nil {
		return err
	}

	return nil
}

func (t *DefaultTransactionDao) GetTransaction(ctx context.Context, txnId string) (*TransactionMsgDTO, error) {
	q := datastore.NewQuery(TXN_KIND).
		Filter("TxnId=", txnId).
		Limit(1)

	txnDtoList := make([]transactionMsgDsDTO, 0, 1)

	keys, err := q.GetAll(ctx, &txnDtoList)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, nil
	}

	return NewTransactionMsgDTOFromDs(txnDtoList[0], keys[0]), nil
}

func (t *DefaultTransactionDao) GetTransactionsByUser(ctx context.Context, parentUserKey *datastore.Key) ([]*TransactionMsgDTO, error) {

	q := datastore.NewQuery(TXN_KIND).
		Ancestor(parentUserKey).
		Order("PaymentDate")

	count, err := q.Count(ctx)
	if err != nil {
		return nil, err
	}

	txnDsDtoList := make([]transactionMsgDsDTO, 0, count)

	keys, err := q.GetAll(ctx, &txnDsDtoList)
	if err != nil {
		return nil, err
	}

	return NewTransactionMsgDTOList(txnDsDtoList, keys), nil
}

func (t *DefaultTransactionDao) GetCurrentTransactionsAfter(ctx context.Context, date time.Time) ([]*TransactionMsgDTO, error) {
	q := datastore.NewQuery(TXN_KIND).
		Filter("PaymentDate>=", date)

	txnDsDtoList := make([]transactionMsgDsDTO, 0)

	keys, err := q.GetAll(ctx, &txnDsDtoList)
	if err != nil {
		return nil, err
	}

	return NewTransactionMsgDTOList(txnDsDtoList, keys), nil
}

func GetTransactionsAboutToExpire(ctx context.Context) ([]*datastore.Key, error) {
	subscriptionDurationInMonth := 6
	warningDeltaDays := 7

	paymentExpiratinDate := time.Now().AddDate(0, -subscriptionDurationInMonth, 0)
	paymentWarningStartDate := paymentExpiratinDate.AddDate(0, 0, -warningDeltaDays)

	log.Infof(ctx, "PaymentDate>=%s", paymentWarningStartDate)
	log.Infof(ctx, "PaymentDate<=%s", paymentExpiratinDate)

	q := datastore.NewQuery(TXN_KIND).
		Filter("PaymentDate>=", paymentWarningStartDate).
		Filter("PaymentDate<=", paymentExpiratinDate).
		KeysOnly()

	return q.GetAll(ctx, nil)
}
