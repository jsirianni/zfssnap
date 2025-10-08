#!/bin/sh
# Build a single binary for specified OS/ARCH
set -e

SCRIPT_DIR="$(dirname "$0")"
. "${SCRIPT_DIR}/build.env"

if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Usage: $0 <GOOS> <GOARCH>" >&2
    exit 1
fi

GOOS="$1"
GOARCH="$2"

mkdir -p "${BUILD_DIR}"

LDFLAGS="$(sh "${SCRIPT_DIR}/build-ldflags.sh")"

GOOS="${GOOS}" GOARCH="${GOARCH}" go build -ldflags "${LDFLAGS}" -o "${BUILD_DIR}/${BINARY_NAME}-${GOOS}-${GOARCH}" "${CMD_DIR}"

