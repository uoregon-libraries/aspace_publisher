package oclc

import(
  "testing"
  "net/http/httptest"
  "net/http"
  "os"
  "fmt"
)

func TestSetHolding(t *testing.T){
   ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/manage/institution/holdings/123456789/set" { t.Fatalf("incorrect url") }
    if r.Method != "POST" { t.Fatalf("request action is not correct") }
    if r.Header.Get("accept") != "application/json" { t.Fatalf("request header is not correct")}
    if r.Header.Get("Content-Type") != "" { t.Fatalf("request header is not correct")}
    fmt.Fprintf(w, `{ "status":"updated", "id":"123456789" }`)
  }))
  defer ts.Close()
  os.Setenv("OCLC_URL", ts.URL)
  _, _ = SetHolding("123456789", "abcdabcdabcd12345")
}

