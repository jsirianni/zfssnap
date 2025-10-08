// Package main implements the zfssnap CLI.
package main

import (
	"flag"

	"github.com/jsirianni/zfssnap/internal/version"
)

func versionSubcommand(args []string) error {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	_ = fs.Parse(args)

	appLogger.Info(version.Version())
	return nil
}
