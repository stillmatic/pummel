package fields

import (
	"encoding/xml"

	"github.com/pkg/errors"
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

var (
	errNoOutputFields            = errors.New("no output fields")
	errNoOutputFieldsWithFeature = errors.New("no output fields with feature")
)

func (o *Outputs) GetFeature(value string) (*OutputField, error) {
	if o.OutputFields == nil {
		return nil, errNoOutputFields
	}
	for _, output := range o.OutputFields {
		if output.Value == value {
			return output, nil
		}
	}
	return nil, errors.Wrapf(errNoOutputFieldsWithFeature, "value: %s", value)
}

func (o *Outputs) GetPredictedValue() (*OutputField, error) {
	for _, output := range o.OutputFields {
		if output.Feature == "predictedValue" {
			return output, nil
		}
	}
	return nil, errors.Wrapf(errNoOutputFieldsWithFeature, "predictedValue")
}
