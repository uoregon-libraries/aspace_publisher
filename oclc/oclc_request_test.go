package oclc

import(
  "testing"
  "io/ioutil"
  "os"
  "net/http"
  "net/http/httptest"
  "fmt"
)

var oreq *OclcRequest

func Test(t *testing.T) {
  url := "/ocm18421997?classificationScheme=LibraryOfCongress&holdingLibraryCode=ABCD"
  marc, _ := ioutil.ReadFile("../files/marc_response.txt")

  server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    if r.URL.String() != url { 
      fmt.Println(r.URL.String())
      t.Errorf("request uri is incorrect")
      fmt.Println(r.URL.String())
    } else { w.Write(marc) }
  }))
  defer server.Close()

  os.Setenv("SCHEMA", "LibraryOfCongress")
  os.Setenv("HOLDING_LIB_CODE", "ABCD")
  os.Setenv("OCLC_METADATA_SERVICE_URL", server.URL)
  oreq = &OclcRequest{}
  _, e := oreq.Initialize("../files/marc.txt")
  if e != nil {
    fmt.Sprintf("%s", e.Error())
  } else {fmt.Println( "oclc request initialize ok") }

  req, e := oreq.RequestPush()
  if e != nil {
    fmt.Printf(e.Error())
  } else { fmt.Println("oclc request push ok") }
  
  oresp, e := DoRequest(req)
  if e != nil {
    fmt.Printf(e.Error())
  } else { fmt.Println("oclc do request ok") }
  
  if oresp.OclcId != "ocm18421997" { t.Errorf("response OCLC id is incorrect") }
  if oresp.OclcDate != "20210720192458.8" { t.Errorf("response OCLC date is incorrect") }
}
