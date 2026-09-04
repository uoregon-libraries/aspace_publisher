package handlers

import(
  "github.com/labstack/echo/v4"
  "aspace_publisher/utils"
  "aspace_publisher/oclc"
  "aspace_publisher/as"
  "fmt"
  "os"
)

func OclcValidateHandler(c echo.Context) error {
  id := c.Param("id")
  repo_id := "2"
  //get session id
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(520, "Cannot retrieve session, try redoing login.") }
  //authenticate with OCLC
  oclc_token, err := oclc.GetToken(c)
  if err != nil { return echo.NewHTTPError(520, err) }
  fm := oclc.FunMap{AsCheckIsPublished: as.CheckIsPublished, AsAcquireMarc: as.AcquireMarc, AsGetOclcId: as.GetOclcId, OclcRequest: oclc.Request, OclcRecord: oclc.Record, AsPost: as.Post}

  fname, err := validateMarc(id, repo_id, session_id, oclc_token, fm)

  base_url := os.Getenv("HOME_URL")
  return c.HTML(200, fmt.Sprintf("<p>Relevant updates will be written to <a href=\"%s/reports/%s\">%s</a></p>", base_url, fname, fname))
}
