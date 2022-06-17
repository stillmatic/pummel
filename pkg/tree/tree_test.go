package tree_test

import (
	"encoding/xml"
	"testing"

	"github.com/stillmatic/pummel/pkg/tree"
	"github.com/stretchr/testify/assert"
)

func TestParseTree(t *testing.T) {
	xmlStr := []byte(`
	<TreeModel modelName="golfing" functionName="classification">
	  <MiningSchema>
		<MiningField name="temperature"/>
		<MiningField name="humidity"/>
		<MiningField name="windy"/>
		<MiningField name="outlook"/>
		<MiningField name="whatIdo" usageType="predicted"/>
	  </MiningSchema>
	  <Node score="will play" id="0">
		<True/>
		<Node score="will play" id="1">
		  <SimplePredicate field="outlook" operator="equal" value="sunny"/>
		  <Node score="will play" id="2">
			<CompoundPredicate booleanOperator="and">
			  <SimplePredicate field="temperature" operator="lessThan" value="90"/>
			  <SimplePredicate field="temperature" operator="greaterThan" value="50"/>
			</CompoundPredicate>
			<Node score="will play" id="3">
			  <SimplePredicate field="humidity" operator="lessThan" value="80"/>
			</Node>
			<Node score="no play" id="4">
			  <SimplePredicate field="humidity" operator="greaterOrEqual" value="80"/>
			</Node>
		  </Node>
		  <Node score="no play" id="5">
			<CompoundPredicate booleanOperator="or">
			  <SimplePredicate field="temperature" operator="greaterOrEqual" value="90"/>
			  <SimplePredicate field="temperature" operator="lessOrEqual" value="50"/>
			</CompoundPredicate>
		  </Node>
		</Node>
		<Node score="may play" id="6">
		  <CompoundPredicate booleanOperator="or">
			<SimplePredicate field="outlook" operator="equal" value="overcast"/>
			<SimplePredicate field="outlook" operator="equal" value="rain"/>
		  </CompoundPredicate>
		  <Node score="may play" id="7">
			<CompoundPredicate booleanOperator="and">
			  <SimplePredicate field="temperature" operator="greaterThan" value="60"/>
			  <SimplePredicate field="temperature" operator="lessThan" value="100"/>
			  <SimplePredicate field="outlook" operator="equal" value="overcast"/>
			  <SimplePredicate field="humidity" operator="lessThan" value="70"/>
			  <SimplePredicate field="windy" operator="equal" value="false"/>
			</CompoundPredicate>
		  </Node>
		  <Node score="no play" id="8">
			<CompoundPredicate booleanOperator="and">
			  <SimplePredicate field="outlook" operator="equal" value="rain"/>
			  <SimplePredicate field="humidity" operator="lessThan" value="70"/>
			</CompoundPredicate>
		  </Node>
		</Node>
	  </Node>
	</TreeModel>
`)
	var tm *tree.TreeModel
	err := xml.Unmarshal(xmlStr, &tm)
	assert.NoError(t, err)
	assert.Equal(t, "golfing", tm.ModelName)
	assert.Equal(t, "classification", tm.FunctionName)
	assert.Equal(t, "temperature", tm.MiningSchema.MiningFields[0].Name)
	assert.Equal(t, "humidity", tm.MiningSchema.MiningFields[1].Name)
	assert.Equal(t, "windy", tm.MiningSchema.MiningFields[2].Name)
	assert.Equal(t, "outlook", tm.MiningSchema.MiningFields[3].Name)
	assert.Equal(t, "whatIdo", tm.MiningSchema.MiningFields[4].Name)

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
