// Package main implements the zfssnap CLI.
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jsirianni/zfssnap/model"
	"github.com/jsirianni/zfssnap/zfs"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [snapshot...]",
	Short: "Get details for ZFS snapshots or list all snapshots",
	Long: `Get details for ZFS snapshots or list all snapshots.

If snapshot names are provided, returns detailed information for those snapshots.
If no snapshot names are provided, lists all snapshots.
If no arguments are provided and stdin is not a terminal, reads snapshot names from stdin (newline-separated).

Examples:
  # List all snapshots
  zfssnap get

  # Get details for specific snapshots
  zfssnap get pool@snapshot1 pool@snapshot2

  # Read snapshot names from stdin
  echo "pool@snapshot1" | zfssnap get`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		s := zfs.NewSnapshot(
			zfs.WithZFSPath(flagZFSPath),
			zfs.WithTimeout(flagTimeout),
		)

		var snapshotNames []string
		if len(args) > 0 {
			// Use provided arguments
			snapshotNames = args
		} else {
			// Check if stdin has data (not a terminal)
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				// stdin is not a terminal, read from stdin
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					name := strings.TrimSpace(scanner.Text())
					if name != "" {
						snapshotNames = append(snapshotNames, name)
					}
				}
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("read from stdin: %w", err)
				}
				if len(snapshotNames) == 0 {
					return fmt.Errorf("no snapshot names provided")
				}
			} else {
				// stdin is a terminal, list all snapshots
				names, err := s.List(ctx)
				if err != nil {
					return fmt.Errorf("list snapshots: %w", err)
				}

				// Use output functions for formatting
				if flagLogType == "json" {
					return outputStringArrayJSON(names, os.Stdout)
				}
				return outputStringArray(names, os.Stdout)
			}
		}

		// Get detailed information for specific snapshots
		var snapshots []*model.Snapshot
		for _, snapshotName := range snapshotNames {
			info, err := s.Get(ctx, snapshotName)
			if err != nil {
				return fmt.Errorf("get snapshot %s: %w", snapshotName, err)
			}
			snapshots = append(snapshots, info)
		}

		// Use output functions for formatting
		if flagLogType == "json" {
			if len(snapshots) == 1 {
				return outputSnapshotJSON(snapshots[0], os.Stdout)
			}
			return outputSnapshotJSONArray(snapshots, os.Stdout)
		}
		return outputSnapshotPlainArray(snapshots, os.Stdout)
	},
}
