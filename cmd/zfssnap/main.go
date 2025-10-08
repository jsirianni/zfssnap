package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jsirianni/zfssnap/internal/logger"
	"github.com/spf13/cobra"
)

var (
	appLogger logger.Logger

	flagZFSPath string
	flagTimeout time.Duration
	flagLogType string
)

var rootCmd = &cobra.Command{
	Use:   "zfssnap",
	Short: "ZFS snapshot utility",
}

func initLogger() error {
	switch strings.ToLower(strings.TrimSpace(flagLogType)) {
	case "", "plain":
		appLogger = logger.PlainLogger{}
	case "json":
		jl, err := logger.NewJSONLogger()
		if err != nil {
			return err
		}
		appLogger = jl
	default:
		return fmt.Errorf("invalid log type: %s", flagLogType)
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagZFSPath, "zfs-bin", "", "Path to zfs binary (default: detect in $PATH)")
	rootCmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", 30*time.Second, "Command timeout")
	rootCmd.PersistentFlags().StringVar(&flagLogType, "output", "plain", "Output format: plain or json")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.ParseFlags(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := initLogger(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		appLogger.Error(err.Error())
		os.Exit(1)
	}
}
