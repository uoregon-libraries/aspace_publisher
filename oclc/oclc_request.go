package oclc

import(
  "net/http"
  "os"
  "io/ioutil"
  "bytes"
  "encoding/xml"
  as "aspace_publisher/aspace_export"
  "strings"
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

func (or *OclcRequest) initialize(marc string) (string, error) {
  or.ServiceUrl = os.Getenv("OCLC_METADATA_SERVICE_URL")
  or.Inst = os.Getenv("INST")
  or.Schema = os.Getenv("SCHEMA")
  or.HoldingCode = os.Getenv("HOLDING_LIB_CODE")
  
  or.AspaceExport = &as.AspaceExport{ FileName: marc }
  str,e := or.AspaceExport.Initialize()
  if e != nil { return str, e }
  return "", nil
}

func (or OclcRequest) build_uri(get bool) string{
  uri := or.ServiceUrl
  uri += "?inst=" + or.Inst
  uri += "&classificationScheme=" + or.Schema
  uri += "&holdingLibraryCode=" + or.HoldingCode
  if get { uri += "&oclcNumber=" + or.AspaceExport.OclcId }
  return uri
}

func manage_oclc_push(marc string) (*OclcResponse, error){
  var oreq OclcRequest
  _, err1 := oreq.initialize(marc); if err1 != nil { return nil, err1 }
  client := authenticated_client()
  response2, err2 := oreq.push(client);  if err2 != nil { return nil, err2 }
  entry, err3 := response_oclc_xml(response2); if err3 != nil { return nil, err3 }
  var oresp OclcResponse
  oresp.set_fields(entry)
  return &oresp, nil
}

//post and put requests to oclc
func (or OclcRequest) push(client *http.Client) (string, error){
  payload, err := ioutil.ReadFile(or.AspaceExport.FileName)
  if err != nil { return or.AspaceExport.AspaceId, err }
  uri := or.build_uri(false) /*******************set the bool properly **************/
  req, err := http.NewRequest(or.AspaceExport.Protocol, uri, bytes.NewBuffer(payload))
  req.Header.Set("content_type", "application/vnd.oclc.marc21+xml")
  if err != nil { return or.AspaceExport.AspaceId, err }
  resp, err := client.Do(req)
  if err != nil{ return or.AspaceExport.AspaceId, err }
  body, err := ioutil.ReadAll(resp.Body) 
  if err != nil { return or.AspaceExport.AspaceId, err }
  resp.Body.Close()
  return string(body), nil
}

//get requests to oclc
//mainly to acquire date
/*func (or OclcRequest) pull(client *http.Client) (string, error){
*
*}
*/

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

