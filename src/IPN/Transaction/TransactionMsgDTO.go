package TransactionDao

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/KrauseStefan/NavitasFitness/AppEngineHelper"
	"github.com/KrauseStefan/NavitasFitness/constants"
)

const (
	TXN_KIND             = "txn"
	TXN_PARENT_STRING_ID = "default_txn"
)

var (
	txnCollectionParentKey = datastore.NameKey(TXN_KIND, TXN_PARENT_STRING_ID, nil)
)

func NewTransactionMsgDTOFromIpn(ipnMessage string) *TransactionMsgDTO {
	return NewTransactionMsgDTOFromIpnWithKey(ipnMessage, nil)
}

func NewTransactionMsgDTOFromIpnWithKey(ipnMessage string, key *datastore.Key) *TransactionMsgDTO {
	t := TransactionMsgDTO{
		key:   key,
		dsDto: new(transactionMsgDsDTO),
	}

	t.AddNewIpnMessage(ipnMessage)

	return &t
}

func NewTransactionMsgDTOFromDs(dto *transactionMsgDsDTO, key *datastore.Key) *TransactionMsgDTO {
	t := TransactionMsgDTO{
		key:   key,
		dsDto: dto,
	}
	t.parseMessage()
	return &t
}

func NewTransactionMsgDTOList(dtos []*transactionMsgDsDTO, keys []*datastore.Key) TransactionList {
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

func (t transactionMsgDsDTO) String() string {
	js, _ := json.MarshalIndent(t, "", "  ")
	return string(js)
}

type TransactionMsgDTO struct {
	dsDto *transactionMsgDsDTO
	key   *datastore.Key

	parsedIpnMessage url.Values
}
type TransactionList []*TransactionMsgDTO

func (txns TransactionList) Filter(filterFn func(*TransactionMsgDTO) bool) TransactionList {
	filteredTxns := make([]*TransactionMsgDTO, 0, len(txns))
	for _, txn := range txns {
		if filterFn(txn) {
			filteredTxns = append(filteredTxns, txn)
		}
	}
	return filteredTxns
}

func (txns TransactionList) GetUserKeys() []*datastore.Key {
	userKeys := make([]*datastore.Key, 0, len(txns))
	for _, txn := range txns {
		userKeys = append(userKeys, txn.GetUser())
	}
	return userKeys
}

func (txns TransactionList) getDatastoreKeyAndDtos() ([]*datastore.Key, []*transactionMsgDsDTO) {
	txnKeys := make([]*datastore.Key, len(txns))
	dsTxns := make([]*transactionMsgDsDTO, len(txns))
	for i, txn := range txns {
		txnKeys[i] = txn.GetKey()
		dsTxns[i] = txn.dsDto
	}
	return txnKeys, dsTxns
}

func (t *TransactionMsgDTO) hasKey() bool {
	return t.key != nil
}

func (t *TransactionMsgDTO) GetKey() *datastore.Key {
	return t.key
}

func (t *TransactionMsgDTO) GetUser() *datastore.Key {
	return t.key.Parent
}

func (t *TransactionMsgDTO) GetTxnId() string {
	return t.dsDto.TxnId
}

func (t *TransactionMsgDTO) GetPayerEmail() string {
	return t.GetField(FIELD_PAYER_EMAIL)
}

func (t *TransactionMsgDTO) parseMessage() *url.Values {
	if t.parsedIpnMessage == nil {
		parsedIpnMessage, _ := url.ParseQuery(t.getLatestIPNMessage())
		t.parsedIpnMessage = parsedIpnMessage
		t.dsDto.PaymentDate = t.GetPaymentDate()
		t.dsDto.TxnId = t.GetField(FIELD_TXN_ID)
	}

	return &t.parsedIpnMessage
}

func (t *TransactionMsgDTO) GetField(field string) string {
	return t.parseMessage().Get(field)
}

func (t *TransactionMsgDTO) getLatestIPNMessage() string {
	if len(t.dsDto.IpnMessages) > 0 {
		return t.dsDto.IpnMessages[0]
	} else {
		return ""
	}
}

func (t *TransactionMsgDTO) AddNewIpnMessage(ipnMessage string) *TransactionMsgDTO {
	t.dsDto.IpnMessages = append([]string{ipnMessage}, t.dsDto.IpnMessages...)
	t.parsedIpnMessage = nil
	t.parseMessage()
	return t
}

func (t *TransactionMsgDTO) GetPaymentStatus() string {
	return t.parseMessage().Get(FIELD_PAYMENT_STATUS)
}

func (t *TransactionMsgDTO) GetPaymentDate() time.Time {
	paymentDate := t.dsDto.PaymentDate

	if paymentDate == (time.Time{}) {
		paymentDate = t.getIpnPaymentDate()
	}

	locCET, err := time.LoadLocation("CET")
	if err != nil {
		panic(err)
	}

	return paymentDate.In(locCET)
}

func (t *TransactionMsgDTO) getIpnPaymentDate() time.Time {
	const layout = "15:04:05 Jan 02, 2006 MST"
	fieldValue := t.parseMessage().Get(FIELD_PAYMENT_DATE)
	splitPoint := strings.LastIndex(fieldValue, " ") + 1
	timeZone := fieldValue[splitPoint:]

	loc, err := AppEngineHelper.LoadLocation(timeZone)
	if err != nil {
		panic(err)
	}

	paymentDate, _ := time.ParseInLocation(layout, fieldValue, loc)
	return paymentDate
}

func (t *TransactionMsgDTO) GetAmount() float64 {
	value, _ := strconv.ParseFloat(t.parseMessage().Get(FIELD_MC_GROSS), 64)
	return value
}

func (t *TransactionMsgDTO) GetCurrency() string {
	return t.parseMessage().Get(FIELD_MC_CURRENCY)
}

func (t *TransactionMsgDTO) PaymentIsCompleted() bool {
	return t.GetPaymentStatus() == STATUS_COMPLEATED
}

func (t *TransactionMsgDTO) GetIpnMessages() []string {
	return t.dsDto.IpnMessages
}

func (t *TransactionMsgDTO) GetReceiverEmail() string {
	return t.parseMessage().Get(FIELD_RECEIVER_EMAIL)
}

func (t TransactionMsgDTO) String() string {

	dsDto := fmt.Sprintf("dsDto: %s\n", t.dsDto.String())
	json, _ := json.MarshalIndent(t.parsedIpnMessage, "", "  ")

	return dsDto + string(json)
}

func (t *TransactionMsgDTO) IsActive() bool {
	endTime := t.GetPaymentDate().AddDate(0, constants.SubscriptionDurationInMonth, 0)
	return t.GetPaymentDate().Before(time.Now()) && time.Now().Before(endTime)
}

func (t *TransactionMsgDTO) ExpirationWarningGiven() bool {
	return t.dsDto.ExpirationWarningGiven
}
