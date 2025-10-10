package daemon

import (
	"context"
	"fmt"

	"github.com/jsirianni/zfssnap/zfs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
func initMetrics(ctx context.Context, config *MetricsConfig, snapshotter *zfs.Snapshot) error {
	if config == nil {
		return fmt.Errorf("config is required")
	}

	if config.ServiceName == "" {
		return fmt.Errorf("service name is required")
	}
	if config.ServiceVersion == "" {
		return fmt.Errorf("service version is required")
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String(config.ServiceVersion),
		),
	)
	if err != nil {
		return fmt.Errorf("create resource: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(meterProvider)

	meter = meterProvider.Meter("zfssnap-daemon")

	snapshotCountGauge, err = meter.Int64ObservableGauge(
		"snapshot_count_current",
		metric.WithDescription("Current number of ZFS snapshots"),
	)
	if err != nil {
		return fmt.Errorf("create snapshot count gauge: %w", err)
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
		return fmt.Errorf("register snapshot count callback: %w", err)
	}

	return nil
}
