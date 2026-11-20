package as

import(
  "testing"
  "net/http/httptest"
  "net/http"
  "fmt"
  "os"
  "log"

)

func LargeTCMap()[]map[string]string{
  var tcmap []map[string]string
  for i := 0; i < 1000; i++ {
  tc := map[string]string { "type": "box", "indicator": "2", "barcode": "35025042622635", "boundwith": "false", "uri":"/repositories/2/top_containers/1234", "ils_holding":"22452547390001852", "ils_item":"23452547370001852", "update_refs":"false"}
  tcmap = append(tcmap, tc)
  }
  return tcmap
}

func BenchmarkLargeTCMap(b *testing.B){
  var tcmap []map[string]string
  for i := 0; i < b.N; i++ {
    tcmap = LargeTCMap()
  }
  fmt.Printf("************************************* %d", len(tcmap))
}

func TestTCListLarge(t * testing.T){
  data1 := `{"page_size":250,"total_hits":1000,"results":[{"id":"/repositories/2/top_containers/1","uri":"/repositories/2/top_containers/1","jsonmodel_type":null}`
  for i := 2; i <= 250; i++ {
    data1 += fmt.Sprintf(`,{"id":"/repositories/2/top_containers/%d", "url":"id":"/repositories/2/top_containers/%d","jsonmodel_type":null}`, i, i)
  }
  data1 += `]}`
  data2 := `{"page_size":250,"total_hits":1000,"results":[{"id":"/repositories/251/top_containers/1","uri":"/repositories/2/top_containers/251","jsonmodel_type":null}`
  for i := 252; i <= 500; i++ {
    data2 += fmt.Sprintf(`,{"id":"/repositories/2/top_containers/%d", "url":"id":"/repositories/2/top_containers/%d","jsonmodel_type":null}`, i, i)
  }
  data2 += `]}`
  data3 := `{"page_size":250,"total_hits":1000,"results":[{"id":"/repositories/2/top_containers/501","uri":"/repositories/2/top_containers/501","jsonmodel_type":null}`
  for i := 502; i <= 750; i++ {
    data3 += fmt.Sprintf(`,{"id":"/repositories/2/top_containers/%d", "url":"id":"/repositories/2/top_containers/%d","jsonmodel_type":null}`, i, i)
  }
  data3 += `]}`
  data4 := `{"page_size":250,"total_hits":1000,"results":[{"id":"/repositories/2/top_containers/751","uri":"/repositories/2/top_containers/751","jsonmodel_type":null}`
  for i := 752; i <= 1000; i++ {
    data4 += fmt.Sprintf(`,{"id":"/repositories/2/top_containers/%d", "url":"id":"/repositories/2/top_containers/%d","jsonmodel_type":null}`, i, i)
  }
  data4 += `]}`

  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Query().Get("page") == "1" { 
      fmt.Fprint(w, data1)
    } else if r.URL.Query().Get("page") == "2" {
      fmt.Fprint(w, data2)
    } else if r.URL.Query().Get("page") == "3" {
      fmt.Fprint(w, data3)
    } else {
      fmt.Fprint(w, data4)
    }
  }))
  defer ts.Close()
  os.Setenv("ASPACE_URL", ts.URL + "/api/")
  list, err := TCList("432143214321", "2", "1234")
  if err != nil { log.Println(err) }
  if len(list) != 1000 { t.Error("list is not correct") }
  if list[999] != "/repositories/2/top_containers/1000" { t.Error("list is not correct") }
}
