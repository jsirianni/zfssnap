// Package daemon provides daemon services with OpenTelemetry metrics.
package daemon

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jsirianni/zfssnap/zfs"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.uber.org/zap"
)

// Daemon represents a daemon service with OpenTelemetry metrics.
type Daemon struct {
	snapshot     *zfs.Snapshot
	promExporter *prometheus.Exporter
	httpServer   *http.Server
	logger       *zap.Logger
}

// New creates a new Daemon instance with OpenTelemetry metrics.
func New(ctx context.Context, serviceName, serviceVersion string, log *zap.Logger) (*Daemon, error) {
	snapshotter := zfs.NewSnapshot()

	promExporter, err := initMetrics(ctx, serviceName, serviceVersion, snapshotter)
	if err != nil {
		return nil, fmt.Errorf("initialize metrics: %w", err)
	}

	daemon := &Daemon{
		snapshot:     snapshotter,
		promExporter: promExporter,
		logger:       log,
	}

	return daemon, nil
}

// Start starts the HTTP server for metrics.
func (d *Daemon) Start(_ context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

		// For now, return a simple metrics response
		// TODO: Implement proper Prometheus metrics collection
		fmt.Fprintf(w, "# HELP zfs_snapshot_count Total number of ZFS snapshots\n")
		fmt.Fprintf(w, "# TYPE zfs_snapshot_count gauge\n")
		fmt.Fprintf(w, "zfs_snapshot_count{source=\"daemon\"} 0\n")
	})

	d.httpServer = &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 30 * time.Second,
	}
	d.logger.Info("HTTP server starting", zap.String("addr", addr+"/metrics"))

	go func() {
		if err := d.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			d.logger.Error("HTTP server error", zap.Error(err))
		}
	}()

	return nil
}

// Stop stops the HTTP server.
func (d *Daemon) Stop(ctx context.Context) error {
	if d.httpServer != nil {
		return d.httpServer.Shutdown(ctx)
	}
	return nil
}
