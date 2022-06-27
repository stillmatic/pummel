package regression

import (
	"math"
)

type Normalizer interface {
	Normalize(map[string]interface{}) map[string]interface{}
}

type SoftMaxNormalizer struct{}
type LogitNormalizer struct{}

func (n SoftMaxNormalizer) Normalize(ys map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{}, len(ys))
	var sum float64
	for _, y := range ys {
		sum += math.Exp(y.(float64))
	}
	for i, y := range ys {
		output[i] = math.Exp(y.(float64)) / sum
	}
	return output
}

func (n LogitNormalizer) Normalize(ys map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{}, len(ys))
	for i, y := range ys {
		output[i] = 1 / (1 + math.Exp(-y.(float64)))
	}
	return output
}
