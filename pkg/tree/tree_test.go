package tree_test

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
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
	assert.Equal(t, "temperature", (*tm.MiningSchema).MiningFields[0].Name)
	assert.Equal(t, "humidity", (*tm.MiningSchema).MiningFields[1].Name)
	assert.Equal(t, "windy", (*tm.MiningSchema).MiningFields[2].Name)
	assert.Equal(t, "outlook", (*tm.MiningSchema).MiningFields[3].Name)
	assert.Equal(t, "whatIdo", (*tm.MiningSchema).MiningFields[4].Name)
	assert.Equal(t, "whatIdo", tm.GetOutputField())

	inputData := map[string]interface{}{
		"outlook":     "overcast",
		"temperature": "75",
		"humidity":    "55",
		"windy":       "false",
	}

	res, err := tm.Evaluate(inputData)
	t.Log(res)
	assert.NoError(t, err)
	assert.Equal(t, "may play", res[tm.GetOutputField()])
}

var complexTreeXML = []byte(`
<TreeModel modelName="Iris" functionName="classification" splitCharacteristic="binarySplit">
<MiningSchema>
  <MiningField name="petal_length" usageType="active"/>
  <MiningField name="petal_width" usageType="active"/>
  <MiningField name="Class" usageType="target"/>
</MiningSchema>
<Output>
  <OutputField dataType="string" feature="predictedValue" name="PredictedClass" optype="categorical"/>
  <OutputField dataType="double" feature="probability" name="Probability_setosa" optype="continuous" value="Iris-setosa"/>
  <OutputField dataType="double" feature="probability" name="Probability_versicolor" optype="continuous" value="Iris-versicolor"/>
  <OutputField dataType="double" feature="probability" name="Probability_virginica" optype="continuous" value="Iris-virginica"/>
</Output>
<Node score="Iris-setosa" recordCount="150">
  <True/>
  <ScoreDistribution value="Iris-setosa" recordCount="50"/>
  <ScoreDistribution value="Iris-versicolor" recordCount="50"/>
  <ScoreDistribution value="Iris-virginica" recordCount="50"/>
  <Node score="Iris-setosa" recordCount="50">
	<SimplePredicate field="petal_length" operator="lessThan" value="2.45"/>
	<ScoreDistribution value="Iris-setosa" recordCount="50"/>
	<ScoreDistribution value="Iris-versicolor" recordCount="0"/>
	<ScoreDistribution value="Iris-virginica" recordCount="0"/>
  </Node>
  <Node score="Iris-versicolor" recordCount="100">
	<SimplePredicate field="petal_length" operator="greaterThan" value="2.45"/>
	<ScoreDistribution value="Iris-setosa" recordCount="0"/>
	<ScoreDistribution value="Iris-versicolor" recordCount="50"/>
	<ScoreDistribution value="Iris-virginica" recordCount="50"/>
	<Node score="Iris-versicolor" recordCount="54">
	  <SimplePredicate field="petal_width" operator="lessThan" value="1.75"/>
	  <ScoreDistribution value="Iris-setosa" recordCount="0"/>
	  <ScoreDistribution value="Iris-versicolor" recordCount="49"/>
	  <ScoreDistribution value="Iris-virginica" recordCount="5"/>
	</Node>
	<Node score="Iris-virginica" recordCount="46">
	  <SimplePredicate field="petal_width" operator="greaterThan" value="1.75"/>
	  <ScoreDistribution value="Iris-setosa" recordCount="0"/>
	  <ScoreDistribution value="Iris-versicolor" recordCount="1"/>
	  <ScoreDistribution value="Iris-virginica" recordCount="45"/>
	</Node>
  </Node>
</Node>
</TreeModel>
`)

func TestComplexTreeXML(t *testing.T) {
	var tm *tree.TreeModel
	err := xml.Unmarshal(complexTreeXML, &tm)
	assert.NoError(t, err)
	assert.Equal(t, len(tm.Output.OutputFields), 4)

	assert.Equal(t, "Iris", tm.ModelName)
	pv, err := tm.Output.GetPredictedValue()
	assert.NoError(t, err)
	assert.Equal(t, "PredictedClass", pv.Name)
	input := map[string]interface{}{
		"petal_length": 2.5,
		"petal_width":  1.5,
		"temperature":  0.5,
		"cloudiness":   0.5,
	}
	res, err := tm.Evaluate(input)
	assert.NoError(t, err)
	t.Log(res)
	assert.Equal(t, "Iris-versicolor", res["PredictedClass"])
}

var TreeTests = []struct {
	features map[string]interface{}
	score    float64
	err      error
}{
	{map[string]interface{}{},
		4.3463944950723456e-4,
		nil,
	},
	{
		map[string]interface{}{"f2": "f2v1"},
		-1.8361380219689046e-4,
		nil,
	},
	{
		map[string]interface{}{"f2": "f2v1", "f1": "f1v3"},
		-6.237581139073701e-4,
		nil,
	},
	{
		map[string]interface{}{"f2": "f2v1", "f1": "f1v3", "f4": 0.08},
		0.0021968294712358194,
		nil,
	},
	{
		map[string]interface{}{"f2": "f2v1", "f1": "f1v3", "f4": 0.09},
		-9.198573460887271e-4,
		nil,
	},
	{
		map[string]interface{}{"f2": "f2v1", "f1": "f1v3", "f4": 0.09, "f3": "f3v2"},
		-0.0021187239505556523,
		nil,
	},
	{
		map[string]interface{}{"f2": "f2v1", "f1": "f1v3", "f4": 0.09, "f3": "f3v4"},
		-3.3516227414227926e-4,
		nil,
	},
	{
		map[string]interface{}{"f2": "f2v1", "f1": "f1v4"},
		0.0011015286521365208,
		nil,
	},
	{
		map[string]interface{}{"f2": "f2v4"},
		0.0022726641744997256,
		nil,
	},
	{
		map[string]interface{}{"f1": "f1v3", "f2": "f2v1", "f3": "f3v7", "f4": 0.09},
		-1,
		errors.New("terminal node without score, Node id: 5"),
	},
}

func TestTreeFixture(t *testing.T) {
	treeXmlIO, err := ioutil.ReadFile("../../testdata/tree.pmml")
	assert.NoError(t, err)
	var tm *tree.TreeModel
	err = xml.Unmarshal(treeXmlIO, &tm)
	assert.NoError(t, err)
	assert.Equal(t, "regression", tm.FunctionName)
	assert.Equal(t, "1", tm.Node.ID)
	assert.Equal(t, 4, len(tm.MiningSchema.MiningFields))
	for i, test := range TreeTests {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			res, err := tm.Evaluate(test.features)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			}
			if err == nil {
				assert.Equal(t, test.score, res[tm.GetOutputField()])
			}
		})
	}
}

//nolint
func BenchmarkTreeFixture(b *testing.B) {
	treeXmlIO, _ := ioutil.ReadFile("../../testdata/tree.pmml")
	var tm *tree.TreeModel
	xml.Unmarshal(treeXmlIO, &tm)
	for i, test := range TreeTests {
		b.Run(fmt.Sprintf("test%d", i), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				tm.Evaluate(test.features)
			}
		})
	}
}
