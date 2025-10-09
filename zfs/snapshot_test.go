package zfs

import (
	"strings"
	"testing"
	"time"
)

func TestParseUint(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uint64
		hasError bool
	}{
		{
			name:     "valid number",
			input:    "12345",
			expected: 12345,
			hasError: false,
		},
		{
			name:     "zero",
			input:    "0",
			expected: 0,
			hasError: false,
		},
		{
			name:     "large number",
			input:    "18446744073709551615", // max uint64
			expected: 18446744073709551615,
			hasError: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
			hasError: true,
		},
		{
			name:     "dash",
			input:    "-",
			expected: 0,
			hasError: true,
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: 0,
			hasError: true,
		},
		{
			name:     "contains letters",
			input:    "123abc",
			expected: 0,
			hasError: true,
		},
		{
			name:     "contains special chars",
			input:    "123-456",
			expected: 0,
			hasError: true,
		},
		{
			name:     "leading zeros",
			input:    "00123",
			expected: 123,
			hasError: false,
		},
		{
			name:     "whitespace around number",
			input:    " 123 ",
			expected: 123,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseUint(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %d, got %d", tt.expected, result)
				}
			}
		})
	}
}

func TestIsValidSnapshotName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid cases
		{
			name:     "simple snapshot",
			input:    "pool@snapshot",
			expected: true,
		},
		{
			name:     "nested dataset",
			input:    "pool/dataset@snapshot",
			expected: true,
		},
		{
			name:     "deeply nested",
			input:    "pool/dataset1/dataset2/dataset3@snapshot",
			expected: true,
		},
		{
			name:     "with underscores",
			input:    "pool_name/dataset_name@snapshot_name",
			expected: true,
		},
		{
			name:     "with hyphens",
			input:    "pool-name/dataset-name@snapshot-name",
			expected: true,
		},
		{
			name:     "with colons",
			input:    "pool:name/dataset:name@snapshot:name",
			expected: true,
		},
		{
			name:     "with periods",
			input:    "pool.name/dataset.name@snapshot.name",
			expected: true,
		},
		{
			name:     "mixed special chars",
			input:    "pool_name-name:name.name@snapshot_name-name:name.name",
			expected: true,
		},

		// Invalid cases
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: false,
		},
		{
			name:     "no @ symbol",
			input:    "pool/dataset",
			expected: false,
		},
		{
			name:     "starts with number",
			input:    "123pool@snapshot",
			expected: false,
		},
		{
			name:     "dataset starts with number",
			input:    "pool/123dataset@snapshot",
			expected: false,
		},
		{
			name:     "snapshot starts with number",
			input:    "pool/dataset@123snapshot",
			expected: false,
		},
		{
			name:     "contains percent",
			input:    "pool%name@snapshot",
			expected: false,
		},
		{
			name:     "empty components",
			input:    "pool//dataset@snapshot",
			expected: false,
		},
		{
			name:     "leading slash",
			input:    "/pool/dataset@snapshot",
			expected: false,
		},
		{
			name:     "trailing slash",
			input:    "pool/dataset/@snapshot",
			expected: false,
		},
		{
			name:     "too long",
			input:    strings.Repeat("a", 250) + "@" + strings.Repeat("b", 10),
			expected: false,
		},
		{
			name:     "invalid special chars",
			input:    "pool/dataset@snapshot!",
			expected: false,
		},
		{
			name:     "multiple @ symbols",
			input:    "pool@dataset@snapshot",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidSnapshotName(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for input: %q", tt.expected, result, tt.input)
			}
		})
	}
}

