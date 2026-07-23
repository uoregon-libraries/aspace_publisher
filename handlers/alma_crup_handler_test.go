package handlers

import (
  "net/http/httptest"
  "net/url"
  "testing"
  "fmt"
  "github.com/labstack/echo/v4"
)

func TestValidHolding(t *testing.T){
  e := echo.New()
  rec := httptest.NewRecorder()
  //set up the query part of the request url
  q := make(url.Values)
  q.Set("holding", "7890123456")
  //add the query to the request like so
  req := httptest.NewRequest(echo.GET, "/?"+q.Encode(), nil)
  c := e.NewContext(req, rec)
  str := validHolding(c)
  fmt.Println(str)
  if str != "7890123456" { t.Errorf("wrong response") }
}

func TestValidHoldingEmpty(t *testing.T){
  e := echo.New()
  rec := httptest.NewRecorder()
  req := httptest.NewRequest(echo.GET, "/", nil)
  c := e.NewContext(req, rec)
  str := validHolding(c)
  if str != "" { t.Errorf("wrong response") }
}
