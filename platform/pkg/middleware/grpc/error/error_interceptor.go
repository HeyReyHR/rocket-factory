package error

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, convertError(ctx, err, info.FullMethod)
		}
		return resp, nil
	}
}

func convertError(ctx context.Context, err error, method string) error {
	if businessErr := GetBusinessError(err); businessErr != nil {
		grpcStatus := BusinessErrorToGRPCStatus(businessErr)
		logger.Error(ctx, "BusinessError:", zap.String("method", method), zap.Int("code", int(businessErr.Code())), zap.String("error", businessErr.Error()))
		return grpcStatus.Err()
	}
	if _, ok := status.FromError(err); ok {
		return err
	}
	logger.Error(ctx, "Unknown error:", zap.String("method", method), zap.Error(err))
	return status.Error(codes.Internal, "internal server error")
}
