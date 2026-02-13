package coinai

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadCandlesFromCSV(t *testing.T) {
	path := writeTempCSV(t, `
date,open,high,low,close,volume
2026-01-01,100,110,95,108,1000
2026-01-02,108,113,101,111,1200
2026-01-03,111,118,109,116,1500
`)

	candles, err := LoadCandlesFromCSV(path, 2)
	if err != nil {
		t.Fatalf("LoadCandlesFromCSV returned error: %v", err)
	}

	if got, want := len(candles), 2; got != want {
		t.Fatalf("len(candles) = %d, want %d", got, want)
	}

	if got, want := candles[0].OpenTime, time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC); !got.Equal(want) {
		t.Fatalf("first open time = %s, want %s", got, want)
	}
	if got, want := candles[1].Close, 116.0; got != want {
		t.Fatalf("last close = %f, want %f", got, want)
	}
	if got, want := candles[1].Volume, 1500.0; got != want {
		t.Fatalf("last volume = %f, want %f", got, want)
	}
}

func TestLoadCandlesFromCSVTimestampMillis(t *testing.T) {
	path := writeTempCSV(t, `
timestamp,open,high,low,close
1735689600000,10,11,9,10.5
1735776000000,10.5,12,10,11.7
`)

	candles, err := LoadCandlesFromCSV(path, 0)
	if err != nil {
		t.Fatalf("LoadCandlesFromCSV returned error: %v", err)
	}
	if got, want := len(candles), 2; got != want {
		t.Fatalf("len(candles) = %d, want %d", got, want)
	}
	if got, want := candles[0].OpenTime, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC); !got.Equal(want) {
		t.Fatalf("first open time = %s, want %s", got, want)
	}
}

func TestLoadCandlesFromCSVMissingOHLC(t *testing.T) {
	path := writeTempCSV(t, `
date,open,close
2026-01-01,100,102
`)

	_, err := LoadCandlesFromCSV(path, 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "missing required OHLC columns") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func writeTempCSV(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "candles.csv")
	trimmed := strings.TrimSpace(content) + "\n"
	if err := os.WriteFile(path, []byte(trimmed), 0o644); err != nil {
		t.Fatalf("write temp csv: %v", err)
	}
	return path
}
