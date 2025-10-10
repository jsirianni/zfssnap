// Package testutil provides testing utilities and mock implementations.
package testutil

import (
	"context"
	"time"

	"github.com/jsirianni/zfssnap/model"
)

// MockSnapshotter is a mock implementation of Snapshotter for testing.
type MockSnapshotter struct {
	ListFunc   func(ctx context.Context) ([]string, error)
	GetFunc    func(ctx context.Context, name string) (*model.Snapshot, error)
	CreateFunc func(ctx context.Context, name, dataset string) error
	DeleteFunc func(ctx context.Context, name string) error
}

// List implements Snapshotter.List.
func (m *MockSnapshotter) List(ctx context.Context) ([]string, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return []string{}, nil
}

// Get implements Snapshotter.Get.
func (m *MockSnapshotter) Get(ctx context.Context, name string) (*model.Snapshot, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, name)
	}
	return nil, nil
}

// Create implements Snapshotter.Create.
func (m *MockSnapshotter) Create(ctx context.Context, name, dataset string) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, name, dataset)
	}
	return nil
}

// Delete implements Snapshotter.Delete.
func (m *MockSnapshotter) Delete(ctx context.Context, name string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, name)
	}
	return nil
}

// NewMockSnapshotter creates a new MockSnapshotter with default implementations.
func NewMockSnapshotter() *MockSnapshotter {
	return &MockSnapshotter{}
}

// WithListFunc sets the List function for the mock.
func (m *MockSnapshotter) WithListFunc(fn func(ctx context.Context) ([]string, error)) *MockSnapshotter {
	m.ListFunc = fn
	return m
}

// WithGetFunc sets the Get function for the mock.
func (m *MockSnapshotter) WithGetFunc(fn func(ctx context.Context, name string) (*model.Snapshot, error)) *MockSnapshotter {
	m.GetFunc = fn
	return m
}

// WithCreateFunc sets the Create function for the mock.
func (m *MockSnapshotter) WithCreateFunc(fn func(ctx context.Context, name, dataset string) error) *MockSnapshotter {
	m.CreateFunc = fn
	return m
}

// WithDeleteFunc sets the Delete function for the mock.
func (m *MockSnapshotter) WithDeleteFunc(fn func(ctx context.Context, name string) error) *MockSnapshotter {
	m.DeleteFunc = fn
	return m
}

// TestData contains real ZFS command outputs for testing.
type TestData struct {
	ListOutput []string
	GetOutput  map[string]*model.Snapshot
}

// NewTestData creates test data based on real ZFS outputs.
func NewTestData() *TestData {
	return &TestData{
		ListOutput: []string{
			"zroot/var/mail@test2",
			"zroot/var/tmp@test",
		},
		GetOutput: map[string]*model.Snapshot{
			"zroot/var/tmp@test": {
				Name:              "zroot/var/tmp@test",
				Dataset:           "zroot/var/tmp",
				Creation:          time.Date(2025, 8, 7, 0, 22, 49, 0, time.UTC),
				Used:              65536,
				Referenced:        114688,
				Clones:            nil,
				DeferDestroy:      false,
				LogicalUsed:       0, // - in output
				LogicalReferenced: 48128,
				GUID:              16532700914722816504,
				UserRefs:          0,
				Written:           114688,
				Type:              "snapshot",
			},
			"zroot/var/mail@test2": {
				Name:              "zroot/var/mail@test2",
				Dataset:           "zroot/var/mail",
				Creation:          time.Date(2025, 8, 7, 0, 22, 49, 0, time.UTC),
				Used:              65536,
				Referenced:        114688,
				Clones:            nil,
				DeferDestroy:      false,
				LogicalUsed:       0, // - in output
				LogicalReferenced: 48128,
				GUID:              16532700914722816504,
				UserRefs:          0,
				Written:           114688,
				Type:              "snapshot",
			},
		},
	}
}

// CreateMockSnapshotter creates a MockSnapshotter with test data.
func (td *TestData) CreateMockSnapshotter() *MockSnapshotter {
	return NewMockSnapshotter().
		WithListFunc(func(_ context.Context) ([]string, error) {
			return td.ListOutput, nil
		}).
		WithGetFunc(func(_ context.Context, name string) (*model.Snapshot, error) {
			if snapshot, exists := td.GetOutput[name]; exists {
				return snapshot, nil
			}
			return nil, &MockError{Message: "snapshot not found: " + name}
		})
}

// MockError represents a mock error for testing.
type MockError struct {
	Message string
}

func (e *MockError) Error() string {
	return e.Message
}
