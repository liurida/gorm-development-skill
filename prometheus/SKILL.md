---
name: gorm-prometheus
description: Use when integrating GORM with Prometheus for database monitoring, collecting DBStats metrics, or creating custom database metrics for observability.
---

# Prometheus

The Prometheus plugin exports GORM metrics to Prometheus for monitoring database health and performance.

**Reference:** https://gorm.io/docs/prometheus.html
**Repository:** https://github.com/go-gorm/prometheus

## Features

- Automatic DBStats collection (connections, wait times)
- Custom metrics support
- Push gateway support
- Built-in HTTP server for metrics endpoint

## Quick Reference

| Config | Purpose |
|--------|---------|
| `DBName` | Label for metrics (required for multi-db) |
| `RefreshInterval` | How often to refresh metrics (seconds) |
| `StartServer` | Start built-in HTTP metrics server |
| `HTTPServerPort` | Port for metrics endpoint |
| `PushAddr` | Prometheus push gateway address |
| `MetricsCollector` | Custom metrics collectors |

## Basic Setup

```go
import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
  "gorm.io/plugin/prometheus"
)

db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

db.Use(prometheus.New(prometheus.Config{
  DBName:          "db1",           // Metrics label
  RefreshInterval: 15,              // Refresh every 15 seconds
  StartServer:     true,            // Start metrics HTTP server
  HTTPServerPort:  8080,            // Metrics available at :8080/metrics
}))
```

## Push Gateway Configuration

For environments where Prometheus cannot scrape directly.

```go
db.Use(prometheus.New(prometheus.Config{
  DBName:          "db1",
  RefreshInterval: 15,
  PushAddr:        "http://pushgateway:9091", // Push metrics here
}))
```

## DBStats Metrics

The plugin automatically collects `database/sql` DBStats:

| Metric | Description |
|--------|-------------|
| `gorm_dbstats_max_open_connections` | Maximum open connections |
| `gorm_dbstats_open_connections` | Current open connections |
| `gorm_dbstats_in_use` | Connections in use |
| `gorm_dbstats_idle` | Idle connections |
| `gorm_dbstats_wait_count` | Total waits for connection |
| `gorm_dbstats_wait_duration` | Total wait time |
| `gorm_dbstats_max_idle_closed` | Closed due to max idle |
| `gorm_dbstats_max_lifetime_closed` | Closed due to max lifetime |

## Custom Metrics Collector

Implement `MetricsCollector` interface for custom metrics.

```go
type MetricsCollector interface {
  Metrics(*Prometheus) []prometheus.Collector
}
```

### MySQL Status Metrics

Built-in collector for MySQL status variables.

```go
db.Use(prometheus.New(prometheus.Config{
  DBName:          "db1",
  RefreshInterval: 15,
  StartServer:     true,
  HTTPServerPort:  8080,
  MetricsCollector: []prometheus.MetricsCollector{
    &prometheus.MySQL{
      Prefix: "gorm_status_",          // Metric name prefix
      Interval: 100,                    // Custom refresh interval
      VariableNames: []string{          // Specific variables to collect
        "Threads_running",
        "Threads_connected",
        "Slow_queries",
        "Questions",
        "Queries",
      },
    },
  },
}))
```

Available MySQL metrics with default prefix `gorm_status_`:
- `gorm_status_Threads_running` - Currently executing queries
- `gorm_status_Threads_connected` - Connected clients
- `gorm_status_Slow_queries` - Slow query count
- `gorm_status_Questions` - Total statements executed
- And many more from `SHOW STATUS`

## Multiple Database Monitoring

```go
// Primary database
db1.Use(prometheus.New(prometheus.Config{
  DBName:         "primary",
  RefreshInterval: 15,
  StartServer:    true,
  HTTPServerPort: 8080,
}))

// Replica database (only first HTTPServerPort is used)
db2.Use(prometheus.New(prometheus.Config{
  DBName:         "replica",
  RefreshInterval: 15,
}))
```

## Custom Metrics Example

```go
type QueryCountCollector struct {
  queryCount *prometheus.CounterVec
}

func (c *QueryCountCollector) Metrics(p *prometheus.Prometheus) []prometheus.Collector {
  c.queryCount = prometheus.NewCounterVec(
    prometheus.CounterOpts{
      Name: "gorm_query_total",
      Help: "Total number of queries",
    },
    []string{"db", "table"},
  )
  return []prometheus.Collector{c.queryCount}
}

// Register with GORM
db.Use(prometheus.New(prometheus.Config{
  DBName:          "db1",
  MetricsCollector: []prometheus.MetricsCollector{
    &QueryCountCollector{},
  },
}))
```

## Prometheus Scrape Config

Add to `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'gorm'
    static_configs:
      - targets: ['app:8080']
    scrape_interval: 15s
```

## Grafana Dashboard Queries

Example PromQL queries:

```promql
# Connection pool utilization
gorm_dbstats_in_use{db="primary"} / gorm_dbstats_max_open_connections{db="primary"}

# Connection wait rate
rate(gorm_dbstats_wait_count{db="primary"}[5m])

# MySQL threads
gorm_status_Threads_running{db="primary"}
```

## When NOT to Use

- **In short-lived applications or serverless functions** - The pull-based model of Prometheus is not a good fit. Use the push gateway (`PushAddr`) instead.
- **If you don't have a Prometheus server** - The plugin only exposes metrics; it doesn't provide a monitoring UI. You need a full Prometheus and Grafana stack to make use of it.
- **For business-level metrics** - This plugin is for monitoring database health. For business metrics (e.g., user sign-ups, orders placed), instrument your application code directly with Prometheus client libraries.
- **When the overhead is too high** - For extremely high-performance, low-latency systems, the small overhead of collecting and serving metrics might be undesirable. Profile first to determine the impact.

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Multiple `HTTPServerPort` configurations | Only first port is used |
| Missing `DBName` with multiple databases | Always set unique `DBName` |
| Refresh interval too low | Keep at 15+ seconds to avoid overhead |
| Not exposing metrics endpoint | Enable `StartServer` or use push gateway |
