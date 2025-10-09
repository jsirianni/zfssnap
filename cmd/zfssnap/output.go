package main

import (
	"encoding/json"
	"fmt"
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

// outputSnapshotPlain writes the snapshot name as plain text to the provided writer.
func outputSnapshotPlain(s *model.Snapshot, w io.Writer) error {
	_, err := fmt.Fprintln(w, s.Name)
	return err
}

// outputSnapshotJSONArray writes an array of snapshots as JSON to the provided writer.
// Single snapshot outputs as object, multiple snapshots as array.
func outputSnapshotJSONArray(snapshots []*model.Snapshot, w io.Writer) error {
	if len(snapshots) == 1 {
		return encodeSnapshotJSON(snapshots[0], w)
	}
	return encodeSnapshotJSONArray(snapshots, w)
}

// outputSnapshotPlainArray writes snapshot names as plain text (one per line) to the provided writer.
func outputSnapshotPlainArray(snapshots []*model.Snapshot, w io.Writer) error {
	for _, snapshot := range snapshots {
		if err := outputSnapshotPlain(snapshot, w); err != nil {
			return err
		}
	}
	return nil
}

// outputStringArray writes strings as plain text (one per line) to the provided writer.
func outputStringArray(strings []string, w io.Writer) error {
	for _, str := range strings {
		if _, err := fmt.Fprintln(w, str); err != nil {
			return err
		}
	}
	return nil
}

// encodeStringArrayJSON writes strings as JSON array to the provided writer.
func encodeStringArrayJSON(strings []string, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(strings)
}

// outputStringArrayJSON writes strings as JSON array to the provided writer.
func outputStringArrayJSON(strings []string, w io.Writer) error {
	return encodeStringArrayJSON(strings, w)
}
