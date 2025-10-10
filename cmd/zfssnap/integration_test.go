package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/jsirianni/zfssnap/model"
	"github.com/jsirianni/zfssnap/testutil"
	"github.com/jsirianni/zfssnap/zfs"
	"github.com/spf13/cobra"
)

func TestGetCommandListAll(t *testing.T) {
	testData := testutil.NewTestData()
	mockSnapshotter := testData.CreateMockSnapshotter()

	tests := []struct {
		name           string
		expectedOutput string
	}{
		{
			name: "json output",
			expectedOutput: `[{"name":"zroot/var/mail@test2","dataset":"zroot/var/mail","creation":"2025-08-07T00:22:49Z","used":65536,"referenced":114688,"defer_destroy":false,"logical_used":0,"logical_referenced":48128,"guid":16532700914722816504,"user_refs":0,"written":114688,"type":"snapshot"},{"name":"zroot/var/tmp@test","dataset":"zroot/var/tmp","creation":"2025-08-07T00:22:49Z","used":65536,"referenced":114688,"defer_destroy":false,"logical_used":0,"logical_referenced":48128,"guid":16532700914722816504,"user_refs":0,"written":114688,"type":"snapshot"}]
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			// Create a test command
			cmd := &cobra.Command{
				Use: "test-get-list",
			}

			// Create a test runner that uses our mock
			runGetListWithMock := func(_ *cobra.Command, _ []string) error {
				ctx := context.Background()
				names, err := mockSnapshotter.List(ctx)
				if err != nil {
					return err
				}

				// Get detailed information for all snapshots
				var snapshots []*model.Snapshot
				for _, snapshotName := range names {
					info, err := mockSnapshotter.Get(ctx, snapshotName)
					if err != nil {
						return err
					}
					snapshots = append(snapshots, info)
				}

				// Use output functions for formatting
				if len(snapshots) == 1 {
					return outputSnapshotJSON(snapshots[0], &buf)
				}
				return outputSnapshotJSONArray(snapshots, &buf)
			}

			err := runGetListWithMock(cmd, []string{})
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			output := buf.String()
			if output != tt.expectedOutput {
				t.Errorf("Expected output:\n%q\nGot:\n%q", tt.expectedOutput, output)
			}
		})
	}
}

func TestGetCommand(t *testing.T) {
	testData := testutil.NewTestData()
	mockSnapshotter := testData.CreateMockSnapshotter()

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name: "single snapshot",
			args: []string{"zroot/var/tmp@test"},
			expectedOutput: `{"name":"zroot/var/tmp@test","dataset":"zroot/var/tmp","creation":"2025-08-07T00:22:49Z","used":65536,"referenced":114688,"defer_destroy":false,"logical_used":0,"logical_referenced":48128,"guid":16532700914722816504,"user_refs":0,"written":114688,"type":"snapshot"}
`,
		},
		{
			name:        "invalid snapshot name",
			args:        []string{"invalid-name"},
			expectError: true,
		},
		{
			name:        "empty snapshot name",
			args:        []string{""},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			// Create a test command
			cmd := &cobra.Command{
				Use: "test-get",
			}

			// Create a test runner that uses our mock
			runGetWithMock := func(_ *cobra.Command, args []string) error {
				var snapshotNames []string
				if len(args) > 0 {
					snapshotNames = args
				}

				ctx := context.Background()
				var snapshots []*model.Snapshot
				for _, snapshotName := range snapshotNames {
					info, err := mockSnapshotter.Get(ctx, snapshotName)
					if err != nil {
						return err
					}
					snapshots = append(snapshots, info)
				}

				// Use output functions for formatting
				if len(snapshots) == 1 {
					return outputSnapshotJSON(snapshots[0], &buf)
				}
				return outputSnapshotJSONArray(snapshots, &buf)
			}

			err := runGetWithMock(cmd, tt.args)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			output := buf.String()
			if output != tt.expectedOutput {
				t.Errorf("Expected output:\n%q\nGot:\n%q", tt.expectedOutput, output)
			}
		})
	}
}

func TestGetCommandStdin(t *testing.T) {
	testData := testutil.NewTestData()
	mockSnapshotter := testData.CreateMockSnapshotter()

	tests := []struct {
		name           string
		stdinInput     string
		expectedOutput string
	}{
		{
			name:       "stdin multiple snapshots",
			stdinInput: "zroot/var/tmp@test\nzroot/var/mail@test2\n",
			expectedOutput: `[{"name":"zroot/var/tmp@test","dataset":"zroot/var/tmp","creation":"2025-08-07T00:22:49Z","used":65536,"referenced":114688,"defer_destroy":false,"logical_used":0,"logical_referenced":48128,"guid":16532700914722816504,"user_refs":0,"written":114688,"type":"snapshot"},{"name":"zroot/var/mail@test2","dataset":"zroot/var/mail","creation":"2025-08-07T00:22:49Z","used":65536,"referenced":114688,"defer_destroy":false,"logical_used":0,"logical_referenced":48128,"guid":16532700914722816504,"user_refs":0,"written":114688,"type":"snapshot"}]
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			// Create a test command
			cmd := &cobra.Command{
				Use: "test-get-stdin",
			}

			// Create a test runner that uses our mock and stdin
			runGetStdinWithMock := func(_ *cobra.Command, _ []string) error {
				var snapshotNames []string
				// Read from stdin
				lines := strings.Split(strings.TrimRight(tt.stdinInput, "\n"), "\n")
				for _, line := range lines {
					name := strings.TrimSpace(line)
					if name != "" {
						snapshotNames = append(snapshotNames, name)
					}
				}

				ctx := context.Background()
				var snapshots []*model.Snapshot
				for _, snapshotName := range snapshotNames {
					info, err := mockSnapshotter.Get(ctx, snapshotName)
					if err != nil {
						return err
					}
					snapshots = append(snapshots, info)
				}

				// Use output functions for formatting
				if len(snapshots) == 1 {
					return outputSnapshotJSON(snapshots[0], &buf)
				}
				return outputSnapshotJSONArray(snapshots, &buf)
			}

			err := runGetStdinWithMock(cmd, []string{})
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			output := buf.String()
			if output != tt.expectedOutput {
				t.Errorf("Expected output:\n%q\nGot:\n%q", tt.expectedOutput, output)
			}
		})
	}
}

