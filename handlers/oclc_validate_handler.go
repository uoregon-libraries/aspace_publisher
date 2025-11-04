package handlers

import(
  "github.com/labstack/echo/v4"
  "aspace_publisher/utils"
  "aspace_publisher/as"
  "aspace_publisher/oclc"
  "aspace_publisher/marc"
  "net/http"
)

func OclcValidateHandler(c echo.Context) error {
  id := c.Param("id")
  repo_id := "2"
  //get session id
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(520, "Cannot retrieve session, try redoing login.") }
  //get aspace resource
  json, err := as.AcquireJson(session_id, repo_id, "resources/" + id)
  if err != nil {
    if len(json) != 0 {
      return echo.NewHTTPError(400, json) } else {
      return echo.NewHTTPError(400, err)
    }
  }

  //is it published?
  published, err := as.IsPublished(json)
  if err != nil { return echo.NewHTTPError(400, err) }
  //get MARC
  marc_rec, err := as.AcquireMarc(session_id, repo_id, id, published)
  if err != nil {
    if marc_rec != "" {
      return echo.NewHTTPError(400, marc_rec) } else {
      return echo.NewHTTPError(400, err)
    }
  }

  //strip outer tag
  marc_stripped, err := marc.StripOuterTags(marc_rec)
  if err != nil { return echo.NewHTTPError(400, err) }
  //authenticate with OCLC
  token, err := oclc.GetToken(c)
  if err != nil { return echo.NewHTTPError(520, err) }
  //push MARC to OCLC
  oclc_resp, err := oclc.Request(token, marc_stripped, "manage/bibs/validate/validateFull", "","json")
  if err != nil {
    if oclc_resp != "" {
      return c.String(http.StatusOK, oclc_resp) } else {
      return echo.NewHTTPError(400, err) }
  }
  return c.String(http.StatusOK, oclc_resp)
}
