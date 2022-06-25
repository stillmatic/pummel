package predicates_test

import (
	"encoding/xml"
	"testing"

	"github.com/stillmatic/pummel/pkg/predicates"
	"github.com/stretchr/testify/assert"
)

func TestTruePredicate(t *testing.T) {
	truePredicateBytes := []byte(`<True/>`)
	var tp predicates.TruePredicate
	err := xml.Unmarshal(truePredicateBytes, &tp)
	assert.NoError(t, err)
	res, err := tp.Evaluate(map[string]interface{}{"age": 30})
	assert.NoError(t, err)
	assert.True(t, res.Valid)
	assert.True(t, res.ValueOrZero())
}

func TestFalsePredicate(t *testing.T) {
	falsePredicateBytes := []byte(`<False/>`)
	var fp predicates.FalsePredicate
	err := xml.Unmarshal(falsePredicateBytes, &fp)
	assert.NoError(t, err)
	res, err := fp.Evaluate(map[string]interface{}{"age": 30})
	assert.NoError(t, err)
	assert.True(t, res.Valid)
	assert.False(t, res.ValueOrZero())
}

type predicateInput struct {
	bytes    []byte
	features map[string]interface{}
}

type predicateTest struct {
	inputs   predicateInput
	expected bool
}

var simplePredicateTests = []predicateTest{
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="lessThan" value="30"/>`), map[string]interface{}{"age": 31}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="lessThan" value="30"/>`), map[string]interface{}{"age": 29.4}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="lessThan" value="30"/>`), map[string]interface{}{"age": 30}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="lessOrEqual" value="30"/>`), map[string]interface{}{"age": 29.4}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="lessOrEqual" value="30"/>`), map[string]interface{}{"age": 31}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="lessOrEqual" value="30"/>`), map[string]interface{}{"age": 30}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="greaterThan" value="30"/>`), map[string]interface{}{"age": 31}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="greaterThan" value="30"/>`), map[string]interface{}{"age": 29.4}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="greaterThan" value="30"/>`), map[string]interface{}{"age": 30}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="greaterThan" value="30"/>`), map[string]interface{}{"age": "29.9"}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="greaterOrEqual" value="30"/>`), map[string]interface{}{"age": 29.4}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="greaterOrEqual" value="30"/>`), map[string]interface{}{"age": 31}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="greaterOrEqual" value="30"/>`), map[string]interface{}{"age": 30}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="equal" value="30"/>`), map[string]interface{}{"age": 32}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="equal" value="30"/>`), map[string]interface{}{"age": 30}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="equal" value="30"/>`), map[string]interface{}{"age": 30.1}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="equal" value="30.0"/>`), map[string]interface{}{"age": 30}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="equal" value="abc"/>`), map[string]interface{}{"age": "abc"}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="young" operator="equal" value="true"/>`), map[string]interface{}{"young": false}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="young" operator="equal" value="true"/>`), map[string]interface{}{"young": true}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="notEqual" value="30.0"/>`), map[string]interface{}{"age": 30}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="notEqual" value="30.0"/>`), map[string]interface{}{"age": 31}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="young" operator="notEqual" value="true"/>`), map[string]interface{}{"young": false}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="isMissing"/>`), map[string]interface{}{"age": 31}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="isMissing"/>`), map[string]interface{}{}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="isMissing"/>`), map[string]interface{}{"height": 61}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="isMissing"/>`), map[string]interface{}{"age": ""}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="isNotMissing"/>`), map[string]interface{}{"age": 31}}, true},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="isNotMissing"/>`), map[string]interface{}{}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="isNotMissing"/>`), map[string]interface{}{"height": 61}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="isNotMissing"/>`), map[string]interface{}{"age": ""}}, false},
}

var simpleSetPredicateTests = []predicateTest{
	{predicateInput{[]byte(`
		<SimpleSetPredicate field="age" booleanOperator="isIn">
			<Array type="string">29 30</Array>
		</SimpleSetPredicate>
	`), map[string]interface{}{"age": "31"}}, false},
	{predicateInput{[]byte(`		
	<SimpleSetPredicate field="age" booleanOperator="isIn">
		<Array type="string">29 30</Array>
	</SimpleSetPredicate>
	`), map[string]interface{}{"age": "30"}}, true},
	{predicateInput{[]byte(`		
	<SimpleSetPredicate field="age" booleanOperator="isNotIn">
		<Array type="string">29 30</Array>
	</SimpleSetPredicate>
	`), map[string]interface{}{"age": "31"}}, true},
	{predicateInput{[]byte(`		
	<SimpleSetPredicate field="age" booleanOperator="isNotIn">
		<Array type="string">29 30</Array>
	</SimpleSetPredicate>
	`), map[string]interface{}{"age": "30"}}, false},
}

func TestSimplePredicates(t *testing.T) {
	for _, test := range simplePredicateTests {
		var sp predicates.SimplePredicate
		err := xml.Unmarshal(test.inputs.bytes, &sp)
		if err != nil {
			t.Fatal("could not unmarshal xml")
		}
		res, err := sp.Evaluate(test.inputs.features)
		assert.NoError(t, err)
		assert.True(t, res.Valid, "expected %s %v to be true", sp.Operator, test.inputs.features)
		assert.Equal(t, test.expected, res.ValueOrZero())
		if res.ValueOrZero() != test.expected {
			t.Errorf("error comparing %v versus %v, with %s", sp.Value, test.inputs.features, sp.Operator)
		}
	}
}

func TestSimpleSetPredicates(t *testing.T) {
	for _, test := range simpleSetPredicateTests {
		var sp predicates.SimpleSetPredicate
		err := xml.Unmarshal(test.inputs.bytes, &sp)
		if err != nil {
			t.Fatal("could not unmarshal xml")
		}
		res, err := sp.Evaluate(test.inputs.features)
		assert.NoError(t, err)
		assert.True(t, res.Valid)
		assert.Equal(t, test.expected, res.ValueOrZero())
		if res.ValueOrZero() != test.expected {
			t.Errorf("error comparing %v versus %v, with %s", sp.Values, test.inputs.features, sp.Operator)
		}
	}
}

var compoundPredicateTests = []predicateTest{
	{predicateInput{[]byte(`
	<CompoundPredicate booleanOperator="or">
		<SimplePredicate field="f" operator="equal" value="A"/>
		<SimplePredicate field="f" operator="equal" value="B"/>
	</CompoundPredicate>
	`), map[string]interface{}{"f": "A"}}, true},
	{predicateInput{[]byte(`
	<CompoundPredicate booleanOperator="or">
		<SimplePredicate field="f" operator="equal" value="A"/>
		<SimplePredicate field="f" operator="equal" value="B"/>
	</CompoundPredicate>
	`), map[string]interface{}{"f": "C"}}, false},
	{predicateInput{[]byte(`
	<CompoundPredicate booleanOperator="and">
		<SimplePredicate field="f" operator="equal" value="A"/>
		<SimplePredicate field="f" operator="equal" value="B"/>
	</CompoundPredicate>
	`), map[string]interface{}{"f": "C"}}, false},
	{predicateInput{[]byte(`
	<CompoundPredicate booleanOperator="and">
		<SimplePredicate field="f" operator="greaterThan" value="10"/>
		<SimplePredicate field="f" operator="lessThan" value="100"/>
	</CompoundPredicate>
	`), map[string]interface{}{"f": 30}}, true},
	{predicateInput{[]byte(`
	<CompoundPredicate booleanOperator="xor">
		<SimplePredicate field="f" operator="greaterThan" value="10"/>
		<SimplePredicate field="f" operator="lessThan" value="100"/>
	</CompoundPredicate>
	`), map[string]interface{}{"f": 30}}, false},
	{predicateInput{[]byte(`
	<CompoundPredicate booleanOperator="xor">
		<SimplePredicate field="f" operator="greaterThan" value="10"/>
		<SimplePredicate field="f" operator="lessThan" value="100"/>
	</CompoundPredicate>
	`), map[string]interface{}{"f": 5}}, true},
	{predicateInput{[]byte(`
	<CompoundPredicate booleanOperator="and">
		<SimpleSetPredicate field="age" booleanOperator="isNotIn">
			<Array type="string">29 30</Array>
		</SimpleSetPredicate>
		<SimplePredicate field="f" operator="lessThan" value="100"/>
	</CompoundPredicate>
	`), map[string]interface{}{"f": 5, "age": "31"}}, true},
	{predicateInput{[]byte(`
	<CompoundPredicate booleanOperator="and">
		<True/>
		<SimplePredicate field="f" operator="lessThan" value="100"/>
	</CompoundPredicate>
	`), map[string]interface{}{"f": 5, "age": "31"}}, true},
	{predicateInput{[]byte(`
	<CompoundPredicate booleanOperator="and">
		<False/>
		<SimplePredicate field="f" operator="lessThan" value="100"/>
	</CompoundPredicate>
	`), map[string]interface{}{"f": 5, "age": "31"}}, false},
}

func TestCompoundPredicates(t *testing.T) {
	for _, test := range compoundPredicateTests {
		var sp predicates.CompoundPredicate
		err := xml.Unmarshal(test.inputs.bytes, &sp)
		if err != nil {
			t.Fatal("could not unmarshal xml", err)
		}
		res, err := sp.Evaluate(test.inputs.features)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, res.ValueOrZero())
		if res.ValueOrZero() != test.expected {
			t.Errorf("error comparing %s with %s", test.inputs.features, sp.Operator)
		}
	}
}

var simplePredicateTestsMissing = []predicateTest{
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="lessThan" value="30"/>`), map[string]interface{}{}}, false},
	{predicateInput{[]byte(`<SimplePredicate field="age" operator="lessThan" value="30"/>`), map[string]interface{}{"alphabet": 29.4}}, false},
}

