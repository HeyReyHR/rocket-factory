package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	serviceName = "assembly-service"
)

var meter = otel.Meter(serviceName)

var (
	AssembleDuration metric.Float64Histogram
)

func InitMetrics() error {
	var err error

	AssembleDuration, err = meter.Float64Histogram(
		"order_assemble_duration_seconds",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
	if err != nil {
		return err
	}

	return nil

}