func TestSnapshotValidation(t *testing.T) {
	tests := []struct {
		name         string
		snapshotName string
		expectValid  bool
	}{
		// Valid snapshot names
		{
			name:         "valid simple snapshot",
			snapshotName: "zroot/var/tmp@test",
			expectValid:  true,
		},
		{
			name:         "valid nested dataset snapshot",
			snapshotName: "pool/dataset1/dataset2@snapshot-2024",
			expectValid:  true,
		},
		{
			name:         "valid with underscores and hyphens",
			snapshotName: "pool_name/dataset-name@snapshot_name",
			expectValid:  true,
		},
		{
			name:         "valid with colons and periods",
			snapshotName: "pool:name/dataset.name@snapshot:name",
			expectValid:  true,
		},

		// Invalid snapshot names
		{
			name:         "invalid - no @",
			snapshotName: "zroot/var/tmp",
			expectValid:  false,
		},
		{
			name:         "invalid - empty name",
			snapshotName: "",
			expectValid:  false,
		},
		{
			name:         "invalid - whitespace only",
			snapshotName: "   ",
			expectValid:  false,
		},
		{
			name:         "invalid - starts with number",
			snapshotName: "123pool@snapshot",
			expectValid:  false,
		},
		{
			name:         "invalid - contains percent",
			snapshotName: "pool%name@snapshot",
			expectValid:  false,
		},
		{
			name:         "invalid - empty components",
			snapshotName: "pool//dataset@snapshot",
			expectValid:  false,
		},
		{
			name:         "invalid - leading slash",
			snapshotName: "/pool/dataset@snapshot",
			expectValid:  false,
		},
		{
			name:         "invalid - trailing slash",
			snapshotName: "pool/dataset/@snapshot",
			expectValid:  false,
		},
		{
			name:         "invalid - snapshot starts with number",
			snapshotName: "pool/dataset@123snapshot",
			expectValid:  false,
		},
		{
			name:         "invalid - too long",
			snapshotName: strings.Repeat("a", 250) + "@" + strings.Repeat("b", 10),
			expectValid:  false,
		},
		{
			name:         "invalid - special characters",
			snapshotName: "pool/dataset@snapshot!",
			expectValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := zfs.IsValidSnapshotName(tt.snapshotName)

			if tt.expectValid {
				if !isValid {
					t.Errorf("Expected valid but got invalid: %s", tt.snapshotName)
				}
			} else {
				if isValid {
					t.Errorf("Expected invalid but got valid: %s", tt.snapshotName)
				}
			}
		})
	}
}
