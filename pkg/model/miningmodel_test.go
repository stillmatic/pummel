package model_test

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stillmatic/pummel/pkg/model"
	"github.com/stillmatic/pummel/pkg/tree"
	"github.com/stretchr/testify/assert"
)

var classificationEnsembleXMLStr = []byte(`
  <MiningModel functionName="classification">
  <MiningSchema>
    <MiningField name="petal_length" usageType="active"/>
    <MiningField name="petal_width" usageType="active"/>
    <MiningField name="day" usageType="active"/>
    <MiningField name="continent" usageType="active"/>
    <MiningField name="sepal_length" usageType="supplementary"/>
    <MiningField name="sepal_width" usageType="supplementary"/>
    <MiningField name="Class" usageType="target"/>
  </MiningSchema>
  <Output>
    <OutputField name="PredictedClass" optype="categorical" dataType="string" feature="predictedValue"/>
    <OutputField name="ProbSetosa" optype="continuous" dataType="double" feature="probability" value="Iris-setosa"/>
    <OutputField name="ProbVeriscolor" optype="continuous" dataType="double" feature="probability" value="Iris-versicolor"/>
    <OutputField name="ProbVirginica" optype="continuous" dataType="double" feature="probability" value="Iris-virginica"/>
  </Output>
  <Segmentation multipleModelMethod="majorityVote">
    <Segment id="1">
      <True/>
      <TreeModel modelName="Iris" functionName="classification" splitCharacteristic="binarySplit">
        <MiningSchema>
          <MiningField name="petal_length" usageType="active"/>
          <MiningField name="petal_width" usageType="active"/>
          <MiningField name="day" usageType="active"/>
          <MiningField name="continent" usageType="active"/>
          <MiningField name="sepal_length" usageType="supplementary"/>
          <MiningField name="sepal_width" usageType="supplementary"/>
          <MiningField name="Class" usageType="target"/>
        </MiningSchema>
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
            <True/>
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
              <True/>
              <ScoreDistribution value="Iris-setosa" recordCount="0"/>
              <ScoreDistribution value="Iris-versicolor" recordCount="1"/>
              <ScoreDistribution value="Iris-virginica" recordCount="45"/>
            </Node>
          </Node>
        </Node>
      </TreeModel>
    </Segment>
    <Segment id="2">
      <True/>
      <TreeModel modelName="Iris" functionName="classification" splitCharacteristic="binarySplit">
        <MiningSchema>
          <MiningField name="petal_length" usageType="active"/>
          <MiningField name="petal_width" usageType="active"/>
          <MiningField name="day" usageType="active"/>
          <MiningField name="continent" usageType="active"/>
          <MiningField name="sepal_length" usageType="supplementary"/>
          <MiningField name="sepal_width" usageType="supplementary"/>
          <MiningField name="Class" usageType="target"/>
        </MiningSchema>
        <Node score="Iris-setosa" recordCount="150">
          <True/>
          <ScoreDistribution value="Iris-setosa" recordCount="50"/>
          <ScoreDistribution value="Iris-versicolor" recordCount="50"/>
          <ScoreDistribution value="Iris-virginica" recordCount="50"/>
          <Node score="Iris-setosa" recordCount="50">
            <SimplePredicate field="petal_length" operator="lessThan" value="2.15"/>
            <ScoreDistribution value="Iris-setosa" recordCount="50"/>
            <ScoreDistribution value="Iris-versicolor" recordCount="0"/>
            <ScoreDistribution value="Iris-virginica" recordCount="0"/>
          </Node>
          <Node score="Iris-versicolor" recordCount="100">
            <True/>
            <ScoreDistribution value="Iris-setosa" recordCount="0"/>
            <ScoreDistribution value="Iris-versicolor" recordCount="50"/>
            <ScoreDistribution value="Iris-virginica" recordCount="50"/>
            <Node score="Iris-versicolor" recordCount="54">
              <SimplePredicate field="petal_width" operator="lessThan" value="1.93"/>
              <ScoreDistribution value="Iris-setosa" recordCount="0"/>
              <ScoreDistribution value="Iris-versicolor" recordCount="49"/>
              <ScoreDistribution value="Iris-virginica" recordCount="5"/>
              <Node score="Iris-versicolor" recordCount="48">
                <SimplePredicate field="continent" operator="equal" value="africa"/>
              </Node>
              <Node score="Iris-virginical" recordCount="6">
                <SimplePredicate field="continent" operator="notEqual" value="africa"/>
              </Node>
            </Node>
            <Node score="Iris-virginica" recordCount="46">
              <True/>
              <ScoreDistribution value="Iris-setosa" recordCount="0"/>
              <ScoreDistribution value="Iris-versicolor" recordCount="1"/>
              <ScoreDistribution value="Iris-virginica" recordCount="45"/>
            </Node>
          </Node>
        </Node>
      </TreeModel>
    </Segment>
    <Segment id="3">
      <True/>
      <TreeModel modelName="Iris" functionName="classification" splitCharacteristic="binarySplit">
        <MiningSchema>
          <MiningField name="petal_length" usageType="active"/>
          <MiningField name="petal_width" usageType="active"/>
          <MiningField name="day" usageType="active"/>
          <MiningField name="continent" usageType="active"/>
          <MiningField name="sepal_length" usageType="supplementary"/>
          <MiningField name="sepal_width" usageType="supplementary"/>
          <MiningField name="Class" usageType="target"/>
        </MiningSchema>
        <Node score="Iris-setosa" recordCount="150">
          <True/>
          <ScoreDistribution value="Iris-setosa" recordCount="50"/>
          <ScoreDistribution value="Iris-versicolor" recordCount="50"/>
          <ScoreDistribution value="Iris-virginica" recordCount="50"/>
          <Node score="Iris-setosa" recordCount="50">
            <SimplePredicate field="petal_width" operator="lessThan" value="2.85"/>
            <ScoreDistribution value="Iris-setosa" recordCount="50"/>
            <ScoreDistribution value="Iris-versicolor" recordCount="0"/>
            <ScoreDistribution value="Iris-virginica" recordCount="0"/>
          </Node>
          <Node score="Iris-versicolor" recordCount="100">
            <True/>
            <ScoreDistribution value="Iris-setosa" recordCount="0"/>
            <ScoreDistribution value="Iris-versicolor" recordCount="50"/>
            <ScoreDistribution value="Iris-virginica" recordCount="50"/>
            <Node score="Iris-versicolor" recordCount="54">
              <SimplePredicate field="continent" operator="equal" value="asia"/>
              <ScoreDistribution value="Iris-setosa" recordCount="0"/>
              <ScoreDistribution value="Iris-versicolor" recordCount="49"/>
              <ScoreDistribution value="Iris-virginica" recordCount="5"/>
            </Node>
            <Node score="Iris-virginica" recordCount="46">
              <SimplePredicate field="continent" operator="notEqual" value="asia"/>
              <ScoreDistribution value="Iris-setosa" recordCount="0"/>
              <ScoreDistribution value="Iris-versicolor" recordCount="1"/>
              <ScoreDistribution value="Iris-virginica" recordCount="45"/>
            </Node>
          </Node>
        </Node>
      </TreeModel>
    </Segment>
  </Segmentation>
</MiningModel>
`)

