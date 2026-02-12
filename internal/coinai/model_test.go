package coinai

import "testing"

func TestLinearModelTrain(t *testing.T) {
	trainX := make([][]float64, 0, 80)
	trainY := make([]float64, 0, 80)

	for i := 0; i < 80; i++ {
		x1 := -1.0 + 2.0*float64(i)/79.0
		x2 := x1 * 0.3
		y := (0.5 * x1) + (-0.2 * x2) + 0.1
		trainX = append(trainX, []float64{x1, x2})
		trainY = append(trainY, y)
	}

	model := NewLinearModel(2)
	stats, err := model.Train(trainX, trainY, TrainConfig{
		Epochs:       1500,
		LearningRate: 0.05,
		L2:           0,
	})
	if err != nil {
		t.Fatalf("Train returned error: %v", err)
	}
	if stats.FinalLoss > 1e-6 {
		t.Fatalf("unexpected final loss: %f", stats.FinalLoss)
	}

	preds := model.PredictBatch(trainX)
	mse := MeanSquaredError(preds, trainY)
	if mse > 1e-6 {
		t.Fatalf("mse = %f, want <= 1e-6", mse)
	}

	acc := DirectionalAccuracy(preds, trainY)
	if acc < 0.95 {
		t.Fatalf("directional accuracy = %f, want >= 0.95", acc)
	}
}
