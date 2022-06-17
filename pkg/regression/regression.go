package regression

import (
	"math"

	"github.com/stillmatic/pummel/pkg/model"
)

type RegressionModel struct {
	model.PMMLModel
	RegressionTable     RegressionTable `xml:"RegressionTable"`
	ModelName           string          `xml:"ModelName,attr"`
	FunctionName        string          `xml:"FunctionName,attr"`
	ModelType           string          `xml:"ModelType,attr"`
	TargetFieldName     string          `xml:"targetFieldName,attr"`
	NormalizationMethod string          `xml:"NormalizationMethod,attr"`
	IsScorable          bool            `xml:"isScorable,attr"`
}

type Normalizer interface {
	Normalize([]float64) []float64
}

type SoftMaxNormalizer struct{}
type LogitNormalizer struct{}

func (n SoftMaxNormalizer) Normalize(ys []float64) []float64 {
	var missing bool
	output := make([]float64, len(ys))
	var sum float64
	for _, y := range ys {
		if math.IsNaN(y) {
			missing = true
			break
		}
		sum += math.Exp(y)
	}
	for i, y := range ys {
		if missing {

		} else {
			output[i] = math.Exp(y) / sum
		}
	}
	return output
}

func (n LogitNormalizer) Normalize(ys []float64) []float64 {
	output := make([]float64, len(ys))
	for i, y := range ys {
		output[i] = math.Log(y / (1 - y))
	}
	return output
}
