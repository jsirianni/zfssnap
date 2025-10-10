package model

import (
	"testing"
	"time"
)

func TestSnapshotStruct(t *testing.T) {
	// Test that the Snapshot struct can be created and accessed
	snapshot := &Snapshot{
		Name:              "pool@snapshot",
		Dataset:           "pool",
		Creation:          time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
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
	}

	// Test field access
	if snapshot.Name != "pool@snapshot" {
		t.Errorf("Expected Name 'pool@snapshot', got %q", snapshot.Name)
	}
	if snapshot.Dataset != "pool" {
		t.Errorf("Expected Dataset 'pool', got %q", snapshot.Dataset)
	}
	if snapshot.Used != 1024 {
		t.Errorf("Expected Used 1024, got %d", snapshot.Used)
	}
	if len(snapshot.Clones) != 2 {
		t.Errorf("Expected 2 clones, got %d", len(snapshot.Clones))
	}
	if !snapshot.DeferDestroy {
		t.Error("Expected DeferDestroy to be true")
	}
}
