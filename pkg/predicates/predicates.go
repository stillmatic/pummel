package predicates

import (
	"encoding/xml"

	"gopkg.in/guregu/null.v4"
)

type Predicate interface {
	True(map[string]interface{}) (null.Bool, error)
}

// TruePredicate always returns true.
type TruePredicate struct {
	XMLName xml.Name `xml:"True"`
}

// FalsePredicate always returns false.
type FalsePredicate struct {
	XMLName xml.Name `xml:"False"`
}

func (p TruePredicate) True(map[string]interface{}) (null.Bool, error) {
	return null.BoolFrom(true), nil
}

func (p FalsePredicate) True(map[string]interface{}) (null.Bool, error) {
	return null.BoolFrom(false), nil
}