func TestClassificationEnsembleModel(t *testing.T) {
	var model model.MiningModel
	err := xml.Unmarshal([]byte(classificationEnsembleXMLStr), &model)
	assert.NoError(t, err)
	assert.Equal(t, "classification", model.FunctionName)
	assert.Equal(t, 3, len(model.Segmentation.Segments))
	// check first segmentation
	seg := model.Segmentation.Segments[0]
	assert.Equal(t, "1", seg.ID)
	assert.Equal(t, 1, len(seg.Predicates))
	pred := seg.Predicates[0]
	dummyInput := make(map[string]interface{})
	predRes, ok, err := pred.Evaluate(dummyInput)
	assert.True(t, predRes)
	assert.True(t, ok)
	assert.NoError(t, err)
	// check first segment's model
	tm := seg.ModelElement
	assert.Equal(t, "Iris", tm.(*tree.TreeModel).ModelName)
	assert.Equal(t, "Class", tm.(*tree.TreeModel).GetOutputField())

	actualInput := map[string]interface{}{
		"petal_length": 2.5,
		"petal_width":  1.0,
		"day":          1.0,
		"continent":    "africa",
		"sepal_length": 1.0,
		"sepal_width":  1.0,
	}
	res, err := tm.Evaluate(actualInput)
	assert.NoError(t, err)
	assert.Equal(t, res["Class"], "Iris-versicolor")
	// evaluate first segment
	res, err = seg.Evaluate(actualInput)
	// t.Log("first seg", res)
	assert.NoError(t, err)
	assert.Equal(t, res["Class"], "Iris-versicolor")
	// evaluate second segment
	res, err = model.Segmentation.Segments[1].Evaluate(actualInput)
	// t.Log("second seg", res)
	assert.NoError(t, err)
	assert.Equal(t, res["Class"], "Iris-versicolor")
	// evaluate third segment
	res, err = model.Segmentation.Segments[2].Evaluate(actualInput)
	// t.Log("third seg", res)
	assert.NoError(t, err)
	assert.Equal(t, res["Class"], "Iris-setosa")

	// check overall segmentation
	res, err = model.Segmentation.Evaluate(actualInput)
	assert.NoError(t, err)
	// t.Log("segmentation", res)
	assert.Equal(t, res["Class"], "Iris-versicolor")
	assert.Equal(t, res["Iris-versicolor"], 2.0)
	assert.Equal(t, res["Iris-setosa"], 1.0)

	// check classification ensemble model
	res, err = model.Evaluate(actualInput)
	assert.NoError(t, err)
	// t.Log("model", res)
	assert.Equal(t, res["Class"], "Iris-versicolor")
	assert.Equal(t, res["ProbVeriscolor"], (2.0 / 3.0))
	assert.Equal(t, res["ProbSetosa"], (1.0 / 3.0))
}

