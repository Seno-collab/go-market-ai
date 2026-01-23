package healthapp

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type stubDB struct {
	err   error
	delay time.Duration
}

func (s stubDB) Ping(ctx context.Context) error {
	if s.delay > 0 {
		select {
		case <-time.After(s.delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return s.err
}

type stubRedis struct {
	err   error
	delay time.Duration
}

func (s stubRedis) Ping(ctx context.Context) *redis.StatusCmd {
	if s.delay > 0 {
		select {
		case <-time.After(s.delay):
		case <-ctx.Done():
			cmd := redis.NewStatusCmd(ctx, "PING")
			cmd.SetErr(ctx.Err())
			return cmd
		}
	}
	cmd := redis.NewStatusCmd(ctx, "PING")
	if s.err != nil {
		cmd.SetErr(s.err)
	} else {
		cmd.SetVal("PONG")
	}
	return cmd
}

func TestExecute_AllHealthy(t *testing.T) {
	uc := NewCheckHealthUseCase(stubDB{}, stubRedis{}, "test")
	resp, ok := uc.Execute(context.Background())

	if !ok {
		t.Fatal("expected health to be healthy")
	}
	if resp.Status != StatusUp {
		t.Fatalf("expected status %q, got %q", StatusUp, resp.Status)
	}
	if resp.Environment != "test" {
		t.Fatalf("expected environment test, got %s", resp.Environment)
	}
	if len(resp.Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(resp.Services))
	}
	for _, svc := range resp.Services {
		if svc.Status != StatusUp {
			t.Fatalf("expected service %s to be up, got %s", svc.Name, svc.Status)
		}
		if svc.LatencyMs < 0 {
			t.Fatalf("expected non-negative latency, got %d", svc.LatencyMs)
		}
	}
	if resp.CheckedAt.IsZero() {
		t.Fatal("expected CheckedAt to be set")
	}
}

func TestExecute_DegradedWhenDependencyFails(t *testing.T) {
	dbErr := errors.New("db down")
	uc := NewCheckHealthUseCase(stubDB{err: dbErr}, stubRedis{}, "test")

	resp, ok := uc.Execute(context.Background())
	if ok {
		t.Fatal("expected health to be degraded when db fails")
	}
	if resp.Status != StatusDegraded {
		t.Fatalf("expected status %q, got %q", StatusDegraded, resp.Status)
	}

	var found bool
	for _, svc := range resp.Services {
		if svc.Name == "postgres" {
			found = true
			if svc.Status != StatusDown {
				t.Fatalf("expected postgres status %q, got %q", StatusDown, svc.Status)
			}
			if svc.Error == "" {
				t.Fatal("expected postgres error to be populated")
			}
		}
	}
	if !found {
		t.Fatal("expected postgres service entry")
	}
}

func TestExecute_DownWhenAllDependenciesFail(t *testing.T) {
	errDB := errors.New("db unavailable")
	errCache := errors.New("redis unavailable")
	uc := NewCheckHealthUseCase(stubDB{err: errDB}, stubRedis{err: errCache}, "test")

	resp, ok := uc.Execute(context.Background())
	if ok {
		t.Fatal("expected health to be down when all dependencies fail")
	}
	if resp.Status != StatusDown {
		t.Fatalf("expected status %q, got %q", StatusDown, resp.Status)
	}
}
