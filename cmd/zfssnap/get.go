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
	Short: "Get details for ZFS snapshots",
	Long:  "Get details for ZFS snapshots. If no snapshot names are provided, reads from stdin (newline-separated).",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		var snapshotNames []string
		if len(args) > 0 {
			snapshotNames = args
		} else {
			// Read from stdin
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
		}

		// Validate all snapshot names
		for _, name := range snapshotNames {
			if !zfs.IsValidSnapshotName(name) {
				return fmt.Errorf("invalid snapshot name format: %s (must contain @)", name)
			}
		}

		ctx := context.Background()
		s := zfs.NewSnapshot(
			zfs.WithZFSPath(flagZFSPath),
			zfs.WithTimeout(flagTimeout),
		)

		var snapshots []*model.Snapshot
		for _, snapshotName := range snapshotNames {
			info, err := s.Get(ctx, snapshotName)
			if err != nil {
				return fmt.Errorf("get snapshot %s: %w", snapshotName, err)
			}
			snapshots = append(snapshots, info)
		}

		// Use model methods for output formatting
		if flagLogType == "json" {
			if len(snapshots) == 1 {
				return snapshots[0].OutputJSON(os.Stdout)
			}
			return snapshots[0].OutputJSONArray(snapshots, os.Stdout)
		}
		return snapshots[0].OutputPlainArray(snapshots, os.Stdout)
	},
}
