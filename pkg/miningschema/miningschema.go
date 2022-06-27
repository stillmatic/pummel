package miningschema

import "encoding/xml"

type MiningSchema struct {
	XMLName      xml.Name       `xml:"MiningSchema"`
	MiningFields []*MiningField `xml:"MiningField"`
}

type MiningField struct {
	XMLName    xml.Name `xml:"MiningField"`
	Name       string   `xml:"name,attr"`
	UsageType  string   `xml:"usageType,attr"`
	OpType     string   `xml:"opType,attr"`
	Importance float64  `xml:"importance,attr"`
	// Outliers determines how outliers are handled by the model.
	// Outliers are valid numeric values which are either greater than the specified
	// highValue or less than the specified lowValue.
	Outliers                string `xml:"outliers,attr"`
	LowValue                int64  `xml:"lowValue,attr"`
	HighValue               int64  `xml:"highValue,attr"`
	MissingValueReplacement string `xml:"missingValueReplacement,attr"`
	MissingValueTreatment   string `xml:"missingValueTreatment,attr"`
	InvalidValueTreatment   string `xml:"invalidValueTreatment,attr"`
	InvalidValueReplacement string `xml:"invalidValueReplacement,attr"`
}

var OutlierTreatmentMethods = struct {
	// AsIs field values treated at face value.
	AsIs string
	// AsMissingValues outlier values are treated as if they were missing.
	AsMissingValues string
	// AsExtremeValues outlier values are changed to a specific high or low value defined in MiningField.
	AsExtremeValues string
}{
	AsIs:            "asIs",
	AsMissingValues: "asMissingValues",
	AsExtremeValues: "asExtremeValues",
}

var MissingValueTreatmentMethods = struct {
	AsIs          string
	AsMean        string
	AsMode        string
	AsMedian      string
	AsValue       string
	ReturnInvalid string
}{
	AsIs:          "asIs",
	AsMean:        "asMean",
	AsMode:        "asMode",
	AsMedian:      "asMedian",
	AsValue:       "asValue",
	ReturnInvalid: "returnInvalid",
}

var InvalidValueTreatmentMethods = struct {
	ReturnInvalid string
	AsIs          string
	AsMissing     string
	AsValue       string
}{
	ReturnInvalid: "returnInvalid",
	AsIs:          "asIs",
	AsMissing:     "asMissing",
	AsValue:       "asValue",
}

func (ms *MiningSchema) GetOutputField() string {
	var out string
	for _, f := range ms.MiningFields {
		// 'predicted' is valid but deprecated as of PMML 4.2
		if f.UsageType == "predicted" || f.UsageType == "target" {
			out = f.Name
			break
		}
	}
	return out
}
