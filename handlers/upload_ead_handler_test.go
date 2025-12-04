package handlers

import (
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"
  "os"
  "fmt"
  "github.com/labstack/echo/v4"
)

func TestUploadEadHandler(t *testing.T){
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/repositories/2/resource_descriptions/3456.xml" {
      w.WriteHeader(418)
      fmt.Fprintf(w, "mayday")
    } else {
      w.WriteHeader(503)
      fmt.Fprintf(w, "not found")
    }
  }))

  os.Setenv("ASPACE_URL", ts.URL + "/")
  os.Setenv("VERBOSE", "true")
  os.Setenv("DEBUG", "true")
  str_session := "123456789"

  e := echo.New()

  rec := httptest.NewRecorder()
  http.SetCookie(rec, &http.Cookie{Name: "as_session", Value: str_session})

  req := httptest.NewRequest(echo.GET, "/", nil)
  req.Header = http.Header{"Cookie": rec.HeaderMap["Set-Cookie"]}

  c := e.NewContext(req, rec)
  c.SetPath("/ead/upload/:id")
  c.SetParamNames("id")
  c.SetParamValues("3456")

  err := UploadEadHandler(c)
  if err != nil {
    if !strings.Contains(err.Error(), "mayday") { t.Errorf("wrong response") }
  }
}
