package node

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/stillmatic/pummel/pkg/predicates"
	"golang.org/x/exp/slices"
)

type Node struct {
	XMLName            xml.Name
	Predicate          predicates.Predicate `xml:"any"`
	ID                 string               `xml:"id,attr"`
	Score              string               `xml:"score,attr"`
	RecordCount        int                  `xml:"recordCount,attr"`
	ScoreDistributions []*ScoreDistribution `xml:"ScoreDistribution"`
	Children           []*Node              `xml:"Node"`
	// DefaultChild gives the id of the child node to use when no predicates can be evaluated due to missing values.
	// Note that only Nodes which are immediate children of the respective Node can be referenced.
	DefaultChild string `xml:"defaultChild,attr"`
}

type Partition struct {
	XMLName xml.Name `xml:"Partition"`
	Name    string   `xml:"name,attr"`
}

type PartitionFieldStats struct {
	XMLName  xml.Name `xml:"PartitionFieldStats"`
	Field    string   `xml:"field,attr"`
	Weighted bool     `xml:"weighted,attr"`
}

type ScoreDistribution struct {
	XMLName     xml.Name `xml:"ScoreDistribution"`
	Value       string   `xml:"value,attr"`
	RecordCount int      `xml:"recordCount,attr"`
	Confidence  float64  `xml:"confidence,attr"`
	Probability float64  `xml:"probability,attr"`
}

func (n *Node) EqualTo(other *Node) bool {
	if n.XMLName != other.XMLName {
		return false
	}
	if n.ID != other.ID {
		return false
	}
	if n.Score != other.Score {
		return false
	}
	if n.RecordCount != other.RecordCount {
		return false
	}
	if len(n.Children) != len(other.Children) {
		return false
	}
	for i, child := range n.Children {
		if !child.EqualTo(other.Children[i]) {
			return false
		}
	}
	return true
}

func (n *Node) String() string {
	return fmt.Sprintf("Node(ID: %s, score: %s)", n.ID, n.Score)
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.XMLName = start.Name
	n.Children = make([]*Node, 0)
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			n.ID = attr.Value
		case "score":
			n.Score = attr.Value
		case "recordCount":
			n.RecordCount, _ = strconv.Atoi(attr.Value)
		case "defaultChild":
			n.DefaultChild = attr.Value
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
			var child Node
			var scoreDistribution ScoreDistribution
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
			// append this child to the list of children
			case "Node":
				err = d.DecodeElement(&child, &tt)
				if err != nil {
					return err
				}
				n.Children = append(n.Children, &child)
			// append this scoredistribution to the list of scoredistributions
			case "ScoreDistribution":
				err = d.DecodeElement(&scoreDistribution, &tt)
				if err != nil {
					return err
				}
				n.ScoreDistributions = append(n.ScoreDistributions, &scoreDistribution)
			default:
				return fmt.Errorf("unknown children type: %s", tt.Name.Local)
			}
			if p != nil {
				if err := d.DecodeElement(&p, &tt); err != nil {
					return err
				}
				n.Predicate = p
			}
		case xml.EndElement:
			return nil
		}
	}
}

func (n *Node) Evaluate(features map[string]interface{}) (bool, bool, error) {
	return n.Predicate.Evaluate(features)
}

func (n *Node) GetDefaultChild() (*Node, error) {
	if n.DefaultChild == "" {
		return nil, fmt.Errorf("no default child specified")
	}
	for _, child := range n.Children {
		if child.ID == n.DefaultChild {
			return child, nil
		}
	}
	return nil, fmt.Errorf("default child not found: %s", n.DefaultChild)
}

func (n *Node) GetRecordCount() int {
	var count int
	for _, sd := range n.ScoreDistributions {
		count += sd.RecordCount
	}
	return count
}

func (n *Node) GetClasses() []string {
	classes := make([]string, len(n.ScoreDistributions))
	for _, sd := range n.ScoreDistributions {
		classes = append(classes, sd.Value)
	}
	return classes
}

// ComputeWeightedConfidence computes confidences for each class from scoring it and each of its sibling Nodes in turn
// (excluding any siblings whose predicates evaluate to FALSE). The confidences returned for each class
// from each sibling Node that was scored are weighted by the proportion of the number of records in that Node,
// then summed to produce a total confidence for each class. The winner is the class with the highest confidence.
// Note that weightedConfidence should be applied recursively to deal with situations where several predicates
// within the tree evaluate to UNKNOWN during the scoring of a case.
func ComputeWeightedConfidence(ns []Node) (string, error) {
	var totalRecords int
	for _, n := range ns {
		totalRecords += n.GetRecordCount()
	}
	confidences := make(map[string]float64, len(ns[0].ScoreDistributions))
	// iterate over each node and add weighted confidence for each class
	for _, n := range ns {
		for _, sd := range n.ScoreDistributions {
			startVal := confidences[sd.Value]
			startVal += sd.Probability * float64(sd.RecordCount) / float64(totalRecords)
			confidences[sd.Value] = startVal
		}
	}
	// sort the confidences in descending order and return the class with the highest confidence
	// TODO: use a heap instead
	var winner string
	winnerConfidence := -1.0
	for class, confidence := range confidences {
		if confidence > winnerConfidence {
			winner = class
			winnerConfidence = confidence
		}
	}
	return winner, nil
}

// CheckScoreDistributionClasses checks that classes are shared among all nodes in the tree.
func (n *Node) CheckScoreDistributionClasses() (bool, error) {
	classes := n.GetClasses()
	for _, n := range n.Children {
		foundClasses := make([]string, 0, len(classes))
		for _, sd := range n.ScoreDistributions {
			if !slices.Contains(classes, sd.Value) {
				return false, fmt.Errorf("class %s not found in node %s", sd.Value, n.ID)
			}
			foundClasses = append(foundClasses, sd.Value)
		}
		if len(foundClasses) != len(classes) {
			return false, fmt.Errorf("classes not shared among all nodes: %v", foundClasses)
		}
	}
	return true, nil
}

func (n *Node) HandleScoreDistributions() (map[string]float64, float64, error) {
	var sum float64
	vals := make(map[string]float64, len(n.ScoreDistributions))
	for _, sd := range n.ScoreDistributions {
		// if sd.Probability > 0 {
		// 	vals[sd.Value] = sd.Probability
		// }
		vals[sd.Value] = float64(sd.RecordCount)
		sum += float64(sd.RecordCount)
	}
	return vals, sum, nil
}