//nolint
func BenchmarkClassificationEnsembleSegmentation(b *testing.B) {
	var model model.MiningModel
	xml.Unmarshal([]byte(classificationEnsembleXMLStr), &model)
	actualInput := map[string]interface{}{
		"petal_length": 2.5,
		"petal_width":  1.0,
		"day":          1.0,
		"continent":    "africa",
		"sepal_length": 1.0,
		"sepal_width":  1.0,
	}
	for i := 0; i < b.N; i++ {
		_, err := model.Segmentation.Evaluate(actualInput)
		assert.NoError(b, err)
	}
}

// func BenchmarkClassificationEnsembleSegmentationConcurrent(b *testing.B) {
// 	var model model.MiningModel
// 	xml.Unmarshal([]byte(classificationEnsembleXMLStr), &model)
// 	actualInput := map[string]interface{}{
// 		"petal_length": 2.5,
// 		"petal_width":  1.0,
// 		"day":          1.0,
// 		"continent":    "africa",
// 		"sepal_length": 1.0,
// 		"sepal_width":  1.0,
// 	}
// 	for i := 0; i < b.N; i++ {
// 		model.Segmentation.EvaluateConcurrently(actualInput)
// 	}
// }

var regressionWeightedAverageXML = []byte(`
<MiningModel functionName="regression">
  <MiningSchema>
    <MiningField name="petal_length" usageType="active"/>
    <MiningField name="petal_width" usageType="active"/>
    <MiningField name="day" usageType="active"/>
    <MiningField name="continent" usageType="active"/>
    <MiningField name="sepal_length" usageType="target"/>
    <MiningField name="sepal_width" usageType="active"/>
  </MiningSchema>
  <Output>
    <OutputField name="PredictedSepalLength" optype="continuous" dataType="double" feature="predictedValue"/>
  </Output>
  <Segmentation multipleModelMethod="weightedAverage">
    <Segment id="1" weight="0.25">
      <True/>
      <TreeModel modelName="Iris" functionName="regression" splitCharacteristic="multiSplit">
        <MiningSchema>
          <MiningField name="petal_length" usageType="active"/>
          <MiningField name="petal_width" usageType="active"/>
          <MiningField name="day" usageType="active"/>
          <MiningField name="continent" usageType="active"/>
          <MiningField name="sepal_length" usageType="target"/>
          <MiningField name="sepal_width" usageType="active"/>
        </MiningSchema>
        <Node score="5.843333" recordCount="150">
          <True/>
          <Node score="5.179452" recordCount="73">
            <SimplePredicate field="petal_length" operator="lessThan" value="4.25"/>
            <Node score="5.005660" recordCount="53">
              <SimplePredicate field="petal_length" operator="lessThan" value="3.40"/>
            </Node>
            <Node score="4.735000" recordCount="20">
              <SimplePredicate field="sepal_width" operator="lessThan" value="3.25"/>
            </Node>
            <Node score="5.169697" recordCount="33">
              <SimplePredicate field="sepal_width" operator="greaterThan" value="3.25"/>
            </Node>
            <Node score="5.640000" recordCount="20">
              <SimplePredicate field="petal_length" operator="greaterThan" value="3.40"/>
            </Node>
          </Node>
          <Node score="6.472727" recordCount="77">
            <SimplePredicate field="petal_length" operator="greaterThan" value="4.25"/>
            <Node score="6.326471" recordCount="68">
              <SimplePredicate field="petal_length" operator="lessThan" value="6.05"/>
              <Node score="6.165116" recordCount="43">
                <SimplePredicate field="petal_length" operator="lessThan" value="5.15"/>
                <Node score="6.054545" recordCount="33">
                  <SimplePredicate field="sepal_width" operator="lessThan" value="3.05"/>
                </Node>
                <Node score="6.530000" recordCount="10">
                  <SimplePredicate field="sepal_width" operator="greaterThan" value="3.05"/>
                </Node>
              </Node>
              <Node score="6.604000" recordCount="25">
                <SimplePredicate field="petal_length" operator="greaterThan" value="5.15"/>
              </Node>
            </Node>
            <Node score="7.577778" recordCount="9">
              <SimplePredicate field="petal_length" operator="greaterThan" value="6.05"/>
            </Node>
          </Node>
        </Node>
      </TreeModel>
    </Segment>
    <Segment id="2" weight="0.25">
      <True/>
      <TreeModel modelName="Iris" functionName="regression" splitCharacteristic="multiSplit">
        <MiningSchema>
          <MiningField name="petal_length" usageType="active"/>
          <MiningField name="petal_width" usageType="active"/>
          <MiningField name="day" usageType="active"/>
          <MiningField name="continent" usageType="active"/>
          <MiningField name="sepal_length" usageType="target"/>
          <MiningField name="sepal_width" usageType="active"/>
        </MiningSchema>
        <Node score="5.843333" recordCount="150">
          <True/>
          <Node score="5.073333" recordCount="60">
            <SimplePredicate field="petal_width" operator="lessThan" value="1.15"/>
            <Node score="4.953659" recordCount="41">
              <SimplePredicate field="petal_width" operator="lessThan" value="0.35"/>
            </Node>
            <Node score="4.688235" recordCount="17">
              <SimplePredicate field="sepal_width" operator="lessThan" value="3.25"/>
            </Node>
            <Node score="5.141667" recordCount="24">
              <SimplePredicate field="sepal_width" operator="greaterThan" value="3.25"/>
            </Node>
            <Node score="5.331579" recordCount="19">
              <SimplePredicate field="petal_width" operator="greaterThan" value="0.35"/>
            </Node>
          </Node>
          <Node score="6.356667" recordCount="90">
            <SimplePredicate field="petal_width" operator="greaterThan" value="1.15"/>
            <Node score="6.160656" recordCount="61">
              <SimplePredicate field="petal_width" operator="lessThan" value="1.95"/>
              <Node score="5.855556" recordCount="18">
                <SimplePredicate field="petal_width" operator="lessThan" value="1.35"/>
              </Node>
              <Node score="6.288372" recordCount="43">
                <SimplePredicate field="petal_width" operator="greaterThan" value="1.35"/>
                <Node score="6.000000" recordCount="13">
                  <SimplePredicate field="sepal_width" operator="lessThan" value="2.75"/>
                </Node>
                <Node score="6.413333" recordCount="30">
                  <SimplePredicate field="sepal_width" operator="greaterThan" value="2.75"/>
                </Node>
              </Node>
            </Node>
            <Node score="6.768966" recordCount="29">
              <SimplePredicate field="petal_width" operator="greaterThan" value="1.95"/>
            </Node>
          </Node>
        </Node>
      </TreeModel>
    </Segment>
    <Segment id="3" weight="0.5">
      <True/>
      <TreeModel modelName="Iris" functionName="regression" splitCharacteristic="multiSplit">
        <MiningSchema>
          <MiningField name="petal_length" usageType="active"/>
          <MiningField name="petal_width" usageType="active"/>
          <MiningField name="day" usageType="active"/>
          <MiningField name="continent" usageType="active"/>
          <MiningField name="sepal_length" usageType="target"/>
          <MiningField name="sepal_width" usageType="active"/>
        </MiningSchema>
        <Node score="5.843333" recordCount="150">
          <True/>
          <Node score="5.179452" recordCount="73">
            <SimplePredicate field="petal_length" operator="lessThan" value="4.25"/>
            <Node score="5.005660" recordCount="53">
              <SimplePredicate field="petal_length" operator="lessThan" value="3.40"/>
            </Node>
            <Node score="5.640000" recordCount="20">
              <SimplePredicate field="petal_length" operator="greaterThan" value="3.40"/>
            </Node>
          </Node>
          <Node score="6.472727" recordCount="77">
            <SimplePredicate field="petal_length" operator="greaterThan" value="4.25"/>
            <Node score="6.326471" recordCount="68">
              <SimplePredicate field="petal_length" operator="lessThan" value="6.05"/>
              <Node score="6.165116" recordCount="43">
                <SimplePredicate field="petal_length" operator="lessThan" value="5.15"/>
              </Node>
              <Node score="6.604000" recordCount="25">
                <SimplePredicate field="petal_length" operator="greaterThan" value="5.15"/>
              </Node>
            </Node>
            <Node score="7.577778" recordCount="9">
              <SimplePredicate field="petal_length" operator="greaterThan" value="6.05"/>
            </Node>
          </Node>
        </Node>
      </TreeModel>
    </Segment>
  </Segmentation>
</MiningModel>
`)

