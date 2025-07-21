package v1

import payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"

type client struct {
	generatedClient payV1.PaymentServiceClient
}

func NewPaymentClient(generatedClient payV1.PaymentServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
