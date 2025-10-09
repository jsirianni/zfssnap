// Package main implements the zfssnap CLI.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jsirianni/zfssnap/zfs"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <snapshot>",
	Short: "Get details for a ZFS snapshot",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		snapshotName := args[0]
		ctx := context.Background()
		s := zfs.NewSnapshot(
			zfs.WithZFSPath(flagZFSPath),
			zfs.WithTimeout(flagTimeout),
		)
		info, err := s.Get(ctx, snapshotName)
		if err != nil {
			return err
		}

		// Respect output mode set at root; plain prints the name for compatibility,
		// json prints the full struct as single-line JSON.
		switch flagLogType {
		case "json":
			enc := json.NewEncoder(os.Stdout)
			enc.SetEscapeHTML(false)
			if err := enc.Encode(info); err != nil {
				return fmt.Errorf("encode json: %w", err)
			}
		default:
			appLogger.Info(info.Name)
		}
		return nil
	},
}
