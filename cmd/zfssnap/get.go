// Package main implements the zfssnap CLI.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jsirianni/zfssnap/zfs"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [snapshot]",
	Short: "Get details for a ZFS snapshot",
	Long:  "Get details for a ZFS snapshot. If no snapshot name is provided, reads from stdin.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var snapshotName string
		if len(args) > 0 {
			snapshotName = args[0]
		} else {
			// Read from stdin
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				snapshotName = strings.TrimSpace(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("read from stdin: %w", err)
			}
			if snapshotName == "" {
				return fmt.Errorf("no snapshot name provided")
			}
		}

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
