package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	api "github.com/HeyReyHR/rocket-factory/inventory/internal/api/inventory/v1"
	repository "github.com/HeyReyHR/rocket-factory/inventory/internal/repository/inventory"
	service "github.com/HeyReyHR/rocket-factory/inventory/internal/service/inventory"
	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

func main() {
	ctx := context.Background()

	dbURI := "mongodb://inventory-service-user:mongo@localhost:27017/database?authSource=admin"

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v\n", err)
	}

	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("failed to disconnect from MongoDB: %v\n", cerr)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping MongoDB: %v\n", err)
	}

	db := client.Database("inventory-service")

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

	newRepository := repository.NewRepository(db)

	newService := service.NewService(newRepository)

	newApi := api.NewApi(newService)

	invV1.RegisterInventoryServiceServer(s, newApi)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 Inventory service is running on port %d", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down inventory service...")
	s.GracefulStop()
	log.Println("✅ Inventory service stopped")
}
