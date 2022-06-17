package model

import "encoding/xml"

type DataDictionary struct {
	XMLName    xml.Name     `xml:"DataDictionary"`
	DataFields []*DataField `xml:"DataField"`
}

type DataField struct {
	XMLName  xml.Name `xml:"DataField"`
	Name     string   `xml:"name,attr"`
	OpType   string   `xml:"optype,attr"`
	DataType string   `xml:"dataType,attr"`
	Values   []Value  `xml:"Value"`
}

type Value struct {
	XMLName xml.Name `xml:"Value"`
	Value   string   `xml:"value,attr"`
}
