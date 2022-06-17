package model_test

import (
	"encoding/xml"
	"testing"

	"github.com/stillmatic/pummel/pkg/model"
	"github.com/stillmatic/pummel/pkg/tree"
	"github.com/stretchr/testify/assert"
)

func TestParseClassificationModel(t *testing.T) {
	var model model.MiningModel
	xmlStr := []byte(`
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
	err := xml.Unmarshal([]byte(xmlStr), &model)
	assert.NoError(t, err)
	assert.Equal(t, "classification", model.FunctionName)
	assert.Equal(t, 3, len(model.Segmentation.Segments))
	seg := model.Segmentation.Segments[0]
	assert.Equal(t, "1", seg.ID)
	assert.Equal(t, 1, len(seg.Predicates))
	pred := *seg.Predicates[0]
	dummyInput := make(map[string]interface{})
	predRes, err := pred.True(dummyInput)
	assert.True(t, predRes.ValueOrZero())
	assert.NoError(t, err)
	tm := seg.ModelElement
	assert.Equal(t, "Iris", tm.(*tree.TreeModel).ModelName)

}
