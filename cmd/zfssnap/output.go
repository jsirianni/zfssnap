package main

import (
	"encoding/json"
	"io"

	"github.com/jsirianni/zfssnap/model"
)

// encodeJSON writes the snapshot as JSON to the provided writer.
func encodeSnapshotJSON(s *model.Snapshot, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(s)
}

// encodeSnapshotJSONArray writes an array of snapshots as JSON to the provided writer.
func encodeSnapshotJSONArray(snapshots []*model.Snapshot, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(snapshots)
}

// outputSnapshotJSON writes the snapshot as JSON to the provided writer.
func outputSnapshotJSON(s *model.Snapshot, w io.Writer) error {
	return encodeSnapshotJSON(s, w)
}

// outputSnapshotJSONArray writes an array of snapshots as JSON to the provided writer.
// Single snapshot outputs as object, multiple snapshots as array.
func outputSnapshotJSONArray(snapshots []*model.Snapshot, w io.Writer) error {
	if len(snapshots) == 1 {
		return encodeSnapshotJSON(snapshots[0], w)
	}
	return encodeSnapshotJSONArray(snapshots, w)
}
