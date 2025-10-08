// Package main implements the zfssnap CLI.
package main

import (
	"context"
	"flag"

	"github.com/jsirianni/zfssnap/zfs"
)

func listSubcommand(args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	ctx := context.Background()
	s := zfs.NewSnapshot(
		zfs.WithZFSPath(flagZFSPath),
		zfs.WithTimeout(flagTimeout),
	)
	names, err := s.List(ctx)
	if err != nil {
		return err
	}
	for _, n := range names {
		appLogger.Info(n)
	}
	return nil
}
