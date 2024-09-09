package oclc

import(
  "testing"
  "strings"
  "regexp"
  "net/http/httptest"
  "net/http"
  "os"
)

func TestAssembleUrl(t *testing.T){
  parts_no_id := []string{"https://blah.org","path",""}
  parts_id := []string{"https://blah.org","path","abcd"}
  url_no_id := assembleUrl(parts_no_id)
  url_id := assembleUrl(parts_id)
  if url_no_id != "https://blah.org/path" {t.Fatalf("assembled url is incorrect")}
  if url_id != "https://blah.org/path/abcd" {t.Fatalf("assembled url is incorrect")}
}

func TestRequestCreate(t *testing.T){
   ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" { t.Fatalf("request action is not correct") }
    if r.Header.Get("accept") != "application/marcxml+xml" { t.Fatalf("request header is not correct")}
    if r.Header.Get("Content-Type") != "application/marcxml+xml" { t.Fatalf("request header is not correct")}
  }))
  defer ts.Close()
  os.Setenv("OCLC_URL", ts.URL + "/")
  marc := `<record></record>`
  _,_ = Request("token", marc, "manage/bibs", "", "marcxml+xml")
}

func TestRequestUpdate(t *testing.T){
   ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Method != "PUT" { t.Fatalf("request action is not correct") }
    if r.Header.Get("accept") != "application/marcxml+xml" { t.Fatalf("request header is not correct")}
    if r.Header.Get("Content-Type") != "application/marcxml+xml" { t.Fatalf("request header is not correct")}
    arr := strings.Split(r.URL.String(), "/")
    id := arr[len(arr)-1]
    matched, _ := regexp.Match(`[0-9]`, []byte(id))
    if matched != true { t.Fatalf("id is not present") }
  }))
  defer ts.Close()
  os.Setenv("OCLC_URL", ts.URL + "/")
  marc := `<record></record>`
  _,_ = Request("token", marc, "manage/bibs", "12345678", "marcxml+xml")
}

func TestRequestValidate(t *testing.T){
   ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" { t.Fatalf("request action is not correct") }
    if r.Header.Get("accept") != "application/json" { t.Fatalf("request header is not correct")}
    if r.Header.Get("Content-Type") != "application/marcxml+xml" { t.Fatalf("request header is not correct")}
  }))
  defer ts.Close()
  os.Setenv("OCLC_URL", ts.URL + "/")
  marc := `<record></record>`
  _,_ = Request("token", marc, "validate/validateFull", "", "json")
}
