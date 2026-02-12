package coinai

import (
	"math"
	"testing"
	"time"
)

func TestBuildDataset(t *testing.T) {
	candles := mockCandles([]float64{100, 101, 102, 103, 104, 105, 106, 107})

	samples, err := BuildDataset(candles)
	if err != nil {
		t.Fatalf("BuildDataset returned error: %v", err)
	}

	if got, want := len(samples), 2; got != want {
		t.Fatalf("len(samples) = %d, want %d", got, want)
	}

	if got, want := len(samples[0].Features), len(featureNames); got != want {
		t.Fatalf("feature length = %d, want %d", got, want)
	}

	gotRange := samples[0].Features[2]
	wantRange := (candles[5].High - candles[5].Low) / candles[5].Close
	if !nearlyEqual(gotRange, wantRange, 1e-12) {
		t.Fatalf("range_ratio = %f, want %f", gotRange, wantRange)
	}

	gotTarget := samples[0].Target
	wantTarget := (candles[6].Close / candles[5].Close) - 1
	if !nearlyEqual(gotTarget, wantTarget, 1e-12) {
		t.Fatalf("target = %f, want %f", gotTarget, wantTarget)
	}
}

func TestBuildLatestFeatures(t *testing.T) {
	candles := mockCandles([]float64{100, 101, 102, 103, 104, 105, 106})

	features, err := BuildLatestFeatures(candles)
	if err != nil {
		t.Fatalf("BuildLatestFeatures returned error: %v", err)
	}

	if got, want := len(features), len(featureNames); got != want {
		t.Fatalf("feature length = %d, want %d", got, want)
	}
}

func TestBuildDatasetRangeRatioDivisionByZero(t *testing.T) {
	candles := mockCandles([]float64{100, 101, 102, 103, 104, 0, 106, 107})

	_, err := BuildDataset(candles)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func mockCandles(closes []float64) []Candle {
	out := make([]Candle, 0, len(closes))
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for i, closePrice := range closes {
		open := closePrice - 1
		high := closePrice + 2
		low := closePrice - 2
		volume := float64(100 + i)
		out = append(out, Candle{
			OpenTime:  base.Add(time.Duration(i) * time.Hour),
			CloseTime: base.Add(time.Duration(i+1) * time.Hour),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     closePrice,
			Volume:    volume,
		})
	}
	return out
}

func nearlyEqual(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}
