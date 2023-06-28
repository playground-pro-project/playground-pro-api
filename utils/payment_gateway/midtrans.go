package paymentgateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/playground-pro-project/playground-pro-api/app/config"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	reservation "github.com/playground-pro-project/playground-pro-api/features/reservation"
)

var log = middlewares.Log()

type PaymetGateway struct {
	Request *coreapi.ChargeReq
}

type PaymentMethod interface {
	Charge(*PaymetGateway) (*ChargeResponse, error)
}

type MyRefund struct{}

type Refund interface {
	RefundTransaction(invoice string, amount int64, reason string) error
}

func (r *MyRefund) RefundTransaction(invoice string, amount int64, reason string) error {
	client := coreapi.Client{}
	client.New(config.MIDTRANS_SERVERKEY, midtrans.Sandbox)

	suffix := uuid.New().String()
	refundKey := config.MIDTRANS_MERCHANT_ID + suffix

	refundRequest := &coreapi.RefundReq{
		RefundKey: refundKey,
		Amount:    amount,
		Reason:    reason,
	}

	_, err := client.RefundTransaction(invoice, refundRequest)
	if err != nil {
		return err
	}

	return nil
}

func ChargeMidtrans(reservationID string, request reservation.PaymentCore) (*ChargeResponse, error) {
	client := coreapi.Client{}
	client.New(config.MIDTRANS_SERVERKEY, midtrans.Sandbox)

	if reservationID == "" {
		log.Error("error reservationID")
		return nil, errors.New("invalid reservationID")
	}

	grandTotal, err := strconv.ParseFloat(request.GrandTotal, 64)
	if err != nil {
		log.Error("error parsing grand_total")
		return nil, err
	}

	pg := PaymetGateway{
		Request: &coreapi.ChargeReq{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  reservationID,
				GrossAmt: int64(grandTotal),
			},
		},
	}

	paymentTypeMap := map[string]PaymentMethod{
		"bri":       &BankPayment{Bank: bri},
		"bca":       &BankPayment{Bank: bca},
		"bni":       &BankPayment{Bank: bni},
		"mandiri":   &BankPayment{Bank: mandiri},
		"permata":   &BankPayment{Bank: permata},
		"indomaret": &ConvStorePayment{Store: indomaret},
		"alfamart":  &ConvStorePayment{Store: alfamart},
		"gopay":     &EWalletPayment{EWallet: gopay},
		"shopeepay": &EWalletPayment{EWallet: shopeepay},
		"qris":      &EWalletPayment{EWallet: qris},
	}

	paymentMethod, ok := paymentTypeMap[request.PaymentType]
	if !ok {
		return nil, errors.New("invalid payment_type")
	}

	return paymentMethod.Charge(&pg)
}

func (pg *PaymetGateway) CustomCharge(request *coreapi.ChargeReq) (*ChargeResponse, error) {
	client := coreapi.Client{}
	client.New(config.MIDTRANS_SERVERKEY, midtrans.Sandbox)

	result := ChargeResponse{}
	jsonRequest, _ := json.Marshal(request)
	err := client.HttpClient.Call(http.MethodPost, fmt.Sprintf("%s/v2/charge", client.Env.BaseUrl()), &client.ServerKey, client.Options, bytes.NewBuffer(jsonRequest), &result)
	if err != nil {
		return nil, err
	}

	switch result.PaymentType {
	case "bank_transfer", "echannel":
		if result.PermataVaNumber != "" {
			result.PaymentCode = result.PermataVaNumber
		} else if result.BillerCode != "" || result.BillKey != "" {
			result.PaymentCode = fmt.Sprintf("BillCode:%s-BillKey:%s", result.BillerCode, result.BillKey)
		} else {
			result.PaymentCode = result.VaNumbers[0].VANumber
		}
	case "gopay", "shopeepay", "qris":
		result.PaymentCode = result.Actions[0].URL
	}

	return &result, nil
}

func IsRefundable(paymentMethod string) bool {
	refundableMethods := []string{"bank_transfer", "cstore", "echannel"}
	for _, method := range refundableMethods {
		if method == paymentMethod {
			return true
		}
	}
	return false
}
