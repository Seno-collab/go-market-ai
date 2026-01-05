package container

import (
	"context"
	"testing"
	"time"
)

func TestConnectRedisSuccess(t *testing.T) {
	client, err := ConnectRedis("redis://localhost:6379/0")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("expected client, got nil")
	}
	if opts := client.Options(); opts == nil || opts.Addr != "localhost:6379" {
		t.Fatalf("unexpected client options: %+v", opts)
	}
	_ = client.Close()
}

func TestConnectRedisInvalidDSN(t *testing.T) {
	client, err := ConnectRedis("://bad-dsn")
	if err == nil {
		t.Fatalf("expected error for invalid DSN, got nil (client=%v)", client)
	}
}

func TestConnectRedisConnectionFailure(t *testing.T) {
	client, err := ConnectRedis("redis://127.0.0.1:1/0")
	if err != nil {
		t.Fatalf("expected client despite unreachable host, got error %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	if err := client.Ping(ctx).Err(); err == nil {
		t.Fatal("expected ping error when Redis is unreachable")
	}
}

func TestConnectRedisParsesConfig(t *testing.T) {
	client, err := ConnectRedis("redis://:s3cr3t@localhost:12345/2")
	if err != nil {
		t.Fatalf("expected no error parsing DSN, got %v", err)
	}
	defer client.Close()

	opts := client.Options()
	if opts == nil {
		t.Fatal("expected options, got nil")
	}
	if opts.Addr != "localhost:12345" {
		t.Fatalf("expected addr localhost:12345, got %s", opts.Addr)
	}
	if opts.DB != 2 {
		t.Fatalf("expected DB 2, got %d", opts.DB)
	}
	if opts.Password != "s3cr3t" {
		t.Fatalf("expected password parsed, got %q", opts.Password)
	}
}
