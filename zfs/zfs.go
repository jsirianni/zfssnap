package zfs

import (
	"context"

	"github.com/jsirianni/zfssnap/model"
)

// Snapshotter defines the contract for managing ZFS snapshots.
type Snapshotter interface {
	// List returns a list of ZFS snapshot names.
	List(ctx context.Context) ([]string, error)

	// Create creates a ZFS snapshot with the given name for the specified dataset.
	Create(ctx context.Context, name, dataset string) error

	// Delete removes the ZFS snapshot with the given name.
	Delete(ctx context.Context, name string) error

	// Get returns detailed information for the specified snapshot name.
	Get(ctx context.Context, name string) (*model.Snapshot, error)
}
