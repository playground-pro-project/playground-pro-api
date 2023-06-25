package paymentgateway

import "github.com/midtrans/midtrans-go/coreapi"

type ConvStorePayment struct {
	Store ConvStore
}

func (cs *ConvStorePayment) Charge(pg *PaymetGateway) (*ChargeResponse, error) {
	pg.Request.PaymentType = "cstore"
	pg.Request.ConvStore = &coreapi.ConvStoreDetails{
		Store: string(cs.Store),
	}

	return pg.CustomCharge(pg.Request)
}
