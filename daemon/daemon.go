// Package daemon provides daemon services with OpenTelemetry metrics.
package daemon

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jsirianni/zfssnap/internal/logger"
	"github.com/jsirianni/zfssnap/zfs"
	"go.opentelemetry.io/otel/exporters/prometheus"
)

// Daemon represents a daemon service with OpenTelemetry metrics.
type Daemon struct {
	snapshot     *zfs.Snapshot
	promExporter *prometheus.Exporter
	httpServer   *http.Server
	logger       logger.Logger
}

// New creates a new Daemon instance with OpenTelemetry metrics.
func New(ctx context.Context, serviceName, serviceVersion string, log logger.Logger) (*Daemon, error) {
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
func (d *Daemon) Start(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		err := d.promExporter.Collect(r.Context(), nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error collecting metrics: %v", err), http.StatusInternalServerError)
			return
		}
	})

	d.httpServer = &http.Server{Addr: addr, Handler: mux}
	d.logger.Info("HTTP server starting", "addr", addr+"/metrics")

	go func() {
		if err := d.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			d.logger.Error("HTTP server error", "error", err)
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
