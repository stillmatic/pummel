package predicates

import (
	"encoding/xml"
)

// Predicates are logical expressions which define the rule for choosing the Node
// If the result is true, evaluate the node
// If the result is false, do not evaluate the node
// If the result is null, use the tree's missing value strategy
// returns val, ok, err
type Predicate interface {
	Evaluate(map[string]interface{}) (bool, bool, error)
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

func (p *TruePredicate) Evaluate(map[string]interface{}) (bool, bool, error) {
	return true, true, nil
}

func (p *FalsePredicate) Evaluate(map[string]interface{}) (bool, bool, error) {
	return false, true, nil
}

func (p *TruePredicate) String() string {
	return "True"
}

func (p *FalsePredicate) String() string {
	return "False"
}