func TestRegressionWeightedAverage(t *testing.T) {
	var model model.MiningModel
	err := xml.Unmarshal([]byte(regressionWeightedAverageXML), &model)
	assert.NoError(t, err)
	assert.Equal(t, "regression", model.FunctionName)
	assert.Equal(t, "PredictedSepalLength", (*model.Output).OutputFields[0].Name)
	assert.Equal(t, "weightedAverage", model.Segmentation.MultipleModelMethod)

	inputs := map[string]interface{}{
		"petal_length": 5.5,
		"petal_width":  1.5,
		"day":          1,
		"continent":    "europe",
		"sepal_length": 5.1,
		"sepal_width":  3.5,
	}
	result, err := model.Segmentation.Evaluate(inputs)
	assert.NoError(t, err)
	assert.InEpsilon(t, 6.55633325, result["sepal_length"], 0.01)
}

//nolint
func BenchmarkRegressionWeightedAverage(b *testing.B) {
	var model model.MiningModel
	xml.Unmarshal([]byte(regressionWeightedAverageXML), &model)
	inputs := map[string]interface{}{
		"petal_length": 5.5,
		"petal_width":  1.5,
		"day":          1,
		"continent":    "Europe",
		"sepal_length": 5.1,
		"sepal_width":  3.5,
	}
	for i := 0; i < b.N; i++ {
		_, err := model.Segmentation.Evaluate(inputs)
		assert.NoError(b, err)
	}
}

