package miningschema_test

import (
	"encoding/xml"
	"testing"

	"github.com/stillmatic/pummel/pkg/miningschema"
	"github.com/stretchr/testify/assert"
)

func TestMiningSchema(t *testing.T) {
	msXML := []byte(`
<MiningSchema>
	<MiningField name="Probability_setosa" usageType="active"/>
	<MiningField name="Probability_versicolor" usageType="active"/>
	<MiningField name="Probability_virginica" usageType="active"/>
	<MiningField name="temperature" usageType="active"/>
	<MiningField name="cloudiness" usageType="active"/>
	<MiningField name="PollenIndex" usageType="target"/>
</MiningSchema>
  `)
	ms := &miningschema.MiningSchema{}
	err := xml.Unmarshal(msXML, ms)
	assert.NoError(t, err)
	res := ms.GetOutputField()
	assert.Equal(t, "PollenIndex", res)
	assert.Equal(t, 6, len(ms.MiningFields))
}
