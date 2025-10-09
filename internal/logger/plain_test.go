package logger

import (
	"testing"
)

func TestPlainLoggerMarshalJSON(t *testing.T) {
	logger := PlainLogger{}
	result, err := logger.MarshalJSON()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := `"plain-logger"`
	if string(result) != expected {
		t.Errorf("Expected %q, got %q", expected, string(result))
	}
}

func TestPlainLoggerInterfaceCompliance(t *testing.T) {
	// Test that PlainLogger implements Logger interface
	var _ Logger = PlainLogger{}

	// Test that we can call all methods without panicking
	logger := PlainLogger{}
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")
	logger.Debug("test")

	// Test MarshalJSON
	_, err := logger.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON should not error: %v", err)
	}
}

func TestPlainLoggerMethods(t *testing.T) {
	logger := PlainLogger{}

	// Test that methods don't panic with various inputs
	tests := []struct {
		name string
		fn   func()
	}{
		{
			name: "Info with no args",
			fn:   func() { logger.Info() },
		},
		{
			name: "Info with string",
			fn:   func() { logger.Info("hello") },
		},
		{
			name: "Info with multiple args",
			fn:   func() { logger.Info("hello", "world", 123) },
		},
		{
			name: "Warn with no args",
			fn:   func() { logger.Warn() },
		},
		{
			name: "Warn with string",
			fn:   func() { logger.Warn("warning") },
		},
		{
			name: "Error with no args",
			fn:   func() { logger.Error() },
		},
		{
			name: "Error with string",
			fn:   func() { logger.Error("error") },
		},
		{
			name: "Debug with no args",
			fn:   func() { logger.Debug() },
		},
		{
			name: "Debug with string",
			fn:   func() { logger.Debug("debug") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the method doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Method panicked: %v", r)
				}
			}()
			tt.fn()
		})
	}
}
