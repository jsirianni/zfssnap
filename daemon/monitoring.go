package daemon

import (
	"context"
	"fmt"

	"github.com/jsirianni/zfssnap/zfs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var (
	meter              metric.Meter
	snapshotCountGauge metric.Int64ObservableGauge
)

// initMetrics initializes metrics and sets up callback registration.
func initMetrics(ctx context.Context, serviceName, serviceVersion string, snapshotter *zfs.Snapshot) (*prometheus.Exporter, error) {
	if serviceName == "" {
		return nil, fmt.Errorf("service name is required")
	}
	if serviceVersion == "" {
		return nil, fmt.Errorf("service version is required")
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	promExporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("create prometheus exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(promExporter),
	)

	otel.SetMeterProvider(meterProvider)

	meter = meterProvider.Meter("zfssnap-daemon")

	snapshotCountGauge, err = meter.Int64ObservableGauge(
		"snapshot_count_current",
		metric.WithDescription("Current number of ZFS snapshots"),
	)
	if err != nil {
		return nil, fmt.Errorf("create snapshot count gauge: %w", err)
	}

	_, err = meter.RegisterCallback(
		func(ctx context.Context, obs metric.Observer) error {
			snapshots, err := snapshotter.List(ctx)
			if err != nil {
				return fmt.Errorf("list snapshots: %w", err)
			}

			count := int64(len(snapshots))

			obs.ObserveInt64(snapshotCountGauge, count, metric.WithAttributes(
				attribute.String("source", "daemon"),
			))

			return nil
		},
		snapshotCountGauge,
	)
	if err != nil {
		return nil, fmt.Errorf("register snapshot count callback: %w", err)
	}

	return promExporter, nil
}
