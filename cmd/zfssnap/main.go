package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jsirianni/zfssnap/internal/logger"
)

var (
	appLogger logger.Logger

	// Global configuration populated from env and flags.
	flagZFSPath string
	flagTimeout time.Duration
	flagLogType string
)

func initDefaultsFromEnv() {
	if v := strings.TrimSpace(os.Getenv("ZFSSNAP_ZFS_PATH")); v != "" {
		flagZFSPath = v
	}
	if v := strings.TrimSpace(os.Getenv("ZFSSNAP_TIMEOUT")); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			flagTimeout = d
		}
	}
	if v := strings.TrimSpace(os.Getenv("ZFSSNAP_LOG_TYPE")); v != "" {
		flagLogType = strings.ToLower(v)
	}
}

func bindGlobalFlags() {
	if flagTimeout <= 0 {
		flagTimeout = 30 * time.Second
	}
	if flagLogType == "" {
		flagLogType = "plain"
	}

	flag.StringVar(&flagZFSPath, "zfs-path", flagZFSPath, "Path to zfs binary (default: detected)")
	flag.DurationVar(&flagTimeout, "timeout", flagTimeout, "Command timeout")
	flag.StringVar(&flagLogType, "log-type", flagLogType, "Logger type: plain or json")
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

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [global flags] <command> [command flags]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nCommands:\n")
	fmt.Fprintf(os.Stderr, "  list\tList ZFS snapshots\n")
	fmt.Fprintf(os.Stderr, "\nGlobal Flags:\n")
	flag.PrintDefaults()
}

func main() {
	initDefaultsFromEnv()
	bindGlobalFlags()
	flag.Usage = usage
	flag.Parse()

	if err := initLogger(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(2)
	}

	subcmd := args[0]
	subArgs := args[1:]

	switch subcmd {
	case "list":
		if err := listSubcommand(subArgs); err != nil {
			appLogger.Error(err.Error())
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", subcmd)
		usage()
		os.Exit(2)
	}
}