var selectFirstXML = []byte(`
<MiningModel functionName="classification">
  <MiningSchema>
    <MiningField name="petal_length" usageType="active"/>
    <MiningField name="petal_width" usageType="active"/>
    <MiningField name="day" usageType="active"/>
    <MiningField name="continent" usageType="active"/>
    <MiningField name="sepal_length" usageType="supplementary"/>
    <MiningField name="sepal_width" usageType="supplementary"/>
    <MiningField name="Class" usageType="target"/>
  </MiningSchema>
  <Output>
    <OutputField name="PredictedClass" optype="categorical" dataType="string" feature="predictedValue"/>
    <OutputField name="ProbSetosa" optype="continuous" dataType="double" feature="probability" value="Iris-setosa"/>
    <OutputField name="ProbVeriscolor" optype="continuous" dataType="double" feature="probability" value="Iris-versicolor"/>
    <OutputField name="ProbVirginica" optype="continuous" dataType="double" feature="probability" value="Iris-virginica"/>
  </Output>
  <Segmentation multipleModelMethod="selectFirst">
    <Segment id="1">
      <CompoundPredicate booleanOperator="and">
        <SimplePredicate field="continent" operator="equal" value="asia"/>
        <SimplePredicate field="day" operator="lessThan" value="60.0"/>
        <SimplePredicate field="day" operator="greaterThan" value="0.0"/>
      </CompoundPredicate>
      <TreeModel modelName="Iris" functionName="classification" splitCharacteristic="binarySplit">
        <MiningSchema>
          <MiningField name="petal_length" usageType="active"/>
          <MiningField name="petal_width" usageType="active"/>
          <MiningField name="day" usageType="active"/>
          <MiningField name="continent" usageType="active"/>
          <MiningField name="sepal_length" usageType="supplementary"/>
          <MiningField name="sepal_width" usageType="supplementary"/>
          <MiningField name="Class" usageType="target"/>
        </MiningSchema>
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
    </Segment>
    <Segment id="2">
      <CompoundPredicate booleanOperator="and">
        <SimplePredicate field="continent" operator="equal" value="africa"/>
        <SimplePredicate field="day" operator="lessThan" value="60.0"/>
        <SimplePredicate field="day" operator="greaterThan" value="0.0"/>
      </CompoundPredicate>
      <TreeModel modelName="Iris" functionName="classification" splitCharacteristic="binarySplit">
        <MiningSchema>
          <MiningField name="petal_length" usageType="active"/>
          <MiningField name="petal_width" usageType="active"/>
          <MiningField name="day" usageType="active"/>
          <MiningField name="continent" usageType="active"/>
          <MiningField name="sepal_length" usageType="supplementary"/>
          <MiningField name="sepal_width" usageType="supplementary"/>
          <MiningField name="Class" usageType="target"/>
        </MiningSchema>
        <Node score="Iris-setosa" recordCount="150">
          <True/>
          <ScoreDistribution value="Iris-setosa" recordCount="50"/>
          <ScoreDistribution value="Iris-versicolor" recordCount="50"/>
          <ScoreDistribution value="Iris-virginica" recordCount="50"/>
          <Node score="Iris-setosa" recordCount="50">
            <SimplePredicate field="petal_length" operator="lessThan" value="2.15"/>
            <ScoreDistribution value="Iris-setosa" recordCount="50"/>
            <ScoreDistribution value="Iris-versicolor" recordCount="0"/>
            <ScoreDistribution value="Iris-virginica" recordCount="0"/>
          </Node>
          <Node score="Iris-versicolor" recordCount="100">
            <SimplePredicate field="petal_length" operator="greaterThan" value="2.15"/>
            <ScoreDistribution value="Iris-setosa" recordCount="0"/>
            <ScoreDistribution value="Iris-versicolor" recordCount="50"/>
            <ScoreDistribution value="Iris-virginica" recordCount="50"/>
            <Node score="Iris-versicolor" recordCount="54">
              <SimplePredicate field="petal_width" operator="lessThan" value="1.93"/>
              <ScoreDistribution value="Iris-setosa" recordCount="0"/>
              <ScoreDistribution value="Iris-versicolor" recordCount="49"/>
              <ScoreDistribution value="Iris-virginica" recordCount="5"/>
            </Node>
            <Node score="Iris-virginica" recordCount="46">
              <SimplePredicate field="petal_width" operator="greaterThan" value="1.93"/>
              <ScoreDistribution value="Iris-setosa" recordCount="0"/>
              <ScoreDistribution value="Iris-versicolor" recordCount="1"/>
              <ScoreDistribution value="Iris-virginica" recordCount="45"/>
            </Node>
          </Node>
        </Node>
      </TreeModel>
    </Segment>
    <Segment id="3">
      <SimplePredicate field="continent" operator="equal" value="africa"/>
      <TreeModel modelName="Iris" functionName="classification" splitCharacteristic="binarySplit">
        <MiningSchema>
          <MiningField name="petal_length" usageType="active"/>
          <MiningField name="petal_width" usageType="active"/>
          <MiningField name="day" usageType="active"/>
          <MiningField name="continent" usageType="active"/>
          <MiningField name="sepal_length" usageType="supplementary"/>
          <MiningField name="sepal_width" usageType="supplementary"/>
          <MiningField name="Class" usageType="target"/>
        </MiningSchema>
        <Node score="Iris-setosa" recordCount="150">
          <True/>
          <ScoreDistribution value="Iris-setosa" recordCount="50"/>
          <ScoreDistribution value="Iris-versicolor" recordCount="50"/>
          <ScoreDistribution value="Iris-virginica" recordCount="50"/>
          <Node score="Iris-setosa" recordCount="50">
            <SimplePredicate field="petal_width" operator="lessThan" value="2.85"/>
            <ScoreDistribution value="Iris-setosa" recordCount="50"/>
            <ScoreDistribution value="Iris-versicolor" recordCount="0"/>
            <ScoreDistribution value="Iris-virginica" recordCount="0"/>
          </Node>
          <Node score="Iris-versicolor" recordCount="100">
            <SimplePredicate field="petal_width" operator="greaterThan" value="2.85"/>
            <ScoreDistribution value="Iris-setosa" recordCount="0"/>
            <ScoreDistribution value="Iris-versicolor" recordCount="50"/>
            <ScoreDistribution value="Iris-virginica" recordCount="50"/>
          </Node>
        </Node>
      </TreeModel>
    </Segment>
  </Segmentation>
</MiningModel>
`)

