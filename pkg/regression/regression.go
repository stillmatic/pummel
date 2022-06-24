package regression

import (
	"encoding/xml"
	"fmt"
)

type RegressionModel struct {
	XMLName          xml.Name           `xml:"RegressionModel"`
	RegressionTables []*RegressionTable `xml:"RegressionTable"`
	ModelName        string             `xml:"modelName,attr"`
	FunctionName     string             `xml:"functionName,attr"`
	ModelType        string             `xml:"modelType,attr"`
	TargetFieldName  string             `xml:"targetFieldName,attr"`
	Normalizer       Normalizer
	IsScorable       bool `xml:"isScorable,attr"`
}

func (rm *RegressionModel) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	rm.XMLName = start.Name
	rm.RegressionTables = make([]*RegressionTable, 0)
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "modelName":
			rm.ModelName = attr.Value
		case "functionName":
			rm.FunctionName = attr.Value
		case "modelType":
			rm.ModelType = attr.Value
		case "targetFieldName":
			rm.TargetFieldName = attr.Value
		case "normalizationMethod":
			switch attr.Value {
			case "softmax":
				rm.Normalizer = SoftMaxNormalizer{}
			case "logit":
				rm.Normalizer = LogitNormalizer{}
			case "":
				rm.Normalizer = nil
			default:
				return fmt.Errorf("unknown normalization method: %s", attr.Value)
			}
		case "isScorable":
			rm.IsScorable = attr.Value == "true"
		}
	}
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch t := t.(type) {
		case xml.StartElement:
			start := t
			if start.Name.Local == "RegressionTable" {
				var rt RegressionTable
				err := d.DecodeElement(&rt, &start)
				if err != nil {
					return err
				}
				rm.RegressionTables = append(rm.RegressionTables, &rt)
			}
		case xml.EndElement:
			end := t
			if end.Name.Local == "RegressionModel" {
				return nil
			}
		}
	}
}

func (rm RegressionModel) Evaluate(inputs map[string]interface{}) (map[string]interface{}, error) {
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
		if rm.Normalizer != nil {
			scores = rm.Normalizer.Normalize(scores)
		}
		return scores, nil
	default:
		return nil, fmt.Errorf("unknown model type: %s", rm.FunctionName)
	}
}
