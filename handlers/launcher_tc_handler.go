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
    e_session, err := utils.Encrypt(session_id)
    if err != nil { return c.String(400, "Could not encrypt session") }
    fname = file.Filename()
    err = alma.CallWorker("startAlmaBWJob", map[string]string{ "id": top_container_id, "filename": fname, "session": e_session } )
    if err != nil { return c.String(400, err.Error()) }
  default:
    return c.String(500, "No workflow submitted")
  }
  base_url := os.Getenv("HOME_URL")
  return c.HTML(status, fmt.Sprintf("<p>Relevant updates will be written to <a href=\"%s/reports/%s\">%s</a></p>", base_url, fname, fname))
}
