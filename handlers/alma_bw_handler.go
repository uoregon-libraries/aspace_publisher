package handlers

import (
  "github.com/labstack/echo/v4"
  "aspace_publisher/utils"
  "aspace_publisher/file"
  "aspace_publisher/alma"
  "os"
  "fmt"
)

func AlmaBWHandler(c echo.Context) error {
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(500, "Cannot retrieve session, try redoing login.") }
  esession, err := utils.Encrypt(session_id)
  if err != nil { return c.String(400, "Could not encrypt session") }
  fname := file.Filename()
  err = alma.CallWorker("startAlmaBWJob", map[string]string{ "id": c.Param("id"), "filename": fname, "session": esession } )
  if err != nil { return c.String(400, err.Error()) }
  base_url := os.Getenv("HOME_URL")
  return c.HTML(200, fmt.Sprintf("<p>Relevant updates will be written to <a href=\"%s/reports/%s\">%s</a></p>", base_url, fname, fname))
}