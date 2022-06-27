package fields_test

import (
	"encoding/xml"
	"testing"

	"github.com/stillmatic/pummel/pkg/fields"
	"github.com/stretchr/testify/assert"
)

func TestParseOutputFields(t *testing.T) {
	outputXML := []byte(`
	<Output>
	<OutputField name="Predicted_Survived" feature="predictedValue"/>
	<OutputField name="Probability_0" optype="continuous" dataType="double" feature="probability" value="0"/>
	<OutputField name="Probability_1" optype="continuous" dataType="double" feature="probability" value="1"/>
</Output>
	`)
	var output *fields.Outputs
	err := xml.Unmarshal(outputXML, &output)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(output.OutputFields))
	assert.Equal(t, "Predicted_Survived", output.OutputFields[0].Name)
}
