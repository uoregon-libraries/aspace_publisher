package handlers

import (
  "github.com/labstack/echo/v4"
  "aspace_publisher/utils"
  "aspace_publisher/as"
  "aspace_publisher/oclc"
  "aspace_publisher/alma"
  "aspace_publisher/file"
  "encoding/json"
  "net/http"
  "os"
  "fmt"
)

func AlmaCrupHandler(c echo.Context) error {
  id := c.Param("id")
  repo_id := "2"
  filename := file.Filename()
  //get session id
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(500, "Cannot retrieve session, try redoing login.") }

  //acquire aspace resource
  rjson, err := as.AcquireJson(session_id, repo_id, "resources/" + id)
  if err != nil { file.WriteReport(filename, []string{ "Could not aquire JSON from aspace: " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}

  oclc_id := as.GetOclcId(rjson)

  //try for mms_id and create based on presence in resource json
  mms_id := as.GetMmsId(rjson)
  create := true
  if mms_id != "" { create = false }
  //needed for holding record, appears as 099 in the aspace MARC but not OCLC's
  id_0 := as.ExtractID0(rjson)

  //authenticate with OCLC
  token, err := oclc.GetToken(c)
  if err != nil { file.WriteReport(filename, []string{ "Could not authenticate with OCLC: " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}

  //get oclc marc
  oclc_marc, err := oclc.Record(token, oclc_id)
  if err != nil { file.WriteReport(filename, []string{ "Could not acquire OCLC MARC " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}

  //create bib, holding, items
  mms_id, err = alma.ProcessBib(mms_id, oclc_marc, create)
  if err != nil { file.WriteReport(filename, []string{ "Could not create Alma Bib " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}

  var holding_id = ""
  if create == false { holding_id = alma.GetHoldingId(mms_id) }
  holding_id, err = alma.ProcessHolding(mms_id, holding_id, oclc_marc, id_0, create)
  if err != nil { file.WriteReport(filename, []string{ "Could not create Alma Holding: " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}

  itemlist := []string{}
  tclist,err := as.TCList(session_id, repo_id, id) //get the top containers
  if err != nil { file.WriteReport(filename, []string{ "Unable to acquire TC list: " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}
  for _,tc_path := range tclist{
    tc_id := as.ExtractID(tc_path)
    jsonTC, err := as.AcquireJson(session_id, repo_id, "top_containers/" + tc_id)
    if err != nil { file.WriteReport(filename, []string{ "Unable to acquire TC json " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}
    item_id, _ := as.GetTCRefs(jsonTC)
    var tc as.TopContainer
    err = json.Unmarshal(jsonTC, &tc)
    if err != nil { file.WriteReport(filename, []string{ "Unable to process TC json: " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}
    item_id, err = alma.ProcessItem(mms_id, holding_id, item_id, tc.Mapify(), create)
    if err != nil { file.WriteReport(filename, []string{ "Unable to process Alma item: " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}
    itemlist = append(itemlist, item_id)
    if create {
      err = as.UpdateTC(repo_id, tc_id, jsonTC, holding_id, item_id, session_id)
      if err != nil { file.WriteReport(filename, []string{ "Unable to update TC in aspace: " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}
    }
  }
  if create {
    //update the aspace resource
    modified, err := as.UpdateUserDefined2(rjson, mms_id)
  if err != nil { file.WriteReport(filename, []string{ err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}
    as.UpdateResource(session_id, "2", id, string(modified))

    //todo: switch to worker.
    alma.LinkToNetwork([]string{ mms_id }, filename)
  }
  base_url := os.Getenv("HOME_URL")
  return c.HTML(http.StatusOK, fmt.Sprintf("<p>Relevant updates will be written to <a href=\"%s/reports/%s\">%s</a></p>", base_url, filename, filename))
}
