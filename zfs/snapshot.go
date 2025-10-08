// Package zfs provides ZFS snapshot management functionality.
package zfs

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// DefaultZFSBinary is the default path to the zfs binary.
const DefaultZFSBinary = "zfs"

// DefaultTimeout is the default timeout for ZFS operations.
const DefaultTimeout = 30 * time.Second

// SnapshotInfo represents a ZFS snapshot and its associated metadata.
// Fields are based on OpenZFS properties commonly exposed by `zfs get`.
type SnapshotInfo struct {
	// Fully qualified snapshot name: pool/dataset@snap
	Name string

	// Parent dataset name without the @ snapshot component
	Dataset string

	// Creation time of the snapshot
	Creation time.Time

	// Space that would be freed if the snapshot were destroyed (bytes)
	Used uint64

	// Space accessible by this snapshot, including shared data (bytes)
	Referenced uint64

	// Datasets that are clones of this snapshot
	Clones []string

	// Whether the snapshot is marked for deferred destroy
	DeferDestroy bool

	// Logical space consumed by this snapshot, ignoring compression/dedup (bytes)
	LogicalUsed uint64

	// Logical space accessible by this snapshot (bytes)
	LogicalReferenced uint64

	// Globally unique identifier for this snapshot
	GUID uint64

	// Number of user holds on this snapshot
	UserRefs uint64

	// Amount of space written to this snapshot since the previous snapshot (bytes)
	Written uint64

	// Dataset type; for snapshots this is typically "snapshot"
	Type string
}

// Option configures the Snapshot implementation.
type Option func(*Snapshot)

// WithZFSPath sets the path to the zfs binary. If not provided, DefaultZFSBinary is used.
func WithZFSPath(path string) Option { return func(s *Snapshot) { s.ZFSPath = path } }

// WithTimeout sets the default timeout for CLI calls.
func WithTimeout(d time.Duration) Option { return func(s *Snapshot) { s.Timeout = d } }

// NewSnapshot creates a new CLI-backed snapshotter with options.
// Defaults: ZFSPath=DefaultZFSBinary, Timeout=DefaultTimeout.
func NewSnapshot(opts ...Option) *Snapshot {
	s := &Snapshot{
		ZFSPath: DefaultZFSBinary,
		Timeout: DefaultTimeout,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}

	// Validate and normalize final values
	s.ZFSPath = strings.TrimSpace(s.ZFSPath)
	if s.ZFSPath == "" {
		s.ZFSPath = DefaultZFSBinary
	}
	if s.Timeout <= 0 {
		s.Timeout = DefaultTimeout
	}
	return s
}

// Snapshot is a concrete implementation of Snapshotter that
// will use the `zfs` command line interface under the hood.
// Methods are currently stubbed and will be implemented to shell out
// to the ZFS CLI in future edits.
type Snapshot struct {
	// Path to the zfs binary, e.g. "/sbin/zfs". If empty, "zfs" on PATH is used.
	ZFSPath string

	// Optional default timeout for CLI calls.
	Timeout time.Duration
}

// Compile-time check that Snapshot implements Snapshotter.
var _ Snapshotter = (*Snapshot)(nil)

// execContext is a small seam for testing and future customization.
func (c *Snapshot) execContext(ctx context.Context, name string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, args...)
}

// List returns the names of ZFS snapshots using the `zfs` CLI.
func (c *Snapshot) List(ctx context.Context) ([]string, error) {
	args := []string{"list", "-H", "-t", "snapshot", "-o", "name"}

	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	cmd := c.execContext(ctx, c.ZFSPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("zfs list failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	out := stdout.String()
	if out == "" {
		return []string{}, nil
	}

	lines := strings.Split(out, "\n")
	snapshots := make([]string, 0, len(lines))
	for _, line := range lines {
		name := strings.TrimSpace(line)
		if name == "" {
			continue
		}
		snapshots = append(snapshots, name)
	}
	return snapshots, nil
}

// Create creates a snapshot. Stub: returns nil for now.
func (c *Snapshot) Create(_ context.Context, _, _ string) error {
	return nil
}

// Delete destroys a snapshot.
func (c *Snapshot) Delete(_ context.Context, _ string) error {
	return nil
}
