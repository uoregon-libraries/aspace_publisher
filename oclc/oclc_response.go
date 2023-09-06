package oclc

import(

)

type OclcResponse struct{
  OclcId string
  OclcDate string
}

func response_oclc_xml(marc string) (*Entry, error){
  xmlfile := strings.NewReader(marc)
  var entry Entry
  byteValue, _ := ioutil.ReadAll(xmlfile)
  xml.Unmarshal(byteValue, &entry)
  return &entry, nil
}

func (or *OclcResponse) set_fields(entry *Entry){
  result := as.Extract_fields(&entry.Content.Response.Record)
  or.OclcId = result[0]
  or.OclcDate = result[1]
}

type Response struct {
  XMLName xml.Name `xml:"response"`
  Record as.Record `xml:"record"`
}
type Content struct {
  XMLName xml.Name `xml:"content"`
  Response Response `xml:"response"`
}
type Entry struct {
  XMLName xml.Name `xml:"entry"`
  Content Content `xml:"content"`
  Id string `xml:"id"`
  Link string `xml:"link"`
}
