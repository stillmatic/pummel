package tree

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/stillmatic/pummel/pkg/fields"
	ms "github.com/stillmatic/pummel/pkg/miningschema"
	"github.com/stillmatic/pummel/pkg/node"
)

type TreeModel struct {
	XMLName              xml.Name         `xml:"TreeModel"`
	Node                 *node.Node       `xml:"Node"`
	MiningSchema         *ms.MiningSchema `xml:"MiningSchema"`
	ModelName            string           `xml:"modelName,attr"`
	FunctionName         string           `xml:"functionName,attr"`
	MissingValueStrategy string           `xml:"missingValueStrategy,attr"`
	MissingValuePenalty  float64          `xml:"missingValuePenalty,attr"`
	NoTrueChildStrategy  string           `xml:"noTrueChildStrategy,attr"`
	SplitCharacteristic  string           `xml:"splitCharacteristic,attr"`
	IsScorable           bool             `xml:"isScorable,attr"`
	Output               *fields.Outputs  `xml:"Output"`
}

// generate an enum struct for MissingValueStrategy
var MissingValueStrategy = struct {
	// evaluation is stopped and the current winner is returned as the final prediction.
	LastPrediction string
	// abort the scoring process and give no prediction.
	NullPrediction string
	// evaluate the attribute defaultChild which gives the child to continue traversing with.
	// Requires the presence of the attribute defaultChild in every non-leaf Node.
	DefaultChild string
	// the confidences for each class is calculated from scoring it and each of its sibling Nodes in
	// turn (excluding any siblings whose predicates evaluate to FALSE). The confidences returned for
	// each class from each sibling Node that was scored are weighted by the proportion of the number of
	// records in that Node, then summed to produce a total confidence for each class.
	// The winner is the class with the highest confidence.
	// Note that weightedConfidence should be applied recursively to deal with situations
	// where several predicates within the tree evaluate to UNKNOWN during the scoring of a case.
	WeightedConfidence string
	AggregateNodes     string
	// Comparisons with missing values other than checks for missing values always evaluate to FALSE.
	// If no rule fires, then use the noTrueChildStrategy to decide on a result.
	// This option requires that missing values be handled after all rules at the Node have been evaluated.
	// Note: In contrast to lastPrediction, evaluation is carried on instead of stopping immediately upon
	// first discovery of a Node who's predicate value cannot be determined due to missing values.
	None string
}{
	LastPrediction:     "lastPrediction",
	NullPrediction:     "nullPrediction",
	DefaultChild:       "defaultChild",
	WeightedConfidence: "weightedConfidence",
	AggregateNodes:     "aggregateNodes",
	None:               "none",
}

func (t *TreeModel) Evaluate(features map[string]interface{}) (map[string]interface{}, error) {
	rootPredRes, err := t.Node.Evaluate(features)
	if err != nil {
		return nil, err
	}
	if !rootPredRes.Valid {
		return nil, nil
	}
	curr, err := t.traverse(features)
	if err != nil {
		return nil, err
	}

	if curr.Score == "" {
		return nil, fmt.Errorf("terminal node without score, Node id: %v", curr.ID)
	}

	// this is the length of the number of outputs
	out := make(map[string]interface{}, 0)
	outputField := t.GetOutputField()

	if curr.ScoreDistributions != nil {
		var sum float64
		vals, sum, err := curr.HandleScoreDistributions()
		if err != nil {
			return nil, err
		}

		for i, val := range vals {
			var nameForField string
			if t.Output != nil {
				fieldName, err := t.Output.GetFeature(i)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to get feature name for index %q", i)
				}
				nameForField = fieldName.Name
			} else {
				nameForField = i
			}
			out[nameForField] = val / sum
		}
	}

	if t.FunctionName == "regression" {
		parsed, err := strconv.ParseFloat(curr.Score, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse score")
		}
		out[outputField] = parsed
	} else {
		out[outputField] = curr.Score
	}
	if t.Output != nil {
		for _, output := range t.Output.OutputFields {
			switch output.Feature {
			case "predictedValue":
				out[output.Name] = curr.Score
			case "probability":
			}
		}
	}
	return out, nil
}

func (t *TreeModel) traverse(features map[string]interface{}) (*node.Node, error) {
	curr := t.Node
	for len(curr.Children) > 0 {
		prevNode := curr
		for _, child := range curr.Children {
			predRes, err := child.Evaluate(features)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to evaluate child %s", child)
			}
			// handle missing value cases
			if !predRes.Valid {
				switch t.MissingValueStrategy {
				case MissingValueStrategy.LastPrediction:
					break
				case MissingValueStrategy.NullPrediction:
					return nil, nil
				case MissingValueStrategy.DefaultChild:
					if curr.DefaultChild == "" {
						break
					}
					curr, err = curr.GetDefaultChild()
					if err != nil {
						return nil, err
					}
				}
			}
			if predRes.Valid && predRes.Bool {
				curr = child
				break
			}
		}
		// could not find a valid child
		if curr == prevNode {
			break
		}
	}
	return curr, nil
}

func (t *TreeModel) GetOutputField() string {
	return t.MiningSchema.GetOutputField()
}
