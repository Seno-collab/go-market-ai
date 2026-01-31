package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Database metrics
	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	DBConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
	)

	DBConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	DBQueryErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_query_errors_total",
			Help: "Total number of database query errors",
		},
		[]string{"operation", "table"},
	)

	// Cache metrics
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type"},
	)

	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type"},
	)

	CacheOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cache_operation_duration_seconds",
			Help:    "Duration of cache operations in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1},
		},
		[]string{"operation", "cache_type"},
	)

	// HTTP metrics (additional to echoprometheus)
	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)

	// Authentication metrics
	AuthAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"status"},
	)

	AuthTokenRefreshes = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_token_refreshes_total",
			Help: "Total number of token refresh operations",
		},
	)

	ActiveSessions = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "auth_active_sessions",
			Help: "Number of active user sessions",
		},
	)

	// Business metrics
	RestaurantsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "restaurants_total",
			Help: "Total number of restaurants",
		},
	)

	MenuItemsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "menu_items_total",
			Help: "Total number of menu items",
		},
	)

	FileUploadsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "file_uploads_total",
			Help: "Total number of file uploads",
		},
		[]string{"type", "status"},
	)

	FileUploadBytes = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "file_upload_bytes_total",
			Help: "Total bytes uploaded",
		},
	)
)

// RecordDBQuery records database query metrics
func RecordDBQuery(operation, table string, durationSeconds float64, err error) {
	DBQueryDuration.WithLabelValues(operation, table).Observe(durationSeconds)
	if err != nil {
		DBQueryErrors.WithLabelValues(operation, table).Inc()
	}
}

// RecordCacheHit records a cache hit
func RecordCacheHit(cacheType string) {
	CacheHits.WithLabelValues(cacheType).Inc()
}

// RecordCacheMiss records a cache miss
func RecordCacheMiss(cacheType string) {
	CacheMisses.WithLabelValues(cacheType).Inc()
}

// RecordCacheOperation records cache operation duration
func RecordCacheOperation(operation, cacheType string, durationSeconds float64) {
	CacheOperationDuration.WithLabelValues(operation, cacheType).Observe(durationSeconds)
}

// RecordAuthAttempt records an authentication attempt
func RecordAuthAttempt(success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	AuthAttempts.WithLabelValues(status).Inc()
}

// RecordFileUpload records a file upload
func RecordFileUpload(fileType string, success bool, sizeBytes int64) {
	status := "success"
	if !success {
		status = "failure"
	}
	FileUploadsTotal.WithLabelValues(fileType, status).Inc()
	if success {
		FileUploadBytes.Add(float64(sizeBytes))
	}
}
