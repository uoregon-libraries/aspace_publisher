package handlers

import(
  "github.com/labstack/echo/v4"
  "aspace_publisher/as"
  "aspace_publisher/utils"
  "net/http"
  "os"
  "fmt"
)

func LauncherTCHandler(c echo.Context) error {
  top_container_id := c.FormValue("top_container_id")
  repo_id := "2"
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(500, "Cannot retrieve session, try redoing login.") }

  status := http.StatusOK
  fname := ""

  err = as.ValidID(top_container_id)
  if err != nil { return c.String(500, "bad tc id") }
  workflow := c.FormValue("workflow")
  switch workflow {
  case "update_single_tc":
    fname = UpdateTopContainerSA(top_container_id, session_id, repo_id)
  case "boundwith":
    //do nothing for now
    return c.String(status, "This process not yet enabled")
  default:
    return c.String(500, "No workflow submitted")
  }
  //if error has been documented in the user report, return the filename
  if err != nil { status = http.StatusInternalServerError }
  base_url := os.Getenv("HOME_URL")
  return c.HTML(status, fmt.Sprintf("<p>Relevant updates will be written to <a href=\"%s/reports/%s\">%s</a></p>", base_url, fname, fname))
}
