package alma

import (
  "aspace_publisher/marc"
  "github.com/beevik/etree"
  "encoding/json"
  "encoding/xml"
  "log"
  "fmt"
  "errors"
  "strings"
)

func ConstructBib(marc_string string)(string, error){
  marc_stripped, err := marc.StripOuterTags(marc_string)
  if err != nil { return "", err }
  return "<bib>" + marc_stripped + "</bib>", nil
}

func ConstructHolding(marc_string string)(string, error){
  marc_xml, err := ParseMarc(marc_string)
  if err != nil { return "", err }
  call_num, err := ExtractCall(marc_xml)
  if err != nil { return "", err }
  link, err := BuildFindingLink(marc_xml)
  if err != nil { return "", err }
  fixed, err := ExtractFixed(marc_xml)
  if err != nil { return "", err }
  var h = Holding{}
  h.Suppress = false
  var rec Record
  rec.Leader, err = ExtractLeader(marc_xml)
  if err != nil { return "", err }

  rec.Cfields = []Controlfield{ Controlfield{Tag:"008", Value: fixed} }
  sfb := Subfield{Code:"b", Value:"Special Collections"}
  sfc := Subfield{Code:"c", Value: "spmanus"}
  sfh := Subfield{Code:"h", Value: call_num}
  df852 := Datafield{Ind1:"8", Ind2:" ", Tag:"852"}
  df852.Sfields = []Subfield{sfb, sfc, sfh}
  sfz := Subfield{Code: "z", Value: link }
  df866 := Datafield{Ind1:"4", Ind2:"1", Tag:"866"}
  df866.Sfields = []Subfield{ sfz }
  rec.Dfields = []Datafield{ df852, df866 }
  h.Rec = rec
  output, err := xml.MarshalIndent(h, "  ", "    ")
  if err != nil { log.Println(err); return "", errors.New("unable to contruct holding json") }
  return string(output), nil
}

func ConstructItem(item_id string, holding_id string, tc_data map[string]string)(string, error){
  var item = Item{}
  item.Holding_data = HoldingData{ Holding_id: holding_id }
  var idata = ItemData{}
  idata.Barcode = tc_data["barcode"]
  idata.Policy = Value{ Val: policy(tc_data["type"]) }
  idata.Description = fmt.Sprintf("%s %s", tc_data["type"], tc_data["indicator"])
  idata.Library = Value{ Val: "SpecColl"}
  idata.Location = Value{ Val: "spmanus"}
  idata.Base_status = Value{ Val: "1" }
  idata.Physical_material_type = Value{ Val: "MANUSCRIPT" }
  item.Item_data = idata
  data, err := json.Marshal(item)
  if err != nil { log.Println(err); return "", errors.New("unable to construct item json") }
  return string(data), nil
}

func policy(_type string)string{
  if strings.Contains(_type, "Unarranged") { return "Unarranged" } else if
    strings.Contains(_type, "Restricted") { return "Restricted" } else {
    return "999"
  }
}
type Record struct{
  Leader string `xml:"leader"`
  Cfields []Controlfield
  Dfields []Datafield
}

type Controlfield struct{
  XMLName xml.Name `xml:"controlfield"`
  Tag string `xml:"tag,attr"`
  Value string `xml:",chardata"`
}

type Datafield struct{
  XMLName xml.Name `xml:"datafield"`
  Tag string `xml:"tag,attr"`
  Ind1 string `xml:"ind1,attr"`
  Ind2 string `xml:"ind2,attr"`
  Sfields []Subfield
  Value string `xml:",chardata"`
}

type Subfield struct{
  XMLName xml.Name `xml:"subfield"`
  Code string `xml:"code,attr"`
  Value string `xml:",chardata"`
}

func ParseMarc(marc_string string)(*etree.Document, error){
  marc_stripped, err := marc.StripOuterTags(marc_string)
  marc_xml := etree.NewDocument()
  err = marc_xml.ReadFromString(marc_stripped)
  if err != nil { log.Println(err); return marc_xml, errors.New("Unable to read XML response from OCLC.") }
  return marc_xml, nil
}

//856->866 which is not on the LOC reference
//marc uses z for the display message, u for the url
//866 is a mashup, z, value is a link
func BuildFindingLink(marc_xml *etree.Document)(string, error){
  url := marc_xml.FindElement("//datafield[@tag='856']/subfield[@code='u']")
  if url == nil { return "", errors.New("unable to extract 856") }
  text := marc_xml.FindElement("//datafield[@tag='856']/subfield[@code='z']")
  if text == nil { return "", errors.New("unable to extract 856") }
  link := fmt.Sprintf("<a href='%s' target='_blank'>%s</a>", url.Text(), text.Text())
  if text.Text() == "Notice of Interest in Unprocessed Collections" {
    link = "UNARRANGED COLLECTION UNAVAILABLE FOR USE. Inquiries regarding these materials should be submitted via the " + link
  }
  return link, nil
}

//leader
func ExtractLeader(marc_xml *etree.Document)(string, error){
  leader := marc_xml.FindElement("//leader")
  if leader == nil { return "", errors.New("unable to extract leader") }
  return leader.Text(), nil
}

//df 099 sf a -> df 852 sf h
func ExtractCall(marc_xml *etree.Document)(string, error){
  call := marc_xml.FindElement("//datafield[@tag='099']/subfield[@code='a']")
  if call == nil { return "", errors.New("unable to extract 099") }
  return call.Text(), nil
}

//CF 008 -> 008
func ExtractFixed(marc_xml *etree.Document)(string, error){
  fixed := marc_xml.FindElement("//controlfield[@tag='008']")
  if fixed == nil { return "", errors.New("unable to extract 008") }
  return fixed.Text(), nil
}
