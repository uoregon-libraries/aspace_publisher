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
  os.Setenv("ASPACE_URL", ts.URL + "/api/")
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

func TestExtractTCData(t *testing.T){
  listdata := "[{\"ref\":\"/repositories/2/top_containers/12345\"},{\"ref\":\"/repositories/2/top_containers/67890\"}]"
  listpath := "/api/repositories/2/resources/987/top_containers"
  tcdata2 := `{"barcode":"35025042622635","ils_holding_id":"22452547390001852","ils_item_id":"23452547370001852","indicator":"[35025042622635]","created_by":"alexagoff","last_modified_by":"alexagoff","create_time":"2023-11-21T20:00:48Z","system_mtime":"2025-12-19T23:32:52Z","user_mtime":"2025-11-03T21:50:00Z","created_for_collection":"/repositories/2/resources/9634","type":"Multiple Collection Box","indicator":"[35025042622635]","collection":[{"ref":"/repositories/2/resources/3340","identifier":"PH 200_333","display_string":"Tony Minthorn collection of Nez Perce photographs"},{"ref":"/repositories/2/resources/3639","identifier":"PH 363","display_string":"Ronda Skubi collection of Lord and Schryver Slides"}],"uri":"/repositories/2/top_containers/67890"}`
  tcpath1 := "/api/repositories/2/top_containers/12345"
  tcpath2 := "/api/repositories/2/top_containers/67890"
   ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == listpath {
      fmt.Fprintf(w, listdata)
    } else if r.URL.Path == tcpath1 {
      w.WriteHeader(http.StatusNotFound)
      fmt.Fprintf(w, "{\"error\": \"Not found\"}")
    } else if r.URL.Path == tcpath2 {
      fmt.Fprintf(w, tcdata2)
    } else {
      t.Errorf("incorrect request url")
    }
  }))
  defer ts.Close()
  os.Setenv("ASPACE_URL", ts.URL + "/api/")
  tcmap,msgs := ExtractTCData("123456789", "2", "987")
  if msgs[0] != "Unable to acquire TC json aspace error exporting record" {
    t.Errorf("incorrect response to error")
  }
  if tcmap[0]["barcode"] != "35025042622635" {
    t.Errorf("incorrect tcmap returned")
  }
  if tcmap[0]["boundwith"] != "true" { t.Errorf("incorrect boundwith value") }
}

func TestIsBoundwith(t *testing.T){
  data1 := `{"barcode":"35025042622635","ils_holding_id":"22452547390001852","ils_item_id":"23452547370001852","indicator":"[35025042622635]","created_by":"alexagoff","last_modified_by":"alexagoff","create_time":"2023-11-21T20:00:48Z","system_mtime":"2025-12-19T23:32:52Z","user_mtime":"2025-11-03T21:50:00Z","created_for_collection":"/repositories/2/resources/9634","type":"Multiple Collection Box","indicator":"[35025042622635]","collection":[{"ref":"/repositories/2/resources/3340","identifier":"PH 200_333","display_string":"Tony Minthorn collection of Nez Perce photographs"},{"ref":"/repositories/2/resources/3639","identifier":"PH 363","display_string":"Ronda Skubi collection of Lord and Schryver Slides"}],"uri":"/repositories/2/top_containers/67890"}`

  data2 := `{"barcode":"35025042622635","ils_holding_id":"22452547390001852","ils_item_id":"23452547370001852","indicator":"[35025042622635]","created_by":"alexagoff","last_modified_by":"alexagoff","create_time":"2023-11-21T20:00:48Z","system_mtime":"2025-12-19T23:32:52Z","user_mtime":"2025-11-03T21:50:00Z","created_for_collection":"/repositories/2/resources/9634","type":"Box","indicator":"[35025042622635]","collection":[{"ref":"/repositories/2/resources/3340","identifier":"PH 200_333","display_string":"Tony Minthorn collection of Nez Perce photographs"}],"uri":"/repositories/2/top_containers/67890"}`
  result1 := IsBoundwith([]byte(data1))
  if result1 != true { t.Errorf("incorrect result") }
  result2 := IsBoundwith([]byte(data2))
  if result2 != false { t.Errorf("incorrect result") }
}
