package aw

import (
  "testing"
  "os"
  "path/filepath"
  "strings"
  "net/http"
  "net/http/httptest"
  "fmt"
)

func TestParseResult(t *testing.T) {
  home_dir := os.Getenv("HOME_DIR")
  html_resp1, err := os.ReadFile(filepath.Join(home_dir, "fixtures/resp1.html"))
  if err != nil { t.Error(err.Error()) }
  r1 := strings.NewReader(string(html_resp1))
  result1, err := ParseResult(r1)
  if err != nil { t.Error(err.Error()) }
  if !strings.Contains(result1, "success") { t.Errorf("Parser failed to find node.") }

  html_resp2, err := os.ReadFile(filepath.Join(home_dir, "fixtures/resp2.html"))
  if err != nil { t.Error(err.Error()) }
  r2 := strings.NewReader(string(html_resp2))
  result2, err := ParseResult(r2)
  if err != nil { t.Error(err.Error()) }
  if !strings.Contains(result2, "errors") { t.Errorf("Parser failed to find node.") }
}

func TestValidArk(t *testing.T){
  if ValidArk("/80444/xv[fill in ARK here]") { t.Error("invalid ark") }
  if ValidArk("80444/xv123") != true { t.Error("incorrect result") }
  if ValidArk("/80444/xv345") != true { t.Error("incorrect result") }
}

func TestExtractArk(t *testing.T){
  ark, err := ExtractArk(ead_fixture1)
  if err != nil { t.Error(err) }
  if ark != "/80444/xv123" { t.Error("incorrect result") }
}

func TestCheckArk(t *testing.T){
  home_dir := os.Getenv("HOME_DIR")
  html_resp, err := os.ReadFile(filepath.Join(home_dir, "fixtures/short_finding.html"))
  if err != nil { t.Error(err) }
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/ark:80444/xv123" { t.Error("incorrect path") }
    fmt.Fprint(w, string(html_resp))
  }))
  defer ts.Close()
  os.Setenv("ARK_URL_BASE", ts.URL + "/ark:")
  result, err := CheckArk("80444/xv123")
  if err != nil { t.Error(err) }
  if result != true { t.Error("incorrect result") }
}
