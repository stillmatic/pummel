package tree

import (
	"encoding/xml"
	"fmt"

	ms "github.com/stillmatic/pummel/pkg/miningschema"
	"github.com/stillmatic/pummel/pkg/node"
	"gopkg.in/guregu/null.v4"
)

type TreeModel struct {
	XMLName              xml.Name         `xml:"TreeModel"`
	Node                 *node.Node       `xml:"Node"`
	MiningSchema         *ms.MiningSchema `xml:"MiningSchema"`
	OutputField          *OutputField     `xml:"OutputField"`
	ModelName            string           `xml:"modelName,attr"`
	FunctionName         string           `xml:"functionName,attr"`
	MissingValueStrategy string           `xml:"missingValueStrategy,attr"`
	MissingValuePenalty  float64          `xml:"missingValuePenalty,attr"`
	NoTrueChildStrategy  string           `xml:"noTrueChildStrategy,attr"`
	SplitCharacteristic  string           `xml:"splitCharacteristic,attr"`
	IsScorable           bool             `xml:"isScorable,attr"`
}

type OutputField struct {
	XMLName     xml.Name `xml:"OutputField"`
	Name        string   `xml:"name,attr"`
	DisplayName string   `xml:"displayName,attr"`
	DataType    string   `xml:"dataType,attr"`
	Feature     string   `xml:"feature,attr"`
}

func (t *TreeModel) Evaluate(features map[string]interface{}) (null.String, error) {
	// check root first
	rootPredRes, err := t.Node.True(features)
	if err != nil {
		return null.StringFromPtr(nil), err
	}
	if !rootPredRes.Valid {
		return null.StringFromPtr(nil), nil
	}
	curr := t.Node
	for len(curr.Children) > 0 {
		for _, child := range curr.Children {
			predRes, err := child.True(features)
			fmt.Println("predRes:", predRes, "child:", child)
			if err != nil {
				return null.StringFromPtr(nil), err
			}
			if predRes.Valid && predRes.Bool {
				curr = &child
				break
			}
		}
	}
	return null.StringFrom(curr.Score), nil
}
