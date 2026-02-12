package coinai

import (
	"fmt"
)

type LinearModel struct {
	Weights []float64 `json:"weights"`
	Bias    float64   `json:"bias"`
}

func NewLinearModel(featureCount int) *LinearModel {
	return &LinearModel{
		Weights: make([]float64, featureCount),
	}
}

func (m *LinearModel) Predict(features []float64) float64 {
	var prediction float64 = m.Bias
	for i, weight := range m.Weights {
		prediction += weight * features[i]
	}
	return prediction
}

func (m *LinearModel) PredictBatch(data [][]float64) []float64 {
	preds := make([]float64, len(data))
	for i, row := range data {
		preds[i] = m.Predict(row)
	}
	return preds
}

func (m *LinearModel) Train(data [][]float64, targets []float64, cfg TrainConfig) (TrainStats, error) {
	if len(data) == 0 {
		return TrainStats{}, fmt.Errorf("empty train data")
	}
	if len(data) != len(targets) {
		return TrainStats{}, fmt.Errorf("train data and targets length mismatch")
	}
	if cfg.Epochs <= 0 {
		cfg.Epochs = 500
	}
	if cfg.LearningRate <= 0 {
		cfg.LearningRate = 0.03
	}
	if cfg.L2 < 0 {
		return TrainStats{}, fmt.Errorf("L2 cannot be negative")
	}

	featureCount := len(data[0])
	if len(m.Weights) != featureCount {
		m.Weights = make([]float64, featureCount)
	}
	for _, row := range data {
		if len(row) != featureCount {
			return TrainStats{}, fmt.Errorf("inconsistent feature dimensions")
		}
	}

	for i := 0; i < cfg.Epochs; i++ {
		gradW := make([]float64, featureCount)
		var gradB float64

		for rowIdx, row := range data {
			pred := m.Predict(row)
			diff := pred - targets[rowIdx]
			gradB += diff
			for j := 0; j < featureCount; j++ {
				gradW[j] += diff * row[j]
			}
		}

		n := float64(len(data))
		for j := 0; j < featureCount; j++ {
			gradW[j] = (gradW[j] / n) + cfg.L2*m.Weights[j]
			m.Weights[j] -= cfg.LearningRate * gradW[j]
		}
		m.Bias -= cfg.LearningRate * (gradB / n)
	}

	preds := m.PredictBatch(data)
	return TrainStats{FinalLoss: MeanSquaredError(preds, targets)}, nil
}

func MeanSquaredError(preds, actuals []float64) float64 {
	if len(preds) == 0 || len(preds) != len(actuals) {
		return 0
	}

	var sum float64
	for i := range preds {
		diff := preds[i] - actuals[i]
		sum += diff * diff
	}
	return sum / float64(len(preds))
}

func DirectionalAccuracy(preds, actuals []float64) float64 {
	if len(preds) == 0 || len(preds) != len(actuals) {
		return 0
	}

	var match int
	for i := range preds {
		if sign(preds[i]) == sign(actuals[i]) {
			match++
		}
	}
	return float64(match) / float64(len(preds))
}

func sign(v float64) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	default:
		return 0
	}
}
