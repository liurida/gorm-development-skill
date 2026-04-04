
package examples

import (
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

// SetupPrometheus configures and registers the GORM Prometheus plugin.
func SetupPrometheus() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Configure the Prometheus plugin
	prom := prometheus.New(prometheus.Config{
		DBName:          "my_app_db", // A name for your database to be used as a label
		RefreshInterval: 15,          // How often to refresh the DBStats (in seconds)
		StartServer:     true,        // Expose metrics on an HTTP endpoint
		HTTPServerPort:  8081,        // A different port to avoid conflicts in tests
	})

	if err := db.Use(prom); err != nil {
		return nil, err
	}

	// The plugin will start an HTTP server on port 8081 in the background.
	// You can now scrape metrics from http://localhost:8081/metrics

	return db, nil
}

// CustomMetricsCollector is an example of a custom metrics collector.	ype CustomMetricsCollector struct{}

// Metrics implements the prometheus.MetricsCollector interface.
func (c *CustomMetricsCollector) Metrics(p *prometheus.Prometheus) []prometheus.Collector {
	// In a real application, you might query the database for custom metrics.
	// For this example, we'll just return a dummy gauge.
	customGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "my_app_custom_metric",
		Help: "An example of a custom application metric.",
	})
	customGauge.Set(42)

	return []prometheus.Collector{customGauge}
}

// SetupPrometheusWithCustomMetrics demonstrates adding a custom metrics collector.
func SetupPrometheusWithCustomMetrics() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	prom := prometheus.New(prometheus.Config{
		DBName:          "custom_metrics_db",
		StartServer:     true,
		HTTPServerPort:  8082,
		MetricsCollector: []prometheus.MetricsCollector{
			&CustomMetricsCollector{},
		},
	})

	if err := db.Use(prom); err != nil {
		return nil, err
	}

	// You can now find 'my_app_custom_metric' at http://localhost:8082/metrics

	return db, nil
}

// Example of how to verify that the metrics server is running.
func VerifyMetricsServer(port string) bool {
	resp, err := http.Get("http://localhost:" + port + "/metrics")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
