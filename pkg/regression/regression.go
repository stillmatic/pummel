package regression

import (
	"encoding/xml"
	"fmt"
	"math"
)

type RegressionModel struct {
	XMLName             xml.Name           `xml:"RegressionModel"`
	RegressionTables    []*RegressionTable `xml:"RegressionTable"`
	ModelName           string             `xml:"modelName,attr"`
	FunctionName        string             `xml:"functionName,attr"`
	ModelType           string             `xml:"modelType,attr"`
	TargetFieldName     string             `xml:"targetFieldName,attr"`
	NormalizationMethod string             `xml:"normalizationMethod,attr"`
	IsScorable          bool               `xml:"isScorable,attr"`
}

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
		output[i] = math.Log(y.(float64) / (1 - y.(float64)))
	}
	return output
}

func (rm RegressionModel) Evaluate(inputs map[string]interface{}) (map[string]interface{}, error) {
	// TODO: move this into the unmarshal step
	var normalizer Normalizer
	switch rm.NormalizationMethod {
	case "softmax":
		normalizer = SoftMaxNormalizer{}
	case "logit":
		normalizer = LogitNormalizer{}
	default:
		// no normalization
	}
	switch rm.FunctionName {
	case "regression":
		// assume only 1 regression table in regression
		val, err := rm.RegressionTables[0].Evaluate(inputs)
		if err != nil {
			return nil, err
		}
		out := map[string]interface{}{
			rm.TargetFieldName: val,
		}
		return out, nil
	case "classification":
		// score each category and return the one with the highest score
		scores := make(map[string]interface{}, len(rm.RegressionTables))
		for _, rt := range rm.RegressionTables {
			val, err := rt.Evaluate(inputs)
			if err != nil {
				return nil, err
			}
			scores[rt.TargetCategory] = val
		}
		if normalizer != nil {
			scores = normalizer.Normalize(scores)
		}
		return scores, nil
	default:
		return nil, fmt.Errorf("unknown model type: %s", rm.ModelType)
	}
}