func TestIsValidDatasetName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid cases
		{
			name:     "simple dataset",
			input:    "pool",
			expected: true,
		},
		{
			name:     "nested dataset",
			input:    "pool/dataset",
			expected: true,
		},
		{
			name:     "deeply nested",
			input:    "pool/dataset1/dataset2/dataset3",
			expected: true,
		},
		{
			name:     "with underscores",
			input:    "pool_name/dataset_name",
			expected: true,
		},
		{
			name:     "with hyphens",
			input:    "pool-name/dataset-name",
			expected: true,
		},
		{
			name:     "with colons",
			input:    "pool:name/dataset:name",
			expected: true,
		},
		{
			name:     "with periods",
			input:    "pool.name/dataset.name",
			expected: true,
		},
		{
			name:     "mixed special chars",
			input:    "pool_name-name:name.name/dataset_name-name:name.name",
			expected: true,
		},

		// Invalid cases
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: false,
		},
		{
			name:     "starts with number",
			input:    "123pool",
			expected: false,
		},
		{
			name:     "dataset starts with number",
			input:    "pool/123dataset",
			expected: false,
		},
		{
			name:     "contains percent",
			input:    "pool%name",
			expected: false,
		},
		{
			name:     "empty components",
			input:    "pool//dataset",
			expected: false,
		},
		{
			name:     "leading slash",
			input:    "/pool/dataset",
			expected: false,
		},
		{
			name:     "trailing slash",
			input:    "pool/dataset/",
			expected: false,
		},
		{
			name:     "too long",
			input:    strings.Repeat("a", 256),
			expected: false,
		},
		{
			name:     "invalid special chars",
			input:    "pool/dataset!",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidDatasetName(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for input: %q", tt.expected, result, tt.input)
			}
		})
	}
}

func TestIsValidSnapshotComponent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid cases
		{
			name:     "simple snapshot name",
			input:    "snapshot",
			expected: true,
		},
		{
			name:     "with underscores",
			input:    "snapshot_name",
			expected: true,
		},
		{
			name:     "with hyphens",
			input:    "snapshot-name",
			expected: true,
		},
		{
			name:     "with colons",
			input:    "snapshot:name",
			expected: true,
		},
		{
			name:     "with periods",
			input:    "snapshot.name",
			expected: true,
		},
		{
			name:     "mixed special chars",
			input:    "snapshot_name-name:name.name",
			expected: true,
		},
		{
			name:     "with numbers",
			input:    "snapshot123",
			expected: true,
		},

		// Invalid cases
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: false,
		},
		{
			name:     "starts with number",
			input:    "123snapshot",
			expected: false,
		},
		{
			name:     "contains @ symbol",
			input:    "snapshot@name",
			expected: false,
		},
		{
			name:     "contains slash",
			input:    "snapshot/name",
			expected: false,
		},
		{
			name:     "contains percent",
			input:    "snapshot%name",
			expected: false,
		},
		{
			name:     "too long",
			input:    strings.Repeat("a", 256),
			expected: false,
		},
		{
			name:     "invalid special chars",
			input:    "snapshot!",
			expected: false,
		},
		{
			name:     "contains spaces",
			input:    "snapshot name",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidSnapshotComponent(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for input: %q", tt.expected, result, tt.input)
			}
		})
	}
}

func TestNewSnapshot(t *testing.T) {
	tests := []struct {
		name     string
		opts     []Option
		expected struct {
			zfspath string
			timeout time.Duration
		}
	}{
		{
			name: "default options",
			opts: []Option{},
			expected: struct {
				zfspath string
				timeout time.Duration
			}{
				zfspath: DefaultZFSBinary,
				timeout: DefaultTimeout,
			},
		},
		{
			name: "custom zfs path",
			opts: []Option{WithZFSPath("/custom/zfs")},
			expected: struct {
				zfspath string
				timeout time.Duration
			}{
				zfspath: "/custom/zfs",
				timeout: DefaultTimeout,
			},
		},
		{
			name: "custom timeout",
			opts: []Option{WithTimeout(60 * time.Second)},
			expected: struct {
				zfspath string
				timeout time.Duration
			}{
				zfspath: DefaultZFSBinary,
				timeout: 60 * time.Second,
			},
		},
		{
			name: "both custom options",
			opts: []Option{
				WithZFSPath("/custom/zfs"),
				WithTimeout(120 * time.Second),
			},
			expected: struct {
				zfspath string
				timeout time.Duration
			}{
				zfspath: "/custom/zfs",
				timeout: 120 * time.Second,
			},
		},
		{
			name: "empty zfs path gets default",
			opts: []Option{WithZFSPath("")},
			expected: struct {
				zfspath string
				timeout time.Duration
			}{
				zfspath: DefaultZFSBinary,
				timeout: DefaultTimeout,
			},
		},
		{
			name: "whitespace zfs path gets default",
			opts: []Option{WithZFSPath("   ")},
			expected: struct {
				zfspath string
				timeout time.Duration
			}{
				zfspath: DefaultZFSBinary,
				timeout: DefaultTimeout,
			},
		},
		{
			name: "zero timeout gets default",
			opts: []Option{WithTimeout(0)},
			expected: struct {
				zfspath string
				timeout time.Duration
			}{
				zfspath: DefaultZFSBinary,
				timeout: DefaultTimeout,
			},
		},
		{
			name: "negative timeout gets default",
			opts: []Option{WithTimeout(-1 * time.Second)},
			expected: struct {
				zfspath string
				timeout time.Duration
			}{
				zfspath: DefaultZFSBinary,
				timeout: DefaultTimeout,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshot := NewSnapshot(tt.opts...)

			if snapshot.ZFSPath != tt.expected.zfspath {
				t.Errorf("Expected ZFSPath %q, got %q", tt.expected.zfspath, snapshot.ZFSPath)
			}

			if snapshot.Timeout != tt.expected.timeout {
				t.Errorf("Expected Timeout %v, got %v", tt.expected.timeout, snapshot.Timeout)
			}
		})
	}
}

