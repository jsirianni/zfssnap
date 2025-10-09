package daemon

import (
	"context"
	"fmt"
	"time"

	"github.com/jsirianni/zfssnap/zfs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

const (
	// DefaultPrometheusPort is the default port for the Prometheus metrics endpoint.
	DefaultPrometheusPort = ":9464"
)

// Daemon represents a daemon service with OpenTelemetry metrics.
type Daemon struct {
	snapshot           *zfs.Snapshot
	meter              metric.Meter
	snapshotCounter    metric.Int64Counter
	snapshotCountGauge metric.Int64Gauge
}

// MetricsConfig holds configuration for metrics exporters.
type MetricsConfig struct {
	PrometheusPort string
	ServiceName    string
	ServiceVersion string
}

// defaultMetricsConfig returns a default configuration for metrics.
func defaultMetricsConfig() *MetricsConfig {
	return &MetricsConfig{
		PrometheusPort: DefaultPrometheusPort,
	}
}

// initMetricsExporter initializes Prometheus exporter.
func initMetricsExporter(ctx context.Context, config *MetricsConfig) (metric.MeterProvider, error) {
	if config == nil {
		return nil, fmt.Errorf("config is required")
	}

	// Validate required fields
	if config.ServiceName == "" {
		return nil, fmt.Errorf("service name is required")
	}
	if config.ServiceVersion == "" {
		return nil, fmt.Errorf("service version is required")
	}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String(config.ServiceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create meter provider without exporter (for now)
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
	)

	// Set as global meter provider
	otel.SetMeterProvider(meterProvider)

	return meterProvider, nil
}

// New creates a new Daemon instance with OpenTelemetry metrics.
func New(serviceName, serviceVersion string) (*Daemon, error) {
	config := defaultMetricsConfig()
	config.ServiceName = serviceName
	config.ServiceVersion = serviceVersion
	return newWithConfig(context.Background(), config)
}

// newWithConfig creates a new Daemon instance with custom metrics configuration.
func newWithConfig(ctx context.Context, config *MetricsConfig) (*Daemon, error) {
	// Initialize metrics exporters
	meterProvider, err := initMetricsExporter(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	// Get meter
	meter := meterProvider.Meter("zfssnap-daemon")

	// Create snapshot counter metric
	snapshotCounter, err := meter.Int64Counter(
		"snapshot_count",
		metric.WithDescription("Total number of ZFS snapshots"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot counter: %w", err)
	}

	// Create snapshot count gauge metric
	snapshotCountGauge, err := meter.Int64Gauge(
		"snapshot_count_current",
		metric.WithDescription("Current number of ZFS snapshots"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot count gauge: %w", err)
	}

	// Create ZFS snapshotter
	snapshotter := zfs.NewSnapshot()

	return &Daemon{
		snapshot:           snapshotter,
		meter:              meter,
		snapshotCounter:    snapshotCounter,
		snapshotCountGauge: snapshotCountGauge,
	}, nil
}

// recordSnapshotCount records the current snapshot count metric.
func (d *Daemon) recordSnapshotCount(ctx context.Context) error {
	// Get snapshot list using the same function CLI get command uses
	snapshots, err := d.snapshot.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list snapshots: %w", err)
	}

	count := int64(len(snapshots))

	// Record the count metric (counter - cumulative)
	d.snapshotCounter.Add(ctx, count, metric.WithAttributes(
		attribute.String("source", "daemon"),
	))

	// Record the current count metric (gauge - current value)
	d.snapshotCountGauge.Record(ctx, count, metric.WithAttributes(
		attribute.String("source", "daemon"),
	))

	return nil
}

// StartMetricsCollection starts periodic collection of snapshot metrics.
func (d *Daemon) StartMetricsCollection(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Record initial count
	if err := d.recordSnapshotCount(ctx); err != nil {
		// Log error but continue
		fmt.Printf("Error recording initial snapshot count: %v\n", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := d.recordSnapshotCount(ctx); err != nil {
				// Log error but continue
				fmt.Printf("Error recording snapshot count: %v\n", err)
			}
		}
	}
}
