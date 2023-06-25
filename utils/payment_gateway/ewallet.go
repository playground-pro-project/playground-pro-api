package paymentgateway

import "github.com/midtrans/midtrans-go/coreapi"

type EWalletPayment struct {
	EWallet Ewallet
}

func (ew *EWalletPayment) Charge(pg *PaymetGateway) (*ChargeResponse, error) {
	pg.Request.PaymentType = coreapi.CoreapiPaymentType(ew.EWallet)
	return pg.CustomCharge(pg.Request)
}
