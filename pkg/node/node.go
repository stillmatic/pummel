package node

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/stillmatic/pummel/pkg/predicates"
	"gopkg.in/guregu/null.v4"
)

type Node struct {
	XMLName     xml.Name
	Predicate   *predicates.Predicate `xml:"any"`
	ID          string                `xml:"id,attr"`
	Score       string                `xml:"score,attr"`
	RecordCount int                   `xml:"recordCount,attr"`
	Children    []Node                `xml:"Node"`
}

func (n Node) EqualTo(other Node) bool {
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

func (n Node) String() string {
	return fmt.Sprintf("Node(ID: %s, score: %s)", n.ID, n.Score)
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.XMLName = start.Name
	n.Children = make([]Node, 0)
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			n.ID = attr.Value
		case "score":
			n.Score = attr.Value
		case "recordCount":
			n.RecordCount, _ = strconv.Atoi(attr.Value)
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
			case "Node":
				err = d.DecodeElement(&child, &tt)
				if err != nil {
					return err
				}
				n.Children = append(n.Children, child)
			default:
				return fmt.Errorf("unknown children type: %s", tt.Name.Local)
			}
			if p != nil {
				if err := d.DecodeElement(&p, &tt); err != nil {
					return err
				}
				n.Predicate = &p
			}
		case xml.EndElement:
			return nil
		}
	}
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

func (n *Node) True(features map[string]interface{}) (null.Bool, error) {
	res, err := (*n.Predicate).True(features)
	return res, err
}
