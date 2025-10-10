package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jsirianni/zfssnap/daemon"
	"github.com/jsirianni/zfssnap/internal/logger"
	"github.com/jsirianni/zfssnap/internal/version"
	"github.com/spf13/cobra"
)

var (
	daemonAddr string
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start the ZFS snapshot daemon with metrics",
	Long: `Start the ZFS snapshot daemon with OpenTelemetry metrics.

The daemon provides Prometheus metrics at /metrics endpoint and automatically
collects ZFS snapshot counts using OpenTelemetry callbacks.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Force JSON logging for daemon
		jsonLogger, err := logger.NewJSONLogger()
		if err != nil {
			return fmt.Errorf("create JSON logger: %w", err)
		}
		defer jsonLogger.Sync()

		// Create daemon instance
		d, err := daemon.New(ctx, "zfssnap-daemon", version.Version(), jsonLogger)
		if err != nil {
			return fmt.Errorf("create daemon: %w", err)
		}

		// Start daemon
		if err := d.Start(ctx, daemonAddr); err != nil {
			return fmt.Errorf("start daemon: %w", err)
		}

		jsonLogger.Info("daemon started successfully", "addr", daemonAddr)

		// Wait for interrupt signal
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-sigChan:
			jsonLogger.Info("received signal, shutting down", "signal", sig.String())
		case <-ctx.Done():
			jsonLogger.Info("context cancelled, shutting down")
		}

		// Graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := d.Stop(shutdownCtx); err != nil {
			jsonLogger.Error("error stopping daemon", "error", err)
			return fmt.Errorf("stop daemon: %w", err)
		}

		jsonLogger.Info("daemon stopped successfully")
		return nil
	},
}

func init() {
	daemonCmd.Flags().StringVarP(&daemonAddr, "addr", "a", ":9464", "Address to bind the metrics server")
	rootCmd.AddCommand(daemonCmd)
}
