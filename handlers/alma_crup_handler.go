package handlers

import (
  "github.com/labstack/echo/v4"
  "aspace_publisher/utils"
  "aspace_publisher/oclc"
  "net/http"
  "os"
  "fmt"
  "regexp"
)

func AlmaCrupHandler(c echo.Context) error {
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(500, "Cannot retrieve session, try redoing login.") }

  //authenticate with OCLC
  oclc_token, err := oclc.GetToken(c)
  if err != nil { return c.String(400, "Could not authenticate with OCLC") }

  status := http.StatusOK
  fname,err := almaCrup(c.Param("id"), validHolding(c), session_id, oclc_token)
  if err != nil { status = http.StatusInternalServerError }
  base_url := os.Getenv("HOME_URL")
  return c.HTML(status, fmt.Sprintf("<p>Relevant updates will be written to <a href=\"%s/reports/%s\">%s</a></p>", base_url, fname, fname))
}

// hopefully unlikely case for this to be useful:
// processing has borked after the holding was created but NO items were created
// thus the holding id has not been added to any top containers
func validHolding(c echo.Context) string{
  h := c.QueryParam("holding")
  if h == "" { return "" }
  re1 := regexp.MustCompile(`[0-9]+`)
  matched1 := re1.Find([]byte(h))
  if string(matched1) == h {
    return h
  } else {
    return ""
  }
}
