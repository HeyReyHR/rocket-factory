package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const grpcPort = 50051

type InventoryService struct {
	invV1.UnimplementedInventoryServiceServer
	inventory InventoryStorage
}
type InventoryStorageInMem struct {
	mu        sync.RWMutex
	inventory map[string]*invV1.Part
}
type InventoryStorage interface {
	Part(uuid string) (*invV1.Part, error)
	Parts(filter *invV1.PartsFilter) ([]*invV1.Part, error)
}

func NewInventoryStorage() *InventoryStorageInMem {
	return &InventoryStorageInMem{
		inventory: map[string]*invV1.Part{
			"engine-001": {
				Uuid:          "engine-001",
				Name:          "Rocket Engine V1",
				Description:   "High-performance rocket engine",
				Price:         15000.50,
				StockQuantity: 10,
				Category:      invV1.Category_ENGINE,
				Manufacturer: &invV1.Manufacturer{
					Name:    "RocketCorp",
					Country: "France",
					Website: "https://rocketcorp.com",
				},
				Tags: []string{"engine", "high-performance", "lol"},
				Dimensions: &invV1.Dimensions{
					Length: 2.5,
					Width:  1.0,
					Height: 1.0,
					Weight: 500.0,
				},
			},
			"fuel-001": {
				Uuid:          "fuel-001",
				Name:          "Liquid Fuel Tank",
				Description:   "High-capacity fuel storage",
				Price:         8500.00,
				StockQuantity: 25,
				Category:      invV1.Category_FUEL,
				Manufacturer: &invV1.Manufacturer{
					Name:    "FuelTech",
					Country: "Germany",
					Website: "https://fueltech.de",
				},
				Tags: []string{"fuel", "liquid", "lol"},
				Dimensions: &invV1.Dimensions{
					Length: 3.0,
					Width:  1.5,
					Height: 1.5,
					Weight: 200.0,
				},
			},
			"wing-001": {
				Uuid:          "wing-001",
				Name:          "Stabilizer Wing",
				Description:   "Aerodynamic stabilizer wing",
				Price:         3200.75,
				StockQuantity: 15,
				Category:      invV1.Category_WING,
				Manufacturer: &invV1.Manufacturer{
					Name:    "AeroWings",
					Country: "France",
					Website: "https://aerowings.fr",
				},
				Tags: []string{"wing", "stabilizer"},
				Dimensions: &invV1.Dimensions{
					Length: 2.0,
					Width:  0.5,
					Height: 0.1,
					Weight: 50.0,
				},
			},
		},
	}
}

func NewInventoryService(inventory InventoryStorage) *InventoryService {
	return &InventoryService{
		inventory: inventory,
	}
}

func (storage *InventoryStorageInMem) Part(uuid string) (*invV1.Part, error) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()
	part, ok := storage.inventory[uuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", uuid)
	}
	return part, nil
}

func (storage *InventoryStorageInMem) Parts(filter *invV1.PartsFilter) ([]*invV1.Part, error) {
	storage.mu.RLock()

	var parts []invV1.Part
	for _, part := range storage.inventory {
		parts = append(parts, invV1.Part{
			Uuid:          part.Uuid,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
			Category:      part.Category,
			Manufacturer:  part.Manufacturer,
			Tags:          part.Tags,
			Dimensions:    part.Dimensions,
		})
	}

	storage.mu.RUnlock()

	filteredParts := storage.filterParts(parts, filter)
	return filteredParts, nil
}

func (s *InventoryService) GetPart(_ context.Context, r *invV1.GetPartRequest) (*invV1.GetPartResponse, error) {
	part, err := s.inventory.Part(r.GetUuid())
	if err != nil {
		log.Printf("GetPart failed for UUID %s: %v", r.GetUuid(), err)
		return nil, status.Errorf(codes.Internal, "GetPart failed for UUID %s: %v", r.GetUuid(), err)
	}

	return &invV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *InventoryService) ListParts(_ context.Context, r *invV1.ListPartsRequest) (*invV1.ListPartsResponse, error) {
	filteredParts, err := s.inventory.Parts(r.Filter)
	if err != nil {
		log.Printf("GetParts failed: %v", err)
		return nil, status.Errorf(codes.Internal, "GetParts failed: %v", err)
	}

	return &invV1.ListPartsResponse{
		Parts: filteredParts,
	}, nil
}

type filter func(part *invV1.Part) bool

func (storage *InventoryStorageInMem) filterParts(parts []invV1.Part, filterReq *invV1.PartsFilter) []*invV1.Part {
	var result []*invV1.Part

	filters := makeFilters(filterReq)

	for i := range parts {
		needAdd := true
		for _, filter := range filters {
			if !filter(&parts[i]) {
				needAdd = false
				break
			}
		}
		if needAdd {
			result = append(result, &parts[i])
		}
	}

	return result
}

func makeFilters(filterReq *invV1.PartsFilter) []filter {
	var filters []filter

	if len(filterReq.GetUuids()) > 0 {
		filters = append(filters, func(part *invV1.Part) bool {
			return slices.Contains(filterReq.GetUuids(), part.Uuid)
		})
	}

	if len(filterReq.GetNames()) > 0 {
		filters = append(filters, func(part *invV1.Part) bool {
			return slices.Contains(filterReq.GetNames(), part.Name)
		})
	}

	if len(filterReq.GetCategories()) > 0 {
		filters = append(filters, func(part *invV1.Part) bool {
			return slices.Contains(filterReq.GetCategories(), part.Category)
		})
	}

	if len(filterReq.GetManufacturerCountries()) > 0 {
		filters = append(filters, func(part *invV1.Part) bool {
			return slices.Contains(filterReq.GetManufacturerCountries(), part.GetManufacturer().GetCountry())
		})
	}

	if len(filterReq.GetTags()) > 0 {
		filters = append(filters, func(part *invV1.Part) bool {
			for _, tag := range part.GetTags() {
				if slices.Contains(filterReq.GetTags(), tag) {
					return true
				}
			}
			return false
		})
	}

	return filters
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

	storage := NewInventoryStorage()

	service := NewInventoryService(storage)

	invV1.RegisterInventoryServiceServer(s, service)

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
