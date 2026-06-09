package as

import (
  "testing"
  "github.com/tidwall/gjson"
  "fmt"
  "os"
  "net/http/httptest"
  "net/http"
)

func TestPost(t *testing.T){
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("X-ArchivesSpace-Session") != "abcde_sessionstring_1234" { t.Errorf("request bungled adding session id to header") }
      fmt.Fprintf(w, `{ "status": "success", "id":"3456"}`)
  }))
  defer ts.Close()
  os.Setenv("ASPACE_URL", ts.URL + "/")

  response := Post("abcde_sessionstring_1234", "9876", "2", "resources/5432", "jsonrecordstring" )
  str_resp := response.ResponseToString()
  value := gjson.Get(str_resp, "id")
  if value.String() != "9876" { t.Errorf("response is wrong") }
  value = gjson.Get(str_resp, "message.id")
  if value.String() != "3456" { t.Errorf("response is wrong") }
}

func TestBuildMessage(t *testing.T){
  mess := `{"id":"12345","warning":"the end of the world is nigh"}`
  resp := Response{ "67890", BuildMessage(mess) }
  r_str := resp.ResponseToString()
  value := gjson.Get(r_str, "id")
  if value.String() != "67890" { t.Errorf("id is wrong") }
  value = gjson.Get(r_str, "message.id")
  if value.String() != "12345" { t.Errorf("message is wrong") }
}
