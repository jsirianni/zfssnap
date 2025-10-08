// Package main implements the zfssnap CLI.
package main

import (
	"flag"

	"github.com/jsirianni/zfssnap/internal/version"
)

func versionSubcommand(args []string) error {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	appLogger.Info(version.Version())
	return nil
}
