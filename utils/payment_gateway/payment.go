package paymentgateway

import (
	"time"

	"github.com/midtrans/midtrans-go/coreapi"
)

// bri, bni, bca, mandiri, permata, alfamart, gopay.
// qris (The bank/payment partner is experiencing some connection issues)

type Bank string
type Ewallet string
type ConvStore string

const (
	bca       Bank      = "bca"
	bni       Bank      = "bni"
	bri       Bank      = "bri"
	mandiri   Bank      = "mandiri"
	permata   Bank      = "permata"
	gopay     Ewallet   = "gopay"
	shopeepay Ewallet   = "shopeepay"
	qris      Ewallet   = "qris"
	indomaret ConvStore = "indomaret"
	alfamart  ConvStore = "alfamart"
)

type ChargeResponse struct {
	TransactionID          string             `json:"transaction_id"`
	OrderID                string             `json:"order_id"`
	GrossAmount            string             `json:"gross_amount"`
	PaymentType            string             `json:"payment_type"`
	TransactionTime        string             `json:"transaction_time"`
	TransactionStatus      string             `json:"transaction_status"`
	FraudStatus            string             `json:"fraud_status"`
	StatusCode             string             `json:"status_code"`
	Bank                   string             `json:"bank"`
	Store                  string             `json:"store"`
	StatusMessage          string             `json:"status_message"`
	ChannelResponseCode    string             `json:"channel_response_code"`
	ChannelResponseMessage string             `json:"channel_response_message"`
	Currency               string             `json:"currency"`
	ValidationMessages     []string           `json:"validation_messages"`
	PermataVaNumber        string             `json:"permata_va_number"`
	VaNumbers              []coreapi.VANumber `json:"va_numbers"`
	BillKey                string             `json:"bill_key"`
	BillerCode             string             `json:"biller_code"`
	Actions                []coreapi.Action   `json:"actions"`
	PaymentCode            string             `json:"payment_code"`
	QRString               string             `json:"qr_string"`
	Expire                 string             `json:"expiry_time"`
	CreatedAt              time.Time          `json:"created_at"`
	UpdatedAt              time.Time          `json:"updated_at"`
	DeletedAt              time.Time          `json:"deleted_at"`
}

func GetPaymentCode(res *ChargeResponse) string {
	if len(res.VaNumbers) > 0 {
		return res.VaNumbers[0].VANumber
	}

	if res.BillKey != "" {
		return res.BillKey
	}

	if res.PaymentCode != "" {
		return res.PaymentCode
	}

	return ""
}

func GetBankType(res *ChargeResponse) string {
	if len(res.VaNumbers) > 0 {
		return string(res.VaNumbers[0].Bank)
	}

	if res.BillKey != "" && res.BillerCode != "" {
		return string(mandiri)
	}

	if res.PermataVaNumber != "" {
		return string(permata)
	}

	if res.PaymentType == "cstore" && res.Store == "indomaret" {
		return string(indomaret)
	}

	if res.PaymentType == "cstore" && res.Store == "alfamart" {
		return string(alfamart)
	}

	if res.PaymentType == "gopay" {
		return string(gopay)
	}

	return ""
}
