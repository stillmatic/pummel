package predicates

import (
	"encoding/xml"
	"fmt"

	"github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
	op "github.com/stillmatic/pummel/pkg/operators"
	"gopkg.in/guregu/null.v4"
)

// SimpleSetPredicate checks whether a field value is element of a set.
// The set of values is specified by the array.
type SimpleSetPredicate struct {
	XMLName  xml.Name `xml:"SimpleSetPredicate"`
	Field    string   `xml:"field,attr"`
	Operator string   `xml:"booleanOperator,attr"`
	Values   string   `xml:"Array"`
}

func (p SimpleSetPredicate) True(features map[string]interface{}) (null.Bool, error) {
	values, err := shellwords.Parse(p.Values)
	if err != nil {
		// returns a null bool if we can't parse this predicate
		return null.BoolFromPtr(nil), errors.Wrapf(err, "failed to parse values for SimpleSetPredicate %s", p.Field)
	}
	featureVal, exists := features[p.Field]
	if !exists {
		return null.BoolFromPtr(nil), nil
	}

	switch p.Operator {
	case op.Operators.IsIn:
		for _, value := range values {
			if value == featureVal {
				return null.BoolFrom(true), nil
			}
		}
		return null.BoolFrom(false), nil
	case op.Operators.IsNotIn:
		for _, value := range values {
			if value == featureVal {
				return null.BoolFrom(false), nil
			}
		}
		return null.BoolFrom(true), nil
	}
	return null.BoolFromPtr(nil), fmt.Errorf("unsupported operator: %s", p.Operator)
}
