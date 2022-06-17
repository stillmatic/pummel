package model_test

import (
	"encoding/xml"
	"testing"

	"github.com/stillmatic/pummel/pkg/model"
	"github.com/stretchr/testify/assert"
)

var simpleXMlStr = []byte(`
<PMML xmlns="https://www.dmg.org/PMML-4_1" version="4.1">
<Header copyright="www.dmg.org" description="A very small binary tree model to show structure."/>
<DataDictionary numberOfFields="5">
  <DataField name="temperature" optype="continuous" dataType="double"/>
  <DataField name="humidity" optype="continuous" dataType="double"/>
  <DataField name="windy" optype="categorical" dataType="string">
	<Value value="true"/>
	<Value value="false"/>
  </DataField>
  <DataField name="outlook" optype="categorical" dataType="string">
	<Value value="sunny"/>
	<Value value="overcast"/>
	<Value value="rain"/>
  </DataField>
  <DataField name="whatIdo" optype="categorical" dataType="string">
	<Value value="will play"/>
	<Value value="may play"/>
	<Value value="no play"/>
  </DataField>
</DataDictionary>
<TreeModel modelName="golfing" functionName="classification">
  <MiningSchema>
	<MiningField name="temperature"/>
	<MiningField name="humidity"/>
	<MiningField name="windy"/>
	<MiningField name="outlook"/>
	<MiningField name="whatIdo" usageType="predicted"/>
  </MiningSchema>
  <Node score="will play">
	<True/>
	<Node score="will play">
	  <SimplePredicate field="outlook" operator="equal" value="sunny"/>
	  <Node score="will play">
		<CompoundPredicate booleanOperator="and">
		  <SimplePredicate field="temperature" operator="lessThan" value="90"/>
		  <SimplePredicate field="temperature" operator="greaterThan" value="50"/>
		</CompoundPredicate>
		<Node score="will play">
		  <SimplePredicate field="humidity" operator="lessThan" value="80"/>
		</Node>
		<Node score="no play">
		  <SimplePredicate field="humidity" operator="greaterOrEqual" value="80"/>
		</Node>
	  </Node>
	  <Node score="no play">
		<CompoundPredicate booleanOperator="or">
		  <SimplePredicate field="temperature" operator="greaterOrEqual" value="90"/>
		  <SimplePredicate field="temperature" operator="lessOrEqual" value="50"/>
		</CompoundPredicate>
	  </Node>
	</Node>
	<Node score="may play">
	  <CompoundPredicate booleanOperator="or">
		<SimplePredicate field="outlook" operator="equal" value="overcast"/>
		<SimplePredicate field="outlook" operator="equal" value="rain"/>
	  </CompoundPredicate>
	  <Node score="may play">
		<CompoundPredicate booleanOperator="and">
		  <SimplePredicate field="temperature" operator="greaterThan" value="60"/>
		  <SimplePredicate field="temperature" operator="lessThan" value="100"/>
		  <SimplePredicate field="outlook" operator="equal" value="overcast"/>
		  <SimplePredicate field="humidity" operator="lessThan" value="70"/>
		  <SimplePredicate field="windy" operator="equal" value="false"/>
		</CompoundPredicate>
	  </Node>
	  <Node score="no play">
		<CompoundPredicate booleanOperator="and">
		  <SimplePredicate field="outlook" operator="equal" value="rain"/>
		  <SimplePredicate field="humidity" operator="lessThan" value="70"/>
		</CompoundPredicate>
	  </Node>
	</Node>
  </Node>
</TreeModel>
</PMML>
`)

func TestParseTreeModel(t *testing.T) {
	var tm *model.PMMLTreeModel
	err := xml.Unmarshal(simpleXMlStr, &tm)
	assert.NoError(t, err)
	assert.Equal(t, "golfing", tm.TreeModel.ModelName)

	inputData := map[string]interface{}{
		"outlook":     "overcast",
		"temperature": "75",
		"humidity":    "55",
		"windy":       "false",
	}

	res, err := tm.Evaluate(inputData)
	assert.NoError(t, err)
	assert.True(t, res.Valid)
	assert.Equal(t, "may play", res.String)
}

type validateFeatureTestCase struct {
	features map[string]interface{}
	valid    bool
}

var featureTestCases = []validateFeatureTestCase{
	{map[string]interface{}{"outlook": "sunny", "temperature": "90", "humidity": "80", "windy": "false"}, true},
	{map[string]interface{}{"outlook": "sunny", "temperature": "90", "humidity": "80"}, true},
	{map[string]interface{}{"gmail": "sunny", "temperature": "90", "humidity": "80", "windy": "false"}, false},
	{map[string]interface{}{"outlook": "notgoodbob", "temperature": "90", "humidity": "80", "windy": "false"}, false},
}

func TestValidateFeatures(t *testing.T) {
	var tm *model.PMMLTreeModel
	err := xml.Unmarshal(simpleXMlStr, &tm)
	assert.NoError(t, err)
	for _, c := range featureTestCases {
		res, _ := tm.ValidateFeatures(c.features)
		assert.Equal(t, c.valid, res)
		if c.valid != res {
			t.Log(c.features)
		}

	}
}

//nolint:errcheck
func BenchmarkTreeModel(b *testing.B) {
	var tm *model.PMMLTreeModel
	xml.Unmarshal(simpleXMlStr, &tm)
	inputData := map[string]interface{}{
		"outlook":     "overcast",
		"temperature": "75",
		"humidity":    "55",
		"windy":       "false",
	}

	for i := 0; i < b.N; i++ {
		tm.Evaluate(inputData)
	}
}
