package predicates

import (
	"encoding/xml"
	"fmt"

	"github.com/pkg/errors"
	op "github.com/stillmatic/pummel/pkg/operators"
)

// CompoundPredicate combines two or more elements.
// booleanOperator can take one of the following logical (boolean) operators: and, or, xor or surrogate.
type CompoundPredicate struct {
	XMLName    xml.Name `xml:"CompoundPredicate"`
	Predicates []Predicate
	Operator   string `xml:"booleanOperator,attr"`
}

func (p *CompoundPredicate) String() string {
	var predicates []string
	for _, predicate := range p.Predicates {
		predicates = append(predicates, predicate.String())
	}
	return fmt.Sprintf("CompoundPredicate(%s(%s))", p.Operator, predicates)
}

func (cp *CompoundPredicate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	cp.XMLName = start.Name
	// set cp.Operator to the attribute with name "booleanOperator"
	for _, attr := range start.Attr {
		if attr.Name.Local == "booleanOperator" {
			cp.Operator = attr.Value
		}
	}
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch tt := t.(type) {
		case xml.StartElement:
			var p Predicate
			// TODO: make these constants
			switch tt.Name.Local {
			case "SimplePredicate":
				p = &SimplePredicate{}
			case "SimpleSetPredicate":
				p = &SimpleSetPredicate{}
			case "True":
				p = &TruePredicate{}
			case "False":
				p = &FalsePredicate{}
			default:
				return fmt.Errorf("unknown predicate type: %s", tt.Name.Local)
			}
			if err := d.DecodeElement(&p, &tt); err != nil {
				return err
			}
			cp.Predicates = append(cp.Predicates, p)
		case xml.EndElement:
			return nil
		}
	}
}

func (p *CompoundPredicate) Evaluate(features map[string]interface{}) (bool, bool, error) {
	count := 0
	for _, predicate := range p.Predicates {
		eval, ok, err := predicate.Evaluate(features)
		if err != nil {
			return false, false, errors.Wrapf(err, "Error when evaluating predicate %s", p)
		}
		switch p.Operator {
		case op.Operators.And:
			// The operator and indicates an evaluation to TRUE if all the predicates evaluate to TRUE.
			// if value is null (and this function returns False), we still return false
			if !eval || !ok {
				return false, true, nil
			}
		case op.Operators.Or, op.Operators.Surrogate:
			// The operator or indicates an evaluation to TRUE if one of the predicates evaluates to TRUE.
			// even if some values are missing, we can still return true
			// --
			// The operator surrogate allows for specifying surrogate predicates.
			// They are used for cases where a missing value appears in the evaluation of the parent predicate such that an alternative predicate is available.
			if eval {
				return true, true, nil
			}
		case op.Operators.Xor:
			// The operator xor indicates an evaluation to TRUE if an odd number of the predicates evaluates to TRUE and all others evaluate to FALSE.
			if eval {
				count++
			}
		default:
			return false, false, fmt.Errorf("unsupported compound predicate operator: %s", p.Operator)
		}
	}
	switch p.Operator {
	case op.Operators.And:
		return true, true, nil
	case op.Operators.Or, op.Operators.Surrogate:
		return false, false, nil
	case op.Operators.Xor:
		return count%2 == 1, true, nil
	default:
		return false, false, fmt.Errorf("unsupported compound predicate operator: %s", p.Operator)
	}
}
