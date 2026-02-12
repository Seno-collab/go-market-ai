package coinai

import (
	"fmt"
	"math"
)

func Backtest(preds, actuals []float64, cfg BacktestConfig) (BacktestResult, error) {
	if len(preds) == 0 {
		return BacktestResult{}, fmt.Errorf("empty predictions")
	}
	if len(preds) != len(actuals) {
		return BacktestResult{}, fmt.Errorf("predictions and actuals length mismatch")
	}
	if cfg.LongThreshold <= cfg.ShortThreshold {
		return BacktestResult{}, fmt.Errorf("long threshold must be greater than short threshold")
	}
	if cfg.FeeRate < 0 {
		return BacktestResult{}, fmt.Errorf("fee rate cannot be negative")
	}

	equity := 1.0
	peakEquity := 1.0
	maxDrawdown := 0.0
	position := 0
	trades := 0
	winBars := 0
	activeBars := 0

	periodReturns := make([]float64, 0, len(preds))
	for i := range preds {
		targetPosition := 0
		if preds[i] >= cfg.LongThreshold {
			targetPosition = 1
		} else if preds[i] <= cfg.ShortThreshold {
			targetPosition = -1
		}

		turnover := absInt(targetPosition - position)
		if turnover > 0 {
			trades++
		}

		fee := cfg.FeeRate * float64(turnover)
		periodReturn := (float64(targetPosition) * actuals[i]) - fee
		equity *= (1 + periodReturn)
		periodReturns = append(periodReturns, periodReturn)

		if targetPosition != 0 {
			activeBars++
			if periodReturn > 0 {
				winBars++
			}
		}

		if equity > peakEquity {
			peakEquity = equity
		}
		drawdown := (peakEquity - equity) / peakEquity
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}

		position = targetPosition
	}

	sharpe := sharpeRatio(periodReturns)
	winRate := 0.0
	if activeBars > 0 {
		winRate = float64(winBars) / float64(activeBars)
	}

	return BacktestResult{
		TotalReturn: equity - 1,
		WinRate:     winRate,
		MaxDrawdown: maxDrawdown,
		Sharpe:      sharpe,
		Trades:      trades,
	}, nil
}

func SignalFromPrediction(predictedReturn, longThreshold, shortThreshold float64) Signal {
	if predictedReturn >= longThreshold {
		return SignalBuy
	}
	if predictedReturn <= shortThreshold {
		return SignalSell
	}
	return SignalHold
}

func sharpeRatio(returns []float64) float64 {
	if len(returns) < 2 {
		return 0
	}

	var mean float64
	for _, r := range returns {
		mean += r
	}
	mean /= float64(len(returns))

	var variance float64
	for _, r := range returns {
		d := r - mean
		variance += d * d
	}
	variance /= float64(len(returns) - 1)
	if variance == 0 {
		return 0
	}

	return (mean / math.Sqrt(variance)) * math.Sqrt(float64(len(returns)))
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
