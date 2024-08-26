package as

import (
  "testing"
  "github.com/tidwall/gjson"
  "fmt"
  "os"
  "net/http/httptest"
  "net/http"
  "strings"
)

func TestExtractIdFromInstance(t *testing.T){
  do_string := `{"jsonmodel_type":"digital_object","linked_instances":[{"ref":"/repositories/2/archival_objects/17801"}],"digital_object_id":"Ax019_b013_f008"}`
  do_json := gjson.Parse(do_string)
  aoid := extractIdFromInstance(do_json)
  if aoid != "17801" { t.Fatalf("aoid is not correct") }

}

func TestExtractIdFromResponse(t *testing.T){
  resp_string := `{"status":"success","id":"12345","errors":[]}`
  r := Response{"abcde",BuildMessage(resp_string)}
  s := r.ResponseToString()
  doid := extractIdFromResponse(s)
  if doid != "12345" { t.Fatalf("doid is not correct") }
}

func TestUpdateWithInstance(t *testing.T){
  ao_string := `{"jsonmodel_type":"archival_object","instances":[{"instance_type":"mixed_materials","jsonmodel_type":"instance"}]}`
  ao_modified,_ := UpdateWithInstance([]byte(ao_string), Instance("/repositories/2/digital_objects/123"))
  do_instances := gjson.Get(string(ao_modified), "instances")
  arr := do_instances.Array()
  if len(arr) != 2 { t.Fatalf("adding to instances fail") }
  inst_type := gjson.Get(arr[1].String(), "instance_type")
  if inst_type.String() != "digital_object" { t.Fatalf("adding to instances fail") }
}

func TestPost(t *testing.T){
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("X-ArchivesSpace-Session") != "abcde_sessionstring_1234" { t.Errorf("request bungled adding session id to header") }
      fmt.Fprintf(w, `{ "status": "success", "id":"3456"}`)
  }))
  defer ts.Close()
  os.Setenv("ASPACE_URL", ts.URL + "/")

  response := Post("abcde_sessionstring_1234", "9876", "2", "resources/5432", "jsonrecordstring" )
  str_resp := response.ResponseToString()
  if !strings.Contains(str_resp, `"id":"9876"`) { t.Errorf("response is wrong") }
  if !strings.Contains(str_resp, `"id":"3456"`) { t.Errorf("response is wrong") }
  if !strings.Contains(str_resp, `"status":"success"`) { t.Errorf("response is wrong") }
}
