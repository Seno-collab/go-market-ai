package coinai

import "time"

type Candle struct {
	OpenTime  time.Time
	CloseTime time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

type Sample struct {
	Time     time.Time
	Features []float64
	Target   float64
}

type TrainConfig struct {
	Epochs       int
	LearningRate float64
	L2           float64
}

type TrainStats struct {
	FinalLoss float64
}

type BacktestConfig struct {
	LongThreshold  float64
	ShortThreshold float64
	FeeRate        float64
}

type BacktestResult struct {
	TotalReturn float64
	WinRate     float64
	MaxDrawdown float64
	Sharpe      float64
	Trades      int
}

type Signal string

const (
	SignalBuy  Signal = "BUY"
	SignalSell Signal = "SELL"
	SignalHold Signal = "HOLD"
)
