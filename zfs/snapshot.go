// Package zfs provides ZFS snapshot management functionality.
package zfs

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/jsirianni/zfssnap/model"
)

// DefaultZFSBinary is the default path to the zfs binary.
const DefaultZFSBinary = "zfs"

// DefaultTimeout is the default timeout for ZFS operations.
const DefaultTimeout = 30 * time.Second

// SnapshotInfo moved to model.Snapshot

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

// Get returns detailed information for a given snapshot using `zfs get`.
func (c *Snapshot) Get(ctx context.Context, name string) (*model.Snapshot, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("snapshot name is required")
	}

	// Query properties in a single call; -H for scriptable, -p for parsable numbers
	props := []string{
		"name", "creation", "used", "referenced", "clones", "defer_destroy",
		"logicalused", "logicalreferenced", "guid", "userrefs", "written", "type",
	}
	args := []string{"get", "-H", "-p", "-o", "property,value", strings.Join(props, ","), name}

	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	cmd := c.execContext(ctx, c.ZFSPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("zfs get failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	out := strings.TrimSpace(stdout.String())
	if out == "" {
		return nil, fmt.Errorf("snapshot not found: %s", name)
	}

	info := &model.Snapshot{}
	// Parse lines of the form: <name>\t<property>\t<value>\t- (with -o property,value we expect: <property>\t<value>)
	// But because we used -o property,value and provided dataset, output is lines: <name>\t<property>\t<value>\t<source>
	// We'll split on newlines and then fields by tabs, taking property and value from positions 1 and 2.
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 3 {
			continue
		}
		// fields[0]=dataset@snap, fields[1]=property, fields[2]=value
		prop := fields[1]
		val := fields[2]
		switch prop {
		case "name":
			info.Name = val
			if at := strings.Index(val, "@"); at > 0 {
				info.Dataset = val[:at]
			}
		case "creation":
			// creation is seconds since epoch with -p
			if v, err := parseUint(val); err == nil {
				info.Creation = time.Unix(int64(v), 0).UTC()
			}
		case "used":
			if v, err := parseUint(val); err == nil {
				info.Used = v
			}
		case "referenced":
			if v, err := parseUint(val); err == nil {
				info.Referenced = v
			}
		case "clones":
			if val == "-" || val == "" {
				info.Clones = nil
			} else {
				info.Clones = strings.Split(val, ",")
			}
		case "defer_destroy":
			info.DeferDestroy = val == "on" || val == "yes" || val == "1"
		case "logicalused":
			if v, err := parseUint(val); err == nil {
				info.LogicalUsed = v
			}
		case "logicalreferenced":
			if v, err := parseUint(val); err == nil {
				info.LogicalReferenced = v
			}
		case "guid":
			if v, err := parseUint(val); err == nil {
				info.GUID = v
			}
		case "userrefs":
			if v, err := parseUint(val); err == nil {
				info.UserRefs = v
			}
		case "written":
			if v, err := parseUint(val); err == nil {
				info.Written = v
			}
		case "type":
			info.Type = val
		}
	}

	// If name was not emitted, fall back to provided
	if info.Name == "" {
		info.Name = name
		if at := strings.Index(name, "@"); at > 0 {
			info.Dataset = name[:at]
		}
	}
	return info, nil
}

// parseUint parses a positive integer, returning 0 on error.
func parseUint(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return 0, fmt.Errorf("empty")
	}
	var v uint64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid")
		}
		v = v*10 + uint64(c-'0')
	}
	return v, nil
}
