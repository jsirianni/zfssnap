# Daemon API

The zfssnap daemon exposes HTTP endpoints for monitoring ZFS snapshots via Prometheus metrics.

## HTTP Endpoints

### `GET /metrics`

Prometheus metrics endpoint. Returns metrics in Prometheus exposition format.

**Sample Output:**
```
# HELP zfs_snapshot_count Total number of ZFS snapshots
# TYPE zfs_snapshot_count gauge
zfs_snapshot_count 3
```

**Metrics:**
- `zfs_snapshot_count`: Current total number of ZFS snapshots

### `GET /health`

Health check endpoint for monitoring daemon status.

**Response:**
- **200 OK**: Daemon is healthy
- **500 Internal Server Error**: Daemon has issues

**Sample Output:**
```
OK
```

## Monitoring Integration

### Prometheus Configuration

Add to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'zfssnap'
    static_configs:
      - targets: ['localhost:9464']
    scrape_interval: 30s
    metrics_path: /metrics
```

### Grafana Dashboard

Use the `zfs_snapshot_count` metric to create dashboards showing:
- Total snapshot count over time
- Snapshot count trends
- Alerting when snapshot count drops unexpectedly

## Daemon Configuration

### Command Line Options

- `-a, --addr string`: Address to bind the metrics server (default: "localhost:9464")

### Examples

```bash
# Start daemon on default port
zfssnap daemon

# Start daemon on custom port
zfssnap daemon --addr ":8080"

# Start daemon on specific interface
zfssnap daemon --addr "192.168.1.100:9464"
```

### Features

- Exposes Prometheus metrics at `/metrics` endpoint
- Health check endpoint at `/health`
- Periodic metric updates (every 30 seconds)
- Graceful shutdown on SIGINT/SIGTERM
- Structured JSON logging
