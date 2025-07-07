package midtrans

import (
	"github.com/CRobinDev/karsa/config/env"
	midtransGateway "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func NewMidtrans() *snap.Client {
	serverKey := env.GetEnv().MidtransServerKey
	client := snap.Client{}
	client.New(serverKey, midtransGateway.Sandbox)

	return &client
}
