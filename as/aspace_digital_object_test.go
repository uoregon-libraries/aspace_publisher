package as

import (
  "testing"
  "github.com/tidwall/gjson"
  "fmt"
  "os"
  "net/http/httptest"
  "net/http"
  "strings"
  "path/filepath"
)

func TestCreateDigitalObjects(t *testing.T){
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
      fmt.Fprintf(w, `{"title":"Ahoy", "instances":[{"instance_type":"thing"}]}`)
    } else if strings.Contains(r.URL.String(), "/repositories/2/digital_objects") {
      fmt.Fprintf(w, `{ "status":"created", "id":"3456"}`)
    } else if strings.Contains(r.URL.String(), "repositories/2/archival_objects/26462") {
      fmt.Fprintf(w, `{ "status":"updated", "id":"26462" }`)
    } else { t.Fatalf("post is not going to correct url") }
  }))
  defer ts.Close()
  os.Setenv("ASPACE_URL", ts.URL + "/")
  home_dir := os.Getenv("HOME_DIR")
  fpath := filepath.Join(home_dir, "fixtures/do_output6.json")
  jstr, err := os.ReadFile(fpath)
  if err != nil { fmt.Println(err) }

  responses := CreateDigitalObjects(string(jstr), "abcde_session_1234")
  if responses.responses[0].id != "Ax078_b012_f009" { t.Fatalf("response id is not correct") }
  if responses.responses[1].id != "archival_objects/26462" { t.Fatalf("response id is not correct") }
}

func TestExtractRefPath(t *testing.T){
  do_string := `{"jsonmodel_type":"digital_object","linked_instances":[{"ref":"/repositories/2/archival_objects/17801"}],"digital_object_id":"Ax019_b013_f008"}`
  do_json := gjson.Parse(do_string)
  aoid := extractRefPathFromResult(do_json)
  if aoid != "/repositories/2/archival_objects/17801" { t.Fatalf("aoid is not correct") }

  aoid2 := extractRefPathFromString(do_string)
  if aoid2 != "/repositories/2/archival_objects/17801" {t.Fatalf("aoid is not correct") }
}

func TestValidateRefId(t *testing.T){
  str := "/repositories/2/archival_objects/17801"
  result := validateRefId(str)
  if result != "archival_objects/17801" { t.Fatalf("ref id is not correct") }

  str = "/repositories/2/17801"
  result = validateRefId(str)
  if result != "" { t.Fatalf("ref id is not correct") }
}

func TestExtractIdFromResponse(t *testing.T){
  resp_string := `{"status":"success","id":"12345","errors":[]}`
  r := Response{"abcde",BuildMessage(resp_string)}
  s := r.ResponseToString()
  doid := extractIdFromResponse(s)
  if doid != "12345" { t.Fatalf("doid is not correct") }
}

func TestModify(t *testing.T){
  home_dir := os.Getenv("HOME_DIR")
  fpath := filepath.Join(home_dir, "fixtures/do_output6.json")
  jstr, err := os.ReadFile(fpath)
  if err != nil { fmt.Println(err) }
  item := gjson.GetBytes(jstr, "digital_objects.0")
  result, err := modify(item)
  if err != nil { t.Fatalf("modify error") }
  fmt.Print(result)
  if !strings.Contains(result, "/repositories/2/archival_objects/26462") {
    t.Fatalf("output is not expected value")
  }
}
