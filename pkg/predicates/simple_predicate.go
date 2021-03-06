package predicates

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	op "github.com/stillmatic/pummel/pkg/operators"
)

// SimplePredicate defines a rule in the form of a simple boolean expression.
// The rule consists of field, operator (booleanOperator) for binary comparison, and value.
type SimplePredicate struct {
	XMLName  xml.Name `xml:"SimplePredicate"`
	Field    string   `xml:"field,attr"`
	Operator string   `xml:"operator,attr"`
	Value    string   `xml:"value,attr"`
}

func (p *SimplePredicate) String() string {
	return fmt.Sprintf("SimplePredicate(%s %s %s)", p.Field, p.Operator, p.Value)
}

func (p *SimplePredicate) Evaluate(features map[string]interface{}) (bool, bool, error) {
	featureValue, exists := features[p.Field]
	// for these, we always return ok = true (even if the feature is missing)
	switch p.Operator {
	case op.Operators.IsMissing:
		return featureValue == "" || featureValue == nil || !exists, true, nil
	case op.Operators.IsNotMissing:
		return featureValue != nil && exists && featureValue != "", true, nil
	}
	if !exists || featureValue == nil {
		// returns ok = false
		return false, false, nil
	}

	switch featureValue := featureValue.(type) {
	case int:
		return p.numericTrue(float64(featureValue))
	case float64:
		return p.numericTrue(featureValue)
	case string:
		return p.stringTrue(featureValue)
	case bool:
		predicateBool, err := strconv.ParseBool(p.Value)
		if err != nil {
			return false, false, errors.Wrapf(err, "failed to parse bool value %s", p.Value)
		}
		return p.boolTrue(featureValue == predicateBool)
	}

	return false, false, errors.Errorf("unsupported simplepredicate operator: %s for type %v", p.Operator, featureValue)
}

func (p *SimplePredicate) stringTrue(featureValue string) (bool, bool, error) {
	predicateValue := p.Value
	var b bool
	switch p.Operator {
	case op.Operators.Equal:
		b = strings.EqualFold(featureValue, predicateValue)
	case op.Operators.NotEqual:
		b = !strings.EqualFold(featureValue, predicateValue)
	case op.Operators.Gt, op.Operators.Gte, op.Operators.Lt, op.Operators.Lte:
		numValue, err := strconv.ParseFloat(predicateValue, 64)
		if err != nil {
			return false, false, errors.Wrapf(err, "failed to parse float value %s", predicateValue)
		}
		return p.numericTrue(numValue)
	default:
		return false, false, fmt.Errorf("unsupported stringTrue operator: %s", p.Operator)
	}
	return b, true, nil
}

func (p SimplePredicate) boolTrue(featureValue bool) (bool, bool, error) {
	predicateValue, _ := strconv.ParseBool(p.Value)
	var b bool
	switch p.Operator {
	case op.Operators.Equal:
		b = featureValue == predicateValue
	case op.Operators.NotEqual:
		b = featureValue != predicateValue
	default:
		return false, false, fmt.Errorf("unsupported boolTrue operator: %s", p.Operator)
	}
	return b, true, nil
}

func (p SimplePredicate) numericTrue(featureValue float64) (bool, bool, error) {
	// NB: could set p.Value to float64 vs parsing each time...
	predicateValue, _ := strconv.ParseFloat(p.Value, 64)
	var b bool
	switch p.Operator {
	case op.Operators.Equal:
		b = featureValue == predicateValue
	case op.Operators.NotEqual:
		b = featureValue != predicateValue
	case op.Operators.Lt:
		b = featureValue < predicateValue
	case op.Operators.Lte:
		b = featureValue <= predicateValue
	case op.Operators.Gt:
		b = featureValue > predicateValue
	case op.Operators.Gte:
		b = featureValue >= predicateValue
	default:
		return false, false, fmt.Errorf("unsupported numericTrue operator: %s", p.Operator)
	}
	return b, true, nil
}
