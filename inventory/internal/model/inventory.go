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
	Metadata      map[string]Value
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

type Value struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}

type Category int

const (
	UNKNOWN Category = iota
	ENGINE
	FUEL
	PORTHOLE
	WING
)

type Filter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
