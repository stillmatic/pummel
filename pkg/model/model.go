package model

import (
	"encoding/xml"
	"fmt"

	"github.com/stillmatic/pummel/pkg/tree"
	"gopkg.in/guregu/null.v4"
)

type Model interface {
	// TODO: Evaluate can return other things.
	Evaluate(map[string]interface{}) (null.String, error)
}

// see https://dmg.org/pmml/v4-3/GeneralStructure.html#xsdGroup_MODEL-ELEMENT
type ModelElement interface {
	// TODO: Evaluate can return other things.
	Evaluate(map[string]interface{}) (null.String, error)
}

type PMMLModel struct {
	XMLName        xml.Name        `xml:"PMML"`
	Header         *Header         `xml:"Header"`
	DataDictionary *DataDictionary `xml:"DataDictionary"`
}

type Header struct {
	XMLName     xml.Name `xml:"Header"`
	Copyright   string   `xml:"copyright,attr"`
	Description string   `xml:"description,attr"`
}

type PMMLTreeModel struct {
	PMMLModel
	TreeModel *tree.TreeModel `xml:"TreeModel"`
}

// TODO: genericize this.
func (ptm *PMMLTreeModel) Evaluate(features map[string]interface{}) (null.String, error) {
	return ptm.TreeModel.Evaluate(features)
}

// ValidateFeatures checks that each input feature is in the data dictionary.
func (pm *PMMLModel) ValidateFeatures(features map[string]interface{}) (bool, error) {
	// iterate over each feature
	for name, value := range features {
		// check if the feature is in the data dictionary
		var foundFeature bool
		for _, field := range pm.DataDictionary.DataFields {
			if field.Name == name {
				foundFeature = true
				// check if feature is a supported value
				if len(field.Values) > 0 {
					var found bool
					for _, validValue := range field.Values {
						if value.(string) == validValue.Value {
							found = true
							break
						}
					}
					if !found {
						fmt.Println("invalid value for feature", name, ":", value)
						return false, fmt.Errorf("invalid value for feature %s: %v", name, value)
					}
				}
			}
		}
		if !foundFeature {
			fmt.Println("unknown feature", name)
			return false, fmt.Errorf("unknown feature %s", name)
		}
	}
	return true, nil
}
