// Package model defines data structures for ZFS snapshots.
package model

import (
	"time"
)

// Snapshot represents a ZFS snapshot and its associated metadata.
// Fields are based on OpenZFS properties commonly exposed by `zfs get`.
type Snapshot struct {
	// Fully qualified snapshot name: pool/dataset@snap
	Name string `json:"name"`

	// Parent dataset name without the @ snapshot component
	Dataset string `json:"dataset"`

	// Creation time of the snapshot
	Creation time.Time `json:"creation"`

	// Space that would be freed if the snapshot were destroyed (bytes)
	Used uint64 `json:"used"`

	// Space accessible by this snapshot, including shared data (bytes)
	Referenced uint64 `json:"referenced"`

	// Datasets that are clones of this snapshot
	Clones []string `json:"clones,omitempty"`

	// Whether the snapshot is marked for deferred destroy
	DeferDestroy bool `json:"defer_destroy"`

	// Logical space consumed by this snapshot, ignoring compression/dedup (bytes)
	LogicalUsed uint64 `json:"logical_used"`

	// Logical space accessible by this snapshot (bytes)
	LogicalReferenced uint64 `json:"logical_referenced"`

	// Globally unique identifier for this snapshot
	GUID uint64 `json:"guid"`

	// Number of user holds on this snapshot
	UserRefs uint64 `json:"user_refs"`

	// Amount of space written to this snapshot since the previous snapshot (bytes)
	Written uint64 `json:"written"`

	// Dataset type; for snapshots this is typically "snapshot"
	Type string `json:"type"`
}
