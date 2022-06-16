package predicates

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/mattn/go-shellwords"
	op "github.com/stillmatic/pummel/pkg/operators"
)

type Predicate interface {
	True(map[string]interface{}) bool
}

// SimplePredicate defines a rule in the form of a simple boolean expression.
// The rule consists of field, operator (booleanOperator) for binary comparison, and value.
type SimplePredicate struct {
	XMLName  xml.Name `xml:"SimplePredicate"`
	Field    string   `xml:"field,attr"`
	Operator string   `xml:"operator,attr"`
	Value    string   `xml:"value,attr"`
}

// SimpleSetPredicate checks whether a field value is element of a set.
// The set of values is specified by the array.
type SimpleSetPredicate struct {
	XMLName  xml.Name `xml:"SimpleSetPredicate"`
	Field    string   `xml:"field,attr"`
	Operator string   `xml:"booleanOperator,attr"`
	Values   string   `xml:"Array"`
}

// TruePredicate always returns true.
type TruePredicate struct {
	XMLName xml.Name `xml:"True"`
}

// FalsePredicate always returns false.
type FalsePredicate struct {
	XMLName xml.Name `xml:"False"`
}

// CompoundPredicate combines two or more elements.
// booleanOperator can take one of the following logical (boolean) operators: and, or, xor or surrogate.
type CompoundPredicate struct {
	XMLName    xml.Name `xml:"CompoundPredicate"`
	Predicates []Predicate
	Operator   string `xml:"booleanOperator,attr"`
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

func (p TruePredicate) True(map[string]interface{}) bool {
	return true
}

func (p FalsePredicate) True(map[string]interface{}) bool {
	return false
}

func (p SimpleSetPredicate) True(features map[string]interface{}) bool {
	values, _ := shellwords.Parse(p.Values)
	featureVal := features[p.Field].(string)
	exists := op.Contains(featureVal, values)

	switch p.Operator {
	case op.Operators.IsIn:
		return exists
	case op.Operators.IsNotIn:
		return !exists
	}
	return false
}

func (p SimplePredicate) True(features map[string]interface{}) bool {
	if p.Operator == op.Operators.IsMissing {
		featureValue, exists := features[p.Field]
		return featureValue == "" || featureValue == nil || !exists
	}
	if p.Operator == op.Operators.IsNotMissing {
		featureValue, exists := features[p.Field]
		return featureValue != nil && exists && featureValue != ""
	}

	switch featureValue := features[p.Field].(type) {
	case int:
		return p.numericTrue(float64(featureValue))
	case float64:
		return p.numericTrue(featureValue)
	case string:
		if p.Operator == op.Operators.Equal {
			return p.Value == features[p.Field]
		}
		numericFeatureValue, err := strconv.ParseFloat(featureValue, 64)
		if err == nil {
			return p.numericTrue(numericFeatureValue)
		}
	case bool:
		if p.Operator == op.Operators.Equal {
			predicateBool, _ := strconv.ParseBool(p.Value)
			return predicateBool == features[p.Field]
		}
		if p.Operator == op.Operators.NotEqual {
			predicateBool, _ := strconv.ParseBool(p.Value)
			return predicateBool != features[p.Field]
		}
	}

	return false
}

func (p SimplePredicate) numericTrue(featureValue float64) bool {
	// NB: could set predicate value to float64 vs parsing each time...
	predicateValue, _ := strconv.ParseFloat(p.Value, 64)

	switch p.Operator {
	case op.Operators.Equal:
		return op.Equal(featureValue, predicateValue)
	case op.Operators.NotEqual:
		return featureValue != predicateValue
	case op.Operators.Lt:
		return featureValue < predicateValue
	case op.Operators.Lte:
		return featureValue <= predicateValue
	case op.Operators.Gt:
		return featureValue > predicateValue
	case op.Operators.Gte:
		return featureValue >= predicateValue
	}
	// NB: consider returning error
	return false
}

func (p CompoundPredicate) True(features map[string]interface{}) bool {
	switch p.Operator {
	case op.Operators.And:
		// The operator and indicates an evaluation to TRUE if all the predicates evaluate to TRUE.
		for _, predicate := range p.Predicates {
			if !predicate.True(features) {
				return false
			}
		}
		return true
	case op.Operators.Or:
		// The operator or indicates an evaluation to TRUE if one of the predicates evaluates to TRUE.
		for _, predicate := range p.Predicates {
			fmt.Println("predicate:", predicate)
			if predicate.True(features) {
				return true
			}
		}
		return false
	case op.Operators.Xor:
		// The operator xor indicates an evaluation to TRUE if an odd number of the predicates evaluates to TRUE and all others evaluate to FALSE.
		count := 0
		for _, predicate := range p.Predicates {
			if predicate.True(features) {
				count++
			}
		}
		return count%2 == 1
	case op.Operators.Surrogate:
		// The operator surrogate allows for specifying surrogate predicates.
		// They are used for cases where a missing value appears in the evaluation of the parent predicate such that an alternative predicate is available.

		// this is not implemented yet, need to handle missing values better.
		return false
	}
	return false
}
