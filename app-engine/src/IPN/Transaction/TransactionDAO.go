package TransactionDao

import (
	"appengine/datastore"
	"appengine"
	"errors"
	"src/User/Dao"
)

var (
	TxnDuplicateTxnMsg = errors.New("Doublicate message recived, this is likely not a programming error")
	//txnUnableToVerify = errors.New("Unable to verify message")
)

func PersistIpnMessage(ctx appengine.Context, ipnTxn *TransactionMsgDTO, userKey string) error {

	var newKey *datastore.Key

	if !ipnTxn.hasKey() {
		if userKey == "" {
			newKey = datastore.NewIncompleteKey(ctx, TXN_KIND, txnCollectionParentKey(ctx))
			//todo log an error here, this is not a normal scenario
		} else {
			newKey = datastore.NewIncompleteKey(ctx, TXN_KIND, UserDao.StringToKey(ctx, userKey))
		}
	} else {
		newKey = ipnTxn.GetDataStoreKey(ctx)
	}

	ipnTxn.PaymentDate = ipnTxn.GetPaymentDate()
	if _, err := datastore.Put(ctx, newKey, ipnTxn); err != nil {
		return err
	}

	return nil
}

func GetTransaction(ctx appengine.Context, txnId string) (*TransactionMsgDTO, error) {
	q := datastore.NewQuery(TXN_KIND).
	Filter("TxnId=", txnId).
	Limit(1)

	txnDtoList := make([]TransactionMsgDTO, 0, 1)

	keys, err := q.GetAll(ctx, &txnDtoList)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, nil
	}

	txnDtoList[0].setKey(keys[0])
	return &txnDtoList[0], nil
}

func GetTransactionsByUser(ctx appengine.Context, parentUserKey *datastore.Key) ([]TransactionMsgDTO, error) {

	q := datastore.NewQuery(TXN_KIND).
		Ancestor(parentUserKey).
		Order("PaymentDate")

	entryCount, err := q.Count(ctx);
	if err != nil {
		return nil, err
	}

	txnDtoList := make([]TransactionMsgDTO, 0, entryCount)

	keys, err := q.GetAll(ctx, &txnDtoList)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		txnDtoList[i].setKey(key)
	}

	return txnDtoList, nil
}