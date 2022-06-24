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
	input := interface{}("10")
	output, err := df.Transform(input)
	assert.NoError(t, err)
	assert.Equal(t, float64(10), output)

}
