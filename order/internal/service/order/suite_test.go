package order

import (
	"testing"

	mocksClient "github.com/HeyReyHR/rocket-factory/order/internal/client/mocks"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	orderRepository *mocks.OrderRepository // ts ain't working on table tests (see inventory)
	inventoryClient *mocksClient.InventoryClient
	paymentClient   *mocksClient.PaymentClient
	service         *service
}

func (s *ServiceSuite) SetupTest() {
	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.inventoryClient = mocksClient.NewInventoryClient(s.T())
	s.paymentClient = mocksClient.NewPaymentClient(s.T())
	s.service = NewService(s.inventoryClient, s.paymentClient, s.orderRepository)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
