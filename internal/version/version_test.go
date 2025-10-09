package version

import (
	"strings"
	"testing"
)

func TestSemver(t *testing.T) {
	result := Semver()
	if result == "" {
		t.Error("Semver should not return empty string")
	}
	// Should return "dev" in test environment
	if result != "dev" {
		t.Errorf("Expected 'dev', got %q", result)
	}
}

func TestCommitHash(t *testing.T) {
	result := CommitHash()
	if result == "" {
		t.Error("CommitHash should not return empty string")
	}
	// Should return "unknown" in test environment
	if result != "unknown" {
		t.Errorf("Expected 'unknown', got %q", result)
	}
}

func TestBuildTime(t *testing.T) {
	result := BuildTime()
	if result == "" {
		t.Error("BuildTime should not return empty string")
	}
	// Should return "unknown" in test environment
	if result != "unknown" {
		t.Errorf("Expected 'unknown', got %q", result)
	}
}

func TestVersion(t *testing.T) {
	result := Version()
	if result == "" {
		t.Error("Version should not return empty string")
	}

	// Should contain semver, commit hash, and build time
	if !strings.Contains(result, "dev") {
		t.Errorf("Version should contain semver 'dev', got %q", result)
	}
	if !strings.Contains(result, "unknown") {
		t.Errorf("Version should contain 'unknown' for commit/build info, got %q", result)
	}

	// Should be in format: "semver (commitHash, buildTime)"
	expectedFormat := "dev (unknown, unknown)"
	if result != expectedFormat {
		t.Errorf("Expected %q, got %q", expectedFormat, result)
	}
}

func TestVersionFormat(t *testing.T) {
	// Test that Version() returns properly formatted string
	result := Version()

	// Should have parentheses and commas
	if !strings.Contains(result, "(") {
		t.Error("Version should contain opening parenthesis")
	}
	if !strings.Contains(result, ")") {
		t.Error("Version should contain closing parenthesis")
	}
	if !strings.Contains(result, ",") {
		t.Error("Version should contain comma separator")
	}

	// Should have exactly one comma
	commaCount := strings.Count(result, ",")
	if commaCount != 1 {
		t.Errorf("Expected exactly 1 comma, got %d", commaCount)
	}
}
