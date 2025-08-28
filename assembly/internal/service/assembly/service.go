package assembly

import (
	"github.com/HeyReyHR/rocket-factory/assembly/internal/repository"
	serviceInterface "github.com/HeyReyHR/rocket-factory/assembly/internal/service"
)

var _ serviceInterface.AssemblyService = (*service)(nil)

type service struct {
	assemblyRepository repository.AssemblyRepository
}

func NewService(assemblyRepository repository.AssemblyRepository) *service {
	return &service{
		assemblyRepository: assemblyRepository,
	}
}
