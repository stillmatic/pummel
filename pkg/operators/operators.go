package operators

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
