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

	"github.com/pkg/errors"
)

type Expression interface {
	Transform(values map[string]interface{}) (interface{}, error)
	RequiredField() string
}

type LocalTransformations struct {
	XMLName       xml.Name `xml:"LocalTransformations"`
	DerivedFields []*DerivedField
}

type DerivedField struct {
	XMLName     xml.Name `xml:"DerivedField"`
	Name        string   `xml:"name,attr"`
	DisplayName string   `xml:"displayName,attr"`
	OpType      string   `xml:"opType,attr"`
	DataType    string   `xml:"dataType,attr"`
	Expression  *Expression
}

type FieldRef struct {
	XMLName      xml.Name `xml:"FieldRef"`
	Field        string   `xml:"field,attr"`
	MapMissingTo string   `xml:"mapMissingTo,attr"`
	DataType     string
}

type Apply struct {
	XMLName  xml.Name `xml:"Apply"`
	Function string   `xml:"function,attr"`
	Children []*Expression
}

type Constant struct {
	XMLName  xml.Name    `xml:"Constant"`
	DataType string      `xml:"dataType,attr"`
	Value    interface{} `xml:"Constant"`
}

func (lt *LocalTransformations) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	lt.XMLName = start.Name
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch tt := t.(type) {
		case xml.StartElement:
			var df DerivedField
			if err := d.DecodeElement(&df, &tt); err != nil {
				return err
			}
			lt.DerivedFields = append(lt.DerivedFields, &df)
		case xml.EndElement:
			return nil
		}
	}
}

// custom XML unmarshal for DerivedField
func (df *DerivedField) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	df.XMLName = start.Name
	// df.RequiredFields = make([]string, 0)
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
			case "Apply":
				expr = &Apply{}
			case "Constant":
				expr = &Constant{}
			default:
				return fmt.Errorf("unexpected element in DerivedField: %s", tt.Name.Local)
			}
			if expr != nil {
				if err := d.DecodeElement(&expr, &tt); err != nil {
					return err
				}
				df.Expression = &expr
				// df.RequiredFields = append(df.RequiredFields, expr.RequiredField())
			}
		case xml.EndElement:
			return nil
		}
	}
}

func (c *Constant) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	c.XMLName = start.Name
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "dataType":
			c.DataType = attr.Value
		}
	}
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}

		switch tt := t.(type) {
		case xml.CharData:
			switch c.DataType {
			case "double":
				parsed, err := strconv.ParseFloat(string(tt), 64)
				if err != nil {
					return err
				}
				c.Value = parsed
			}
		case xml.EndElement:
			return nil
		}
	}
}

func (a *Apply) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	a.XMLName = start.Name
	a.Children = make([]*Expression, 0)
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "function":
			a.Function = attr.Value
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
			case "Constant":
				expr = &Constant{}
			case "FieldRef":
				expr = &FieldRef{}
			case "Apply":
				expr = &Apply{}
			default:
				return fmt.Errorf("unexpected element in Apply: %s", tt.Name.Local)
			}
			if expr != nil {
				if err := d.DecodeElement(&expr, &tt); err != nil {
					return err
				}
				a.Children = append(a.Children, &expr)
			}
		case xml.EndElement:
			return nil
		}
	}
}

func (df *DerivedField) Transform(values map[string]interface{}) (interface{}, error) {
	return (*df.Expression).Transform(values)
}

// NB hack, this is not actually the required field but the output
func (df *DerivedField) RequiredField() string {
	return df.Name
}

func (fr *FieldRef) Transform(values map[string]interface{}) (interface{}, error) {
	value, ok := values[fr.Field]
	if !ok {
		return nil, errors.New("missing field " + fr.Field)
	}
	switch fr.DataType {
	case "float":
		return strconv.ParseFloat(value.(string), 64)
	case "double":
		return strconv.ParseFloat(value.(string), 64)
	default:
		return value, nil
	}
}

func (fr *FieldRef) RequiredField() string {
	return fr.Field
}

func (a *Apply) RequiredField() string {
	return ""
}

func (c *Constant) Transform(values map[string]interface{}) (interface{}, error) {
	switch c.DataType {
	case "double":
		return c.Value.(float64), nil
	}
	return c.Value, nil
}

func (c *Constant) RequiredField() string {
	return ""
}

func (a *Apply) Transform(values map[string]interface{}) (interface{}, error) {
	if len(a.Children) < 2 {
		return nil, errors.New("Apply requires at least two children")
	}
	// assume only 2 values for each of these...dont know if correct
	l, err := (*a.Children[0]).Transform(values)
	if err != nil {
		return nil, errors.Wrap(err, "error transforming left")
	}
	r, err := (*a.Children[1]).Transform(values)
	if err != nil {
		return nil, errors.Wrap(err, "error transforming right")
	}
	var lf, rf float64
	// TODO: this isn't right if these aren't actually numbers but yolo
	switch l := l.(type) {
	case float64:
		lf = l
	case int:
		lf = float64(l)
	case string:
		lf, err = strconv.ParseFloat(l, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "error parsing left %v", l)
		}
	default:
		return nil, errors.Errorf("left is not a number %v", l)
	}
	switch r := r.(type) {
	case float64:
		rf = r
	case int:
		rf = float64(r)
	case string:
		rf, err = strconv.ParseFloat(r, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "error parsing right %v", r)
		}
	default:
		return nil, errors.Errorf("right is not a number %v", r)
	}
	switch a.Function {
	case "-":
		return interface{}(lf - rf), nil
	case "*":
		return interface{}(lf * rf), nil
	case "/":
		return interface{}(lf / rf), nil
	}
	return nil, nil
}
