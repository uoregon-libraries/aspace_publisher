package as

import (
  "testing"
  "os"
  "net/http"
  "net/http/httptest"
  "io/ioutil"
  "encoding/json"
  "slices"
  "fmt"
  "path/filepath"
)

func TestLightweightPreflight(t *testing.T){

  rpath := "/api/repositories/2/resources/987"
  resource := `{ "user_defined": { "string_1": "123123123", "string_2":"987987987" }`
  listpath := "/api/repositories/2/top_containers/search"
  listdata := `{"response":{"docs":[{"id":"/repositories/2/top_containers/123", "collection_uri_u_sstr":["/repositories/2/resources/987"]},{"id":"/repositories/2/top_containers/456", "collection_uri_u_sstr":["/repositories/2/resources/987"]},{"id":"/repositories/2/top_containers/789", "collection_uri_u_sstr":["/repositories/2/resources/987"]}]}}`
  tcpath123 := "/api/repositories/2/top_containers/123"
  tcpath456 := "/api/repositories/2/top_containers/456"
  tcpath789 := "/api/repositories/2/top_containers/789"
  tcdata123 := `{"barcode":"35025041622635","ils_holding_id":"22452547380001852","ils_item_id":"23452547370001852","uri":"/repositories/2/top_containers/123","collection":[{"ref":"/repositories/2/resources/3340","identifier":"PH 200_333","display_string":"Tony Minthorn collection of Nez Perce photographs"},{"ref":"/repositories/2/resources/3639","identifier":"PH 363","display_string":"Ronda Skubi collection of Lord and Schryver Slides"}]}`
  tcdata456 := `{"barcode":"35025042622635","ils_holding_id":"","ils_item_id":"","uri":"/repositories/2/top_containers/456"}`
  tcdata789 := `{"barcode":"35025043622635","ils_holding_id":"22452547390001852","ils_item_id":"23452547380001852","uri":"/repositories/2/top_containers/789"}`

  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("path " + r.URL.Path)
    switch r.URL.Path {
    case rpath:
      fmt.Fprint(w, resource)
    case listpath:
      fmt.Fprint(w, listdata)
    case tcpath123:
      fmt.Fprint(w, tcdata123)
    case tcpath456:
      fmt.Fprint(w, tcdata456)
    case tcpath789:
      fmt.Fprint(w, tcdata789)
    }
  }))
  defer ts.Close()
  os.Setenv("ASPACE_URL", ts.URL + "/api/")
  homedir := os.Getenv("HOME_DIR")
  filename := LightweightPreflight("1234abcde", "987", "2")
  
  var p Preflight
  jstr,err := ioutil.ReadFile(filepath.Join(homedir, "views/reports/" + filename))
  if err != nil { t.Error(err) }
  err = json.Unmarshal(jstr, &p)
  if err != nil { t.Error(err) }
  if p.CreateBib { t.Error("incorrect bib val") }
  if p.CreateHolding { t.Error("incorrect hold val") }
  if !slices.Contains(p.BoundwithBarcodes, "35025041622635") { t.Error("incorrect boundwith val") }
  if !slices.Contains(p.CreateItemsBarcodes, "35025042622635") { t.Error("incorrect create item val") }
  if !slices.Contains(p.UpdateItemsBarcodesPids, "35025043622635/23452547380001852") { t.Error("incorrect update item val") }
}