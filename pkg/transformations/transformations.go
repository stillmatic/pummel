// These transforms fields in the input dictioanry to the model, before the model is evaluated.
// For performance reasons, we expect the model to have something like `PrepareDictionary`
// which calls each of these transforms *concurrently*, with only the expected field.
// Hence, each DerivedField only operates on a single field value (not map[string]interface{}).
// We also skip checking if the field names are identical, etc.
package transformations

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

type Expression interface {
	Transform(value interface{}) (interface{}, error)
	RequiredField() string
}

type DerivedField struct {
	XMLName        xml.Name `xml:"DerivedField"`
	Name           string   `xml:"name,attr"`
	DisplayName    string   `xml:"displayName,attr"`
	OpType         string   `xml:"opType,attr"`
	DataType       string   `xml:"dataType,attr"`
	Expression     *Expression
	RequiredFields []string
}

type FieldRef struct {
	XMLName      xml.Name `xml:"FieldRef"`
	Field        string   `xml:"field,attr"`
	MapMissingTo string   `xml:"mapMissingTo,attr"`
	DataType     string
}

// custom XML unmarshal for DerivedField
func (df *DerivedField) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	df.XMLName = start.Name
	df.RequiredFields = make([]string, 0)
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "name":
			df.Name = attr.Value
		case "displayName":
			df.DisplayName = attr.Value
		case "opType":
			df.OpType = attr.Value
		case "dataType":
			df.DataType = attr.Value
		}
	}
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch tt := t.(type) {
		case xml.StartElement:
			var expr Expression
			switch tt.Name.Local {
			case "FieldRef":
				expr = &FieldRef{DataType: df.DataType}
			default:
				return fmt.Errorf("unexpected element in DerivedField: %s", tt.Name.Local)
			}
			if expr != nil {
				if err := d.DecodeElement(&expr, &tt); err != nil {
					return err
				}
				df.Expression = &expr
				df.RequiredFields = append(df.RequiredFields, expr.RequiredField())
			}
		case xml.EndElement:
			return nil
		}
	}
}

func (df *DerivedField) Transform(value interface{}) (interface{}, error) {
	return (*df.Expression).Transform(value)
}

func (fr *FieldRef) Transform(value interface{}) (interface{}, error) {
	switch fr.DataType {
	case "float":
		return strconv.ParseFloat(value.(string), 64)
	case "double":
		return strconv.ParseFloat(value.(string), 64)
	default:
		return nil, fmt.Errorf("unexpected data type: %s", fr.DataType)
	}
}

func (fr *FieldRef) RequiredField() string {
	return fr.Field
}
