package coinai

import (
	"math"
	"testing"
)

func TestBacktestAndSignal(t *testing.T) {
	preds := []float64{0.02, -0.03, 0.0}
	actuals := []float64{0.01, 0.02, -0.01}

	result, err := Backtest(preds, actuals, BacktestConfig{
		LongThreshold:  0.01,
		ShortThreshold: -0.01,
		FeeRate:        0,
	})
	if err != nil {
		t.Fatalf("Backtest returned error: %v", err)
	}

	wantTotalReturn := (1.01 * 0.98 * 1.0) - 1
	if !closeEnough(result.TotalReturn, wantTotalReturn, 1e-12) {
		t.Fatalf("total return = %f, want %f", result.TotalReturn, wantTotalReturn)
	}

	if got, want := result.Trades, 3; got != want {
		t.Fatalf("trades = %d, want %d", got, want)
	}

	if got, want := result.WinRate, 0.5; !closeEnough(got, want, 1e-12) {
		t.Fatalf("win rate = %f, want %f", got, want)
	}

	if got, want := result.MaxDrawdown, 0.02; !closeEnough(got, want, 1e-10) {
		t.Fatalf("max drawdown = %f, want %f", got, want)
	}

	if got, want := SignalFromPrediction(0.02, 0.01, -0.01), SignalBuy; got != want {
		t.Fatalf("signal buy = %s, want %s", got, want)
	}
	if got, want := SignalFromPrediction(-0.02, 0.01, -0.01), SignalSell; got != want {
		t.Fatalf("signal sell = %s, want %s", got, want)
	}
	if got, want := SignalFromPrediction(0.001, 0.01, -0.01), SignalHold; got != want {
		t.Fatalf("signal hold = %s, want %s", got, want)
	}
}

func closeEnough(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}
