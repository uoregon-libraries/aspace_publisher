package handlers

import (
  "github.com/labstack/echo/v4"
  "aspace_publisher/utils"
  "aspace_publisher/as"
  "aspace_publisher/oclc"
  "aspace_publisher/alma"
  "encoding/json"
  "net/http"
)

func AlmaCrupHandler(c echo.Context) error {
  id := c.Param("id")
  repo_id := "2"
  //get session id
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(520, "Cannot retrieve session, try redoing login.") }

  //acquire aspace resource
  rjson, err := as.AcquireJson(session_id, repo_id, "resources/" + id)
  if err != nil { return echo.NewHTTPError(400,  err) }

  oclc_id := as.GetOclcId(rjson)

  //try for mms_id and create based on presence in resource json
  mms_id := as.GetMmsId(rjson)
  create := true
  if mms_id != "" { create = false }

  //authenticate with OCLC
  token, err := oclc.GetToken(c)
  if err != nil { return echo.NewHTTPError(520, err) }

  //get oclc marc
  oclc_marc, err_ := oclc.Record(token, oclc_id)
  if err_ != nil{ return echo.NewHTTPError(400, err_) }

  //create bib, holding, items
  mms_id, err = alma.ProcessBib(mms_id, oclc_marc, create)
  if err != nil { return echo.NewHTTPError(400, err) }

  var holding_id = ""
  if create == false { holding_id = alma.GetHoldingId(mms_id) }
  holding_id, err = alma.ProcessHolding(mms_id, holding_id, oclc_marc, create)

  itemlist := []string{}
  tclist,err := as.TCList(session_id, repo_id, mms_id) //get the top containers
  if err != nil { return echo.NewHTTPError(400, err) }
  for _,tc_path := range tclist{
    tc_id := as.ExtractID(tc_path)
    jsonTC, err := as.AcquireJson(session_id, repo_id, "top_containers/" + tc_id)
    if err != nil { return echo.NewHTTPError(400, err) }
    item_id, _ := as.GetTCRefs(jsonTC)
    var tc as.TopContainer
    err = json.Unmarshal(jsonTC, &tc)
    if err != nil { return echo.NewHTTPError(400, err) }
    if err != nil { return echo.NewHTTPError(400, err) }
    item_id, err = alma.ProcessItem(mms_id, holding_id, item_id, tc.Mapify(), create)
    itemlist = append(itemlist, item_id)
    if create {
      err = as.UpdateTC(repo_id, tc_id, jsonTC, holding_id, item_id, session_id)
      if err != nil { return echo.NewHTTPError(400, err) }
    }
  //use itemlist in reporting
  itemlist = append(itemlist, item_id)
  }
  if create {
    //update the aspace resource
    modified, err := as.UpdateUserDefined2(rjson, mms_id)
    if err != nil { return echo.NewHTTPError(400, err) }
    as.UpdateResource(session_id, "2", id, string(modified))
    filename := ""//todo: add filename and reporting
    // alma job, will take a bit longer so run last; todo: switch to worker.
    alma.LinkToNetwork([]string{ mms_id }, filename)
  }
  return c.String(http.StatusOK, "ok")
}
