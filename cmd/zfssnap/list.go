// Package main implements the zfssnap CLI.
package main

import (
	"context"

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
		for _, n := range names {
			appLogger.Info(n)
		}
		return nil
	},
}
