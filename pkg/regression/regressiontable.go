package regression

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type RegressionTable struct {
	XMLName        xml.Name `xml:"RegressionTable"`
	Predictors     []Predictor
	Intercept      float64 `xml:"intercept,attr"`
	TargetCategory string  `xml:"targetCategory,attr"`
}

func (r *RegressionTable) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	r.XMLName = start.Name
	r.Predictors = make([]Predictor, 0)
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "intercept":
			r.Intercept, _ = strconv.ParseFloat(attr.Value, 64)
		case "targetCategory":
			r.TargetCategory = attr.Value
		}
	}
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch tt := t.(type) {
		case xml.StartElement:
			var p Predictor
			switch tt.Name.Local {
			case "NumericPredictor":
				// Note that the exponent defaults to 1, hence it is not always necessary to specify.
				p = &NumericPredictor{Exponent: 1}
			case "CategoricalPredictor":
				p = &CategoricalPredictor{}
			case "PredictorTerm":
				p = &PredictorTerm{}
			default:
				return fmt.Errorf("unknown element type: %s", tt.Name.Local)
			}
			if p != nil {
				if err = d.DecodeElement(&p, &tt); err != nil {
					return errors.Wrap(err, "error decoding predictor")
				}
				r.Predictors = append(r.Predictors, p)
			}
		case xml.EndElement:
			return nil
		}
	}
}

func (r *RegressionTable) Evaluate(inputs map[string]interface{}) (float64, error) {
	result := r.Intercept

	for _, predictor := range r.Predictors {
		var value float64
		value, err := predictor.Evaluate(inputs)

		if err != nil {
			return 0, err
		}
		result += value
	}
	return result, nil
}
