package inventory

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/middleware/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error) {
	_, ok := grpc.GetUserFromContext(ctx)
	if !ok {
		return []model.Part{}, status.Error(codes.Unauthenticated, "user not found")
	}
	
	parts, err := s.inventoryRepository.ListParts(ctx, filter)
	if err != nil {
		return []model.Part{}, err
	}
	return parts, nil
}