func TestSelectFirst(t *testing.T) {
	var model model.MiningModel
	err := xml.Unmarshal(selectFirstXML, &model)
	assert.NoError(t, err)
	inputs := map[string]interface{}{
		"petal_length": 5.5,
		"petal_width":  1.5,
		"day":          1,
		"continent":    "asia",
		"sepal_length": 5.1,
		"sepal_width":  3.5,
	}
	assert.Equal(t, "classification", model.FunctionName)

	result, err := model.Segmentation.Evaluate(inputs)
	assert.NoError(t, err)
	t.Log(result)
	assert.Equal(t, "Iris-versicolor", result["Class"])

	// assert second segment
	inputs = map[string]interface{}{
		"petal_length": 1.5,
		"petal_width":  1.5,
		"day":          1,
		"continent":    "africa",
		"sepal_length": 5.1,
		"sepal_width":  3.5,
	}
	assert.Equal(t, "classification", model.FunctionName)

	result, err = model.Segmentation.Evaluate(inputs)
	assert.NoError(t, err)
	assert.Equal(t, "Iris-setosa", result["Class"])
}

var modelChainXML = []byte(`
<MiningModel functionName="regression">
  <MiningSchema>
    <MiningField name="petal_length" usageType="active"/>
    <MiningField name="petal_width" usageType="active"/>
    <MiningField name="temperature" usageType="active"/>
    <MiningField name="cloudiness" usageType="active"/>
    <MiningField name="sepal_length" usageType="supplementary"/>
    <MiningField name="sepal_width" usageType="supplementary"/>
    <MiningField name="Class" usageType="target"/>
    <MiningField name="PollenIndex" usageType="target"/>
  </MiningSchema>
  <Output>
    <OutputField dataType="string" feature="predictedValue" name="PredictedClass" optype="categorical" targetField="Class" segmentId="1"/>
    <OutputField dataType="double" feature="probability" name="Probability_setosa" optype="continuous" targetField="Class" value="Iris-setosa" segmentId="1"/>
    <OutputField dataType="double" feature="probability" name="Probability_versicolor" optype="continuous" targetField="Class" value="Iris-versicolor" segmentId="1"/>
    <OutputField dataType="double" feature="probability" name="Probability_virginica" optype="continuous" targetField="Class" value="Iris-virginica" segmentId="1"/>
    <OutputField dataType="double" feature="predictedValue" name="Pollen Index" optype="continuous" targetField="PollenIndex"/>
  </Output>
  <Segmentation multipleModelMethod="modelChain">
    <Segment id="1">
      <True/>
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
    </Segment>
    <Segment id="2">
      <True/>
      <RegressionModel modelName="PollenIndex" functionName="regression">
        <MiningSchema>
          <MiningField name="Probability_setosa" usageType="active"/>
          <MiningField name="Probability_versicolor" usageType="active"/>
          <MiningField name="Probability_virginica" usageType="active"/>
          <MiningField name="temperature" usageType="active"/>
          <MiningField name="cloudiness" usageType="active"/>
          <MiningField name="PollenIndex" usageType="target"/>
        </MiningSchema>
        <Output>
          <OutputField dataType="double" feature="predictedValue" name="Pollen Index" optype="continuous"/>
        </Output>
        <RegressionTable intercept="0.3">
          <NumericPredictor coefficient="0.8" exponent="1" name="Probability_setosa"/>
          <NumericPredictor coefficient="0.7" exponent="1" name="Probability_versicolor"/>
          <NumericPredictor coefficient="0.9" exponent="1" name="Probability_virginica"/>
          <NumericPredictor coefficient="0.02" exponent="1" name="temperature"/>
          <NumericPredictor coefficient="-0.1" exponent="1" name="cloudiness"/>
        </RegressionTable>
      </RegressionModel>
    </Segment>
  </Segmentation>
</MiningModel>
`)

