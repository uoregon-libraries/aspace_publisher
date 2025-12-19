package as

import(
  "testing"
  "net/http/httptest"
  "net/http"
  "os"
  "fmt"
)

func TestAcquireJson(t *testing.T){
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("X-ArchivesSpace-Session") != "123456789" { t.Errorf("prob with sessionid") }

    if r.URL.Path == "/repositories/2/resource/1234" {
      fmt.Fprintf(w, "hello hello")
    } else if r.URL.Path == "/repositories/2/resource/3456" {
      w.WriteHeader(418)
      fmt.Fprintf(w, "mayday")
    }
  }))

  os.Setenv("ASPACE_URL", ts.URL + "/")
  os.Setenv("VERBOSE", "true")
  os.Setenv("DEBUG", "true")
  str_session := "123456789"
  str_rec1 := "resource/1234"
  str_rec2 := "resource/2345"
  str_rec3 := "resource/3456"

  response, err := AcquireJson(str_session, "2", str_rec1)
  if err != nil { t.Errorf(err.Error()) }
  if string(response) != "hello hello" { t.Errorf("wrong response") }


  response, err = AcquireJson(str_session, "2", str_rec3)
  if err.Error() != "aspace error exporting record" { t.Errorf(err.Error()) }
  if string(response) != "mayday" { t.Errorf("wrong response") }

  ts.Close()

  response, err = AcquireJson(str_session, "2", str_rec2)
  if err == nil {
    t.Errorf("there should be an error") 
  } else if err.Error() != "unable to complete request to archivesspace." {
    t.Errorf("problem with error") }
  if string(response) != "" { t.Errorf("wrong response") }
}

func TestAcquireMarc(t *testing.T){
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("X-ArchivesSpace-Session") != "123456789" { t.Errorf("prob with sessionid") }
    if v := r.URL.Query().Get("include_unpublished_marc"); v != "true" { t.Errorf("params wrong") }
    if r.URL.Path == "/repositories/2/resources/marc21/1234.xml" {
      fmt.Fprintf(w, "hello hello")
    } else if r.URL.Path == "/repositories/2/resources/marc21/3456.xml" {
      w.WriteHeader(418)
      fmt.Fprintf(w, "mayday")
    }
  }))

  os.Setenv("ASPACE_URL", ts.URL + "/")
  fmt.Println(ts.URL)
  os.Setenv("VERBOSE", "true")
  os.Setenv("DEBUG", "true")
  str_session := "123456789"
  str_rec1 := "1234"
  str_rec2 := "2345"
  str_rec3 := "3456"

  response, err := AcquireMarc(str_session, "2", str_rec1, "false")
  if err != nil { t.Errorf(err.Error()) }
  if string(response) != "hello hello" { t.Errorf("wrong response") }

  response, err = AcquireMarc(str_session, "2", str_rec3, "false")
  if err.Error() != "aspace error exporting MARC" { t.Errorf(err.Error()) }
  if string(response) != "mayday" { t.Errorf("wrong response") }

  ts.Close()

  response, err = AcquireMarc(str_session, "2", str_rec2, "false")
  if err == nil {
    t.Errorf("there should be an error") 
  } else if err.Error() != "unable to complete request to archivesspace" {
    t.Errorf("problem with error") }
  if string(response) != "" { t.Errorf("wrong response") }
}
