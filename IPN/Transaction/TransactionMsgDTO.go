package TransactionDao

import (
	"net/url"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"

	"AppEngineHelper"
)

const (
	TXN_KIND             = "txn"
	TXN_PARENT_STRING_ID = "default_txn"
)

var (
	txnCollectionParentKey = AppEngineHelper.CollectionParentKeyGetFnGenerator(TXN_KIND, TXN_PARENT_STRING_ID, 0)
	txnIntIDToKeyInt64     = AppEngineHelper.IntIDToKeyInt64(TXN_KIND, txnCollectionParentKey)
)

func NewTransactionMsgDTOFromIpn(ipnMessage string) *TransactionMsgDTO {
	t := TransactionMsgDTO{
		key:   nil,
		dsDto: *new(transactionMsgDsDTO),
	}

	t.AddNewIpnMessage(ipnMessage)

	return &t
}

func NewTransactionMsgDTOFromDs(dto transactionMsgDsDTO, key *datastore.Key) *TransactionMsgDTO {
	t := TransactionMsgDTO{
		key:   key,
		dsDto: dto,
	}
	t.parseMessage()
	return &t
}

func NewTransactionMsgDTOList(dtos []transactionMsgDsDTO, keys []*datastore.Key) []*TransactionMsgDTO {
	txnDtoList := make([]*TransactionMsgDTO, len(dtos))

	for i, dto := range dtos {
		var key *datastore.Key = nil
		if len(keys) > i && keys[i] != nil {
			key = keys[i]
		}
		txnDtoList[i] = NewTransactionMsgDTOFromDs(dto, key)
	}

	return txnDtoList
}

type transactionMsgDsDTO struct {
	IpnMessages []string //History of IpnMessages

	PaymentActivationDate      time.Time // not used ?
	PaymentDate                time.Time
	SubscriptionActivationDate time.Time // Decides if the subscription is active
	TxnId                      string
}

type TransactionMsgDTO struct {
	dsDto transactionMsgDsDTO
	key   *datastore.Key

	parsedIpnMessage url.Values
}

func (txDto *TransactionMsgDTO) hasKey() bool {
	return txDto.key != nil
}

func (txDto *TransactionMsgDTO) GetDataStoreKey(ctx appengine.Context) *datastore.Key {
	return txDto.key
}

func (txDto *TransactionMsgDTO) parseMessage() *url.Values {
	if txDto.parsedIpnMessage == nil {
		parsedIpnMessage, _ := url.ParseQuery(txDto.getLatestIPNMessage())
		txDto.parsedIpnMessage = parsedIpnMessage
		txDto.dsDto.PaymentDate = txDto.GetPaymentDate()
		txDto.dsDto.TxnId = txDto.GetField(FIELD_TXN_ID)
		if txDto.PaymentIsComplected() && txDto.dsDto.PaymentActivationDate.IsZero() {
			txDto.dsDto.PaymentActivationDate = time.Now()
		}
	}

	return &txDto.parsedIpnMessage
}

func (txDto *TransactionMsgDTO) GetField(field string) string {
	return txDto.parseMessage().Get(field)
}

func (txDto *TransactionMsgDTO) getLatestIPNMessage() string {
	if len(txDto.dsDto.IpnMessages) > 0 {
		return txDto.dsDto.IpnMessages[0]
	} else {
		return ""
	}
}

func (txDto *TransactionMsgDTO) AddNewIpnMessage(ipnMessage string) *TransactionMsgDTO {
	txDto.dsDto.IpnMessages = append([]string{ipnMessage}, txDto.dsDto.IpnMessages...)
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

func (txDto *TransactionMsgDTO) PaymentIsComplected() bool {
	return txDto.GetPaymentStatus() == STATUS_COMPLEATED && txDto.GetAmount() == 200
}

func (txDto *TransactionMsgDTO) GetPaymentActivationDate() time.Time {
	return txDto.dsDto.PaymentActivationDate
}

func (txDto *TransactionMsgDTO) IsActive() bool {
	endTime := txDto.GetPaymentActivationDate().AddDate(0, 6, 0)
	return txDto.GetPaymentActivationDate().Before(time.Now()) && time.Now().Before(endTime)
}