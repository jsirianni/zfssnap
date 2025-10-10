#!/bin/sh
# Generate ldflags for build
set -e

SCRIPT_DIR="$(dirname "$0")"
VERSION="$(sh "${SCRIPT_DIR}/get-version.sh")"
COMMIT="$(sh "${SCRIPT_DIR}/get-commit.sh")"
BUILDTIME="$(sh "${SCRIPT_DIR}/get-buildtime.sh")"

printf -- "-X github.com/jsirianni/zfssnap/internal/version.semver=%s -X github.com/jsirianni/zfssnap/internal/version.commitHash=%s -X github.com/jsirianni/zfssnap/internal/version.buildTime=%s" \
    "${VERSION}" "${COMMIT}" "${BUILDTIME}"

