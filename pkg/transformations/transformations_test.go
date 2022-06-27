package transformations_test

import (
	"encoding/xml"
	"testing"

	"github.com/stillmatic/pummel/pkg/transformations"
	"github.com/stretchr/testify/assert"
)

func TestDerivedFieldFloat(t *testing.T) {
	derivedFieldFloatXML := []byte(` <DerivedField dataType="double" optype="continuous" name="float(price)">
	<FieldRef field="price" />
</DerivedField>`)
	var df transformations.DerivedField
	err := xml.Unmarshal(derivedFieldFloatXML, &df)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "float(price)", df.Name)
	input := map[string]interface{}{"price": "10"}
	output, err := df.Transform(input)
	assert.NoError(t, err)
	assert.Equal(t, float64(10), output)
}

func TestDerivedFieldComplex(t *testing.T) {
	derivedFieldFloatXML := []byte(`<DerivedField name="standardScaler(Age)" optype="continuous" dataType="double">
	<Apply function="/">
			<Apply function="-">
					<FieldRef field="Age"/>
					<Constant dataType="double">38.30279094260137</Constant>
			</Apply>
			<Constant dataType="double">13.010323102003973</Constant>
	</Apply>
</DerivedField>`)
	var df transformations.DerivedField
	err := xml.Unmarshal(derivedFieldFloatXML, &df)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "standardScaler(Age)", df.Name)
	input := map[string]interface{}{"Age": "10"}
	output, err := df.Transform(input)
	assert.NoError(t, err)
	assert.InEpsilon(t, float64(-2.17535337387), output, 0.01)
}

func BenchmarkDerivedFieldFloat(b *testing.B) {
	derivedFieldFloatXML := []byte(`<DerivedField name="standardScaler(Age)" optype="continuous" dataType="double">
	<Apply function="/">
			<Apply function="-">
					<FieldRef field="Age"/>
					<Constant dataType="double">38.30279094260137</Constant>
			</Apply>
			<Constant dataType="double">13.010323102003973</Constant>
	</Apply>
</DerivedField>`)
	var df transformations.DerivedField
	err := xml.Unmarshal(derivedFieldFloatXML, &df)
	if err != nil {
		b.Error(err)
	}
	input := map[string]interface{}{"Age": "10"}
	for i := 0; i < b.N; i++ {
		_, err := df.Transform(input)
		if err != nil {
			b.Error(err)
		}
	}
}

func TestConstant(t *testing.T) {
	constantXML := []byte(` 
	<Constant dataType="double">38.30279094260137</Constant>
	`)
	var c transformations.Constant
	err := xml.Unmarshal(constantXML, &c)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, c.Value)
	assert.Equal(t, 38.30279094260137, c.Value.(float64))
	assert.Equal(t, "double", c.DataType)
	input := map[string]interface{}{"price": "10"}
	output, err := c.Transform(input)
	assert.NoError(t, err)
	assert.Equal(t, 38.30279094260137, output.(float64))
}

var ApplyTestCases = []struct {
	XML      string
	Input    map[string]interface{}
	Function string
	Output   interface{}
	Error    error
}{
	{
		XML: `<Apply function="*">
			<Constant dataType="double">2</Constant>
			<FieldRef field="price" />
		</Apply>`,
		Input:    map[string]interface{}{"price": "10"},
		Function: "*",
		Output:   float64(20),
		Error:    nil,
	},
	{
		XML: `<Apply function="-">
			<Constant dataType="double">38.3027</Constant>
			<FieldRef field="Age" />
		</Apply>`,
		Input:    map[string]interface{}{"Age": "10"},
		Function: "-",
		Output:   float64(28.3027),
		Error:    nil,
	},
}

func TestApply(t *testing.T) {
	for _, tc := range ApplyTestCases {
		var a transformations.Apply
		err := xml.Unmarshal([]byte(tc.XML), &a)
		if err != nil {
			t.Error(err)
		}
		t.Log((*a.Children[0]).RequiredField())
		assert.Equal(t, 2, len(a.Children))
		assert.Equal(t, tc.Function, a.Function)
		output, err := a.Transform(tc.Input)
		assert.Equal(t, tc.Error, err)
		assert.Equal(t, tc.Output, output)
	}
}

func TestTransforms(t *testing.T) {
	localTransformsXMLA := []byte(`<DerivedField name="standardScaler(Age)" optype="continuous" dataType="double">
	<Apply function="/">
			<Apply function="-">
					<FieldRef field="Age"/>
					<Constant dataType="double">38.30279094260137</Constant>
			</Apply>
			<Constant dataType="double">13.010323102003973</Constant>
	</Apply>
</DerivedField>`)
	localTransformsXMLB := []byte(`<DerivedField name="standardScaler(Hours)" optype="continuous" dataType="double">
	<Apply function="/">
			<Apply function="-">
					<FieldRef field="Hours"/>
					<Constant dataType="double">40.56714060031596</Constant>
			</Apply>
			<Constant dataType="double">11.656262333704255</Constant>
	</Apply>
</DerivedField>`)
	localTransformsXMLC := []byte(`<DerivedField name="standardScaler(Income)" optype="continuous" dataType="double">
	<Apply function="/">
			<Apply function="-">
					<FieldRef field="Income"/>
					<Constant dataType="double">84404.87069510268</Constant>
			</Apply>
			<Constant dataType="double">69670.62788525566</Constant>
	</Apply>
</DerivedField>`)
	var transformA transformations.DerivedField
	var transformB transformations.DerivedField
	var transformC transformations.DerivedField
	err := xml.Unmarshal(localTransformsXMLA, &transformA)
	assert.NoError(t, err)
	assert.NotNil(t, transformA)
	assert.Equal(t, "standardScaler(Age)", transformA.RequiredField())
	err = xml.Unmarshal(localTransformsXMLB, &transformB)
	assert.NoError(t, err)
	assert.Equal(t, "standardScaler(Hours)", transformB.RequiredField())
	err = xml.Unmarshal(localTransformsXMLC, &transformC)
	assert.NoError(t, err)
	assert.Equal(t, "standardScaler(Income)", transformC.RequiredField())
	transforms := []transformations.DerivedField{transformA, transformB, transformC}
	assert.Equal(t, 3, len(transforms))

	input := map[string]interface{}{"Age": 10, "Hours": 20, "Income": 30}
	for i, tr := range transforms {
		t.Log(i)
		output, err := tr.Transform(input)
		assert.NoError(t, err)
		input[tr.RequiredField()] = output
		assert.NoError(t, err)
		assert.NotNil(t, output)
	}
	assert.Equal(t, 6, len(input))
}
