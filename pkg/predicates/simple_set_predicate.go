package predicates

import (
	"encoding/xml"
	"fmt"

	"github.com/mattn/go-shellwords"
	op "github.com/stillmatic/pummel/pkg/operators"
	"gopkg.in/guregu/null.v4"
)

// SimpleSetPredicate checks whether a field value is element of a set.
// The set of values is specified by the array.
type SimpleSetPredicate struct {
	XMLName  xml.Name `xml:"SimpleSetPredicate"`
	Field    string   `xml:"field,attr"`
	Operator string   `xml:"booleanOperator,attr"`
	Values   []string `xml:"Array"`
}

// Custom XML Unmarshal for SimpleSetPredicate
func (p *SimpleSetPredicate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	p.XMLName = start.Name
	p.Values = make([]string, 0)
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "field":
			p.Field = attr.Value
		case "booleanOperator":
			p.Operator = attr.Value
		}
	}
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch tt := t.(type) {
		case xml.StartElement:
			var v string
			if err := d.DecodeElement(&v, &tt); err != nil {
				return err
			}
			// TODO: see if this is really necessary...
			res, err := shellwords.Parse(v)
			if err != nil {
				return err
			}
			p.Values = res
		case xml.EndElement:
			return nil
		}
	}

}

func (p *SimpleSetPredicate) String() string {
	return fmt.Sprintf("SimpleSetPredicate(%s %s %s)", p.Field, p.Operator, p.Values)
}

func (p *SimpleSetPredicate) Evaluate(features map[string]interface{}) (null.Bool, error) {
	featureVal, exists := features[p.Field]
	if !exists {
		return null.BoolFromPtr(nil), nil
	}

	switch p.Operator {
	case op.Operators.IsIn:
		for _, value := range p.Values {
			if value == featureVal {
				return null.BoolFrom(true), nil
			}
		}
		return null.BoolFrom(false), nil
	case op.Operators.IsNotIn:
		for _, value := range p.Values {
			if value == featureVal {
				return null.BoolFrom(false), nil
			}
		}
		return null.BoolFrom(true), nil
	}
	return null.BoolFromPtr(nil), fmt.Errorf("unsupported operator: %s", p.Operator)
}
