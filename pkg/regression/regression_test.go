package regression_test

import (
	"encoding/xml"
	"testing"

	"github.com/stillmatic/pummel/pkg/regression"
	"github.com/stretchr/testify/assert"
)

var linearRegressionXML = []byte(`<PMML xmlns="https://www.dmg.org/PMML-4_1" version="4.1">
<Header copyright="DMG.org"/>
<DataDictionary numberOfFields="4">
  <DataField name="age" optype="continuous" dataType="double"/>
  <DataField name="salary" optype="continuous" dataType="double"/>
  <DataField name="car_location" optype="categorical" dataType="string">
	<Value value="carpark"/>
	<Value value="street"/>
  </DataField>
  <DataField name="number_of_claims" optype="continuous" dataType="integer"/>
</DataDictionary>
<RegressionModel modelName="Sample for linear regression" functionName="regression" algorithmName="linearRegression" targetFieldName="number_of_claims">
  <MiningSchema>
	<MiningField name="age"/>
	<MiningField name="salary"/>
	<MiningField name="car_location"/>
	<MiningField name="number_of_claims" usageType="predicted"/>
  </MiningSchema>
  <RegressionTable intercept="132.37">
	<NumericPredictor name="age" exponent="1" coefficient="7.1"/>
	<NumericPredictor name="salary" exponent="1" coefficient="0.01"/>
	<CategoricalPredictor name="car_location" value="carpark" coefficient="41.1"/>
	<CategoricalPredictor name="car_location" value="street" coefficient="325.03"/>
  </RegressionTable>
</RegressionModel>
</PMML>`)

func TestRegression(t *testing.T) {
	var rm regression.RegressionModel
	err := xml.Unmarshal(linearRegressionXML, &rm)
	assert.NoError(t, err)
}
