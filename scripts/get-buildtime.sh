#!/bin/sh
# Get current UTC timestamp in ISO 8601 format
date -u "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || printf "unknown"

