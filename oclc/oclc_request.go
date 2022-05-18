package oclc

import(
  "net/http"
  "os"
  "io/ioutil"
  "bytes"
  "encoding/xml"
  as "aspace_publisher/aspace_export"
  "strings"
  "errors"
  "time"
)

type OclcRequest struct{
  ServiceUrl string
  Inst string
  Schema string
  HoldingCode string
  AspaceExport *as.AspaceExport
}

type OclcResponse struct{
  OclcId string
  OclcDate string
}

func (or *OclcRequest) Initialize(marc string) (string, error) {
  or.ServiceUrl = os.Getenv("OCLC_METADATA_SERVICE_URL")
  or.Schema = os.Getenv("SCHEMA")
  or.HoldingCode = os.Getenv("HOLDING_LIB_CODE")
  
  or.AspaceExport = &as.AspaceExport{ FileName: marc }
  str,e := or.AspaceExport.Initialize()
  if e != nil { return str, e }
  return "", nil
}

func (or OclcRequest) build_get_uri() (string){
  uri := or.ServiceUrl
  uri += "/" + or.AspaceExport.OclcId
  uri += "?classificationScheme=" + or.Schema
  uri += "&holdingLibraryCode=" + or.HoldingCode

  return uri
}

func (or OclcRequest) build_uri() (string){
  uri := or.ServiceUrl
  if or.AspaceExport.Protocol == "PUT" {
    uri += "/" + or.AspaceExport.OclcId
  }
  uri += "?classificationScheme=" + or.Schema
  uri += "&holdingLibraryCode=" + or.HoldingCode

  return uri
}

func (or OclcRequest) RequestPush() (*http.Request, error){
  payload, err := ioutil.ReadFile(or.AspaceExport.FileName); if err != nil { return nil, err }
  uri := or.build_uri()
  req, err := http.NewRequest(or.AspaceExport.Protocol, uri, bytes.NewBuffer(payload))
  if err != nil { return nil, err }
  req.Header.Set("content_type", "application/vnd.oclc.marc21+xml")
  return req, nil
}

func (or OclcRequest) RequestPull() (*http.Request, error){
  if or.AspaceExport.OclcId == "" { return nil, errors.New("no oclc to pull") }
  uri := or.build_get_uri()
  req, err := http.NewRequest("GET", uri, nil)
  return req, err
}

func DoRequest(req *http.Request)(*OclcResponse, error){
  client := &http.Client{
    Timeout: time.Second * 10,
  }
  response, err := client.Do(req); if err != nil { return nil, err }
  body, err := ioutil.ReadAll(response.Body); if err != nil { return nil, err }
  response.Body.Close()
  entry, err := response_oclc_xml(string(body)); if err != nil { return nil, err }
  var oresp OclcResponse
  oresp.set_fields(entry)
  return &oresp, nil
}

func AddToken(req *http.Request) (*http.Request, error){
  var ot OclcToken
  err := ot.GetToken(); if err != nil { return nil, err }
  req.Header.Set("Authorization", "Bearer " + ot.AccessToken)
  return req, nil
}

// requires path to file and xml type
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

