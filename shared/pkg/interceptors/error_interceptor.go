package interceptors

import (
	"context"
	"github.com/HeyReyHR/rocket-factory/shared/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, convertError(err, info.FullMethod)
		}
		return resp, nil
	}
}

func convertError(err error, method string) error {
	if businessErr := errors.GetBusinessError(err); businessErr != nil {
		grpcStatus := errors.BusinessErrorToGRPCStatus(businessErr)
		log.Printf("BusinessError in method %s: code=%d, message=%s",
			method, businessErr.Code(), businessErr.Error())
		return grpcStatus.Err()
	}
	if _, ok := status.FromError(err); ok {
		return err
	}
	log.Printf("Unknown error in method %s: %v", method, err)
	return status.Error(codes.Internal, "internal server error")
}
