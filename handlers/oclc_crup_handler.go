package handlers

import(
  "github.com/labstack/echo/v4"
  "aspace_publisher/marc"
  "aspace_publisher/utils"
  "aspace_publisher/as"
  "aspace_publisher/oclc"
  "net/http"
)

func OclcCrupHandler(c echo.Context) error {

  id := c.Param("id")
  repo_id := "2"
  //get session id
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(520, "Aspace authorization is in progress, please wait a moment and try request again.") }
  //acquire aspace resource, which is in json
  json, err := as.AcquireJson(session_id, repo_id, id)
  if err != nil { return echo.NewHTTPError(400,  err) }
  //is it published?
  published, err := as.IsPublished(json)
  if err != nil { return echo.NewHTTPError(400, err) }

  //is it a new record?
  oclc_id := as.GetOclcId(json)
  
  //get MARC
  marc_rec, err := as.AcquireMarc(session_id, repo_id, id, published)
  if err != nil { return echo.NewHTTPError(400, err) }
  //strip outer tag
  marc_stripped, err := marc.StripOuterTags(marc_rec)
  if err != nil { return echo.NewHTTPError(500, err) }

  //authenticate with OCLC
  token, err := oclc.GetToken(c)
  if err != nil { return echo.NewHTTPError(520, err) }

  //edit the marc if updating
  if oclc_id != "" {
    oclc_marc, err := oclc.Record(token, oclc_id)
    if err != nil{ return echo.NewHTTPError(400, err) }
    edited_marc, err := marc.EditMarcForOCLC(as_marc, oclc_marc)
    if err != nil{ return echo.NewHTTPError(400, err) }
  }
  //post or put the marc
  var oclc_resp string
  if oclc_id == "" {
    oclc_resp, err = oclc.Request(token, as_marc, "manage/bibs", "", "marcxml+xml")
  } else {
    oclc_resp, err = oclc.Request(token, edited_marc, "manage/bibs", oclc_id, "marcxml+xml")
  }
  if err != nil {
    if oclc_resp != "" { return c.String(http.StatusOK, oclc_resp) } else {
      return echo.NewHTTPError(400, err) }
  }

  //if updating, done
  if oclc_id != "" {
    return c.String(http.StatusOK, oclc_resp)
  }

  oclc_id, err = marc.ExtractOclc(string(oclc_resp))
  if err != nil { return echo.NewHTTPError(500, err) }
  //insert oclc
  modified, err := as.UpdateUserDefined1(json, oclc_id)
  if err != nil { return echo.NewHTTPError(500,  err) }

  //post resource json back to aspace
  as_resp := as.Post(session_id, id, repo_id, id, string(modified))

  //print response to user
  return c.String(http.StatusOK, as_resp.ResponseToString())
}
