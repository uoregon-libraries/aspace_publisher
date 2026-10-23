package as

import (
  "testing"
  "net/http"
  "net/http/httptest"
  "os"
  "fmt"
  "reflect"
  "io/ioutil"
  "encoding/json"
)

func TestProcessEvent(t *testing.T){
  expected := fmt.Sprintf(event_fixture1, ConstructTime())
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    fmt.Println(expected)
    if err != nil { t.Errorf("error reading request body") }
    if compareEvents(body, []byte(expected)) != true { t.Errorf("incorrectly constructed event") }
    fmt.Fprint(w, `{ "yolo": "" }`)
  }))
  defer ts.Close()
  os.Setenv("ASPACE_URL", ts.URL + "/api/")
  _ = ProcessEvent("abcde123456", "/agents/people/99", "123", "2", "EAD published")
}

func compareEvents(e1, e2 []byte)bool{
  var event1 Event
  var event2 Event
  json.Unmarshal(e1, &event1)
  json.Unmarshal(e2, &event2)
  return reflect.DeepEqual(event1, event2)
}
