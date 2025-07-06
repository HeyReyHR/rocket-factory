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

	mu        sync.RWMutex
	inventory map[string]*invV1.Part
}

func (s *inventoryService) GetPart(_ context.Context, r *invV1.GetPartRequest) (*invV1.GetPartResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	part, ok := s.inventory[r.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", r.GetUuid())
	}
	return &invV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *inventoryService) ListParts(_ context.Context, r *invV1.ListPartsRequest) (*invV1.ListPartsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var parts []*invV1.Part
	for _, part := range s.inventory {
		parts = append(parts, part)
	}

	filteredParts := s.filterParts(parts, r.Filter)

	return &invV1.ListPartsResponse{
		Parts: filteredParts,
	}, nil
}

func (s *inventoryService) filterParts(parts []*invV1.Part, filter *invV1.PartsFilter) []*invV1.Part {
	var result []*invV1.Part

	for _, part := range parts {
		if s.matchesFilter(part, filter) {
			result = append(result, part)
		}
	}

	return result
}

func (s *inventoryService) matchesFilter(part *invV1.Part, filter *invV1.PartsFilter) bool {
	return s.matchesUUIDs(part, filter.Uuids) &&
		s.matchesNames(part, filter.Names) &&
		s.matchesCountries(part, filter.ManufacturerCountries) &&
		s.matchesCategories(part, filter.Categories) &&
		s.matchesTags(part, filter.Tags)
}

func (s *inventoryService) matchesUUIDs(part *invV1.Part, uuids []string) bool {
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

func (s *inventoryService) matchesNames(part *invV1.Part, names []string) bool {
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

func (s *inventoryService) matchesCountries(part *invV1.Part, countries []string) bool {
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

func (s *inventoryService) matchesCategories(part *invV1.Part, categories []invV1.Category) bool {
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

func (s *inventoryService) matchesTags(part *invV1.Part, filterTags []string) bool {
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
	service := &inventoryService{
		inventory: map[string]*invV1.Part{
			"1": {
				Uuid:          "1",
				Name:          "Raptor Engine V2",
				Description:   "Метановый ракетный двигатель полного сгорания",
				Price:         2500000.0,
				StockQuantity: 15,
				Category:      invV1.Category_ENGINE,
				Manufacturer: &invV1.Manufacturer{
					Name:    "SpaceX",
					Country: "USA",
					Website: "https://spacex.com",
				},
				Tags: []string{"methane", "reusable", "high-performance"},
				Metadata: map[string]*invV1.Value{
					"thrust": {
						ValueType: &invV1.Value_StringValue{StringValue: "2300 kN"},
					},
					"isp_vacuum": {
						ValueType: &invV1.Value_Int64Value{Int64Value: 380},
					},
					"mass": {
						ValueType: &invV1.Value_DoubleValue{DoubleValue: 1600.0},
					},
					"reusable": {
						ValueType: &invV1.Value_BoolValue{BoolValue: true},
					},
				},
			},
			"2": {
				Uuid:          "2",
				Name:          "RP-1 Fuel Tank",
				Description:   "Алюминиевый топливный бак для керосина RP-1",
				Price:         450000.0,
				StockQuantity: 8,
				Category:      invV1.Category_FUEL,
				Manufacturer: &invV1.Manufacturer{
					Name:    "Blue Origin",
					Country: "USA",
					Website: "https://blueorigin.com",
				},
				Tags: []string{"fuel-tank", "aluminum", "rp1"},
				Metadata: map[string]*invV1.Value{
					"capacity": {
						ValueType: &invV1.Value_StringValue{StringValue: "150000 L"},
					},
					"max_pressure": {
						ValueType: &invV1.Value_DoubleValue{DoubleValue: 2.5},
					},
					"empty_weight": {
						ValueType: &invV1.Value_DoubleValue{DoubleValue: 2500.0},
					},
				},
			},
			"3": {
				Uuid:          "3",
				Name:          "Dragon Porthole Window",
				Description:   "Круглое окно для космического корабля Dragon",
				Price:         75000.0,
				StockQuantity: 25,
				Category:      invV1.Category_PORTHOLE,
				Manufacturer: &invV1.Manufacturer{
					Name:    "SpaceX",
					Country: "CHINA",
					Website: "https://spacex.com",
				},
				Tags: []string{"window", "transparent", "pressurized", "lol"},
				Metadata: map[string]*invV1.Value{
					"diameter": {
						ValueType: &invV1.Value_StringValue{StringValue: "45 cm"},
					},
					"thickness": {
						ValueType: &invV1.Value_DoubleValue{DoubleValue: 5.0},
					},
					"material": {
						ValueType: &invV1.Value_StringValue{StringValue: "Borosilicate glass"},
					},
					"pressure_rating": {
						ValueType: &invV1.Value_DoubleValue{DoubleValue: 1.5},
					},
				},
			},
			"4": {
				Uuid:          "4",
				Name:          "Grid Fin Assembly",
				Description:   "Аэродинамическое крыло для управления при возвращении",
				Price:         185000.0,
				StockQuantity: 12,
				Category:      invV1.Category_WING,
				Manufacturer: &invV1.Manufacturer{
					Name:    "SpaceX",
					Country: "USA",
					Website: "https://spacex.com",
				},
				Tags: []string{"grid-fin", "titanium", "reentry", "lol"},
				Metadata: map[string]*invV1.Value{
					"material": {
						ValueType: &invV1.Value_StringValue{StringValue: "Titanium"},
					},
					"fin_count": {
						ValueType: &invV1.Value_Int64Value{Int64Value: 4},
					},
					"max_temperature": {
						ValueType: &invV1.Value_DoubleValue{DoubleValue: 1200.0},
					},
					"deployable": {
						ValueType: &invV1.Value_BoolValue{BoolValue: true},
					},
				},
			},
		},
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
