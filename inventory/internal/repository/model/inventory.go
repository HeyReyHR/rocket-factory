package model

import (
	"time"
)

type Part struct {
	Uuid          string           `bson:"uuid,omitempty"`
	Name          string           `bson:"name"`
	Description   string           `bson:"description,omitempty"`
	Price         float64          `bson:"price"`
	Category      Category         `bson:"category"`
	StockQuantity int64            `bson:"stock_quantity"`
	Manufacturer  Manufacturer     `bson:"manufacturer"`
	Tags          []string         `bson:"tags"`
	Metadata      map[string]Value `bson:"metadata,omitempty"`
	Dimensions    Dimensions       `bson:"dimensions"`
	CreatedAt     time.Time        `bson:"created_at"`
	UpdatedAt     time.Time        `bson:"updated_at,omitempty"`
}
type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type Value struct {
	StringValue *string  `bson:"string_value,omitempty"`
	Int64Value  *int64   `bson:"int64_value,omitempty"`
	DoubleValue *float64 `bson:"double_value,omitempty"`
	BoolValue   *bool    `bson:"bool_value,omitempty"`
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
