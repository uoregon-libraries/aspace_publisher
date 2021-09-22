package aspace_export

import (
  "encoding/xml"
  "os"
  "io/ioutil"
  "errors"
) 

type AspaceExport struct{
  FileName string
  AspaceId string
  Protocol string
  OclcId string
  OclcDate string
  RecordXml *Record
}

func (ae AspaceExport) Initialize() (string, error){
  err := ae.set_xml(); if err != nil { return ae.FileName, err }
  err = ae.set_AspaceId(); if err != nil { return ae.FileName, err }
  ae.set_oclc_fields()
  ae.set_protocol()
  return "", nil
}

func (ae AspaceExport) set_xml() error{
  xmlfile, err1 := os.Open(ae.FileName); if err1 != nil { return err1 }
  defer xmlfile.Close()
  byteValue, err2 := ioutil.ReadAll(xmlfile); if err2 != nil { return err2 }
  err3 := xml.Unmarshal(byteValue, ae.RecordXml); if err3 != nil { return err3 }
  return nil
}

func (ae AspaceExport) set_AspaceId() error{
  datafield := ae.get_datafield("856")
  if datafield == nil { return errors.New("Cannot set ApaceId") }
  subfield := get_subfield("u", datafield)
  if subfield == nil { return errors.New("Cannot set ApaceId") }
  ae.AspaceId = subfield.Value
  return nil
}

func (ae AspaceExport) get_datafield(tag string) *DataField{
  dfs := ae.RecordXml.DataFields
  for i := len(dfs)-1; i > 0; i-- {
    if dfs[i].Tag == tag { return &dfs[i] }
  }
  return nil
}

func get_subfield(code string, datafield *DataField) *SubField{
  sfs := datafield.SubFields
  for i := 0; i < len(sfs); i++ {
    if sfs[i].Code == code { return &sfs[i] }   
  }
  return nil
}

func (ae AspaceExport) set_oclc_fields(){
  cfs := ae.RecordXml.ControlFields
  for i := 0; i < len(cfs); i++ {
    if cfs[i].Tag == "001"{ 
      ae.OclcId = cfs[i].Value
    } else if cfs[i].Tag == "005"{
      ae.OclcDate = cfs[i].Value
    }
  }
}

func (ae AspaceExport) set_protocol() {
  ae.Protocol = "POST"
  if ae.OclcId == "" {
    ae.Protocol = "PUT"
  }
}

type SubField struct {
  XMLName xml.Name `xml:"subfield"`
  Code string `xml:"tag,attr"`
  Value string `xml:",chardata"`
}
type DataField struct {
  XMLName xml.Name `xml:"datafield"`
  Tag string `xml:"tag,attr"`
  SubFields []SubField `xml:"subfield"`
}
type ControlField struct {
  XMLName xml.Name `xml:"controlfield"`
  Tag string `xml:"tag,attr"`
  Value string `xml:",chardata"`
}
type Record struct {
  XMLName xml.Name `xml:"record"`
  Leader string `xml:"leader"`
  ControlFields []ControlField `xml:"controlfield"`
  DataFields []DataField `xml:"datafield"`
}
