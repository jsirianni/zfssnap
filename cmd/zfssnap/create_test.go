package main

import (
	"strings"
	"testing"

	"github.com/jsirianni/zfssnap/zfs"
)

func TestCreateCommandValidation(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "valid dataset and snapshot",
			args:        []string{"pool/dataset", "backup"},
			expectError: false,
		},
		{
			name:        "simple dataset",
			args:        []string{"pool", "snapshot"},
			expectError: false,
		},
		{
			name:        "multiple datasets",
			args:        []string{"pool/dataset1", "pool/dataset2", "backup"},
			expectError: false,
		},
		{
			name:        "invalid dataset format",
			args:        []string{"123pool/dataset", "backup"},
			expectError: true,
		},
		{
			name:        "invalid snapshot name format",
			args:        []string{"pool/dataset", "123backup"},
			expectError: true,
		},
		{
			name:        "missing snapshot name",
			args:        []string{"pool/dataset"},
			expectError: true,
		},
		{
			name:        "missing arguments",
			args:        []string{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test argument validation logic
			if len(tt.args) < 2 {
				if !tt.expectError {
					t.Errorf("Expected no error for insufficient args")
				}
				return
			}

			// Last argument is the snapshot name
			snapshotName := tt.args[len(tt.args)-1]
			datasets := tt.args[:len(tt.args)-1]

			// Test dataset validation
			for _, dataset := range datasets {
				if !zfs.IsValidDatasetName(dataset) {
					if !tt.expectError {
						t.Errorf("Expected no error for invalid dataset: %s", dataset)
					}
					return
				}
			}

			// Test snapshot name validation
			if !zfs.IsValidSnapshotComponent(snapshotName) {
				if !tt.expectError {
					t.Errorf("Expected no error for invalid snapshot name: %s", snapshotName)
				}
				return
			}

			// If we get here, validation passed
			if tt.expectError {
				t.Errorf("Expected error but validation passed")
			}
		})
	}
}

func TestApplyNamingTransformations(t *testing.T) {
	tests := []struct {
		name           string
		snapshotName   string
		prefix         string
		suffix         string
		timestamp      bool
		expectedPrefix string
	}{
		{
			name:           "no transformations",
			snapshotName:   "backup",
			prefix:         "",
			suffix:         "",
			timestamp:      false,
			expectedPrefix: "backup",
		},
		{
			name:           "with prefix",
			snapshotName:   "backup",
			prefix:         "daily",
			suffix:         "",
			timestamp:      false,
			expectedPrefix: "daily-backup",
		},
		{
			name:           "with suffix",
			snapshotName:   "backup",
			prefix:         "",
			suffix:         "manual",
			timestamp:      false,
			expectedPrefix: "backup-manual",
		},
		{
			name:           "with prefix and suffix",
			snapshotName:   "backup",
			prefix:         "daily",
			suffix:         "manual",
			timestamp:      false,
			expectedPrefix: "daily-backup-manual",
		},
		{
			name:           "with timestamp",
			snapshotName:   "backup",
			prefix:         "",
			suffix:         "",
			timestamp:      true,
			expectedPrefix: "backup-",
		},
		{
			name:           "all transformations",
			snapshotName:   "backup",
			prefix:         "daily",
			suffix:         "manual",
			timestamp:      true,
			expectedPrefix: "daily-backup-manual-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set flags
			flagPrefix = tt.prefix
			flagSuffix = tt.suffix
			flagTimestamp = tt.timestamp

			result := applyNamingTransformations(tt.snapshotName)

			// Check that result starts with expected prefix
			if !strings.HasPrefix(result, tt.expectedPrefix) {
				t.Errorf("Expected result to start with %q, got %q", tt.expectedPrefix, result)
			}

			// If timestamp is enabled, check that it ends with timestamp format
			if tt.timestamp {
				// Should end with timestamp format (YYYYMMDD-HHMMSS)
				if len(result) < len(tt.expectedPrefix)+15 {
					t.Errorf("Expected timestamp to be added, got %q", result)
				}
			}

			// Reset flags
			flagPrefix = ""
			flagSuffix = ""
			flagTimestamp = false
		})
	}
}
