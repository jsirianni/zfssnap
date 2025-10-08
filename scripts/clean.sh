#!/bin/sh
# Clean build artifacts
SCRIPT_DIR="$(dirname "$0")"
. "${SCRIPT_DIR}/build.env"

rm -rf "${BUILD_DIR}"

