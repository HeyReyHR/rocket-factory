package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const grpcPort = 50051

type inventoryService struct {
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
	storage.mu.Lock()
	defer storage.mu.Unlock()
	var parts []*invV1.Part
	for _, part := range storage.inventory {
		parts = append(parts, part)
	}

	filteredParts := storage.filterParts(parts, filter)
	return filteredParts, nil
}

func (s *inventoryService) GetPart(_ context.Context, r *invV1.GetPartRequest) (*invV1.GetPartResponse, error) {
	part, err := s.inventory.Part(r.GetUuid())
	if err != nil {
		log.Printf("GetPart failed for UUID %s: %v", r.GetUuid(), err)
		return nil, err
	}

	return &invV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *inventoryService) ListParts(_ context.Context, r *invV1.ListPartsRequest) (*invV1.ListPartsResponse, error) {
	filteredParts, err := s.inventory.Parts(r.Filter)
	if err != nil {
		log.Printf("GetParts failed: %v", err)
		return nil, err
	}

	return &invV1.ListPartsResponse{
		Parts: filteredParts,
	}, nil
}

func (storage *InventoryStorageInMem) filterParts(parts []*invV1.Part, filter *invV1.PartsFilter) []*invV1.Part {
	var result []*invV1.Part
	for _, part := range parts {
		if storage.matchesFilter(part, filter) {
			result = append(result, part)
		}
	}

	return result
}

func (storage *InventoryStorageInMem) matchesFilter(part *invV1.Part, filter *invV1.PartsFilter) bool {
	return storage.matchesUUIDs(part, filter.Uuids) &&
		storage.matchesNames(part, filter.Names) &&
		storage.matchesCountries(part, filter.ManufacturerCountries) &&
		storage.matchesCategories(part, filter.Categories) &&
		storage.matchesTags(part, filter.Tags)
}

func (storage *InventoryStorageInMem) matchesUUIDs(part *invV1.Part, uuids []string) bool {
	if uuids == nil {
		return true
	}

	for _, uuid := range uuids {
		if strings.EqualFold(part.Uuid, uuid) {
			return true
		}
	}

	return false
}

func (storage *InventoryStorageInMem) matchesNames(part *invV1.Part, names []string) bool {
	if names == nil {
		return true
	}

	partNameLower := strings.ToLower(part.Name)
	for _, name := range names {
		if strings.Contains(partNameLower, strings.ToLower(name)) {
			return true
		}
	}

	return false
}

func (storage *InventoryStorageInMem) matchesCountries(part *invV1.Part, countries []string) bool {
	if countries == nil {
		return true
	}

	for _, country := range countries {
		if strings.EqualFold(part.Manufacturer.Country, country) {
			return true
		}
	}

	return false
}

func (storage *InventoryStorageInMem) matchesCategories(part *invV1.Part, categories []invV1.Category) bool {
	if categories == nil {
		return true
	}

	for _, category := range categories {
		if strings.EqualFold(part.Category.String(), category.String()) {
			return true
		}
	}

	return false
}

func (storage *InventoryStorageInMem) matchesTags(part *invV1.Part, filterTags []string) bool {
	if filterTags == nil {
		return true
	}

	for _, filterTag := range filterTags {
		for _, partTag := range part.Tags {
			if strings.EqualFold(partTag, filterTag) {
				return true
			}
		}
	}

	return false
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

	storage := &InventoryStorageInMem{
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
					Country: "USA",
					Website: "https://rocketcorp.com",
				},
				Tags: []string{"engine", "high-performance"},
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
				Tags: []string{"fuel", "liquid"},
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

	service := &inventoryService{
		inventory: storage,
	}

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