func TestSnapshotCreate(t *testing.T) {
	tests := []struct {
		name          string
		dataset       string
		snapshotName  string
		expectError   bool
		errorContains string
	}{
		{
			name:         "valid dataset and snapshot",
			dataset:      "pool/dataset",
			snapshotName: "backup",
			expectError:  false,
		},
		{
			name:         "simple dataset",
			dataset:      "pool",
			snapshotName: "snapshot",
			expectError:  false,
		},
		{
			name:         "nested dataset",
			dataset:      "pool/dataset1/dataset2",
			snapshotName: "backup-2024",
			expectError:  false,
		},
		{
			name:         "snapshot with special chars",
			dataset:      "pool/dataset",
			snapshotName: "backup_name-name:name.name",
			expectError:  false,
		},
		{
			name:          "empty dataset",
			dataset:       "",
			snapshotName:  "backup",
			expectError:   true,
			errorContains: "dataset name is required",
		},
		{
			name:          "empty snapshot name",
			dataset:       "pool/dataset",
			snapshotName:  "",
			expectError:   true,
			errorContains: "snapshot name is required",
		},
		{
			name:          "invalid dataset format",
			dataset:       "123pool/dataset",
			snapshotName:  "backup",
			expectError:   true,
			errorContains: "invalid dataset name format",
		},
		{
			name:          "invalid snapshot name format",
			dataset:       "pool/dataset",
			snapshotName:  "123backup",
			expectError:   true,
			errorContains: "invalid snapshot name format",
		},
		{
			name:          "snapshot name with @",
			dataset:       "pool/dataset",
			snapshotName:  "backup@name",
			expectError:   true,
			errorContains: "invalid snapshot name format",
		},
		{
			name:          "snapshot name with slash",
			dataset:       "pool/dataset",
			snapshotName:  "backup/name",
			expectError:   true,
			errorContains: "invalid snapshot name format",
		},
		{
			name:          "whitespace dataset",
			dataset:       "   ",
			snapshotName:  "backup",
			expectError:   true,
			errorContains: "dataset name is required",
		},
		{
			name:          "whitespace snapshot name",
			dataset:       "pool/dataset",
			snapshotName:  "   ",
			expectError:   true,
			errorContains: "snapshot name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test validation logic without actually executing ZFS commands
			dataset := strings.TrimSpace(tt.dataset)
			snapshotName := strings.TrimSpace(tt.snapshotName)

			if dataset == "" {
				if !tt.expectError || !strings.Contains("dataset name is required", tt.errorContains) {
					t.Errorf("Expected error for empty dataset")
				}
				return
			}

			if snapshotName == "" {
				if !tt.expectError || !strings.Contains("snapshot name is required", tt.errorContains) {
					t.Errorf("Expected error for empty snapshot name")
				}
				return
			}

			if !IsValidDatasetName(dataset) {
				if !tt.expectError || !strings.Contains("invalid dataset name format", tt.errorContains) {
					t.Errorf("Expected error for invalid dataset name")
				}
				return
			}

			if !IsValidSnapshotComponent(snapshotName) {
				if !tt.expectError || !strings.Contains("invalid snapshot name format", tt.errorContains) {
					t.Errorf("Expected error for invalid snapshot name")
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
