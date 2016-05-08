package TransActionDao

import (
	"time"
	"net/url"
	"strconv"
)

type TransactionMsgDTO struct {
	IpnMessages      []string //History of IpnMessages
	StatusResp       string

	PaymentDate      time.Time
	TxnId            string
	parsedIpnMessage url.Values `datastore:"-"`
}

func (txDto *TransactionMsgDTO) parseMessage() *url.Values {
	if txDto.parsedIpnMessage == nil {
		parsedIpnMessage, _ := url.ParseQuery(txDto.GetLatestIPNMessage())
		txDto.parsedIpnMessage = parsedIpnMessage
		txDto.PaymentDate = txDto.GetPaymentDate()
		txDto.TxnId = txDto.GetTxnId()
	}

	return &txDto.parsedIpnMessage
}

func (txDto *TransactionMsgDTO) GetTxnId() string {
	return txDto.parseMessage().Get(FIELD_TXN_ID)
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