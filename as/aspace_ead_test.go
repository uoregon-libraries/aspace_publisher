package as

import(
  "testing"
  "net/http/httptest"
  "net/http"
  "os"
  "fmt"
  "strings"
)

func TestAcquireEad(t *testing.T){
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("X-ArchivesSpace-Session") != "123456789" { t.Errorf("prob with sessionid") }
    if v := r.URL.Query().Get("include_unpublished"); v != "False" { t.Errorf("params wrong") }
    if v := r.URL.Query().Get("include_daos"); v != "True" { t.Errorf("params wrong") }
    if v := r.URL.Query().Get("numbered_cs"); v != "True" { t.Errorf("params wrong") }
    if v := r.URL.Query().Get("ead3"); v != "False" { t.Errorf("params wrong") }
    if r.URL.Path == "/repositories/2/resource_descriptions/1234.xml" {
      fmt.Fprintf(w, "hello hello")
    } else if r.URL.Path == "/repositories/2/resource_descriptions/3456.xml" {
      w.WriteHeader(418)
      fmt.Fprintf(w, "mayday")
    }
  }))

  os.Setenv("ASPACE_URL", ts.URL + "/")
  os.Setenv("VERBOSE", "true")
  os.Setenv("DEBUG", "true")
  str_session := "123456789"
  str_rec1 := "1234"
  str_rec2 := "2345"
  str_rec3 := "3456"

  response, err := AcquireEad(str_session, "2", str_rec1, "true")
  if err != nil { t.Errorf(err.Error()) }
  if string(response) != "hello hello" { t.Errorf("wrong response") }

  response, err = AcquireEad(str_session, "2", str_rec3, "false")
  if err.Error() != "problem retrieving ead" { t.Errorf(err.Error()) }
  if string(response) != "mayday" { t.Errorf("wrong response") }

  ts.Close()

  response, err = AcquireEad(str_session, "2", str_rec2, "true")
  if err == nil {
    t.Errorf("there should be an error") 
  } else if strings.Contains(err.Error(), "connect: connection refused") == false {
    t.Errorf("problem with error") }
  if string(response) != "" { t.Errorf("wrong response") }
}
