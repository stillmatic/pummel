package regression

import "encoding/xml"

type Predictor interface{}

type NumericPredictor struct {
	XMLName     xml.Name `xml:"NumericPredictor"`
	Name        string   `xml:"name,attr"`
	Exponent    int64    `xml:"exponent,attr"`
	Coefficient float64  `xml:"coefficient,attr"`
}

type CategoricalPredictor struct {
	XMLName     xml.Name `xml:"CategoricalPredictor"`
	Name        string   `xml:"name,attr"`
	Value       string   `xml:"value,attr"`
	Coefficient float64  `xml:"coefficient,attr"`
}

type PredictorTerm struct {
	XMLName     xml.Name `xml:"PredictorTerm"`
	Name        string   `xml:"name,attr"`
	Coefficient float64  `xml:"coefficient,attr"`
}
