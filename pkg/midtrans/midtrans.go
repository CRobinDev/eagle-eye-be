package midtrans

import (
	"errors"
	"fmt"

	"github.com/CRobinDev/karsa/config/env"
	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/google/uuid"
	midtransGateway "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type IMidtrans interface {
	NewTransactionToken(req dto.PaymentRequest) (*snap.Response, error)
	GetServerKey() string
}

type midtrans struct {
	Client    *snap.Client
	serverKey string
}

func NewMidtransService(client *snap.Client) IMidtrans {
	return &midtrans{
		Client:    client,
		serverKey: env.GetEnv().MidtransServerKey,
	}
}

func (m *midtrans) NewTransactionToken(req dto.PaymentRequest) (*snap.Response, error) {
	request := &snap.Request{
		TransactionDetails: midtransGateway.TransactionDetails{
			OrderID:  req.OrderID,
			GrossAmt: req.Amount,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		Items: &[]midtransGateway.ItemDetails{
			{
				ID:    uuid.NewString(),
				Name:  fmt.Sprintf("Konsultasi Pembelian Paket %s di EagleEye 😊", req.TierOrder),
				Price: req.Amount,
				Qty:   1,
			},
		},
		EnabledPayments: snap.AllSnapPaymentType,
		CustomerDetail: &midtransGateway.CustomerDetails{
			FName: req.CustomerName,
			Email: req.CustomerEmail,
		},
		Expiry: &snap.ExpiryDetails{
			Duration: 30,
			Unit:     "minute",
		},
	}

	snapResp, err := m.Client.CreateTransaction(request)
	var midtransErr *midtransGateway.Error
	if errors.As(err, &midtransErr) && midtransErr == nil {
		return snapResp, nil
	}

	return snapResp, err
}

func (m *midtrans) GetServerKey() string {
	return m.serverKey
}
