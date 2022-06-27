package regression

import (
	"encoding/xml"
	"fmt"

	"github.com/stillmatic/pummel/pkg/fields"
	"github.com/stillmatic/pummel/pkg/miningschema"
	"github.com/stillmatic/pummel/pkg/transformations"
)

type RegressionModel struct {
	XMLName              xml.Name                   `xml:"RegressionModel"`
	RegressionTables     []*RegressionTable         `xml:"RegressionTable"`
	ModelName            string                     `xml:"modelName,attr"`
	FunctionName         string                     `xml:"functionName,attr"`
	ModelType            string                     `xml:"modelType,attr"`
	TargetFieldName      string                     `xml:"targetFieldName,attr"`
	MiningSchema         *miningschema.MiningSchema `xml:"MiningSchema"`
	Normalizer           Normalizer
	IsScorable           bool                                 `xml:"isScorable,attr"`
	Output               *fields.Outputs                      `xml:"Output>OutputField"`
	LocalTransformations transformations.LocalTransformations `xml:"LocalTransformations"`
}

func (rm *RegressionModel) GetOutputField() string {
	// todo: combine with output field name
	return rm.MiningSchema.GetOutputField()
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
		switch tt := t.(type) {
		case xml.StartElement:
			switch tt.Name.Local {
			case "RegressionTable":
				var rt RegressionTable
				err := d.DecodeElement(&rt, &tt)
				if err != nil {
					return err
				}
				rm.RegressionTables = append(rm.RegressionTables, &rt)
			case "MiningSchema":
				var ms miningschema.MiningSchema
				err := d.DecodeElement(&ms, &tt)
				if err != nil {
					return err
				}
				rm.MiningSchema = &ms
			case "Output":
				var out fields.Outputs
				err := d.DecodeElement(&out, &tt)
				if err != nil {
					return err
				}
				rm.Output = &out
			case "LocalTransformations":
				var lt transformations.LocalTransformations
				err := d.DecodeElement(&lt, &tt)
				if err != nil {
					return err
				}
				rm.LocalTransformations = lt
			default:
				return fmt.Errorf("unknown element: %s", tt.Name.Local)
			}

		case xml.EndElement:
			return nil
		}
	}
}

func (rm *RegressionModel) Evaluate(inputs map[string]interface{}) (map[string]interface{}, error) {
	if len(rm.LocalTransformations.DerivedFields) > 0 {
		for _, tr := range rm.LocalTransformations.DerivedFields {
			val, err := tr.Transform(inputs)
			if err != nil {
				return nil, err
			}
			inputs[tr.RequiredField()] = val
		}
	}

	switch rm.FunctionName {
	case "regression":
		return rm.EvaluateRegression(inputs)
	case "classification":
		return rm.EvaluateClassification(inputs)
	default:
		return nil, fmt.Errorf("unknown model type: %s", rm.FunctionName)
	}
}

func (rm *RegressionModel) EvaluateRegression(inputs map[string]interface{}) (map[string]interface{}, error) {
	// assume only 1 regression table in regression
	val, err := rm.RegressionTables[0].Evaluate(inputs)
	if err != nil {
		return nil, err
	}
	targetFieldName := rm.GetOutputField()

	out := map[string]interface{}{
		targetFieldName: val,
	}
	return out, nil
}

func (rm *RegressionModel) EvaluateClassification(inputs map[string]interface{}) (map[string]interface{}, error) {
	// score each category and return the one with the highest score
	scores := make(map[string]interface{}, len(rm.RegressionTables))
	var topCategory string
	var topScore float64
	for _, rt := range rm.RegressionTables {
		val, err := rt.Evaluate(inputs)
		if err != nil {
			return nil, err
		}
		if val > topScore {
			topScore = val
			topCategory = rt.TargetCategory
		}
		// check if we have a output field for this value
		if rm.Output != nil {
			tc, err := rm.Output.GetFeature(rt.TargetCategory)
			if err != nil {
				scores[rt.TargetCategory] = val
			}
			scores[tc.Name] = val
		} else {
			scores[rt.TargetCategory] = val
		}
	}
	if rm.Normalizer != nil {
		scores = rm.Normalizer.Normalize(scores)
	}
	scores[rm.GetOutputField()] = topCategory
	return scores, nil
}
