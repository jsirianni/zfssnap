#!/bin/sh
# Get git commit hash
git rev-parse --short HEAD 2>/dev/null || printf "unknown"

