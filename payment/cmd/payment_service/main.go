package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	v1 "github.com/HeyReyHR/rocket-factory/payment/internal/api/payment/v1"
	"github.com/HeyReyHR/rocket-factory/payment/internal/service/payment"
	"github.com/HeyReyHR/rocket-factory/shared/pkg/interceptors"
	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

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

	s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.UnaryErrorInterceptor()))

	service := payment.NewService()

	api := v1.NewPaymentApi(service)

	payV1.RegisterPaymentServiceServer(s, api)

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
