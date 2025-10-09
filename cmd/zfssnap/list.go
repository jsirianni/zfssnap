// Package main implements the zfssnap CLI.
package main

import (
	"context"
	"os"

	"github.com/jsirianni/zfssnap/model"
	"github.com/jsirianni/zfssnap/zfs"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List ZFS snapshots",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		s := zfs.NewSnapshot(
			zfs.WithZFSPath(flagZFSPath),
			zfs.WithTimeout(flagTimeout),
		)
		names, err := s.List(ctx)
		if err != nil {
			return err
		}

		// Use model methods for output formatting
		if flagLogType == "json" {
			return model.OutputStringArrayJSON(names, os.Stdout)
		}
		return model.OutputStringArray(names, os.Stdout)
	},
}
