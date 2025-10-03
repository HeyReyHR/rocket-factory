package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	serviceName = "order-service"
)

var meter = otel.Meter(serviceName)

var (
	OrdersTotal        metric.Int64Counter
	OrdersRevenueTotal metric.Float64Counter
)

func InitMetrics() error {
	var err error

	OrdersTotal, err = meter.Int64Counter(
		"orders_total",
		metric.WithDescription("Total orders"))
	if err != nil {
		return err
	}

	OrdersRevenueTotal, err = meter.Float64Counter(
		"orders_revenue",
		metric.WithUnit("money"),
		metric.WithDescription("Total order revenue"))
	if err != nil {
		return err
	}

	return nil

}
