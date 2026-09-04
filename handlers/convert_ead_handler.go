package handlers

import(
  "github.com/labstack/echo/v4"
  "os"
  "aspace_publisher/utils"
  "aspace_publisher/aw"
  "fmt"
)


func ConvertEadHandler(c echo.Context) error {
  resource_id := c.Param("id")
  repo_id := "2"

  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(520, "Cannot retrieve session. Try redoing login.") }

    //get session for aw
  aw_session, err := aw.GetSession(c)
  if err != nil { return c.String(400, "could not authenticate with AWest") }

  fname,err := processEad(resource_id, session_id, repo_id, aw_session, "convert")

  base_url := os.Getenv("HOME_URL")
  return c.HTML(200, fmt.Sprintf("<p>Relevant updates will be written to <a href=\"%s/reports/%s\">%s</a></p>", base_url, fname, fname))

}
