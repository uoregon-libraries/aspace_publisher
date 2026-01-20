package alma

import (
  "aspace_publisher/marc"
  "github.com/beevik/etree"
  "encoding/xml"
  "log"
  "fmt"
  "errors"
  "strings"
)

func ConstructBib(marc_string string, suppress bool)(Bib){
  var bib = Bib{}
  bib.SuppressPublish = suppress
  bib.SuppressExternal = true
  var rec = Record{}
  xml.Unmarshal([]byte(marc_string), &rec)
  bib.Rec = rec
  return bib
}

func ConstructBoundwith(boundwith_bib []byte, resource_marc string, resource_mmsid string, tcmap map[string]string)(Bib, error){
  var bwbib = Bib{}
  xml.Unmarshal(boundwith_bib, &bwbib)
  bw_xml, err := ParseMarc(string(boundwith_bib))
  if err != nil { return bwbib, err }
  if df774Exists(bw_xml, resource_mmsid) { return bwbib, nil }

  bib_xml, err := ParseMarc(resource_marc)
  if err != nil { return bwbib, err }

  title, err := ExtractTitle(bib_xml)
  if err != nil { return bwbib, err }
  sft := Subfield{Code: "t", Value: title}//title from the new coll/bib
  sfw := Subfield{Code: "w", Value: resource_mmsid }//mms_id of the new coll/bib
  d774 := Datafield{Ind1:"1", Ind2:" ", Tag:"774"}
  d774.Subfield = []Subfield{sft, sfw}
  bwbib.Rec.Datafield = append(bwbib.Rec.Datafield, d774)
  return bwbib, nil
}

func ConstructHolding(marc_string string, h Holding, id_0 string)(Holding, error){
  marc_xml, err := ParseMarc(marc_string)
  if err != nil { return h, err }
  link, err := BuildFindingLink(marc_xml)
  if err != nil { return h, err }
  fixed, err := ExtractFixed(marc_xml)
  if err != nil { return h, err }
  h.Suppress = false
  var rec Record
  rec.Leader, err = ExtractLeader(marc_xml)
  if err != nil { return h, err }

  rec.Controlfield = []Controlfield{ Controlfield{Tag:"008", Value: fixed} }
  sfb := Subfield{Code:"b", Value:"SpecColl"}
  sfc := Subfield{Code:"c", Value: "spmanus"}
  sfh := Subfield{Code:"h", Value: id_0}
  df852 := Datafield{Ind1:"8", Ind2:" ", Tag:"852"}
  df852.Subfield = []Subfield{sfb, sfc, sfh}
  sfz := Subfield{Code: "z", Value: link }
  df866 := Datafield{Ind1:"4", Ind2:"1", Tag:"866"}
  df866.Subfield = []Subfield{ sfz }
  rec.Datafield = []Datafield{ df852, df866 }
  h.Rec = rec
  return h, nil
}

func ConstructItem(holding_id string, item Item, tc_data map[string]string)(Item, error){
  item.Holding_data = HoldingData{ Holding_id: holding_id, Copy_id: "1" }
  var idata = ItemData{}
  idata.Barcode = tc_data["barcode"]
  idata.Policy = Value{ Val: policy(tc_data["type"]) }
  idata.Description = fmt.Sprintf("%s %s", tc_data["type"], tc_data["indicator"])
  idata.Library = Value{ Val: "SpecColl"}
  idata.Location = Value{ Val: "spmanus"}
  idata.Base_status = Value{ Val: "1" }
  idata.Physical_material_type = Value{ Val: "MANUSCRIPT" }
  item.Item_data = idata
  return item, nil
}

func policy(_type string)string{
  if strings.Contains(_type, "Unarranged") { return "Unarranged" } else if
    strings.Contains(_type, "Restricted") { return "Restricted" } else {
    return "999"
  }
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
  message := text.Text()
  if strings.Contains(message, "Connect to the online") { message = strings.ToUpper(text.Text()) }
  link := fmt.Sprintf("<a href='%s' target='_blank'>%s</a>", url.Text(), message)
  if message == "Notice of Interest in Unprocessed Collections" {
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

func ExtractTitle(marc_xml *etree.Document)(string, error){
  title := marc_xml.FindElement("//datafield[@tag='245']/subfield[@code='a']")
  if title == nil { return "", errors.New("unable to extract 245") }
  return title.Text(), nil
}

func df774Exists(marc_xml *etree.Document, resource_mmsid string) bool{
  df774 := marc_xml.FindElements(fmt.Sprintf("//datafield[@tag='774']/[subfield='%s']", resource_mmsid))
  if len(df774) == 0 { return false }
  return true
}
