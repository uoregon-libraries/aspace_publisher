package handlers

import(
  "github.com/labstack/echo/v4"
  "net/http"
  "io/ioutil"
)

func VersionHandler(c echo.Context) error {
  return c.String(http.StatusOK, read_version())
}

func read_version() string{
  content, err := ioutil.ReadFile("version.txt")
  if err != nil { return "could not read version" }
  return string(content)
}