var simpleSetPredicateTestsMissing = []predicateTest{
	{predicateInput{[]byte(`
		<SimpleSetPredicate field="age" booleanOperator="isIn">
			<Array type="string">29 30</Array>
		</SimpleSetPredicate>
	`), map[string]interface{}{}}, false},
	{predicateInput{[]byte(`		
	<SimpleSetPredicate field="age" booleanOperator="isIn">
		<Array type="string">29 30</Array>
	</SimpleSetPredicate>
	`), map[string]interface{}{"height": "30"}}, false},
}

var compoundPredicateTestsMissing = []predicateTest{
	{predicateInput{[]byte(`
	<CompoundPredicate booleanOperator="or">
		<SimplePredicate field="f" operator="equal" value="A"/>
		<SimplePredicate field="f" operator="equal" value="B"/>
	</CompoundPredicate>
	`), map[string]interface{}{"g": "A"}}, false},
	{predicateInput{[]byte(`
	<CompoundPredicate booleanOperator="or">
		<SimplePredicate field="g" operator="equal" value="A"/>
		<SimplePredicate field="f" operator="equal" value="B"/>
	</CompoundPredicate>
	`), map[string]interface{}{"g": "A"}}, true},
	{predicateInput{[]byte(`
	<CompoundPredicate booleanOperator="surrogate">
		<SimplePredicate field="f" operator="equal" value="A"/>
		<SimplePredicate field="f" operator="equal" value="B"/>
		<True/>
	</CompoundPredicate>
	`), map[string]interface{}{"g": "A"}}, true},
	{predicateInput{[]byte(`
	<CompoundPredicate booleanOperator="surrogate">
		<SimplePredicate field="f" operator="equal" value="A"/>
		<SimplePredicate field="f" operator="equal" value="B"/>
		<False/>
	</CompoundPredicate>
	`), map[string]interface{}{"g": "A"}}, false},
}

