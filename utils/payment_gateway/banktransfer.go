package paymentgateway

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type BankPayment struct {
	Bank Bank
}

func (bp *BankPayment) Charge(pg *PaymetGateway) (*ChargeResponse, error) {
	if bp.Bank != mandiri {
		pg.Request.PaymentType = "bank_transfer"
		pg.Request.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtrans.Bank(bp.Bank),
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
