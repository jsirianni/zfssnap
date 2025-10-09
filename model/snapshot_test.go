package model

import (
	"bytes"
	"testing"
	"time"
)

func TestSnapshotOutputPlain(t *testing.T) {
	tests := []struct {
		name     string
		snapshot *Snapshot
		expected string
	}{
		{
			name: "simple snapshot",
			snapshot: &Snapshot{
				Name: "pool@snapshot",
			},
			expected: "pool@snapshot\n",
		},
		{
			name: "nested dataset snapshot",
			snapshot: &Snapshot{
				Name: "pool/dataset@snapshot",
			},
			expected: "pool/dataset@snapshot\n",
		},
		{
			name: "empty name",
			snapshot: &Snapshot{
				Name: "",
			},
			expected: "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.snapshot.OutputPlain(&buf)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if buf.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestSnapshotOutputPlainArray(t *testing.T) {
	tests := []struct {
		name      string
		snapshots []*Snapshot
		expected  string
	}{
		{
			name: "single snapshot",
			snapshots: []*Snapshot{
				{Name: "pool@snapshot1"},
			},
			expected: "pool@snapshot1\n",
		},
		{
			name: "multiple snapshots",
			snapshots: []*Snapshot{
				{Name: "pool@snapshot1"},
				{Name: "pool@snapshot2"},
				{Name: "pool@snapshot3"},
			},
			expected: "pool@snapshot1\npool@snapshot2\npool@snapshot3\n",
		},
		{
			name:      "empty array",
			snapshots: []*Snapshot{},
			expected:  "",
		},
		{
			name:      "nil array",
			snapshots: nil,
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			snapshot := &Snapshot{}
			err := snapshot.OutputPlainArray(tt.snapshots, &buf)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if buf.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestOutputStringArray(t *testing.T) {
	tests := []struct {
		name     string
		strings  []string
		expected string
	}{
		{
			name:     "single string",
			strings:  []string{"hello"},
			expected: "hello\n",
		},
		{
			name:     "multiple strings",
			strings:  []string{"hello", "world", "test"},
			expected: "hello\nworld\ntest\n",
		},
		{
			name:     "empty array",
			strings:  []string{},
			expected: "",
		},
		{
			name:     "nil array",
			strings:  nil,
			expected: "",
		},
		{
			name:     "empty strings",
			strings:  []string{"", "", ""},
			expected: "\n\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := OutputStringArray(tt.strings, &buf)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if buf.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestSnapshotOutputJSONArray(t *testing.T) {
	tests := []struct {
		name      string
		snapshots []*Snapshot
		expected  string
	}{
		{
			name: "single snapshot - should output as object",
			snapshots: []*Snapshot{
				{
					Name:     "pool@snapshot",
					Dataset:  "pool",
					Creation: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					Used:     1024,
					Type:     "snapshot",
				},
			},
			expected: `{"name":"pool@snapshot","dataset":"pool","creation":"2025-01-01T00:00:00Z","used":1024,"referenced":0,"defer_destroy":false,"logical_used":0,"logical_referenced":0,"guid":0,"user_refs":0,"written":0,"type":"snapshot"}
`,
		},
		{
			name: "multiple snapshots - should output as array",
			snapshots: []*Snapshot{
				{
					Name:    "pool@snapshot1",
					Dataset: "pool",
					Used:    1024,
					Type:    "snapshot",
				},
				{
					Name:    "pool@snapshot2",
					Dataset: "pool",
					Used:    2048,
					Type:    "snapshot",
				},
			},
			expected: `[{"name":"pool@snapshot1","dataset":"pool","creation":"0001-01-01T00:00:00Z","used":1024,"referenced":0,"defer_destroy":false,"logical_used":0,"logical_referenced":0,"guid":0,"user_refs":0,"written":0,"type":"snapshot"},{"name":"pool@snapshot2","dataset":"pool","creation":"0001-01-01T00:00:00Z","used":2048,"referenced":0,"defer_destroy":false,"logical_used":0,"logical_referenced":0,"guid":0,"user_refs":0,"written":0,"type":"snapshot"}]
`,
		},
		{
			name:      "empty array",
			snapshots: []*Snapshot{},
			expected: `[]
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			snapshot := &Snapshot{}
			err := snapshot.OutputJSONArray(tt.snapshots, &buf)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if buf.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestOutputStringArrayJSON(t *testing.T) {
	tests := []struct {
		name     string
		strings  []string
		expected string
	}{
		{
			name:    "single string",
			strings: []string{"hello"},
			expected: `["hello"]
`,
		},
		{
			name:    "multiple strings",
			strings: []string{"hello", "world", "test"},
			expected: `["hello","world","test"]
`,
		},
		{
			name:    "empty array",
			strings: []string{},
			expected: `[]
`,
		},
		{
			name:    "nil array",
			strings: nil,
			expected: `null
`,
		},
		{
			name:    "empty strings",
			strings: []string{"", "", ""},
			expected: `["","",""]
`,
		},
		{
			name:    "special characters",
			strings: []string{"hello\"world", "test\nline"},
			expected: `["hello\"world","test\nline"]
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := OutputStringArrayJSON(tt.strings, &buf)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if buf.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestSnapshotEncodeJSON(t *testing.T) {
	tests := []struct {
		name     string
		snapshot *Snapshot
		expected string
	}{
		{
			name: "complete snapshot",
			snapshot: &Snapshot{
				Name:              "pool@snapshot",
				Dataset:           "pool",
				Creation:          time.Date(2025, 1, 1, 12, 30, 45, 0, time.UTC),
				Used:              1024,
				Referenced:        2048,
				Clones:            []string{"clone1", "clone2"},
				DeferDestroy:      true,
				LogicalUsed:       512,
				LogicalReferenced: 1536,
				GUID:              123456789,
				UserRefs:          2,
				Written:           256,
				Type:              "snapshot",
			},
			expected: `{"name":"pool@snapshot","dataset":"pool","creation":"2025-01-01T12:30:45Z","used":1024,"referenced":2048,"clones":["clone1","clone2"],"defer_destroy":true,"logical_used":512,"logical_referenced":1536,"guid":123456789,"user_refs":2,"written":256,"type":"snapshot"}
`,
		},
		{
			name: "minimal snapshot",
			snapshot: &Snapshot{
				Name: "pool@snapshot",
			},
			expected: `{"name":"pool@snapshot","dataset":"","creation":"0001-01-01T00:00:00Z","used":0,"referenced":0,"defer_destroy":false,"logical_used":0,"logical_referenced":0,"guid":0,"user_refs":0,"written":0,"type":""}
`,
		},
		{
			name: "nil clones",
			snapshot: &Snapshot{
				Name:   "pool@snapshot",
				Clones: nil,
			},
			expected: `{"name":"pool@snapshot","dataset":"","creation":"0001-01-01T00:00:00Z","used":0,"referenced":0,"defer_destroy":false,"logical_used":0,"logical_referenced":0,"guid":0,"user_refs":0,"written":0,"type":""}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.snapshot.EncodeJSON(&buf)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if buf.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}
