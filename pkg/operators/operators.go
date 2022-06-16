package operators

import "github.com/stillmatic/pummel/pkg/constraints"

var Operators = struct {
	IsIn         string
	IsNotIn      string
	IsMissing    string
	IsNotMissing string
	Equal        string
	NotEqual     string
	Lt           string
	Lte          string
	Gt           string
	Gte          string
	Or           string
	And          string
	Xor          string
	Surrogate    string
}{
	IsIn:         "isIn",
	IsNotIn:      "isNotIn",
	IsMissing:    "isMissing",
	IsNotMissing: "isNotMissing",
	Equal:        "equal",
	NotEqual:     "notEqual",
	Lt:           "lessThan",
	Lte:          "lessOrEqual",
	Gt:           "greaterThan",
	Gte:          "greaterOrEqual",
	Or:           "or",
	And:          "and",
	Xor:          "xor",
	Surrogate:    "surrogate",
}

func Equal[T comparable](a, b T) bool {
	return a == b
}

func NotEqual[T comparable](a, b T) bool {
	return a != b
}

func LessThan[T constraints.Ordered](a, b T) bool {
	return a < b
}

func LessOrEqual[T constraints.Ordered](a, b T) bool {
	return a <= b
}

func GreaterThan[T constraints.Ordered](a, b T) bool {
	return a > b
}

func GreaterOrEqual[T constraints.Ordered](a, b T) bool {
	return a >= b
}

// Contains checks if an element a is in an array b
func Contains[T comparable](a T, b []T) bool {
	for _, v := range b {
		if v == a {
			return true
		}
	}
	return false
}
