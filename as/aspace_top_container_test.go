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
  data := `{"page_size":250,"total_hits":2,"results":[{"id":"/repositories/2/top_containers/70040","uri":"/repositories/2/top_containers/70040","jsonmodel_type":null},{"id":"/repositories/2/top_containers/70041","uri":"/repositories/2/top_containers/70041","jsonmodel_type":null}`
  path := "/api/repositories/2/search"
   ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == path {
      if v1 := r.URL.Query().Get("type[]"); v1 != "top_container" { t.Errorf("params wrong") }
      if v2 := r.URL.Query().Get("q"); v2 != `collection_uri_u_sstr:"/repositories/2/resources/1234"` { t.Errorf("params wrong") }
      if v3 := r.URL.Query().Get("page"); v3 != "1" { t.Error("wrong page val") }
      fmt.Fprint(w, data)
    } else {
      t.Error("incorrect request url: " + r.URL.Path)
    }
  }))
  defer ts.Close()
  os.Setenv("ASPACE_URL", ts.URL + "/api/")
  list, err := TCList("432143214321", "2", "1234")
  if err != nil { log.Println(err) }
  if slices.Contains(list, "/repositories/2/top_containers/70040") != true { t.Errorf("incorrect top container list") }
  if slices.Contains(list, "/repositories/2/top_containers/70041") != true { t.Errorf("incorrect top container list") }
}

func TestCalculateRuns(t *testing.T){
  data := `{"page_size":250,"total_hits":33,"results":[{"id":"/repositories/2/top_containers/70040","uri":"/repositories/2/top_containers/70040","jsonmodel_type":null}]}`
  result, err := calculateRuns([]byte(data))
  if err != nil { t.Error(err) }
  if result != 1 { t.Error("incorrect num runs") }
  data2 := `{"page_size":250,"total_hits":383,"results":[{"id":"/repositories/2/top_containers/70040","uri":"/repositories/2/top_containers/70040","jsonmodel_type":null}]}`
  result, err = calculateRuns([]byte(data2))
  if err != nil { t.Error(err) }
  if result != 2 { t.Error("incorrect num runs") }
  data3 := `{"page_size":250,"total_hits":250,"results":[{"id":"/repositories/2/top_containers/70040","uri":"/repositories/2/top_containers/70040","jsonmodel_type":null}]}`
  result, err = calculateRuns([]byte(data3))
  if result != 1 { t.Error("incorrect num runs") }
  data4 := `{"error":"we don't understand what you're asking for here"}`
  result, err = calculateRuns([]byte(data4))
  if err == nil { t.Error("there should be an error") }

}

func TestProcessList(t *testing.T){
  data := `{"page_size":250,"total_hits":383,"results":[{"id":"/repositories/2/top_containers/70040","uri":"/repositories/2/top_containers/70040","jsonmodel_type":null},{"id":"/repositories/2/top_containers/70041","uri":"/repositories/2/top_containers/70041","jsonmodel_type":null}`
  list := processList([]byte(data))
  if slices.Contains(list, "/repositories/2/top_containers/70040") != true { t.Error("incorrect parse of list") }
  if slices.Contains(list, "/repositories/2/top_containers/70041") != true { t.Error("incorrect parse of list") }
}
func TestSearchString(t *testing.T){
  sstring := search_string("2", "1234", "1")
  expected := `search?type[]=top_container&q=collection_uri_u_sstr%3A%22%2Frepositories%2F2%2Fresources%2F1234%22&fields[]=id&page=1&page_size=250`
  if sstring != expected { t.Error("incorrect search string returned") }
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
  listdata := `{"page_size":250,"total_hits":2,"results":[{"id":"/repositories/2/top_containers/12345","uri":"/repositories/2/top_containers/12345","jsonmodel_type":null},{"id":"/repositories/2/top_containers/67890","uri":"/repositories/2/top_containers/67890","jsonmodel_type":null}]}`
  listpath := "/api/repositories/2/search"
  tcdata2 := topcontainer_fixture1
  tcpath1 := "/api/repositories/2/top_containers/12345"
  tcpath2 := "/api/repositories/2/top_containers/67890"
   ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    //If TCList is the requester
    if r.URL.Path == listpath {
      if v1 := r.URL.Query().Get("type[]"); v1 != "top_container" { t.Errorf("params wrong") }
      if v2 := r.URL.Query().Get("q"); v2 != `collection_uri_u_sstr:"/repositories/2/resources/987"` { t.Errorf("params wrong") }
      fmt.Fprint(w, listdata)
    } else if r.URL.Path == tcpath1 { //Iterating through the tclist
      w.WriteHeader(http.StatusNotFound)
      fmt.Fprintf(w, "{\"error\": \"Not found\"}")
    } else if r.URL.Path == tcpath2 {
      fmt.Fprint(w, tcdata2)
    } else {
      t.Errorf("incorrect request url")
    }
  }))
  defer ts.Close()
  os.Setenv("ASPACE_URL", ts.URL + "/api/")
  tcmap,msgs := ExtractTCData("123456789", "2", "987")
  if msgs[0] != "Unable to acquire 12345: aspace error exporting record" {
    t.Errorf("incorrect response to error")
  }
  if tcmap[0]["barcode"] != "35025042622635" {
    t.Errorf("incorrect tcmap returned")
  }
  if tcmap[0]["boundwith"] != "true" { t.Errorf("incorrect boundwith value") }
}

func TestIsBoundwith(t *testing.T){
  data1 := topcontainer_fixture1
  data2 := topcontainer_fixture2
  result1 := IsBoundwith([]byte(data1))
  if result1 != true { t.Errorf("incorrect result") }
  result2 := IsBoundwith([]byte(data2))
  if result2 != false { t.Errorf("incorrect result") }
}

func TestValidateTCData(t *testing.T){
  tcdata := TopContainer{ Barcode: "RONCO3200" }
  if ValidateTCData(tcdata) != true { t.Errorf("incorrect validation") }
}

func TestValidateBarcode(t *testing.T){
  if validateBarcode("ronco3200") != true { t.Errorf("incorrect validation") }
  if validateBarcode("ronco 3200") != false { t.Errorf("incorrect validation") }
}
