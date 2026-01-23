package healthapp

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	StatusUp       = "up"
	StatusDegraded = "degraded"
	StatusDown     = "down"
)

type databasePinger interface {
	Ping(ctx context.Context) error
}

type cachePinger interface {
	Ping(ctx context.Context) *redis.StatusCmd
}

// CheckHealthUseCase probes critical dependencies to derive overall service health.
type CheckHealthUseCase struct {
	db          databasePinger
	cache       cachePinger
	environment string
}

func NewCheckHealthUseCase(db databasePinger, cache cachePinger, environment string) *CheckHealthUseCase {
	return &CheckHealthUseCase{
		db:          db,
		cache:       cache,
		environment: environment,
	}
}

// Execute returns the current health snapshot and whether the service is fully healthy.
func (uc *CheckHealthUseCase) Execute(ctx context.Context) (HealthResponse, bool) {
	services := make([]ServiceStatus, 0, 2)

	services = append(services, uc.check(ctx, "postgres", func(ctx context.Context) error {
		return uc.db.Ping(ctx)
	}))

	services = append(services, uc.check(ctx, "redis", func(ctx context.Context) error {
		return uc.cache.Ping(ctx).Err()
	}))

	overall := deriveOverallStatus(services)

	return HealthResponse{
		Status:      overall,
		Environment: uc.environment,
		CheckedAt:   time.Now().UTC(),
		Services:    services,
	}, overall == StatusUp
}

func (uc *CheckHealthUseCase) check(ctx context.Context, name string, probe func(context.Context) error) ServiceStatus {
	start := time.Now()
	err := probe(ctx)
	latency := time.Since(start).Milliseconds()

	status := StatusUp
	errMsg := ""
	if err != nil {
		status = StatusDown
		errMsg = err.Error()
	}

	return ServiceStatus{
		Name:      name,
		Status:    status,
		LatencyMs: latency,
		Error:     errMsg,
	}
}

func deriveOverallStatus(services []ServiceStatus) string {
	if len(services) == 0 {
		return StatusDown
	}

	allDown := true
	for _, s := range services {
		if s.Status == StatusUp {
			allDown = false
		}
	}

	if allDown {
		return StatusDown
	}

	for _, s := range services {
		if s.Status == StatusDown {
			return StatusDegraded
		}
	}

	return StatusUp
}
