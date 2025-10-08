// Package version provides build-time version information.
package version

var (
	// semver is the semantic version, injected at build time.
	semver = "dev"

	// commitHash is the git commit hash, injected at build time.
	commitHash = "unknown"

	// buildTime is the ISO 8601 UTC timestamp, injected at build time.
	buildTime = "unknown"
)

// Semver returns the semantic version.
func Semver() string {
	return semver
}

// CommitHash returns the git commit hash.
func CommitHash() string {
	return commitHash
}

// BuildTime returns the build timestamp.
func BuildTime() string {
	return buildTime
}

// Version returns a formatted version string combining semver, commit hash, and build time.
func Version() string {
	return semver + " (" + commitHash + ", " + buildTime + ")"
}
