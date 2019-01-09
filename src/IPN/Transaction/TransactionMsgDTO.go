package TransactionDao

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"AppEngineHelper"
	"strings"
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
	return NewTransactionMsgDTOFromIpnWithKey(ipnMessage, nil)
}

func NewTransactionMsgDTOFromIpnWithKey(ipnMessage string, key *datastore.Key) *TransactionMsgDTO {
	t := TransactionMsgDTO{
		key:   key,
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
	IpnMessages []string `json:"ipn_messages" datastore:",noindex"` // History of IpnMessages

	PaymentDate time.Time `json:"payment_date"` // Included so that it can be indexed
	TxnId       string    `json:"txn_id"`

	ExpirationWarningGiven bool
}

type TransactionMsgDTO struct {
	dsDto transactionMsgDsDTO
	key   *datastore.Key

	parsedIpnMessage url.Values
}

func (txDto *TransactionMsgDTO) hasKey() bool {
	return txDto.key != nil
}

func (txDto *TransactionMsgDTO) GetDataStoreKey(ctx context.Context) *datastore.Key {
	return txDto.key
}

func (txDto *TransactionMsgDTO) GetUser() *datastore.Key {
	return txDto.key.Parent()
}

func (txDto *TransactionMsgDTO) GetTxnId() string {
	return txDto.dsDto.TxnId
}

func (txDto *TransactionMsgDTO) GetPayerEmail() string {
	return txDto.GetField(FIELD_PAYER_EMAIL)
}

func (txDto *TransactionMsgDTO) parseMessage() *url.Values {
	if txDto.parsedIpnMessage == nil {
		parsedIpnMessage, _ := url.ParseQuery(txDto.getLatestIPNMessage())
		txDto.parsedIpnMessage = parsedIpnMessage
		txDto.dsDto.PaymentDate = txDto.GetPaymentDate()
		txDto.dsDto.TxnId = txDto.GetField(FIELD_TXN_ID)
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
	const layout = "15:04:05 Jan 02, 2006 MST"
	fieldValue := txDto.parseMessage().Get(FIELD_PAYMENT_DATE)
	splitPoint := strings.LastIndex(fieldValue, " ") + 1
	timeZone := fieldValue[splitPoint:]

	loc, err := AppEngineHelper.LoadLocation(timeZone)
	if err != nil {
		panic(err)
	}

	t, _ := time.ParseInLocation(layout, fieldValue, loc)
	locCET, err := time.LoadLocation("CET")
	if err != nil {
		panic(err)
	}

	return t.In(locCET)
}

func (txDto *TransactionMsgDTO) GetAmount() float64 {
	value, _ := strconv.ParseFloat(txDto.parseMessage().Get(FIELD_MC_GROSS), 64)
	return value
}

func (txDto *TransactionMsgDTO) GetCurrency() string {
	return txDto.parseMessage().Get(FIELD_MC_CURRENCY)
}

func (txDto *TransactionMsgDTO) PaymentIsCompleted() bool {
	return txDto.GetPaymentStatus() == STATUS_COMPLEATED
}

func (txDto *TransactionMsgDTO) GetIpnMessages() []string {
	return txDto.dsDto.IpnMessages
}

func (txDto *TransactionMsgDTO) GetReceiverEmail() string {
	return txDto.parseMessage().Get(FIELD_RECEIVER_EMAIL)
}

func (txDto TransactionMsgDTO) String() string {

	dsDto := fmt.Sprintf("dsDto: %s\n", txDto.dsDto.String())
	json, _ := json.MarshalIndent(txDto.parsedIpnMessage, "", "  ")

	return dsDto + string(json)
}

func (txDto transactionMsgDsDTO) String() string {
	js, _ := json.MarshalIndent(txDto, "", "  ")
	return string(js)
}

func (txDto *TransactionMsgDTO) IsActive() bool {
	endTime := txDto.GetPaymentDate().AddDate(0, 6, 0)
	return txDto.GetPaymentDate().Before(time.Now()) && time.Now().Before(endTime)
}
