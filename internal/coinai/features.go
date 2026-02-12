package coinai

import (
	"fmt"
	"math"
)

const (
	featureWindow = 5
	minCandles    = featureWindow + 2
)

var featureNames = []string{
	"ret_1",
	"mom_3",
	"range_ratio",
	"vol_change",
	"volatility_5",
}

func FeatureNames() []string {
	names := make([]string, len(featureNames))
	copy(names, featureNames)
	return names
}

func BuildDataset(candles []Candle) ([]Sample, error) {
	if len(candles) < minCandles {
		return nil, fmt.Errorf("need at least %d candles, got %d", minCandles, len(candles))
	}

	samples := make([]Sample, 0, len(candles)-minCandles+1)
	for i := featureWindow; i < len(candles)-1; i++ {
		features, err := featureAt(candles, i)
		if err != nil {
			return nil, err
		}
		target, err := pctChange(candles[i+1].Close, candles[i].Close)
		if err != nil {
			return nil, fmt.Errorf("target at index %d: %w", i, err)
		}

		samples = append(samples, Sample{
			Time:     candles[i].CloseTime,
			Features: features,
			Target:   target,
		})
	}

	return samples, nil
}

func BuildLatestFeatures(candles []Candle) ([]float64, error) {
	if len(candles) < minCandles-1 {
		return nil, fmt.Errorf("need at least %d candles, got %d", minCandles-1, len(candles))
	}
	return featureAt(candles, len(candles)-1)
}

func featureAt(candles []Candle, i int) ([]float64, error) {
	if i < featureWindow || i >= len(candles) {
		return nil, fmt.Errorf("feature index out of range")
	}

	ret1, err := pctChange(candles[i].Close, candles[i-1].Close)
	if err != nil {
		return nil, fmt.Errorf("ret_1 index %d: %w", i, err)
	}
	mom3, err := pctChange(candles[i].Close, candles[i-3].Close)
	if err != nil {
		return nil, fmt.Errorf("mom_3 index %d: %w", i, err)
	}
	if candles[i].Close == 0 {
		return nil, fmt.Errorf("range_ratio index %d: division by zero", i)
	}
	rangeRatio := (candles[i].High - candles[i].Low) / candles[i].Close
	volChg, err := pctChange(candles[i].Volume, candles[i-1].Volume)
	if err != nil {
		return nil, fmt.Errorf("vol_change index %d: %w", i, err)
	}
	volatility5, err := rollingVolatility(candles, i, featureWindow)
	if err != nil {
		return nil, err
	}

	return []float64{ret1, mom3, rangeRatio, volChg, volatility5}, nil
}

func rollingVolatility(candles []Candle, endIdx, window int) (float64, error) {
	if endIdx < window || endIdx >= len(candles) {
		return 0, fmt.Errorf("volatility index out of range")
	}

	rets := make([]float64, 0, window)
	for i := endIdx - window + 1; i <= endIdx; i++ {
		r, err := pctChange(candles[i].Close, candles[i-1].Close)
		if err != nil {
			return 0, fmt.Errorf("volatility return index %d: %w", i, err)
		}
		rets = append(rets, r)
	}
	return stddev(rets), nil
}

func stddev(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	var mean float64
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))

	var variance float64
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values))

	return math.Sqrt(variance)
}

func pctChange(newVal, oldVal float64) (float64, error) {
	if oldVal == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return (newVal / oldVal) - 1, nil
}
