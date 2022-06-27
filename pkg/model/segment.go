package model

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/stillmatic/pummel/pkg/predicates"
	"github.com/stillmatic/pummel/pkg/regression"
	"github.com/stillmatic/pummel/pkg/tree"
)

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
			case "RegressionModel":
				var rm regression.RegressionModel
				err = d.DecodeElement(&rm, &tt)
				if err != nil {
					return err
				}
				s.ModelElement = &rm
			case "MiningModel":
				var mm MiningModel
				err = d.DecodeElement(&mm, &tt)
				if err != nil {
					return err
				}
				s.ModelElement = &mm
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
	res, err := s.ModelElement.Evaluate(values)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to evaluate model element")
	}
	return res, nil
}

// Evaluate aggregates results from each segmentation
func (sg *Segmentation) Evaluate(values map[string]interface{}) (map[string]interface{}, error) {
	switch sg.MultipleModelMethod {
	case MultipleModelMethod.MajorityVote:
		return sg.EvaluateMajorityVote(values)
	case MultipleModelMethod.WeightedAverage:
		return sg.EvaluateWeightedAverage(values)
	case MultipleModelMethod.SelectFirst:
		return sg.EvaluateSelectFirst(values)
	case MultipleModelMethod.ModelChain:
		return sg.EvaluateModelChain(values)
	case MultipleModelMethod.Average:
		return sg.EvaluateAverage(values)
	default:
		return nil, fmt.Errorf("unknown multiple model method: %s", sg.MultipleModelMethod)
	}
}

func (sg *Segmentation) EvaluateSum(values map[string]interface{}, targets []Target) (map[string]interface{}, error) {
	// assume only one target
	target := targets[0]
	out := make(map[string]interface{}, 1)
	score := target.RescaleConstant
	var outputName string
	for _, s := range sg.Segments {
		res, err := s.Evaluate(values)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to evaluate segment")
		}
		for k, v := range res {
			if outputName == "" {
				outputName = k
			}
			score += v.(float64)
		}
	}
	// I suppose you could set it to 0. But why?!
	if target.RescaleFactor != 0.0 {
		score *= target.RescaleFactor
	}
	out[outputName] = score
	return out, nil
}

func (sg *Segmentation) EvaluateModelChain(values map[string]interface{}) (map[string]interface{}, error) {
	for _, s := range sg.Segments {
		res, err := s.Evaluate(values)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to evaluate segment")
		}
		for k, v := range res {
			values[k] = v
		}
	}
	return values, nil
}

func (sg *Segmentation) EvaluateSelectFirst(values map[string]interface{}) (map[string]interface{}, error) {
	for _, s := range sg.Segments {
		res, err := s.Evaluate(values)
		if err != nil {
			return nil, errors.Wrap(err, "failed to evaluate segment")
		}
		if res != nil {
			return res, nil
		}
	}
	return nil, nil
}

func (sg *Segmentation) EvaluateMajorityVote(values map[string]interface{}) (map[string]interface{}, error) {
	outputName := sg.Segments[0].ModelElement.GetOutputField()
	var topCount float64
	var topCategory string
	out := make(map[string]interface{})
	count := make(map[string]float64)
	for _, s := range sg.Segments {
		res, err := s.Evaluate(values)
		// fmt.Println("segment", i, "result", res)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to evaluate segment")
		}
		for k, v := range res {
			if k == outputName {
				newCount := count[v.(string)] + 1.0
				if newCount > topCount {
					topCount = newCount
					topCategory = v.(string)
				}
				count[v.(string)] = newCount
			}
		}
	}

	for k, v := range count {
		out[k] = v
	}
	// fmt.Println("majority vote result", out)

	out[outputName] = topCategory
	return out, nil
}

func (sg *Segmentation) EvaluateWeightedAverage(values map[string]interface{}) (map[string]interface{}, error) {
	outputName := sg.Segments[0].ModelElement.GetOutputField()

	var totalValue float64
	for _, s := range sg.Segments {
		res, err := s.Evaluate(values)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to evaluate segment")
		}
		ret := res[outputName].(float64) * s.Weight
		totalValue += ret
	}
	out := map[string]interface{}{outputName: totalValue}
	return out, nil
}

func (sg *Segmentation) EvaluateAverage(values map[string]interface{}) (map[string]interface{}, error) {
	outputName := sg.Segments[0].ModelElement.GetOutputField()
	nSegments := float64(len(sg.Segments))
	out := make(map[string]interface{})
	count := make(map[string]float64)
	for _, s := range sg.Segments {
		res, err := s.Evaluate(values)
		// fmt.Println("segment", i, "result", res)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to evaluate segment")
		}
		for k, v := range res {
			if k == outputName {
				newCount := count[k] + v.(float64)
				count[k] = newCount
			}
		}
	}

	for k, v := range count {
		out[k] = v / nSegments
	}
	fmt.Println("average vote result", out)

	return out, nil
}
