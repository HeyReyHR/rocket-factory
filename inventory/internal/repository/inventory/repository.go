package inventory

import (
	"sync"
	"time"

	def "github.com/HeyReyHR/rocket-factory/inventory/internal/repository"
	"github.com/HeyReyHR/rocket-factory/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]model.Part
}

func NewRepository() *repository {
	return &repository{
		data: map[string]model.Part{
			// Двигатели
			"engine-001": {
				Uuid:          "engine-001",
				Name:          "Rocket Engine V1",
				Description:   "High-performance rocket engine",
				Price:         15000.50,
				StockQuantity: 0,
				Category:      model.ENGINE,
				Manufacturer: model.Manufacturer{
					Name:    "RocketCorp",
					Country: "France",
					Website: "https://rocketcorp.com",
				},
				Tags: []string{"engine", "high-performance", "liquid"},
				Dimensions: model.Dimensions{
					Length: 2.5,
					Width:  1.0,
					Height: 1.0,
					Weight: 500.0,
				},
				Metadata: map[string]model.Value{
					"max_thrust": {Int64Value: &[]int64{25000}[0]},
					"fuel_type":  {StringValue: &[]string{"liquid"}[0]},
					"efficiency": {DoubleValue: &[]float64{0.95}[0]},
					"tested":     {BoolValue: &[]bool{true}[0]},
				},
				CreatedAt: time.Date(2024, 11, 3, 14, 27, 19, 0, time.UTC),
				UpdatedAt: time.Date(2024, 11, 3, 14, 27, 19, 0, time.UTC),
			},
			"engine-002": {
				Uuid:          "engine-002",
				Name:          "Advanced Turbo Engine",
				Description:   "Next-generation turbo rocket engine",
				Price:         28500.00,
				StockQuantity: 5,
				Category:      model.ENGINE,
				Manufacturer: model.Manufacturer{
					Name:    "TurboTech",
					Country: "USA",
					Website: "https://turbotech.com",
				},
				Tags: []string{"engine", "turbo", "advanced", "high-thrust"},
				Dimensions: model.Dimensions{
					Length: 3.2,
					Width:  1.5,
					Height: 1.5,
					Weight: 750.0,
				},
				Metadata: map[string]model.Value{
					"max_thrust": {Int64Value: &[]int64{45000}[0]},
					"fuel_type":  {StringValue: &[]string{"hybrid"}[0]},
					"efficiency": {DoubleValue: &[]float64{0.98}[0]},
					"tested":     {BoolValue: &[]bool{true}[0]},
				},
				CreatedAt: time.Date(2024, 11, 1, 10, 15, 30, 0, time.UTC),
				UpdatedAt: time.Date(2024, 11, 1, 10, 15, 30, 0, time.UTC),
			},
			"engine-003": {
				Uuid:          "engine-003",
				Name:          "Compact Engine",
				Description:   "Small but powerful engine for lightweight rockets",
				Price:         8900.00,
				StockQuantity: 20,
				Category:      model.ENGINE,
				Manufacturer: model.Manufacturer{
					Name:    "MiniProp",
					Country: "Japan",
					Website: "https://miniprop.jp",
				},
				Tags: []string{"engine", "compact", "lightweight"},
				Dimensions: model.Dimensions{
					Length: 1.8,
					Width:  0.8,
					Height: 0.8,
					Weight: 250.0,
				},
				Metadata: map[string]model.Value{
					"max_thrust": {Int64Value: &[]int64{12000}[0]},
					"fuel_type":  {StringValue: &[]string{"solid"}[0]},
					"efficiency": {DoubleValue: &[]float64{0.85}[0]},
					"tested":     {BoolValue: &[]bool{false}[0]},
				},
				CreatedAt: time.Date(2024, 10, 28, 16, 42, 10, 0, time.UTC),
				UpdatedAt: time.Date(2024, 10, 28, 16, 42, 10, 0, time.UTC),
			},

			// Топливные баки
			"fuel-001": {
				Uuid:          "fuel-001",
				Name:          "Liquid Fuel Tank",
				Description:   "High-capacity fuel storage",
				Price:         8500.00,
				StockQuantity: 25,
				Category:      model.FUEL,
				Manufacturer: model.Manufacturer{
					Name:    "FuelTech",
					Country: "Germany",
					Website: "https://fueltech.de",
				},
				Tags: []string{"fuel", "liquid", "storage"},
				Dimensions: model.Dimensions{
					Length: 3.0,
					Width:  1.5,
					Height: 1.5,
					Weight: 200.0,
				},
				Metadata: map[string]model.Value{
					"capacity":  {Int64Value: &[]int64{1000}[0]},
					"fuel_type": {StringValue: &[]string{"liquid_oxygen"}[0]},
					"pressure":  {DoubleValue: &[]float64{150.5}[0]},
					"insulated": {BoolValue: &[]bool{true}[0]},
				},
				CreatedAt: time.Date(2024, 11, 3, 14, 27, 19, 0, time.UTC),
				UpdatedAt: time.Date(2024, 11, 3, 14, 27, 19, 0, time.UTC),
			},
			"fuel-002": {
				Uuid:          "fuel-002",
				Name:          "Cryogenic Fuel Tank",
				Description:   "Ultra-cold fuel storage system",
				Price:         15750.00,
				StockQuantity: 8,
				Category:      model.FUEL,
				Manufacturer: model.Manufacturer{
					Name:    "CryoSystems",
					Country: "Russia",
					Website: "https://cryosystems.ru",
				},
				Tags: []string{"fuel", "cryogenic", "ultra-cold"},
				Dimensions: model.Dimensions{
					Length: 4.0,
					Width:  2.0,
					Height: 2.0,
					Weight: 400.0,
				},
				Metadata: map[string]model.Value{
					"capacity":  {Int64Value: &[]int64{1500}[0]},
					"fuel_type": {StringValue: &[]string{"liquid_hydrogen"}[0]},
					"min_temp":  {DoubleValue: &[]float64{-253.0}[0]},
					"insulated": {BoolValue: &[]bool{true}[0]},
				},
				CreatedAt: time.Date(2024, 10, 30, 11, 20, 45, 0, time.UTC),
				UpdatedAt: time.Date(2024, 10, 30, 11, 20, 45, 0, time.UTC),
			},
			"fuel-003": {
				Uuid:          "fuel-003",
				Name:          "Solid Fuel Module",
				Description:   "Solid propellant fuel module",
				Price:         6200.00,
				StockQuantity: 35,
				Category:      model.FUEL,
				Manufacturer: model.Manufacturer{
					Name:    "SolidProp",
					Country: "Italy",
					Website: "https://solidprop.it",
				},
				Tags: []string{"fuel", "solid", "propellant"},
				Dimensions: model.Dimensions{
					Length: 2.5,
					Width:  1.0,
					Height: 1.0,
					Weight: 150.0,
				},
				Metadata: map[string]model.Value{
					"burn_time": {Int64Value: &[]int64{120}[0]},
					"fuel_type": {StringValue: &[]string{"solid_composite"}[0]},
					"thrust":    {DoubleValue: &[]float64{15000.0}[0]},
					"reusable":  {BoolValue: &[]bool{false}[0]},
				},
				CreatedAt: time.Date(2024, 11, 2, 9, 18, 30, 0, time.UTC),
				UpdatedAt: time.Date(2024, 11, 2, 9, 18, 30, 0, time.UTC),
			},

			// Крылья
			"wing-001": {
				Uuid:          "wing-001",
				Name:          "Stabilizer Wing",
				Description:   "Aerodynamic stabilizer wing",
				Price:         3200.75,
				StockQuantity: 15,
				Category:      model.WING,
				Manufacturer: model.Manufacturer{
					Name:    "AeroWings",
					Country: "France",
					Website: "https://aerowings.fr",
				},
				Tags: []string{"wing", "stabilizer", "aerodynamic"},
				Dimensions: model.Dimensions{
					Length: 2.0,
					Width:  0.5,
					Height: 0.1,
					Weight: 50.0,
				},
				Metadata: map[string]model.Value{
					"material":    {StringValue: &[]string{"carbon_fiber"}[0]},
					"max_speed":   {Int64Value: &[]int64{3000}[0]},
					"lift_coeff":  {DoubleValue: &[]float64{0.7}[0]},
					"retractable": {BoolValue: &[]bool{false}[0]},
				},
				CreatedAt: time.Date(2024, 11, 3, 14, 27, 19, 0, time.UTC),
				UpdatedAt: time.Date(2024, 11, 3, 14, 27, 19, 0, time.UTC),
			},
			"wing-002": {
				Uuid:          "wing-002",
				Name:          "Control Wing",
				Description:   "Precision flight control wing",
				Price:         4850.00,
				StockQuantity: 12,
				Category:      model.WING,
				Manufacturer: model.Manufacturer{
					Name:    "PrecisionFlight",
					Country: "Sweden",
					Website: "https://precisionflight.se",
				},
				Tags: []string{"wing", "control", "precision"},
				Dimensions: model.Dimensions{
					Length: 1.8,
					Width:  0.4,
					Height: 0.08,
					Weight: 35.0,
				},
				Metadata: map[string]model.Value{
					"material":    {StringValue: &[]string{"titanium_alloy"}[0]},
					"max_speed":   {Int64Value: &[]int64{4000}[0]},
					"lift_coeff":  {DoubleValue: &[]float64{0.8}[0]},
					"retractable": {BoolValue: &[]bool{true}[0]},
				},
				CreatedAt: time.Date(2024, 10, 25, 13, 55, 20, 0, time.UTC),
				UpdatedAt: time.Date(2024, 10, 25, 13, 55, 20, 0, time.UTC),
			},
			"wing-003": {
				Uuid:          "wing-003",
				Name:          "Delta Wing",
				Description:   "High-speed delta configuration wing",
				Price:         5650.00,
				StockQuantity: 8,
				Category:      model.WING,
				Manufacturer: model.Manufacturer{
					Name:    "DeltaAero",
					Country: "Canada",
					Website: "https://deltaaero.ca",
				},
				Tags: []string{"wing", "delta", "high-speed"},
				Dimensions: model.Dimensions{
					Length: 3.5,
					Width:  0.6,
					Height: 0.12,
					Weight: 75.0,
				},
				Metadata: map[string]model.Value{
					"material":    {StringValue: &[]string{"composite"}[0]},
					"max_speed":   {Int64Value: &[]int64{5000}[0]},
					"lift_coeff":  {DoubleValue: &[]float64{0.9}[0]},
					"retractable": {BoolValue: &[]bool{false}[0]},
				},
				CreatedAt: time.Date(2024, 10, 22, 8, 30, 15, 0, time.UTC),
				UpdatedAt: time.Date(2024, 10, 22, 8, 30, 15, 0, time.UTC),
			},

			// Иллюминаторы
			"porthole-001": {
				Uuid:          "porthole-001",
				Name:          "Observation Porthole",
				Description:   "High-strength observation window",
				Price:         1850.00,
				StockQuantity: 30,
				Category:      model.PORTHOLE,
				Manufacturer: model.Manufacturer{
					Name:    "ViewTech",
					Country: "Netherlands",
					Website: "https://viewtech.nl",
				},
				Tags: []string{"porthole", "observation", "transparent"},
				Dimensions: model.Dimensions{
					Length: 0.5,
					Width:  0.5,
					Height: 0.05,
					Weight: 15.0,
				},
				Metadata: map[string]model.Value{
					"material":       {StringValue: &[]string{"reinforced_glass"}[0]},
					"pressure_limit": {Int64Value: &[]int64{200}[0]},
					"transparency":   {DoubleValue: &[]float64{0.98}[0]},
					"heated":         {BoolValue: &[]bool{false}[0]},
				},
				CreatedAt: time.Date(2024, 11, 1, 15, 12, 45, 0, time.UTC),
				UpdatedAt: time.Date(2024, 11, 1, 15, 12, 45, 0, time.UTC),
			},
			"porthole-002": {
				Uuid:          "porthole-002",
				Name:          "Reinforced Porthole",
				Description:   "Extra-strong porthole for extreme conditions",
				Price:         3200.00,
				StockQuantity: 18,
				Category:      model.PORTHOLE,
				Manufacturer: model.Manufacturer{
					Name:    "StrongView",
					Country: "Germany",
					Website: "https://strongview.de",
				},
				Tags: []string{"porthole", "reinforced", "extreme"},
				Dimensions: model.Dimensions{
					Length: 0.6,
					Width:  0.6,
					Height: 0.08,
					Weight: 25.0,
				},
				Metadata: map[string]model.Value{
					"material":       {StringValue: &[]string{"sapphire_crystal"}[0]},
					"pressure_limit": {Int64Value: &[]int64{500}[0]},
					"transparency":   {DoubleValue: &[]float64{0.95}[0]},
					"heated":         {BoolValue: &[]bool{true}[0]},
				},
				CreatedAt: time.Date(2024, 10, 29, 12, 22, 30, 0, time.UTC),
				UpdatedAt: time.Date(2024, 10, 29, 12, 22, 30, 0, time.UTC),
			},
			"porthole-003": {
				Uuid:          "porthole-003",
				Name:          "Emergency Porthole",
				Description:   "Quick-release emergency exit porthole",
				Price:         4750.00,
				StockQuantity: 10,
				Category:      model.PORTHOLE,
				Manufacturer: model.Manufacturer{
					Name:    "SafeExit",
					Country: "Norway",
					Website: "https://safeexit.no",
				},
				Tags: []string{"porthole", "emergency", "quick-release"},
				Dimensions: model.Dimensions{
					Length: 0.8,
					Width:  0.8,
					Height: 0.1,
					Weight: 40.0,
				},
				Metadata: map[string]model.Value{
					"material":       {StringValue: &[]string{"titanium_glass"}[0]},
					"pressure_limit": {Int64Value: &[]int64{300}[0]},
					"transparency":   {DoubleValue: &[]float64{0.92}[0]},
					"heated":         {BoolValue: &[]bool{true}[0]},
				},
				CreatedAt: time.Date(2024, 10, 26, 14, 40, 55, 0, time.UTC),
				UpdatedAt: time.Date(2024, 10, 26, 14, 40, 55, 0, time.UTC),
			},
		},
	}
}
