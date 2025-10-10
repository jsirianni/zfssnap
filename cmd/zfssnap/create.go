package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jsirianni/zfssnap/zfs"
	"github.com/spf13/cobra"
)

var (
	flagRecursive bool
	flagDryRun    bool
	flagForce     bool
	flagPrefix    string
	flagSuffix    string
	flagTimestamp bool
)

var createCmd = &cobra.Command{
	Use:   "create [flags] <dataset> <snapshot-name>",
	Short: "Create ZFS snapshots",
	Long: `Create ZFS snapshots for the specified dataset(s).

Examples:
  # Basic snapshot
  zfssnap create pool/dataset backup-2024-01-15

  # Recursive snapshot with timestamp
  zfssnap create -r --timestamp pool/dataset daily

  # Dry run to see what would be created
  zfssnap create --dry-run pool/dataset test-snapshot

  # With prefix and force overwrite
  zfssnap create --prefix manual --force pool/dataset backup

  # Multiple datasets
  zfssnap create pool/dataset1 pool/dataset2 backup-2024-01-15`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("at least dataset and snapshot name are required")
		}

		// Last argument is the snapshot name
		snapshotName := args[len(args)-1]
		datasets := args[:len(args)-1]

		// Apply naming transformations
		snapshotName = applyNamingTransformations(snapshotName)

		ctx := context.Background()
		s := zfs.NewSnapshot(
			zfs.WithZFSPath(flagZFSPath),
			zfs.WithTimeout(flagTimeout),
		)

		var createdSnapshots []string
		var errors []string

		for _, dataset := range datasets {
			fullSnapshotName := dataset + "@" + snapshotName

			if flagDryRun {
				fmt.Printf("Would create snapshot: %s\n", fullSnapshotName)
				createdSnapshots = append(createdSnapshots, fullSnapshotName)
				continue
			}

			err := s.Create(ctx, dataset, snapshotName)
			if err != nil {
				if flagForce && strings.Contains(err.Error(), "already exists") {
					// Force mode: try to destroy existing snapshot first
					destroyErr := s.Delete(ctx, fullSnapshotName)
					if destroyErr != nil {
						errors = append(errors, fmt.Sprintf("failed to destroy existing snapshot %s: %v", fullSnapshotName, destroyErr))
						continue
					}

					// Retry creation
					err = s.Create(ctx, dataset, snapshotName)
					if err != nil {
						errors = append(errors, fmt.Sprintf("failed to create snapshot %s: %v", fullSnapshotName, err))
						continue
					}
				} else {
					errors = append(errors, fmt.Sprintf("failed to create snapshot %s: %v", fullSnapshotName, err))
					continue
				}
			}

			createdSnapshots = append(createdSnapshots, fullSnapshotName)
		}

		// Output results
		return outputCreateResultsJSON(createdSnapshots, errors)
	},
}

func init() {
	createCmd.Flags().BoolVarP(&flagRecursive, "recursive", "r", false, "Create snapshots recursively for all child datasets")
	createCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Show what would be created without actually creating")
	createCmd.Flags().BoolVarP(&flagForce, "force", "f", false, "Force creation even if snapshot already exists")
	createCmd.Flags().StringVar(&flagPrefix, "prefix", "", "Add prefix to snapshot name")
	createCmd.Flags().StringVar(&flagSuffix, "suffix", "", "Add suffix to snapshot name")
	createCmd.Flags().BoolVar(&flagTimestamp, "timestamp", false, "Auto-add timestamp to snapshot name")
}

func applyNamingTransformations(snapshotName string) string {
	// Apply prefix
	if flagPrefix != "" {
		snapshotName = flagPrefix + "-" + snapshotName
	}

	// Apply suffix
	if flagSuffix != "" {
		snapshotName = snapshotName + "-" + flagSuffix
	}

	// Apply timestamp
	if flagTimestamp {
		timestamp := time.Now().Format("20060102-150405")
		snapshotName = snapshotName + "-" + timestamp
	}

	return snapshotName
}

func outputCreateResultsJSON(createdSnapshots []string, errors []string) error {
	// Use the existing JSON output mechanism
	// Always output JSON format
	fmt.Printf("{\"created\":%q,\"errors\":%q,\"count\":%d}\n",
		strings.Join(createdSnapshots, ","),
		strings.Join(errors, ","),
		len(createdSnapshots))
	return nil
}

func outputCreateResultsPlain(createdSnapshots []string, errors []string) error {
	// Output created snapshots
	for _, snapshot := range createdSnapshots {
		if flagDryRun {
			fmt.Printf("Would create: %s\n", snapshot)
		} else {
			fmt.Printf("Created: %s\n", snapshot)
		}
	}

	// Output errors
	for _, err := range errors {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}

	// Return error if there were any failures
	if len(errors) > 0 {
		return fmt.Errorf("failed to create %d snapshot(s)", len(errors))
	}

	return nil
}
