package as

import (
  "testing"
  "os"
  "net/http"
  "net/http/httptest"
  "fmt"
  "log"
  "slices"
)

func TestTCList(t *testing.T){
  data := "[{\"ref\":\"/repositories/2/top_containers/12345\"},{\"ref\":\"/repositories/2/top_containers/67890\"}]"
  path := "/api/repositories/2/resources/987/top_containers"
   ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == path {
      fmt.Fprintf(w, data)
    } else {
      t.Errorf("incorrect request url")
    }
  }))
  defer ts.Close()
  os.Setenv("ASPACE_URL", ts.URL + "/api")
  list, err := TCList("432143214321", "2", "987")
  if err != nil { log.Println(err) }
  if slices.Contains(list, "/repositories/2/top_containers/12345") != true { t.Errorf("incorrect top container list") }
  if slices.Contains(list, "/repositories/2/top_containers/67890") != true { t.Errorf("incorrect top container list") }
}

func TestGetTCRefs(t *testing.T){
  data := "{ \"ils_holding_id\":\"22299219100001852\",\"ils_item_id\":\"23299219080001852\" }"
  item, holding := GetTCRefs([]byte(data))
  if item != "23299219080001852" { t.Errorf("incorrect item id") }
  if holding != "22299219100001852" { t.Errorf("incorrect_holding_id") }
}

func TestUpdateIlsIds(t *testing.T){
  data := "{ \"barcode\":\"123456789\"}"
  modified, err := UpdateIlsIds([]byte(data), "321321321", "456456456")
  if err != nil { log.Println(err) }
  if string(modified) != "{ \"barcode\":\"123456789\",\"ils_holding_id\":\"321321321\",\"ils_item_id\":\"456456456\"}" { t.Errorf("incorrect ils refs") }
}

func TestMapify(t *testing.T){
  tc := TopContainer{Barcode: "123456", Indicator: "yellow", Type: "square"}
  tmap := tc.Mapify()
  if tmap["barcode"] != "123456" { t.Errorf("barcode is incorrect") }
  if tmap["indicator"] != "yellow" { t.Errorf("indicator is incorrect") }
  if tmap["type"] != "square" { t.Errorf("type is incorrect") }
}
