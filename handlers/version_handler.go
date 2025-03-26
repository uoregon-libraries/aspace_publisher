package handlers

import(
  "github.com/labstack/echo/v4"
  "net/http"
  "io/ioutil"
  "log"
  "os"
)

func VersionHandler(c echo.Context) error {
  return c.String(http.StatusOK, read_version())
}

func read_version() string{
  path := os.Getenv("HOME_DIR")
  content, err := ioutil.ReadFile(path + "version.txt")
  if err != nil { log.Println(err); return "could not read version" }
  return string(content)
}
