package model

import (
	"time"
)

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	Category      Category
	StockQuantity int64
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]interface{}
	Dimensions    Dimensions
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Category string

const (
	UNKNOWN  Category = "unknown"
	ENGINE   Category = "engine"
	FUEL     Category = "fuel"
	PORTHOLE Category = "porthole"
	WING     Category = "wing"
)

type Filter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
