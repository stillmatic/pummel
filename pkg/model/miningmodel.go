package model

import (
	"encoding/xml"

	"github.com/pkg/errors"
	"github.com/stillmatic/pummel/pkg/fields"
	"github.com/stillmatic/pummel/pkg/miningschema"
	"github.com/stillmatic/pummel/pkg/transformations"
)

type MiningModel struct {
	XMLName              xml.Name                       `xml:"MiningModel"`
	MiningSchema         *miningschema.MiningSchema     `xml:"MiningSchema"`
	Output               *fields.Outputs                `xml:"Output"`
	Segmentation         Segmentation                   `xml:"Segmentation"`
	FunctionName         string                         `xml:"functionName,attr"`
	ModelName            string                         `xml:"modelName,attr"`
	AlgorithmName        string                         `xml:"algorithmName,attr"`
	LocalTransformations []transformations.DerivedField `xml:"LocalTransformations"`
	IsScorable           bool                           `xml:"isScorable,attr"`
	Targets              []Target                       `xml:"Targets>Target"`
}

type Target struct {
	XMLName         xml.Name `xml:"Target"`
	A               string   `xml:"a,attr"`
	RescaleConstant float64  `xml:"rescaleConstant,attr"`
	RescaleFactor   float64  `xml:"rescaleFactor,attr"`
}

var MultipleModelMethod = struct {
	MajorityVote         string
	WeightedMajorityVote string
	Average              string
	WeightedAverage      string
	Median               string
	Max                  string
	Sum                  string
	SelectFirst          string
	SelectAll            string
	ModelChain           string
}{
	MajorityVote:         "majorityVote",
	WeightedMajorityVote: "weightedMajorityVote",
	Average:              "average",
	WeightedAverage:      "weightedAverage",
	Median:               "median",
	Max:                  "max",
	Sum:                  "sum",
	SelectFirst:          "selectFirst",
	SelectAll:            "selectAll",
	ModelChain:           "modelChain",
}

func (mm *MiningModel) Evaluate(values map[string]interface{}) (map[string]interface{}, error) {
	var sum float64
	var res map[string]interface{}
	var err error
	switch mm.Segmentation.MultipleModelMethod {
	case MultipleModelMethod.Sum:
		res, err = mm.Segmentation.EvaluateSum(values, mm.Targets)
	case MultipleModelMethod.SelectFirst, MultipleModelMethod.ModelChain, MultipleModelMethod.MajorityVote:
		res, err = mm.Segmentation.Evaluate(values)

	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to evaluate segmentation")
	}
	for i, v := range res {
		outputName, err := mm.Output.GetFeature(i)
		if err != nil {
			continue
		}
		switch outputName.OpType {
		case "continuous":
			res[outputName.Name] = v.(float64)
			sum += v.(float64)
		case "categorical":
			res[outputName.Name] = v.(string)
		}
	}
	for _, v := range mm.Output.OutputFields {
		if v.Feature == "probability" {
			val, ok := res[v.Name]
			if !ok {
				val = 0.0
			}
			res[v.Name] = val.(float64) / sum
		}
	}
	return res, nil
}

func (mm *MiningModel) GetOutputField() string {
	return mm.MiningSchema.GetOutputField()
}
