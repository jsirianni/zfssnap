# zfssnap

A ZFS snapshot utility with CLI commands and Prometheus metrics daemon for managing and monitoring ZFS snapshots.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
  - [Build from Source](#build-from-source)
  - [Dependencies](#dependencies)
- [CLI Usage](#cli-usage)
  - [Global Flags](#global-flags)
  - [Commands](#commands)
    - [`get` - List or Get Snapshot Details](#get---list-or-get-snapshot-details)
    - [`create` - Create Snapshots](#create---create-snapshots)
    - [`version` - Show Version Information](#version---show-version-information)
    - [`daemon` - Run as Prometheus Metrics Daemon](#daemon---run-as-prometheus-metrics-daemon)
- [Daemon API](docs/api.md)
- [Data Models](#data-models)
  - [Snapshot Object](#snapshot-object)
- [Examples](#examples)
  - [Complete Workflow](#complete-workflow)
  - [Integration with Scripts](#integration-with-scripts)
- [License](#license)

## Features

- **CLI Commands**: List, get details, and create ZFS snapshots
- **Prometheus Metrics**: Daemon mode with HTTP endpoint for monitoring
- **JSON Output**: Structured output for easy parsing and integration
- **Input Validation**: Robust validation of ZFS dataset and snapshot names
- **Structured Logging**: JSON logging with zap for production use

## Installation

### Build from Source

```bash
git clone https://github.com/jsirianni/zfssnap.git
cd zfssnap
go build ./cmd/zfssnap
```

### Dependencies

- Go 1.25 or later
- ZFS utilities (`zfs` command) in PATH
- Linux/Unix system with ZFS support

## CLI Usage

### Global Flags

All commands support these global flags:

- `--zfs-bin string`: Path to zfs binary (default: detect in $PATH)
- `--timeout duration`: Command timeout (default: 30s)

### Commands

#### `get` - List or Get Snapshot Details

```bash
zfssnap get [snapshot...]
```

**Behavior:**
- **No arguments**: Lists all snapshots with full details
- **With arguments**: Returns detailed information for specified snapshots
- **Stdin input**: Reads newline-separated snapshot names from stdin when no arguments provided and stdin is not a terminal

**Examples:**
```bash
# List all snapshots
zfssnap get

# Get details for specific snapshots
zfssnap get pool@snapshot1 pool@snapshot2

# Read snapshot names from stdin
echo -e "pool@snap1\npool@snap2" | zfssnap get
```

**Output Format:**
- **Single snapshot**: JSON object
- **Multiple snapshots**: JSON array
- **No snapshots found**: Empty array `[]`

**Sample Output:**
```json
{
  "name": "pool/dataset@backup-2024-01-15",
  "dataset": "pool/dataset",
  "creation": "2024-01-15T10:30:00Z",
  "used": 1048576,
  "referenced": 2097152,
  "clones": null,
  "defer_destroy": false,
  "logical_used": 1048576,
  "logical_referenced": 2097152,
  "guid": 12345678901234567890,
  "user_refs": 0,
  "written": 1048576,
  "type": "snapshot"
}
```

#### `create` - Create Snapshots

```bash
zfssnap create [flags] <dataset> <snapshot-name>
```

**Flags:**
- `-r, --recursive`: Create snapshots recursively for all child datasets
- `--dry-run`: Show what would be created without actually creating snapshots
- `-f, --force`: Force creation even if snapshot already exists
- `--prefix string`: Add prefix to snapshot name
- `--suffix string`: Add suffix to snapshot name
- `--timestamp`: Add timestamp to snapshot name (format: YYYY-MM-DD-HHMMSS)

**Examples:**
```bash
# Basic snapshot
zfssnap create pool/dataset backup-2024-01-15

# Recursive snapshot with timestamp
zfssnap create -r --timestamp pool/dataset backup

# Dry run to see what would be created
zfssnap create --dry-run pool/dataset test-snapshot

# Force creation with prefix
zfssnap create -f --prefix "daily-" pool/dataset backup
```

**Output Format:**
```json
{
  "created": "pool/dataset@backup-2024-01-15",
  "errors": "",
  "count": 1
}
```

#### `version` - Show Version Information

```bash
zfssnap version
```

**Output:**
- Version number
- Build information
- Git commit hash

#### `daemon` - Run as Prometheus Metrics Daemon

```bash
zfssnap daemon [flags]
```

**Flags:**
- `-a, --addr string`: Address to bind the metrics server (default: "localhost:9464")

**Examples:**
```bash
# Start daemon on default port
zfssnap daemon

# Start daemon on custom port
zfssnap daemon --addr ":8080"

# Start daemon on specific interface
zfssnap daemon --addr "192.168.1.100:9464"
```

**Features:**
- Exposes Prometheus metrics at `/metrics` endpoint
- Health check endpoint at `/health`
- Periodic metric updates (every 30 seconds)
- Graceful shutdown on SIGINT/SIGTERM
- Structured JSON logging

## Data Models

### Snapshot Object

The `Snapshot` struct represents a ZFS snapshot with the following fields:

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Fully qualified snapshot name (pool/dataset@snap) |
| `dataset` | string | Parent dataset name without snapshot component |
| `creation` | time.Time | Creation timestamp (RFC3339 format) |
| `used` | uint64 | Space that would be freed if destroyed (bytes) |
| `referenced` | uint64 | Space accessible by this snapshot (bytes) |
| `clones` | []string | Datasets that are clones of this snapshot |
| `defer_destroy` | bool | Whether marked for deferred destroy |
| `logical_used` | uint64 | Logical space consumed (bytes) |
| `logical_referenced` | uint64 | Logical space accessible (bytes) |
| `guid` | uint64 | Globally unique identifier |
| `user_refs` | uint64 | Number of user holds |
| `written` | uint64 | Space written since previous snapshot (bytes) |
| `type` | string | Dataset type (typically "snapshot") |

## Examples

### Complete Workflow

```bash
# List all snapshots
zfssnap get

# Create a new snapshot
zfssnap create pool/dataset backup-$(date +%Y-%m-%d)

# Verify the snapshot was created
zfssnap get pool/dataset@backup-$(date +%Y-%m-%d)

# Start monitoring daemon
zfssnap daemon --addr ":9464"

# Check metrics
curl http://localhost:9464/metrics
```

### Integration with Scripts

```bash
#!/bin/bash
# Create snapshot and get details
SNAP_NAME="backup-$(date +%Y-%m-%d-%H%M%S)"
zfssnap create pool/dataset "$SNAP_NAME"

# Get snapshot details and extract size
SNAP_INFO=$(zfssnap get "pool/dataset@$SNAP_NAME")
SNAP_SIZE=$(echo "$SNAP_INFO" | jq -r '.used')
echo "Created snapshot $SNAP_NAME with size $SNAP_SIZE bytes"
```

## License

MIT License - see LICENSE file for details.
