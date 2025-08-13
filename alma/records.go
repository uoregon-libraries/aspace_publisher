package alma

import (
  "encoding/xml"
  "github.com/tidwall/gjson"
  "log"
)

func ExtractBibID(data []byte)string{
  var b Bib
  err := xml.Unmarshal(data, &b)
  if err != nil { log.Println(err); return "" }
  return b.Mms_id
}

// GET returns json
// POST requires XML in and out
type Bib struct {
  XMLName xml.Name 	`xml:"bib"`	
  Mms_id string 	`xml:"mms_id"`
}

// pulls holding list for a existing bib from alma api
// expects only one holding per bib
func GetHoldingId(mms_id string)string{
  _url := BuildUrl( []string{"bibs", mms_id, "holdings"} )
  params := []string { ApiKey() }
  body,err := Get(_url, params)
  if err != nil { log.Println(err); return "" }
  holding_id := gjson.GetBytes(body, "holding.0.holding_id")
  return holding_id.String()
}

// pulls holding from the create output operation
func ExtractHoldingID(data []byte)string{
  var h Holding
  err := xml.Unmarshal(data, &h)
  if err != nil {}
  return h.HoldingId
}

// GET returns json
// POST requires XML in and out
type Holding struct {
  XMLName xml.Name `xml:"holding"`
  HoldingId string `xml:"holding_id,omitempty"`
  Suppress bool `xml:"suppress_from_publishing"`
  Rec Record `xml:"record"`
}

type Value struct {
  Val string `json:"value"`
}
