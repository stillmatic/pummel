package node

import (
	"encoding/xml"

	"github.com/stillmatic/pummel/pkg/predicates"
)

type Node struct {
	XMLName     xml.Name `xml:"Node"`
	Predicate   predicates.Predicate
	ID          string `xml:"id,attr"`
	Score       string `xml:"score,attr"`
	RecordCount int    `xml:"recordCount,attr"`
}

type Partition struct {
	XMLName xml.Name `xml:"Partition"`
	Name    string   `xml:"name,attr"`
}

type PartitionFieldStats struct {
	XMLName  xml.Name `xml:"PartitionFieldStats"`
	Field    string   `xml:"field,attr"`
	Weighted bool     `xml:"weighted,attr"`
}
