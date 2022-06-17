package node_test

import (
	"encoding/xml"
	"testing"

	"github.com/stillmatic/pummel/pkg/node"
	"github.com/stretchr/testify/assert"
)

func TestParseComplexNode(t *testing.T) {
	complexNodeXML := []byte(`    
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
`)
	var node *node.Node
	err := xml.Unmarshal(complexNodeXML, &node)
	assert.NoError(t, err)
	assert.Equal(t, "will play", node.Score)
	assert.Equal(t, 2, len(node.Children))
}

func TestSimpleNode(t *testing.T) {
	nodeXML := []byte(`
<Node id="N1" score="0">
  <True/>
  <Node id="T1" score="1">
    <SimplePredicate field="prob1" operator="greaterThan" value="0.33"/>
  </Node>
</Node>
`)
	var node *node.Node
	err := xml.Unmarshal(nodeXML, &node)
	assert.NoError(t, err)
	assert.Equal(t, "0", node.Score)
	assert.Equal(t, 1, len(node.Children))

}
