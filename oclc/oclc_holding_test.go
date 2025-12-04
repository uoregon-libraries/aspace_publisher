package oclc

import(
  "testing"
  "strings"
  "regexp"
  "net/http/httptest"
  "net/http"
  "os"
)

func TestSetHolding(t *testing.T){
   ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.path != "manage/institution/holdings/abcdabcdabcd12345/set" { t.Fatalf("incorrect url") }
    if r.Method != "POST" { t.Fatalf("request action is not correct") }
    if r.Header.Get("accept") != "json" { t.Fatalf("request header is not correct")}
    if r.Header.Get("Content-Type") != "" { t.Fatalf("request header is not correct")}
  }))
  defer ts.Close()
  os.Setenv("OCLC_URL", ts.URL + "/")
  _,_ := SetHolding('123456789', 'abcdabcdabcd12345')
}

