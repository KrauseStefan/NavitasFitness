package TransActionDao

import (
	"src/Common"
	"appengine/datastore"
	"appengine"
	"errors"
	"src/User/Dao"
)

const (
	TXN_KIND = "txn"
	TXN_PARENT_STRING_ID = "default_txn"
)

var (
	txnNotFoundError = errors.New("Transaction does not exist in datastore")
	txnDuplicateTxnMsg = errors.New("Doublicate message recived, this is likely not a programming error")
	txnUnableToVerify = errors.New("Unable to verify message")
)

var (
	txnCollectionParentKey = Common.CollectionParentKeyGetFnGenerator(TXN_KIND, TXN_PARENT_STRING_ID, 0)
	txnIntIDToKeyInt64 = Common.IntIDToKeyInt64(TXN_KIND, txnCollectionParentKey)
)

func PersistIpnMessage(ctx appengine.Context, ipnTxn *TransactionMsgDTO, userKey string) error {

	key, dbTxn, err := GetTransaction(ctx, ipnTxn.GetTxnId())

	if err != nil && err != txnNotFoundError {
		return err
	}

	if dbTxn != nil && ipnTxn.GetPaymentStatus() == dbTxn.GetPaymentStatus() {
		//Verify that the IPN is not a duplicate. To do this, save the transaction ID and last payment status in each IPN message in a database and verify that the current IPN's values for these fields are not already in this database.
		//Duplicate txnMsg
		//Persist anyway?, with status duplicate?
		return txnDuplicateTxnMsg
	}

	if dbTxn != nil {
		//Add new ipn message to current DTO
		ipnTxn = dbTxn.AddNewIpnMessage(ipnTxn.GetLatestIPNMessage())
	}

	if key == nil {
		if userKey == "" {
			key = datastore.NewIncompleteKey(ctx, TXN_KIND, txnCollectionParentKey(ctx))
		} else {
			key = datastore.NewIncompleteKey(ctx, TXN_KIND, UserDao.StringToKey(ctx, userKey))
		}
	}

	ipnTxn.PaymentDate = ipnTxn.GetPaymentDate()
	if _, err := datastore.Put(ctx, key, ipnTxn); err != nil {
		return err
	}

	return nil
}

func GetTransaction(ctx appengine.Context, txnId string) (*datastore.Key, *TransactionMsgDTO, error) {
	q := datastore.NewQuery(TXN_KIND).
	Filter("TxnId=", txnId).
	Limit(1)

	txnDtoList := make([]TransactionMsgDTO, 0, 1)

	keys, err := q.GetAll(ctx, &txnDtoList)
	if err != nil {
		return nil, nil, err
	}

	if len(keys) == 0 {
		return nil, nil, txnNotFoundError
	}

	return keys[0], &txnDtoList[0], nil
}

func GetTransactionsByUser(ctx appengine.Context, parentUserKey *datastore.Key) ([]TransactionMsgDTO, error) {

	var ( err error; entryCount int )

	q := datastore.NewQuery(TXN_KIND).
		Ancestor(parentUserKey).
		Order("PaymentDate")

	entryCount, err = q.Count(ctx);
	if err != nil {
		return nil, err
	}

	txnDtoList := make([]TransactionMsgDTO, 0, entryCount)

	_, err = q.GetAll(ctx, &txnDtoList)
	if err != nil {
		return nil, err
	}

	return txnDtoList, nil
}