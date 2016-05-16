package TransactionDao

import (
	"time"
	"net/url"
	"strconv"
	"appengine"
	"appengine/datastore"
	"src/Common"
)

const (
	TXN_KIND = "txn"
	TXN_PARENT_STRING_ID = "default_txn"
)

var (
	txnCollectionParentKey = Common.CollectionParentKeyGetFnGenerator(TXN_KIND, TXN_PARENT_STRING_ID, 0)
	txnIntIDToKeyInt64 = Common.IntIDToKeyInt64(TXN_KIND, txnCollectionParentKey)
)

type TransactionMsgDTO struct {
	Key              string
	IpnMessages      []string //History of IpnMessages
	StatusResp       string

	PaymentDate      time.Time
	TxnId            string
	parsedIpnMessage url.Values `datastore:"-"`
}

func (txDto *TransactionMsgDTO) hasKey() bool {
	return len(txDto.Key) > 0
}

func (txDto *TransactionMsgDTO) GetDataStoreKey(ctx appengine.Context) *datastore.Key {
	return StringToKey(ctx, txDto.Key)
}

func(txDto *TransactionMsgDTO) setKey(key *datastore.Key) *TransactionMsgDTO {
	txDto.Key = strconv.FormatInt(key.IntID(), 10)
	return txDto
}

func StringToKey(ctx appengine.Context, key string) *datastore.Key {
	return txnIntIDToKeyInt64(ctx, key)
}

func (txDto *TransactionMsgDTO) parseMessage() *url.Values {
	if txDto.parsedIpnMessage == nil {
		parsedIpnMessage, _ := url.ParseQuery(txDto.GetLatestIPNMessage())
		txDto.parsedIpnMessage = parsedIpnMessage
		txDto.PaymentDate = txDto.GetPaymentDate()
		txDto.TxnId = txDto.GetField(FIELD_TXN_ID)
	}

	return &txDto.parsedIpnMessage
}

func (txDto *TransactionMsgDTO) GetField(field string) string {
	return txDto.parseMessage().Get(field)
}

func (txDto *TransactionMsgDTO) GetLatestIPNMessage() string {
	if len(txDto.IpnMessages) > 0 {
		return txDto.IpnMessages[0]
	} else {
		return ""
	}
}

func (txDto *TransactionMsgDTO) AddNewIpnMessage(ipnMessage string) *TransactionMsgDTO {
	txDto.IpnMessages = append([]string{ipnMessage}, txDto.IpnMessages...)
	txDto.parsedIpnMessage = nil
	txDto.parseMessage()
	return txDto
}

func (txDto *TransactionMsgDTO) GetPaymentStatus() string {
	return txDto.parseMessage().Get(FIELD_PAYMENT_STATUS)
}

func (txDto *TransactionMsgDTO) GetPaymentDate() time.Time {
	value := txDto.parseMessage().Get(FIELD_PAYMENT_DATE)
	const layout = "15:04:05 Jan 02, 2006 MST" //Reference time Mon Jan 2 15:04:05 -0700 MST 2006
	t, _ := time.Parse(layout, value)
	return t
}

func (txDto *TransactionMsgDTO) GetAmount() float64 {
	value, _ := strconv.ParseFloat(txDto.parseMessage().Get(FIELD_MC_GROSS), 64)
	return value
}

func (txDto *TransactionMsgDTO) GetCurrency() string {
	return txDto.parseMessage().Get(FIELD_MC_CURRENCY)
}

func (txDto *TransactionMsgDTO) IsVerified() bool {
	return txDto.StatusResp == "VERIFIED" // other option is "INVALID"
}

func (txDto *TransactionMsgDTO) PaymentIsComplected() bool {
	return txDto.IsVerified() && txDto.GetPaymentStatus() == STATUS_COMPLEATED && txDto.GetAmount() == 200
}