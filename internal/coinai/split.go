package coinai

import "fmt"

func SplitSequential(samples []Sample, trainRatio float64) (train []Sample, test []Sample, err error) {
	if len(samples) < 2 {
		return nil, nil, fmt.Errorf("need at least 2 samples")
	}
	if trainRatio <= 0 || trainRatio >= 1 {
		return nil, nil, fmt.Errorf("train ratio must be in (0,1)")
	}

	splitIdx := int(float64(len(samples)) * trainRatio)
	if splitIdx <= 0 || splitIdx >= len(samples) {
		return nil, nil, fmt.Errorf("invalid split index %d", splitIdx)
	}

	train = append(train, samples[:splitIdx]...)
	test = append(test, samples[splitIdx:]...)
	return train, test, nil
}
