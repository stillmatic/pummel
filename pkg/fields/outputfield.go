package fields

import (
	"encoding/xml"
	"fmt"
)

type Outputs struct {
	OutputFields []*OutputField `xml:"OutputField"`
}

type OutputField struct {
	XMLName     xml.Name `xml:"OutputField"`
	Name        string   `xml:"name,attr"`
	DisplayName string   `xml:"displayName,attr"`
	OpType      string   `xml:"optype,attr"`
	DataType    string   `xml:"dataType,attr"`
	Feature     string   `xml:"feature,attr"`
	Value       string   `xml:"value,attr"`
}

func (o *Outputs) GetFeature(value string) (*OutputField, error) {
	for _, output := range (*o).OutputFields {
		if output.Value == value {
			return output, nil
		}
	}
	return nil, fmt.Errorf("no output field with value %s", value)
}

func (o *Outputs) GetPredictedValue() (*OutputField, error) {
	for _, output := range (*o).OutputFields {
		if output.Feature == "predictedValue" {
			return output, nil
		}
	}
	return nil, fmt.Errorf("no output field with feature predictedValue")
}
