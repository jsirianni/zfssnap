package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	appLogger *zap.Logger

	flagZFSPath string
	flagTimeout time.Duration
)

var rootCmd = &cobra.Command{
	Use:   "zfssnap",
	Short: "ZFS snapshot utility",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		if err := initLogger(); err != nil {
			fmt.Fprintf(os.Stderr, "initialize logger: %v\n", err)
			os.Exit(1)
		}
	},
}

func initLogger() error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "ts"
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	var err error
	appLogger, err = config.Build()
	return err
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagZFSPath, "zfs-bin", "", "Path to zfs binary (default: detect in $PATH)")
	rootCmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", 30*time.Second, "Command timeout")

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
