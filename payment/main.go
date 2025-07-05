package main

import (
	"context"
	"fmt"
	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const grpcPort = 50052

type paymentService struct {
	payV1.UnimplementedPaymentServiceServer

	mu sync.RWMutex
}

func (s *paymentService) PayOrder(_ context.Context, r *payV1.PayOrderRequest) (*payV1.PayOrderResponse, error) {
	transactionUuid := uuid.NewString()
	fmt.Printf("Payment has been succeded, transaction_uuid: %s", transactionUuid)
	return &payV1.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}, nil
}
func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()
	service := &paymentService{}
	payV1.RegisterPaymentServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 Payment service is running on port %d", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down payment service...")
	s.GracefulStop()
	log.Println("✅ Payment service stopped")

}
