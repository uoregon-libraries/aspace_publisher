package handlers

import(
  "github.com/labstack/echo/v4"
  "aspace_publisher/as"
  "aspace_publisher/utils"
  "aspace_publisher/oclc"
  "net/http"
  "os"
  "fmt"
)

func LauncherHandler(c echo.Context) error {
  resource_id := c.FormValue("resource_id")
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(500, "Cannot retrieve session, try redoing login.") }

  status := http.StatusOK

  err = as.ValidID(resource_id)
  fname := ""
  if err != nil { return c.String(400, "bad resource id") }
  workflow := c.FormValue("workflow")
  switch workflow {
  case "validate_ead":
    fname,err = validateEad(resource_id, session_id)
  case "upload_ead":
    fname,err = uploadEad(resource_id, session_id)
  case "validate_marc":
    //authenticate with OCLC
    oclc_token, err := oclc.GetToken(c)
    if err != nil { return c.String(400, "Could not authenticate with OCLC") }
    fname,err = validateMarc(resource_id, session_id, oclc_token)
  case "publish_marc":
    //authenticate with OCLC
    oclc_token, err := oclc.GetToken(c)
    if err != nil { return c.String(400, "Could not authenticate with OCLC") }
    fname,err = oclcCrup(resource_id, session_id, oclc_token)
  case "publish_alma":
    //authenticate with OCLC
    oclc_token, err := oclc.GetToken(c)
    if err != nil { return c.String(400, "Could not authenticate with OCLC") }
    fname,err = almaCrup(resource_id, validHolding(c), session_id, oclc_token)
  default:
    return c.String(500, "No workflow submitted")
  }
  if err != nil { status = http.StatusInternalServerError }
  base_url := os.Getenv("HOME_URL")
  return c.HTML(status, fmt.Sprintf("<p>Relevant updates will be written to <a href=\"%s/reports/%s\">%s</a></p>", base_url, fname, fname))
}
