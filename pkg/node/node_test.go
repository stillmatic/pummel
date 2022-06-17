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

type NodeEqualityTest struct {
	NodeA    *node.Node
	NodeB    *node.Node
	Expected bool
}

var NodeEqualityTests = []NodeEqualityTest{
	{
		NodeA: &node.Node{
			Score: "0",
			Children: []node.Node{
				{
					Score:    "1",
					Children: make([]node.Node, 0),
				},
				{
					Score:    "1",
					Children: make([]node.Node, 0),
				},
			},
		},
		NodeB: &node.Node{
			Score: "0",
			Children: []node.Node{
				{
					Score:    "1",
					Children: make([]node.Node, 0),
				},
				{
					Score:    "1",
					Children: make([]node.Node, 0),
				},
			},
		},
		Expected: true,
	},
	{
		NodeA: &node.Node{
			Score:    "0",
			Children: []node.Node{},
		},
		NodeB: &node.Node{
			Score:    "0",
			Children: []node.Node{},
		},
		Expected: true,
	},
	{
		NodeA: &node.Node{
			Score:    "1",
			Children: []node.Node{},
		},
		NodeB: &node.Node{
			Score:    "0",
			Children: []node.Node{},
		},
		Expected: false,
	},
	{
		NodeA: &node.Node{
			Score:    "1",
			ID:       "N1",
			Children: []node.Node{},
		},
		NodeB: &node.Node{
			Score: "1",
			ID:    "N1",
			Children: []node.Node{
				{
					Score:    "1",
					ID:       "T1",
					Children: make([]node.Node, 0),
				},
			},
		},
		Expected: false,
	},
	{
		NodeA: &node.Node{
			Score: "1",
			ID:    "N1",
			Children: []node.Node{
				{
					Score:    "2",
					ID:       "T1",
					Children: make([]node.Node, 0),
				},
			},
		},
		NodeB: &node.Node{
			Score: "1",
			ID:    "N1",
			Children: []node.Node{
				{
					Score:    "1",
					ID:       "T1",
					Children: make([]node.Node, 0),
				},
			},
		},
		Expected: false,
	},
}

func TestNodeEquality(t *testing.T) {
	for _, test := range NodeEqualityTests {
		assert.Equal(t, test.Expected, test.NodeA.EqualTo(*test.NodeB))
	}
}

func TestComputeWeightedConfidence(t *testing.T) {
	xmlStr := []byte(`
	<Node id="2" score="will play" recordCount="50" defaultChild="3">
	<SimplePredicate field="outlook" operator="equal" value="sunny"/>
	<ScoreDistribution value="will play" recordCount="40" confidence="0.8"/>
	<ScoreDistribution value="may play" recordCount="2" confidence="0.04"/>
	<ScoreDistribution value="no play" recordCount="8" confidence="0.16"/>
	<Node id="3" score="will play" recordCount="40">
	  <CompoundPredicate booleanOperator="surrogate">
		<SimplePredicate field="temperature" operator="greaterOrEqual" value="50"/>
		<SimplePredicate field="humidity" operator="lessThan" value="80"/>
	  </CompoundPredicate>
	  <ScoreDistribution value="will play" recordCount="36" confidence="0.9"/>
	  <ScoreDistribution value="may play" recordCount="2" confidence="0.05"/>
	  <ScoreDistribution value="no play" recordCount="2" confidence="0.05"/>
	</Node>
	<Node id="4" score="no play" recordCount="10">
	  <CompoundPredicate booleanOperator="surrogate">
		<SimplePredicate field="temperature" operator="lessThan" value="50"/>
		<SimplePredicate field="humidity" operator="greaterOrEqual" value="80"/>
	  </CompoundPredicate>
	  <ScoreDistribution value="will play" recordCount="4" confidence="0.4"/>
	  <ScoreDistribution value="may play" recordCount="0" confidence="0.0"/>
	  <ScoreDistribution value="no play" recordCount="6" confidence="0.6"/>
	</Node>
  </Node>
  `)
	var nd *node.Node
	err := xml.Unmarshal(xmlStr, &nd)
	assert.NoError(t, err)
	score, err := node.ComputeWeightedConfidence(nd.Children)
	assert.Equal(t, "will play", score)
	assert.NoError(t, err)
}