func TestModelChain(t *testing.T) {
	var model model.MiningModel
	err := xml.Unmarshal(modelChainXML, &model)
	assert.NoError(t, err)
	assert.NotNil(t, model.Output)
	input := map[string]interface{}{
		"petal_length": 2.5,
		"petal_width":  1.5,
		"temperature":  0.5,
		"cloudiness":   0.5,
	}
	res, err := model.Segmentation.Evaluate(input)
	assert.InEpsilon(t, 0.9785185185185183, res["PollenIndex"].(float64), 0.01)
	assert.NoError(t, err)
}

var RFFixtureCases = []struct {
	name          string
	features      map[string]interface{}
	expectedScore float64
	expectedErr   error
}{
	{
		"low",
		map[string]interface{}{
			"Sex":      "male",
			"Parch":    0,
			"Age":      30,
			"Fare":     9.6875,
			"Pclass":   2,
			"SibSp":    0,
			"Embarked": "Q"},
		(2.0 / 15.0),
		nil,
	},
	{
		"high",
		map[string]interface{}{
			"Sex":      "female",
			"Parch":    0,
			"Age":      38,
			"Fare":     71.2833,
			"Pclass":   2,
			"SibSp":    1,
			"Embarked": "C",
		},
		(14.0 / 15.0),
		nil,
	},
	{
		"error",
		map[string]interface{}{
			"Sex":      "female",
			"Parch":    0,
			"Age":      38,
			"Fare":     71.2833,
			"Pclass":   2,
			"SibSp":    1,
			"Embarked": "UnknownCategory",
		},
		0.0,
		errors.New("failed to evaluate segmentation: failed to evaluate segment: failed to evaluate segment: terminal node without score, Node id: 1"),
	},
}

