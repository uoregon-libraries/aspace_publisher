package handlers

import (
  "github.com/labstack/echo/v4"
  "fmt"
  "path/filepath"
  "path"
)

func DisplayReportsHandler(c echo.Context) error {
  page := "<h3>Reports</h3>"
  reports, err := filepath.Glob("views/reports/*")
  if err != nil { return c.HTML(500, err.Error())}
  for _, r := range reports {
    report := path.Base(r)
    page += fmt.Sprintf("<p><a href=\"/reports/%s\">%s</a></p>", report, report)
  }
  return c.HTML(200, page)
}