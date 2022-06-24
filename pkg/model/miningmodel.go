package model

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/stillmatic/pummel/pkg/miningschema"
	"github.com/stillmatic/pummel/pkg/predicates"
	"github.com/stillmatic/pummel/pkg/transformations"
	"github.com/stillmatic/pummel/pkg/tree"
)

type MiningModel struct {
	XMLName              xml.Name                       `xml:"MiningModel"`
	MiningSchema         *miningschema.MiningSchema     `xml:"MiningSchema"`
	Segmentation         Segmentation                   `xml:"Segmentation"`
	FunctionName         string                         `xml:"functionName,attr"`
	ModelName            string                         `xml:"modelName,attr"`
	AlgorithmName        string                         `xml:"algorithmName,attr"`
	LocalTransformations []transformations.DerivedField `xml:"LocalTransformations"`
	IsScorable           bool                           `xml:"isScorable,attr"`
}

type Segmentation struct {
	XMLName             xml.Name  `xml:"Segmentation"`
	MultipleModelMethod string    `xml:"multipleModelMethod,attr"`
	Segments            []Segment `xml:"Segment"`
}

type Segment struct {
	XMLName      xml.Name                `xml:"Segment"`
	Predicates   []*predicates.Predicate `xml:"Predicate"`
	ModelElement ModelElement
	ID           string  `xml:"id,attr"`
	Weight       float64 `xml:"weight,attr"`
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

// custom xml unmarshaler for Segment
func (s *Segment) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	s.XMLName = start.Name
	s.Predicates = make([]*predicates.Predicate, 0)
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			s.ID = attr.Value
		case "weight":
			s.Weight, _ = strconv.ParseFloat(attr.Value, 64)
		}
	}
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch tt := t.(type) {
		case xml.StartElement:
			var p predicates.Predicate
			switch tt.Name.Local {
			case "SimplePredicate":
				p = &predicates.SimplePredicate{}
			case "SimpleSetPredicate":
				p = &predicates.SimpleSetPredicate{}
			case "True":
				p = &predicates.TruePredicate{}
			case "False":
				p = &predicates.FalsePredicate{}
			case "CompoundPredicate":
				p = &predicates.CompoundPredicate{}
			case "TreeModel":
				var tm tree.TreeModel
				err = d.DecodeElement(&tm, &tt)
				if err != nil {
					return err
				}
				s.ModelElement = &tm
			default:
				return fmt.Errorf("unknown children type: %s", tt.Name.Local)
			}
			if p != nil {
				if err := d.DecodeElement(&p, &tt); err != nil {
					return err
				}
				s.Predicates = append(s.Predicates, &p)
			}
		case xml.EndElement:
			return nil
		}
	}
}

func (s *Segment) Evaluate(values map[string]interface{}) (map[string]interface{}, error) {
	// out := make(map[string]interface{})
	for _, p := range s.Predicates {
		// skip if predicate is not satisfied
		if predEval, _ := (*p).Evaluate(values); !predEval.ValueOrZero() {
			return nil, nil
		}
	}
	// res, err := s.ModelElement.Evaluate(values)
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "failed to evaluate model element")
	// }
	return nil, nil

}
