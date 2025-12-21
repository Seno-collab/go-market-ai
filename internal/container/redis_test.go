package container

import (
	"testing"
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
