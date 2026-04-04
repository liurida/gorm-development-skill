
# Key Concepts for GORM Prometheus Plugin

This document provides key concepts for using the GORM Prometheus plugin to monitor database metrics.

## Overview

The `gorm.io/plugin/prometheus` plugin allows you to collect and expose database statistics and custom metrics in a format that can be scraped by a Prometheus server. This is essential for monitoring the health and performance of your database in a production environment.

## Basic Configuration

You register the Prometheus plugin with your GORM DB instance and provide a configuration.

```go
import "gorm.io/plugin/prometheus"

db.Use(prometheus.New(prometheus.Config{
    DBName:          "my_app_db", // A name for your database to be used as a label
    RefreshInterval: 15,          // How often to refresh the DBStats (in seconds)
    StartServer:     true,        // Expose metrics on an HTTP endpoint
    HTTPServerPort:  8080,        // The port for the metrics server
}))
```

### Key Configuration Options

- **`DBName`**: A label to identify the database in your metrics. This is crucial if you are monitoring multiple databases.
- **`RefreshInterval`**: The frequency at which the plugin collects `sql.DBStats`.
- **`StartServer`**: If `true`, the plugin will start an HTTP server to expose the metrics on a `/metrics` endpoint.
- **`HTTPServerPort`**: The port for the metrics server.
- **`PushAddr`**: If you are using a Prometheus Pushgateway, you can specify its address here.

## Default Metrics

By default, the plugin collects and exposes the standard `sql.DBStats` metrics, including:
- `gorm_dbstats_max_open_connections`: Maximum number of open connections to the database.
- `gorm_dbstats_open_connections`: The number of established connections both in use and idle.
- `gorm_dbstats_in_use`: The number of connections currently in use.
- `gorm_dbstats_idle`: The number of idle connections.
- `gorm_dbstats_wait_count`: The total number of connections waited for.
- `gorm_dbstats_wait_duration`: The total time blocked waiting for a new connection.
- `gorm_dbstats_max_idle_closed`: The total number of connections closed due to `SetMaxIdleConns`.
- `gorm_dbstats_max_lifetime_closed`: The total number of connections closed due to `SetConnMaxLifetime`.

These metrics are automatically labeled with the `db_name` you provide in the configuration.

## Custom Metrics

You can define your own metrics to be collected by implementing the `prometheus.MetricsCollector` interface.

```go
type MetricsCollector interface {
    Metrics(*Prometheus) []prometheus.Collector
}
```

This is useful for collecting database-specific metrics (e.g., from `SHOW STATUS` in MySQL) or application-level metrics related to database usage.

### Example: Custom MySQL Metrics

The plugin provides a built-in collector for MySQL status variables.

```go
db.Use(prometheus.New(prometheus.Config{
    // ... other config
    MetricsCollector: []prometheus.MetricsCollector{
        &prometheus.MySQL{
            VariableNames: []string{"Threads_running", "Connections"},
        },
    },
}))
```

This will expose metrics like `gorm_status_Threads_running` and `gorm_status_Connections`.

Using the Prometheus plugin is a powerful way to gain insight into how your application is interacting with the database and to proactively identify potential issues.
