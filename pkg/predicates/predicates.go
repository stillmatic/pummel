package predicates

import (
	"encoding/xml"

	"gopkg.in/guregu/null.v4"
)

type Predicate interface {
	Evaluate(map[string]interface{}) (null.Bool, error)
	String() string
}

// TruePredicate always returns true.
type TruePredicate struct {
	XMLName xml.Name `xml:"True"`
}

// FalsePredicate always returns false.
type FalsePredicate struct {
	XMLName xml.Name `xml:"False"`
}

func (p *TruePredicate) Evaluate(map[string]interface{}) (null.Bool, error) {
	return null.BoolFrom(true), nil
}

func (p *FalsePredicate) Evaluate(map[string]interface{}) (null.Bool, error) {
	return null.BoolFrom(false), nil
}

func (p *TruePredicate) String() string {
	return "True"
}

func (p *FalsePredicate) String() string {
	return "False"
}
