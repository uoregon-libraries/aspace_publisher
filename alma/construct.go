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
// create should include the bib.suppress fields
// no updates
func ConstructBib(mms_id string, marc_string string, suppress string)(Bib){
  var bib = Bib{}
  bib.SuppressPublish = suppress
  bib.SuppressExternal = "true"
  var rec = Record{}
  xml.Unmarshal([]byte(marc_string), &rec)
  bib.Rec = rec
  return bib
}
// only used in the regular bib workflow!
// updates the boundwith when folders with mms_ids are added to the associated top container
func UpdateBoundwith(boundwith_bib []byte, resource_marc string, resource_mmsid string, tcmap map[string]string)(Bib, error){
  var bwbib = Bib{}
  xml.Unmarshal(boundwith_bib, &bwbib)
  if df774Exists(bwbib, resource_mmsid) { return bwbib, errors.New("skip") }

  bib_xml, err := ParseMarc(resource_marc)
  if err != nil { return bwbib, err }
  title, err := ExtractTitle(bib_xml)
  if err != nil { return bwbib, err }
  sft := Subfield{Code: "t", Value: title}//title from the new coll/bib
  sfw := Subfield{Code: "w", Value: resource_mmsid }//mms_id of the new coll/bib
  d774 := Datafield{Ind1:"1", Ind2:" ", Tag:"774"}
  d774.Subfield = []Subfield{sft, sfw}
  var newbwbib = Bib{}
  newbwbib.Rec = bwbib.Rec
  newbwbib.Rec.Datafield = append(newbwbib.Rec.Datafield, d774)
  return newbwbib, nil
}

func UpdateHolding(marc_string string, hold string)(string, error){
  marc_xml, err := ParseMarc(marc_string)
  if err != nil { return "", err }
  link, err := BuildFindingLink(marc_xml)
  if err != nil { return hold, err }
  holding, err := ParseXML(hold)
  if err != nil { return "", err }
  // acc to the API docs, remove the holdingId
  id_ptr := holding.FindElement("//holding_id")
  parent := id_ptr.Parent()
  parent.RemoveChild(id_ptr)
  holding.FindElement("//subfield[@code='z']").SetText(link)
  str, err := holding.WriteToString()
  return str, err
}
// create and update should both only include suppress_from_publishing and record fields
func ConstructHolding(marc_string string, id_0 string)(Holding, error){
  var h = Holding{}
  marc_xml, err := ParseMarc(marc_string)
  if err != nil { return h, err }
  link, err := BuildFindingLink(marc_xml)
  if err != nil { return h, err }
  fixed, err := ExtractFixed(marc_xml)
  if err != nil { return h, err }
  h.Suppress = false
  h.Rec.Leader, err = ExtractLeader(marc_xml)
  if err != nil { return h, err }
  h.Rec.Controlfield = []Controlfield{ Controlfield{Tag:"008", Value: fixed} }
  sfb := Subfield{Code:"b", Value:"SpecColl"}
  sfc := Subfield{Code:"c", Value: "spmanus"}
  sfh := Subfield{Code:"h", Value: id_0}
  df852 := Datafield{Ind1:"8", Ind2:" ", Tag:"852"}
  df852.Subfield = []Subfield{sfb, sfc, sfh}
  sfz := Subfield{Code: "z", Value: link }
  df866 := Datafield{Ind1:"4", Ind2:"1", Tag:"866"}
  df866.Subfield = []Subfield{ sfz }
  h.Rec.Datafield = []Datafield{ df852, df866 }

  return h, nil
}

//expects holding_data, item_data
func ConstructItem(holding_id string, item Item, tc_data map[string]string)(Item, error){
  if item.Item_data.Item_pid == "" {
    item.Holding_data.Holding_id = holding_id
    item.Holding_data.Copy_id = "1"
    item.Item_data.Barcode = tc_data["barcode"]
    item.Item_data.Library = Value{ Val: "SpecColl"}
    item.Item_data.Location = Value{ Val: "spmanus"}
    item.Item_data.Base_status = Value{ Val: "1" }
    item.Item_data.Physical_material_type = Value{ Val: "MANUSCRIPT" }
  }//note that the item pid SHOULD be left in the item_data
  item.Item_data.Policy = Value{ Val: policy(tc_data["type"]) }
  item.Item_data.Description = fmt.Sprintf("%s %s", tc_data["type"], tc_data["indicator"])

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

// This does not strip out <xml>, however holding xml doesn't appear to have xml tags
func ParseXML(xml_string string)(*etree.Document, error){
  xml_doc := etree.NewDocument()
  err := xml_doc.ReadFromString(xml_string)
  if err != nil { log.Println(err); return xml_doc, errors.New("Unable to read XML") }
  return xml_doc, nil
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

func df774Exists(bw Bib, mms_id string) bool {
  for _, d := range bw.Rec.Datafield{
    for _, s := range d.Subfield {
      if s.Value == mms_id { return true }
    }
  }
  return false
}
