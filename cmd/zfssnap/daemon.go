package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jsirianni/zfssnap/daemon"
	"github.com/jsirianni/zfssnap/internal/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "ts"
		config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
		zapLogger, err := config.Build()
		if err != nil {
			return fmt.Errorf("create zap logger: %w", err)
		}
		defer zapLogger.Sync()

		// Create daemon instance
		d, err := daemon.New(ctx, "zfssnap-daemon", version.Version(), zapLogger)
		if err != nil {
			return fmt.Errorf("create daemon: %w", err)
		}

		// Start daemon
		if err := d.Start(ctx, daemonAddr); err != nil {
			return fmt.Errorf("start daemon: %w", err)
		}

		zapLogger.Info("daemon started successfully", zap.String("addr", daemonAddr))

		// Wait for interrupt signal
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-sigChan:
			zapLogger.Info("received signal, shutting down", zap.String("signal", sig.String()))
		case <-ctx.Done():
			zapLogger.Info("context cancelled, shutting down")
		}

		// Graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := d.Stop(shutdownCtx); err != nil {
			zapLogger.Error("error stopping daemon", zap.Error(err))
			return fmt.Errorf("stop daemon: %w", err)
		}

		zapLogger.Info("daemon stopped successfully")
		return nil
	},
}

func init() {
	daemonCmd.Flags().StringVarP(&daemonAddr, "addr", "a", "localhost:9464", "Address to bind the metrics server")
	rootCmd.AddCommand(daemonCmd)
}
