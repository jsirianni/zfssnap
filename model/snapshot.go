package model

import (
	"encoding/json"
	"fmt"
	"io"
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

// encodeJSON writes the snapshot as JSON to the provided writer.
func (s *Snapshot) encodeJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(s)
}

// encodeJSONArray writes an array of snapshots as JSON to the provided writer.
func encodeJSONArray(snapshots []*Snapshot, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(snapshots)
}

// EncodeJSON writes the snapshot as JSON to the provided writer.
func (s *Snapshot) EncodeJSON(w io.Writer) error {
	return s.encodeJSON(w)
}

// EncodeJSONArray writes an array of snapshots as JSON to the provided writer.
func EncodeJSONArray(snapshots []*Snapshot, w io.Writer) error {
	return encodeJSONArray(snapshots, w)
}

// OutputJSON writes the snapshot as JSON to the provided writer.
func (s *Snapshot) OutputJSON(w io.Writer) error {
	return s.EncodeJSON(w)
}

// OutputPlain writes the snapshot name as plain text to the provided writer.
func (s *Snapshot) OutputPlain(w io.Writer) error {
	_, err := fmt.Fprintln(w, s.Name)
	return err
}

// OutputJSONArray writes an array of snapshots as JSON to the provided writer.
// Single snapshot outputs as object, multiple snapshots as array.
func (s *Snapshot) OutputJSONArray(snapshots []*Snapshot, w io.Writer) error {
	if len(snapshots) == 1 {
		return snapshots[0].encodeJSON(w)
	}
	return encodeJSONArray(snapshots, w)
}

// OutputPlainArray writes snapshot names as plain text (one per line) to the provided writer.
func (s *Snapshot) OutputPlainArray(snapshots []*Snapshot, w io.Writer) error {
	for _, snapshot := range snapshots {
		if err := snapshot.OutputPlain(w); err != nil {
			return err
		}
	}
	return nil
}

// OutputStringArray writes strings as plain text (one per line) to the provided writer.
func OutputStringArray(strings []string, w io.Writer) error {
	for _, str := range strings {
		if _, err := fmt.Fprintln(w, str); err != nil {
			return err
		}
	}
	return nil
}

// outputStringArrayJSON writes strings as JSON array to the provided writer.
func outputStringArrayJSON(strings []string, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(strings)
}

// OutputStringArrayJSON writes strings as JSON array to the provided writer.
func OutputStringArrayJSON(strings []string, w io.Writer) error {
	return outputStringArrayJSON(strings, w)
}
