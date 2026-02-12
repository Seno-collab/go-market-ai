package coinai

import (
	"fmt"
	"math"
)

type StandardScaler struct {
	Means []float64 `json:"means"`
	Stds  []float64 `json:"stds"`
}

func NewStandardScaler(featureCount int) *StandardScaler {
	return &StandardScaler{
		Means: make([]float64, featureCount),
		Stds:  make([]float64, featureCount),
	}
}

func (s *StandardScaler) Fit(data [][]float64) error {
	if len(data) == 0 {
		return fmt.Errorf("empty data")
	}
	if len(data[0]) == 0 {
		return fmt.Errorf("zero feature dimensions")
	}

	featureCount := len(data[0])
	if len(s.Means) != featureCount || len(s.Stds) != featureCount {
		s.Means = make([]float64, featureCount)
		s.Stds = make([]float64, featureCount)
	}

	for _, row := range data {
		if len(row) != featureCount {
			return fmt.Errorf("inconsistent feature dimensions")
		}
		for j, value := range row {
			s.Means[j] += value
		}
	}
	for j := range s.Means {
		s.Means[j] /= float64(len(data))
	}

	for _, row := range data {
		for j, value := range row {
			diff := value - s.Means[j]
			s.Stds[j] += diff * diff
		}
	}
	for j := range s.Stds {
		s.Stds[j] /= float64(len(data))
		if s.Stds[j] == 0 {
			s.Stds[j] = 1
			continue
		}
		s.Stds[j] = math.Sqrt(s.Stds[j])
	}

	return nil
}

func (s *StandardScaler) TransformBatch(data [][]float64) ([][]float64, error) {
	out := make([][]float64, 0, len(data))
	for _, row := range data {
		normalized, err := s.Transform(row)
		if err != nil {
			return nil, err
		}
		out = append(out, normalized)
	}
	return out, nil
}

func (s *StandardScaler) Transform(row []float64) ([]float64, error) {
	if len(row) != len(s.Means) || len(row) != len(s.Stds) {
		return nil, fmt.Errorf("row dimensions mismatch scaler")
	}

	normalized := make([]float64, len(row))
	for i, value := range row {
		normalized[i] = (value - s.Means[i]) / s.Stds[i]
	}
	return normalized, nil
}
