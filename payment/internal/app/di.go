package app

import (
	"context"
	payV1API "github.com/HeyReyHR/rocket-factory/payment/internal/api/payment/v1"

	"github.com/HeyReyHR/rocket-factory/payment/internal/service"
	paymentService "github.com/HeyReyHR/rocket-factory/payment/internal/service/payment"

	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
	"go.mongodb.org/mongo-driver/mongo"
)

type diContainer struct {
	paymentV1API payV1.PaymentServiceServer

	paymentService service.PaymentService

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API(ctx context.Context) payV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = payV1API.NewApi(d.PaymentService(ctx))
	}

	return d.paymentV1API
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}

	return d.paymentService
}
