package regression

import (
	"encoding/xml"
	"strconv"
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
				err = d.DecodeElement(&p, &start)
				if err != nil {
					return err
				}
				r.Predictors = append(r.Predictors, p)
			case "CategoricalPredictor":
				err = d.DecodeElement(&p, &start)
				if err != nil {
					return err
				}
				r.Predictors = append(r.Predictors, p)
			case "PredictorTerm":
				err = d.DecodeElement(&p, &start)
				if err != nil {
					return err
				}
				r.Predictors = append(r.Predictors, p)
			}
		case xml.EndElement:
			end := t.(xml.EndElement)
			if end.Name.Local == "RegressionTable" {
				return nil
			}
		}
	}
}
