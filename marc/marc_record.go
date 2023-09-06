package marc

import (
  "encoding/xml"
  "os"
  "io/ioutil"
  "errors"
) 

type MarcRecord struct{
  OclcId string
  OclcDate string
  RecordXml *Record
}

func (mr *MarcRecord) Initialize(marc string) (string, error){
  mr.RecordXml = &Record{}
  err := mr.set_xml_from_string(marc); if err != nil { return "", err }
  mr.set_oclc_fields()

  return "", nil
}

func (mr *MarcRecord) set_xml_from_file(filename string) error{
  xmlfile, err1 := os.Open(filename); if err1 != nil { return err1 }
  byteValue, err2 := ioutil.ReadAll(xmlfile); if err2 != nil { return err2 }
  err3 := xml.Unmarshal(byteValue, mr.RecordXml); if err3 != nil { return err3 }
  xmlfile.Close()
  return nil
}

func (mr *MarcRecord) set_xml_from_string(marc string) error{
  err := xml.Unmarshal(marc, mr.RecordXml); if err3 != nil { return err3 }
  return nil
}

func (mr *MarcRecord) get_datafield(tag string) *DataField{
  dfs := mr.RecordXml.DataFields
  for i := len(dfs)-1; i > 0; i-- {
    if dfs[i].Tag == tag { return &dfs[i] }
  }
  return nil
}

func Get_subfield(code string, datafield *DataField) *SubField{
  sfs := datafield.SubFields
  for i := 0; i < len(sfs); i++ {
    if sfs[i].Code == code { return &sfs[i] }   
  }
  return nil
}

func (mr *MarcRecord) set_oclc_fields(){
  results := Extract_fields(mr.RecordXml)
  mr.OclcId = results[0]
  mr.OclcDate = results[1]
}

func Extract_fields(record *Record) [2]string{
  cfs := record.ControlFields
  var result [2]string
  for i := 0; i < len(cfs); i++ {
    if cfs[i].Tag == "001"{ 
      result[0] = cfs[i].Value
    } else if cfs[i].Tag == "005"{
      result[1] = cfs[i].Value
    }
  }
  return result
}

type SubField struct {
  XMLName xml.Name `xml:"subfield"`
  Code string `xml:"code,attr"`
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

