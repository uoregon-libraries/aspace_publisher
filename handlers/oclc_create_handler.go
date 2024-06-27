package handlers

import(
  "github.com/labstack/echo/v4"
  "aspace_publisher/marc"
  "aspace_publisher/utils"
  "aspace_publisher/as"
  "aspace_publisher/oclc"
  "net/http"
)

func OclcCreateHandler(c echo.Context) error {

  id := c.Param("id")
  repo_id := "2"
  //get session id
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(520, "Aspace authorization is in progress, please wait a moment and try request again.") }
  //get MARC
  marc_rec, err := as.AcquireMarc(session_id, repo_id, id)
  if err != nil { return echo.NewHTTPError(400, err) }
  //strip outer tag
  marc_stripped, err := marc.StripOuterTags(marc_rec)
  if err != nil { return echo.NewHTTPError(500, err) }

  //authenticate with OCLC
  token, err := oclc.GetToken(c)
  if err != nil { return echo.NewHTTPError(520, "Oclc authorization is in progress, please wait a moment and try request again.") }
  //push MARC to OCLC
  oclc_resp, err := oclc.Create(token, marc_stripped)
  if err != nil { return echo.NewHTTPError(400, err) }

  //extract oclc number
  oclc, err := marc.ExtractOclc(string(oclc_resp))
  if err != nil { return echo.NewHTTPError(500, err) }
  //acquire aspace resource, which is in json
  json, err := as.AcquireJson(session_id, repo_id, id)
    if err != nil { return echo.NewHTTPError(400,  err) }
  //insert oclc
  modified, err := as.UpdateUserDefined1(json, oclc)
  if err != nil { return echo.NewHTTPError(500,  err) }

  //post resource json back to aspace
  as_resp, err := as.UpdateResource(session_id, repo_id, id, string(modified))
  if err != nil { return echo.NewHTTPError(400, err) }
  //print response to user
  return c.String(http.StatusOK, as_resp)
}
