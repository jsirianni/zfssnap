// Package daemon provides daemon services with OpenTelemetry metrics.
package daemon

import (
	"context"
	"fmt"

	"github.com/jsirianni/zfssnap/zfs"
)

const (
	// DefaultPrometheusPort is the default port for the Prometheus metrics endpoint.
	DefaultPrometheusPort = ":9464"
)

// Daemon represents a daemon service with OpenTelemetry metrics.
type Daemon struct {
	snapshot *zfs.Snapshot
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

// New creates a new Daemon instance with OpenTelemetry metrics.
func New(serviceName, serviceVersion string) (*Daemon, error) {
	config := defaultMetricsConfig()
	config.ServiceName = serviceName
	config.ServiceVersion = serviceVersion
	return newWithConfig(context.Background(), config)
}

// newWithConfig creates a new Daemon instance with custom metrics configuration.
func newWithConfig(ctx context.Context, config *MetricsConfig) (*Daemon, error) {
	snapshotter := zfs.NewSnapshot()

	if err := initMetrics(ctx, config, snapshotter); err != nil {
		return nil, fmt.Errorf("initialize metrics: %w", err)
	}

	daemon := &Daemon{
		snapshot: snapshotter,
	}

	return daemon, nil
}
