package paymentgateway

import (
	"time"

	"github.com/midtrans/midtrans-go/coreapi"
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
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              time.Time
}

type Bank string

const (
	bca     Bank = "bca"
	bni     Bank = "bni"
	bri     Bank = "bri"
	mandiri Bank = "mandiri"
)

type ConvStore string

const (
	indomaret ConvStore = "indomart"
	alfamart  ConvStore = "alfamart"
)

type Ewallet string

const (
	gopay     Ewallet = "gopay"
	shopeepay Ewallet = "shopeepay"
	qris      Ewallet = "qris"
)

// // Charge Midtrans with Bank Transfer
// var bank = map[string]coreapi.BankTransferDetails{
// 	"bni": {
// 		Bank: midtrans.BankBni,
// 	},
// 	"bca": {
// 		Bank: midtrans.BankBca,
// 	},
// 	"bri": {
// 		Bank: midtrans.BankBri,
// 	},
// 	"mandiri": {
// 		Bank: midtrans.BankMandiri,
// 	},
// }

// func ChargeMidtrans(reservationId string, request reservation.PaymentCore) (*ChargeResponse, error) {
// 	var c = coreapi.Client{}
// 	c.New(config.MIDTRANS_SERVERKEY, midtrans.Sandbox)
// 	bankTransfer, ok := bank[request.PaymentType]
// 	if !ok {
// 		return nil, errors.New("bank error")
// 	}

// 	if reservationId == "" {
// 		log.Error("error reservation ID")
// 		return nil, errors.New("invalid reservation ID")
// 	}
// 	log.Sugar().Infoln("reservationnya : ", request.Reservation.ReservationID)

// 	req := &coreapi.ChargeReq{
// 		PaymentType:  coreapi.PaymentTypeBankTransfer,
// 		BankTransfer: &bankTransfer,
// 		TransactionDetails: midtrans.TransactionDetails{
// 			OrderID:  reservationId, // Gunakan reservationID dari parameter
// 			GrossAmt: int64(request.GrandTotal),
// 		},
// 	}

// 	resp, _ := c.ChargeTransaction(req)
// 	banks := make([]string, len(resp.VaNumbers))
// 	for i, bank := range resp.VaNumbers {
// 		banks[i] = bank.Bank
// 	}
// 	banksStr := strings.Join(banks, ",")

// 	vaNumbers := make([]string, len(resp.VaNumbers))
// 	for i, vaNumber := range resp.VaNumbers {
// 		vaNumbers[i] = vaNumber.VANumber
// 	}
// 	vaNumbersStr := strings.Join(vaNumbers, ",")
// 	grandTotal, err := strconv.ParseFloat(resp.GrossAmount, 64)
// 	if err != nil {
// 		log.Error("error GrossAmount")
// 		return nil, err
// 	}

// 	log.Sugar().Infoln(resp)
// 	return &ChargeResponse{
// 		TransactionID:     resp.TransactionID,
// 		VaNumbers:         vaNumbersStr,
// 		PaymentType:       resp.PaymentType,
// 		Bank:              banksStr,
// 		GrossAmount:       grandTotal,
// 		TransactionStatus: resp.TransactionStatus,
// 		CreatedAt:         time.Time{},
// 		UpdatedAt:         time.Time{},
// 		DeletedAt:         time.Time{},
// 	}, nil
// }
