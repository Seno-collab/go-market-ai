package healthapp

import "time"

// ServiceStatus captures the health of an external dependency.
type ServiceStatus struct {
	Name      string `json:"name" example:"postgres"`
	Status    string `json:"status" example:"up"`
	LatencyMs int64  `json:"latency_ms,omitempty" example:"12"`
	Error     string `json:"error,omitempty" example:"dial tcp 127.0.0.1:5432: connect: connection refused"`
}

// HealthResponse aggregates the status of the application and its dependencies.
type HealthResponse struct {
	Status      string          `json:"status" example:"up"`
	Environment string          `json:"environment,omitempty" example:"development"`
	CheckedAt   time.Time       `json:"checked_at" example:"2026-01-22T12:34:56Z"`
	Services    []ServiceStatus `json:"services"`
}