func TestSimplePredicatesMissing(t *testing.T) {
	for _, test := range simplePredicateTestsMissing {
		var sp predicates.SimplePredicate
		err := xml.Unmarshal(test.inputs.bytes, &sp)
		if err != nil {
			t.Fatal("could not unmarshal xml")
		}
		res, err := sp.Evaluate(test.inputs.features)
		assert.NoError(t, err)
		assert.Equal(t, res.Valid, test.expected)
	}
}

func TestSimpleSetPredicatesMissing(t *testing.T) {
	for _, test := range simpleSetPredicateTestsMissing {
		var sp predicates.SimpleSetPredicate
		err := xml.Unmarshal(test.inputs.bytes, &sp)
		if err != nil {
			t.Fatal("could not unmarshal xml")
		}
		res, err := sp.Evaluate(test.inputs.features)
		assert.NoError(t, err)
		assert.Equal(t, res.Valid, test.expected)
	}
}

func TestCompoundPredicatesMissing(t *testing.T) {
	for _, test := range compoundPredicateTestsMissing {
		var sp predicates.CompoundPredicate
		err := xml.Unmarshal(test.inputs.bytes, &sp)
		if err != nil {
			t.Fatal("could not unmarshal xml", err)
		}
		res, err := sp.Evaluate(test.inputs.features)
		assert.NoError(t, err)
		assert.True(t, res.Valid)
		assert.Equal(t, res.ValueOrZero(), test.expected)
	}
}

//nolint
func BenchmarkCompoundPredicates(b *testing.B) {
	var sp predicates.CompoundPredicate
	err := xml.Unmarshal([]byte(`
	<CompoundPredicate booleanOperator="and">
		<SimplePredicate field="f" operator="equal" value="A"/>
		<SimplePredicate field="f" operator="equal" value="B"/>
	</CompoundPredicate>
	`), &sp)
	if err != nil {
		b.Fatal("could not unmarshal xml", err)
	}
	for i := 0; i < b.N; i++ {
		sp.Evaluate(map[string]interface{}{"f": "A"})
	}
}
