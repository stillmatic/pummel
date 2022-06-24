package regression_test

import (
	"encoding/xml"
	"testing"

	"github.com/stillmatic/pummel/pkg/model"
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

var polynomialRegressionXML = []byte(`<PMML xmlns="https://www.dmg.org/PMML-4_4" version="4.4">
<Header copyright="DMG.org"/>
<DataDictionary numberOfFields="3">
  <DataField name="salary" optype="continuous" dataType="double"/>
  <DataField name="car_location" optype="categorical" dataType="string">
	<Value value="carpark"/>
	<Value value="street"/>
  </DataField>
  <DataField name="number_of_claims" optype="continuous" dataType="integer"/>
</DataDictionary>
<RegressionModel functionName="regression" modelName="Sample for stepwise polynomial regression" algorithmName="stepwisePolynomialRegression" targetFieldName="number_of_claims">
  <MiningSchema>
	<MiningField name="salary"/>
	<MiningField name="car_location"/>
	<MiningField name="number_of_claims" usageType="target"/>
  </MiningSchema>
  <RegressionTable intercept="3216.38">
	<NumericPredictor name="salary" exponent="1" coefficient="-0.08"/>
	<NumericPredictor name="salary" exponent="2" coefficient="9.54E-7"/>
	<NumericPredictor name="salary" exponent="3" coefficient="-2.67E-12"/>
	<CategoricalPredictor name="car_location" value="carpark" coefficient="93.78"/>
	<CategoricalPredictor name="car_location" value="street" coefficient="288.75"/>
  </RegressionTable>
</RegressionModel>
</PMML>
`)

var logisticRegressionXML = []byte(`<PMML xmlns="https://www.dmg.org/PMML-4_4" version="4.4">
<Header copyright="DMG.org"/>
<DataDictionary numberOfFields="3">
  <DataField name="x1" optype="continuous" dataType="double"/>
  <DataField name="x2" optype="continuous" dataType="double"/>
  <DataField name="y" optype="categorical" dataType="string">
	<Value value="yes"/>
	<Value value="no"/>
  </DataField>
</DataDictionary>
<RegressionModel functionName="classification" modelName="Sample for stepwise polynomial regression" algorithmName="stepwisePolynomialRegression" normalizationMethod="softmax" targetFieldName="y">
  <MiningSchema>
	<MiningField name="x1"/>
	<MiningField name="x2"/>
	<MiningField name="y" usageType="target"/>
  </MiningSchema>
  <RegressionTable targetCategory="no" intercept="125.56601826">
	<NumericPredictor name="x1" coefficient="-28.6617384"/>
	<NumericPredictor name="x2" coefficient="-20.42027426"/>
  </RegressionTable>
  <RegressionTable targetCategory="yes" intercept="0"/>
</RegressionModel>
</PMML>
`)

var complexClassificationXML = []byte(`
<PMML xmlns="https://www.dmg.org/PMML-4_4" version="4.4">
  <Header copyright="DMG.org"/>
  <DataDictionary numberOfFields="5">
    <DataField name="age" optype="continuous" dataType="double"/>
    <DataField name="work" optype="continuous" dataType="double"/>
    <DataField name="sex" optype="categorical" dataType="string">
      <Value value="0"/>
      <Value value="1"/>
    </DataField>
    <DataField name="minority" optype="categorical" dataType="integer">
      <Value value="0"/>
      <Value value="1"/>
    </DataField>
    <DataField name="jobcat" optype="categorical" dataType="string">
      <Value value="clerical"/>
      <Value value="professional"/>
      <Value value="trainee"/>
      <Value value="skilled"/>
    </DataField>
  </DataDictionary>
  <RegressionModel modelName="Sample for logistic regression" functionName="classification" algorithmName="logisticRegression" normalizationMethod="softmax" targetFieldName="jobcat">
    <MiningSchema>
      <MiningField name="age"/>
      <MiningField name="work"/>
      <MiningField name="sex"/>
      <MiningField name="minority"/>
      <MiningField name="jobcat" usageType="target"/>
    </MiningSchema>
    <RegressionTable intercept="46.418" targetCategory="clerical">
      <NumericPredictor name="age" exponent="1" coefficient="-0.132"/>
      <NumericPredictor name="work" exponent="1" coefficient="7.867E-02"/>
      <CategoricalPredictor name="sex" value="0" coefficient="-20.525"/>
      <CategoricalPredictor name="sex" value="1" coefficient="0.5"/>
      <CategoricalPredictor name="minority" value="0" coefficient="-19.054"/>
      <CategoricalPredictor name="minority" value="1" coefficient="0"/>
    </RegressionTable>
    <RegressionTable intercept="51.169" targetCategory="professional">
      <NumericPredictor name="age" exponent="1" coefficient="-0.302"/>
      <NumericPredictor name="work" exponent="1" coefficient="0.155"/>
      <CategoricalPredictor name="sex" value="0" coefficient="-21.389"/>
      <CategoricalPredictor name="sex" value="1" coefficient="0.1"/>
      <CategoricalPredictor name="minority" value="0" coefficient="-18.443"/>
      <CategoricalPredictor name="minority" value="1" coefficient="0"/>
    </RegressionTable>
    <RegressionTable intercept="25.478" targetCategory="trainee">
      <NumericPredictor name="age" exponent="1" coefficient="-0.154"/>
      <NumericPredictor name="work" exponent="1" coefficient="0.266"/>
      <CategoricalPredictor name="sex" value="0" coefficient="-2.639"/>
      <CategoricalPredictor name="sex" value="1" coefficient="0.8"/>
      <CategoricalPredictor name="minority" value="0" coefficient="-19.821"/>
      <CategoricalPredictor name="minority" value="1" coefficient="0.2"/>
    </RegressionTable>
    <RegressionTable intercept="0.0" targetCategory="skilled"/>
  </RegressionModel>
</PMML>`)

var interactionTermsXML = []byte(`
<PMML xmlns="https://www.dmg.org/PMML-4_4" version="4.4">
  <Header copyright="DMG.org"/>
  <DataDictionary numberOfFields="4">
    <DataField name="age" optype="continuous" dataType="double"/>
    <DataField name="work" optype="continuous" dataType="double"/>
    <DataField name="sex" optype="categorical" dataType="string">
      <Value value="male"/>
      <Value value="female"/>
    </DataField>
    <DataField name="y" optype="continuous" dataType="double"/>
  </DataDictionary>
  <RegressionModel modelName="Sample for interaction terms" functionName="regression" targetFieldName="y">
    <MiningSchema> 
      <MiningField name="age" optype="continuous"/> 
      <MiningField name="work" optype="continuous"/> 
      <MiningField name="sex" optype="categorical"/> 
      <MiningField name="y" optype="continuous" usageType="target"/> 
    </MiningSchema>
    <RegressionTable intercept="2.1">
      <CategoricalPredictor name="sex" value="female" coefficient="-20.525"/>
      <PredictorTerm coefficient="-0.1">
        <FieldRef field="age"/>
        <FieldRef field="work"/>
      </PredictorTerm>
    </RegressionTable>
  </RegressionModel>
</PMML>
`)

func TestLinearRegression(t *testing.T) {
	var prm model.PMMLRegressionModel
	err := xml.Unmarshal(linearRegressionXML, &prm)
	assert.Equal(t, 4, len(prm.DataDictionary.DataFields))
	rm := prm.RegressionModel
	assert.NoError(t, err)
	assert.Equal(t, "Sample for linear regression", rm.ModelName)
	assert.Equal(t, "regression", rm.FunctionName)
	assert.Equal(t, float64(132.37), rm.RegressionTables[0].Intercept)
	inputs := map[string]interface{}{
		"age":          float64(30),
		"salary":       float64(1000),
		"car_location": "carpark",
	}
	out, err := rm.Evaluate(inputs)
	assert.NoError(t, err)
	assert.Equal(t, float64(396.47), out["number_of_claims"])
}

func TestPolynomialRegression(t *testing.T) {
	var prm model.PMMLRegressionModel
	err := xml.Unmarshal(polynomialRegressionXML, &prm)
	assert.Equal(t, 3, len(prm.DataDictionary.DataFields))
	rm := prm.RegressionModel
	assert.NoError(t, err)
	assert.Equal(t, "Sample for stepwise polynomial regression", rm.ModelName)
	assert.Equal(t, "regression", rm.FunctionName)
	assert.Equal(t, float64(3216.38), rm.RegressionTables[0].Intercept)
	inputs := map[string]interface{}{
		"salary":       float64(1000),
		"car_location": "carpark",
	}
	out, err := rm.Evaluate(inputs)
	assert.NoError(t, err)
	assert.InEpsilon(t, float64(3231.11133), out["number_of_claims"], 0.01)
}

func TestLogisticRegression(t *testing.T) {
	var prm model.PMMLRegressionModel
	err := xml.Unmarshal(logisticRegressionXML, &prm)
	assert.Equal(t, 3, len(prm.DataDictionary.DataFields))
	rm := prm.RegressionModel
	assert.NoError(t, err)
	assert.Equal(t, "Sample for stepwise polynomial regression", rm.ModelName)
	assert.Equal(t, "classification", rm.FunctionName)
	assert.Equal(t, "no", rm.RegressionTables[0].TargetCategory)
	assert.Equal(t, float64(125.56601826), rm.RegressionTables[0].Intercept)
	inputs := map[string]interface{}{
		"x1": float64(1),
		"x2": float64(2),
	}
	out, err := rm.Evaluate(inputs)
	assert.NoError(t, err)
	assert.InEpsilon(t, 1, out["no"], 0.01)
	assert.InEpsilon(t, 4.485e-25, out["yes"], 0.01)
	// should add to 1
	assert.InEpsilon(t, 1, out["no"].(float64)+out["yes"].(float64), 0.01)
}

func TestComplexClassifcation(t *testing.T) {
	var prm model.PMMLRegressionModel
	err := xml.Unmarshal(complexClassificationXML, &prm)
	assert.Equal(t, 5, len(prm.DataDictionary.DataFields))
	rm := prm.RegressionModel
	assert.NoError(t, err)
	assert.Equal(t, "Sample for logistic regression", rm.ModelName)
	assert.Equal(t, "jobcat", rm.TargetFieldName)
	assert.Equal(t, "classification", rm.FunctionName)
	assert.Equal(t, "clerical", rm.RegressionTables[0].TargetCategory)
	assert.Equal(t, float64(46.418), rm.RegressionTables[0].Intercept)
	inputs := map[string]interface{}{
		"age":      float64(30),
		"work":     float64(0.1),
		"sex":      "0",
		"minority": "0",
	}
	out, err := rm.Evaluate(inputs)
	assert.NoError(t, err)
	// hand verified
	assert.InEpsilon(t, 0.6175894658256361, out["clerical"], 0.01)
	assert.InEpsilon(t, 0.34085492331618605, out["professional"], 0.01)
	assert.InEpsilon(t, 0.034430986935813215, out["skilled"], 0.01)
	assert.InEpsilon(t, 0.007124623922364555, out["trainee"], 0.01)
	sum := out["clerical"].(float64) + out["professional"].(float64) + out["skilled"].(float64) + out["trainee"].(float64)
	assert.InEpsilon(t, 1, sum, 0.01)
}

func TestInteractionTerms(t *testing.T) {
	var prm model.PMMLRegressionModel
	err := xml.Unmarshal(interactionTermsXML, &prm)
	assert.Equal(t, 4, len(prm.DataDictionary.DataFields))
	rm := prm.RegressionModel
	assert.NoError(t, err)
	assert.Equal(t, "Sample for interaction terms", rm.ModelName)
	assert.Equal(t, "regression", rm.FunctionName)
	assert.Equal(t, float64(2.1), rm.RegressionTables[0].Intercept)
	inputs := map[string]interface{}{
		"age":  float64(30),
		"work": float64(0.01),
		"sex":  "female",
	}
	out, err := rm.Evaluate(inputs)
	assert.NoError(t, err)
	t.Log(out)
	assert.InEpsilon(t, float64(-18.455), out["y"], 0.01)

}