func TestRFFixture(t *testing.T) {
	rfXMLIO, err := ioutil.ReadFile("../../testdata/rf.pmml")
	assert.NoError(t, err)
	var mm model.PMMLMiningModel
	err = xml.Unmarshal(rfXMLIO, &mm)
	assert.NoError(t, err)
	assert.Equal(t, len(mm.DataDictionary.DataFields), 12)
	assert.Equal(t, len(mm.MiningModel.MiningSchema.MiningFields), 12)
	assert.Equal(t, len(mm.MiningModel.Output.OutputFields), 3)
	assert.Equal(t, len(mm.MiningModel.Segmentation.Segments), 15)

	for _, tc := range RFFixtureCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := mm.MiningModel.Evaluate(tc.features)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			}
			if err == nil {
				assert.InEpsilon(t, tc.expectedScore, res["Probability_1"].(float64), 0.01)
			}
		})
	}
}

var GBMFixtureCases = []struct {
	name          string
	features      map[string]interface{}
	expectedScore float64
	expectedErr   error
}{
	{
		"low",
		map[string]interface{}{
			"Sex":      "male",
			"Parch":    0,
			"Age":      30,
			"Fare":     9.6875,
			"Pclass":   2,
			"SibSp":    0,
			"Embarked": "Q"},
		0.3652639329522468,
		nil,
	},
	{
		"high",
		map[string]interface{}{
			"Sex":      "female",
			"Parch":    0,
			"Age":      38,
			"Fare":     71.2833,
			"Pclass":   2,
			"SibSp":    1,
			"Embarked": "C",
		},
		0.4178155014037758,
		nil,
	},
}

func TestGBMFixture(t *testing.T) {
	gbmXMLIO, err := ioutil.ReadFile("../../testdata/gbm.pmml")
	assert.NoError(t, err)
	var mm model.PMMLMiningModel
	err = xml.Unmarshal(gbmXMLIO, &mm)
	assert.NoError(t, err)
	assert.Equal(t, len(mm.DataDictionary.DataFields), 2)
	assert.Equal(t, len(mm.MiningModel.MiningSchema.MiningFields), 2)
	assert.Equal(t, len(mm.MiningModel.Segmentation.Segments), 2)

	for _, tc := range GBMFixtureCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := mm.MiningModel.Evaluate(tc.features)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			}
			if err == nil {
				assert.InEpsilon(t, tc.expectedScore, res["probability(1)"].(float64), 0.01)
			}
		})
	}
}

//nolint
func BenchmarkGBMFixture(b *testing.B) {
	gbmXMLIO, _ := ioutil.ReadFile("../../testdata/gbm.pmml")
	var mm model.PMMLMiningModel
	xml.Unmarshal(gbmXMLIO, &mm)

	for _, tc := range GBMFixtureCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := mm.MiningModel.Evaluate(tc.features)
				assert.NoError(b, err)
			}
		})
	}
}

//nolint
func BenchmarkRFFixture(b *testing.B) {
	rfXMLIO, _ := ioutil.ReadFile("../../testdata/rf.pmml")
	var mm model.PMMLMiningModel
	xml.Unmarshal(rfXMLIO, &mm)

	for _, tc := range RFFixtureCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := mm.MiningModel.Evaluate(tc.features)
				assert.NoError(b, err)
			}
		})
	}
}
