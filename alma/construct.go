package alma

import (
  "aspace_publisher/marc"
  "github.com/beevik/etree"
  "encoding/xml"
  "log"
  "fmt"
  "errors"
)

func ConstructBib(marc_string string)string{
  marc_stripped, err := marc.StripOuterTags(marc_string)
  if err != nil { log.Println(err); return "" }
  return "<bib>" + marc_stripped + "</bib>"
}

func ConstructHolding(marc_string string)string{
  marc_xml, err := ParseMarc(marc_string)
  if err != nil { log.Println(err); return "" }
  call_num := ExtractCall(marc_xml)
  link := BuildFindingLink(marc_xml)
  fixed := ExtractFixed(marc_xml)
  var h = Holding{}
  h.Suppress = false
  var rec Record
  rec.Leader = ExtractLeader(marc_xml)
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
  if err != nil { log.Println(err); return "" }
  return string(output)
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
func BuildFindingLink(marc_xml *etree.Document)string{
  url := marc_xml.FindElement("//datafield[@tag='856']/subfield[@code='u']").Text()
  text := marc_xml.FindElement("//datafield[@tag='856']/subfield[@code='z']").Text()
  link := fmt.Sprintf("<a href='%s' target='_blank'>%s</a>", url, text)
  if text == "Notice of Interest in Unprocessed Collections" {
    link = "UNARRANGED COLLECTION UNAVAILABLE FOR USE. Inquiries regarding these materials should be submitted via the " + link
  }
  return link
}

//leader
func ExtractLeader(marc_xml *etree.Document) string{
  leader := marc_xml.FindElement("//leader").Text()
  return leader
}

//df 099 sf a -> df 852 sf h
func ExtractCall(marc_xml *etree.Document)string{
  call := marc_xml.FindElement("//datafield[@tag='099']/subfield[@code='a']").Text()
  return call
}

//CF 008 -> 008
func ExtractFixed(marc_xml *etree.Document)string{
  fixed := marc_xml.FindElement("//controlfield[@tag='008']").Text()
  return fixed
}
