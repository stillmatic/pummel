package regression

import (
	"encoding/xml"
	"math"
)

type Predictor interface {
	Evaluate(map[string]interface{}) (float64, error)
}

type NumericPredictor struct {
	XMLName     xml.Name `xml:"NumericPredictor"`
	Name        string   `xml:"name,attr"`
	Exponent    float64  `xml:"exponent,attr"`
	Coefficient float64  `xml:"coefficient,attr"`
}

type CategoricalPredictor struct {
	XMLName     xml.Name `xml:"CategoricalPredictor"`
	Name        string   `xml:"name,attr"`
	Value       string   `xml:"value,attr"`
	Coefficient float64  `xml:"coefficient,attr"`
}

type PredictorTerm struct {
	XMLName     xml.Name   `xml:"PredictorTerm"`
	Name        string     `xml:"name,attr"`
	Coefficient float64    `xml:"coefficient,attr"`
	FieldRefs   []FieldRef `xml:"FieldRef"`
}

type FieldRef struct {
	XMLName xml.Name `xml:"FieldRef"`
	Field   string   `xml:"field,attr"`
}

func (cp CategoricalPredictor) Evaluate(inputs map[string]interface{}) (float64, error) {
	if value, ok := inputs[cp.Name]; ok {
		if value == cp.Value {
			return cp.Coefficient, nil
		}
	}
	// If the input value is missing then variable_name(v) yields 0 for any v.
	return 0, nil
}

func (np NumericPredictor) Evaluate(inputs map[string]interface{}) (float64, error) {
	if value, ok := inputs[np.Name]; ok {
		return np.Coefficient * math.Pow(value.(float64), np.Exponent), nil
	}
	// if the input value is missing, the result evaluates to a missing value.
	return 0, nil
}

func (pt PredictorTerm) Evaluate(inputs map[string]interface{}) (float64, error) {
	result := pt.Coefficient
	for _, fieldRef := range pt.FieldRefs {
		if value, ok := inputs[fieldRef.Field]; ok {
			result *= value.(float64)
		}
	}
	return result, nil
}
