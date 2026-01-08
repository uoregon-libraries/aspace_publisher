package alma

import (
  "encoding/xml"
  "github.com/tidwall/gjson"
  "encoding/json"
  "log"
  "errors"
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

func(b Bib)Stringify()(string, error){
  output, err := xml.Marshal(b)
  if err != nil { log.Println(err); return "", errors.New("unable to construct bib xml") }
  return string(output), nil
}

type FetchBibIDFun func(string)string
// retrieves bib id by doing GET on item barcode
// for use with boundwith process
func FetchBibID(barcode string)string{
  path := []string{"items?item_barcode=" + barcode}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  item,err := Get(_url, params, "application/json")
  if err != nil { log.Println(err); return ""}
  mms_id := gjson.GetBytes(item, "bib_data.mms_id")
  return mms_id.String()
}

// pulls holding list for a existing bib from alma api
// expects only one holding per bib
func GetHoldingId(mms_id string)string{
  _url := BuildUrl( []string{"bibs", mms_id, "holdings"} )
  params := []string { ApiKey() }
  body,err := Get(_url, params, "application/xml")
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
func(h Holding)Stringify()(string, error){
  output, err := xml.Marshal(h)
  if err != nil { log.Println(err); return "", errors.New("unable to construct holding xml") }
  return string(output), nil
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

func (i Item)Stringify()(string, error){
  output, err := json.Marshal(i)
  if err != nil { log.Println(err); return "", errors.New("unable to construct item json") }
  return string(output), nil

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

type Record struct{
  Leader string `xml:"leader"`
  Controlfield []Controlfield `xml:"controlfield"`
  Datafield []Datafield `xml:"datafield"`
}

type Controlfield struct{
  Tag string `xml:"tag,attr"`
  Value string `xml:",chardata"`
}

type Datafield struct{
  Tag string `xml:"tag,attr"`
  Ind1 string `xml:"ind1,attr"`
  Ind2 string `xml:"ind2,attr"`
  Subfield []Subfield `xml:"subfield"`
  Value string `xml:",chardata"`
}

type Subfield struct{
  Code string `xml:"code,attr"`
  Value string `xml:",chardata"`
}
