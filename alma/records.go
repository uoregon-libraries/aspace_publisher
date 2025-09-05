package alma

import (
  "encoding/xml"
  "github.com/tidwall/gjson"
  "encoding/json"
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
  Mms_id string 	`xml:"mms_id,omitempty"`
  SuppressPublish bool	`xml:"suppress_from_publishing"`
  SuppressExternal bool	`xml:"suppress_from_external_search"`
  Rec Record `xml:"record"`
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

func ExtractItemID(data []byte)string{
  var i Item
  err := json.Unmarshal(data, &i)
  if err != nil {}
  return i.Item_data.Item_pid
}

type Item struct{
  Holding_data HoldingData `json:"holding_data"`
  Item_data ItemData `json:"item_data"`
}

type BibData struct{
  Mms_id string `json:"mms_id,omitempty"`
}

type HoldingData struct{
  Holding_id string `json:"holding_id"`
  Copy_id string `json:"copy_id"`
}

type ItemData struct{
  Item_pid string `json:"pid,omitempty"`
  Barcode string `json:"barcode"`
  Policy Value `json:"policy"`
  Description string `json:"description"`
  Base_status Value `json:"base_status"`
  Library Value `json:"library"`
  Location Value `json:"location"`
  Physical_material_type Value `json:"physical_material_type"`
}

type Value struct {
  Val string `json:"value"`
}
