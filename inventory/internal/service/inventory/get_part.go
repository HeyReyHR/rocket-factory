package inventory

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/middleware/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	_, ok := grpc.GetUserFromContext(ctx)
	if !ok {
		return model.Part{}, status.Error(codes.Unauthenticated, "user not found")
	}
	
	part, err := s.inventoryRepository.GetPart(ctx, uuid)
	if err != nil {
		return model.Part{}, err
	}
	return part, nil
}
