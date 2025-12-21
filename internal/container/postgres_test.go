package container

import (
	"testing"
	"time"
)

func withPostgresTestConfig(t *testing.T, retries int, backoff time.Duration) func() {
	t.Helper()
	oldRetries := postgresMaxRetries
	oldBackoff := postgresInitialBackoff
	oldTimeout := postgresDialTimeout

	postgresMaxRetries = retries
	postgresInitialBackoff = backoff
	postgresDialTimeout = time.Second

	return func() {
		postgresMaxRetries = oldRetries
		postgresInitialBackoff = oldBackoff
		postgresDialTimeout = oldTimeout
	}
}

func TestConnectPostgresInvalidDSN(t *testing.T) {
	t.Cleanup(withPostgresTestConfig(t, 1, time.Millisecond))

	pool, err := ConnectPostgres("://invalid-dsn")
	if err == nil {
		t.Fatal("expected error for invalid DSN")
	}
	if pool != nil {
		t.Fatalf("expected nil pool, got %v", pool)
	}
}

func TestConnectPostgresConnectionFailure(t *testing.T) {
	t.Cleanup(withPostgresTestConfig(t, 1, time.Millisecond))

	pool, err := ConnectPostgres("postgres://user:pass@127.0.0.1:1/testdb?sslmode=disable&connect_timeout=1")
	if err == nil {
		t.Fatal("expected error when connection cannot be established")
	}
	if pool != nil {
		t.Fatalf("expected nil pool, got %v", pool)
	}
}
