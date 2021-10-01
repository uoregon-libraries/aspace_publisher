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
var oresp *OclcResponse

func Test(t *testing.T) {
  url := "/?inst=1234&classificationScheme=LibraryOfCongress&holdingLibraryCode=ABCD"
  marc, _ := ioutil.ReadFile("../files/marc_response.txt")

  server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    if r.URL.String() != url { 
      t.Errorf("request uri is incorrect")
      fmt.Println(r.URL.String())
    } else { w.Write(marc) }
  }))
  defer server.Close()
  c := server.Client()

  os.Setenv("INST", "1234")
  os.Setenv("SCHEMA", "LibraryOfCongress")
  os.Setenv("HOLDING_LIB_CODE", "ABCD")
  os.Setenv("OCLC_METADATA_SERVICE_URL", server.URL)
  oreq = &OclcRequest{}
  _, e := oreq.initialize("../files/marc.txt")
  if e != nil { 
    fmt.Sprintf("%s", e.Error())
  } else {fmt.Println( "oclc request initialize ok") }

  response, e2 := oreq.push(c)
  if e2 != nil { 
    fmt.Printf(e2.Error()) 
  } else { fmt.Println("oclc push ok") }
  
  entry, _ := response_oclc_xml(response)
  oresp = &OclcResponse{}
  oresp.set_fields(entry)
  
  if oresp.OclcId != "ocm18421997" { t.Errorf("response OCLC id is incorrect") }
  if oresp.OclcDate != "20210720192458.8" { t.Errorf("response OCLC date is incorrect") }
}
