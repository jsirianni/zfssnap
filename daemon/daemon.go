// Package daemon provides daemon services with Prometheus metrics.
package daemon

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jsirianni/zfssnap/zfs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	// snapshotCountGauge tracks the total number of ZFS snapshots
	snapshotCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "zfs_snapshot_count",
		Help: "Total number of ZFS snapshots",
	})
)

func init() {
	// Register the Prometheus metrics
	prometheus.MustRegister(snapshotCountGauge)
}

// Daemon represents a daemon service with Prometheus metrics.
type Daemon struct {
	snapshot   *zfs.Snapshot
	httpServer *http.Server
	logger     *zap.Logger
}

// New creates a new Daemon instance with Prometheus metrics.
func New(_ context.Context, _, _ string, log *zap.Logger) (*Daemon, error) {
	snapshotter := zfs.NewSnapshot()

	daemon := &Daemon{
		snapshot: snapshotter,
		logger:   log,
	}

	return daemon, nil
}

// updateSnapshotCount updates the Prometheus gauge with the current snapshot count
func (d *Daemon) updateSnapshotCount() {
	ctx := context.Background()
	snapshots, err := d.snapshot.List(ctx)
	if err != nil {
		d.logger.Error("list snapshots", zap.Error(err))
		return
	}

	snapshotCountGauge.Set(float64(len(snapshots)))
}

// startMetricUpdates starts a goroutine that periodically updates metrics
func (d *Daemon) startMetricUpdates(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.updateSnapshotCount()
		}
	}
}

// Start starts the HTTP server for metrics.
func (d *Daemon) Start(ctx context.Context, addr string) error {
	// Update metrics before starting server
	d.updateSnapshotCount()

	// Start periodic metric updates
	go d.startMetricUpdates(ctx)

	mux := http.NewServeMux()

	// Use the proper Prometheus HTTP handler
	mux.Handle("/metrics", promhttp.Handler())

	// Add a health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
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
