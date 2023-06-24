package paymentgateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

func ChargeMidtrans(reservationID string, request reservation.PaymentCore) (*ChargeResponse, error) {
	client := coreapi.Client{}
	client.New(config.MIDTRANS_SERVERKEY, midtrans.Sandbox)

	if reservationID == "" {
		log.Error("error reservationID")
		return nil, errors.New("invalid reservationID")
	}

	grandTotal, err := strconv.ParseInt(request.GrandTotal, 10, 64)
	if err != nil {
		log.Error("error parsing grand_total")
		return nil, err
	}

	pg := PaymetGateway{
		Request: &coreapi.ChargeReq{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  reservationID,
				GrossAmt: grandTotal,
			},
		},
	}

	paymentTypeMap := map[string]func() (*ChargeResponse, error){
		"bca": func() (*ChargeResponse, error) {
			return pg.ChargeWithBank(bca)
		},
		"mandiri": func() (*ChargeResponse, error) {
			return pg.ChargeWithBank(mandiri)
		},
		"bni": func() (*ChargeResponse, error) {
			return pg.ChargeWithBank(bni)
		},
		"indomaret": func() (*ChargeResponse, error) {
			return pg.ChargeWithConvStore(indomaret)
		},
		"alfamart": func() (*ChargeResponse, error) {
			return pg.ChargeWithConvStore(alfamart)
		},
		"gopay": func() (*ChargeResponse, error) {
			return pg.ChargeWithEWallet(gopay)
		},
		"shopeepay": func() (*ChargeResponse, error) {
			return pg.ChargeWithEWallet(shopeepay)
		},
		"qris": func() (*ChargeResponse, error) {
			return pg.ChargeWithEWallet(qris)
		},
	}

	chargeFunc, ok := paymentTypeMap[request.PaymentType]
	if !ok {
		return nil, errors.New("invalid payment_type")
	}

	return chargeFunc()
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

func (pg *PaymetGateway) ChargeWithBank(b Bank) (*ChargeResponse, error) {
	if b != mandiri {
		pg.Request.PaymentType = "bank_transfer"
		pg.Request.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtrans.Bank(b),
		}
	} else {
		pg.Request.PaymentType = coreapi.PaymentTypeEChannel
		pg.Request.EChannel = &coreapi.EChannelDetail{
			BillInfo1: "pembayaran",
			BillInfo2: "pembayaran",
		}
	}

	return pg.CustomCharge(pg.Request)
}

func (pg *PaymetGateway) ChargeWithEWallet(w Ewallet) (*ChargeResponse, error) {
	pg.Request.PaymentType = coreapi.CoreapiPaymentType(w)
	return pg.CustomCharge(pg.Request)
}

func (pg *PaymetGateway) ChargeWithConvStore(c ConvStore) (*ChargeResponse, error) {
	pg.Request.PaymentType = "cstore"
	pg.Request.ConvStore = &coreapi.ConvStoreDetails{
		Store: string(c),
	}

	return pg.CustomCharge(pg.Request)
}

func (pg *PaymetGateway) Refund(request *coreapi.RefundReq, invoice string) error {
	client := coreapi.Client{}
	client.New(config.MIDTRANS_SERVERKEY, midtrans.Sandbox)
	_, err := client.RefundTransaction(invoice, request)
	if err != nil {
		return err
	}

	return nil
}
