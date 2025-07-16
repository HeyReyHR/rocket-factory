package inventory

import (
	"testing"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	inventoryRepository *mocks.InventoryRepository
	service             *service
}

func (s *ServiceSuite) SetupTest() {
	s.inventoryRepository = mocks.NewInventoryRepository(s.T())

	s.service = NewService(s.inventoryRepository)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
