package handlers

import (
  "github.com/labstack/echo/v4"
  "aspace_publisher/utils"
  "aspace_publisher/as"
  "aspace_publisher/oclc"
  "aspace_publisher/alma"
  "aspace_publisher/file"
  "net/http"
  "os"
  "fmt"
)

func AlmaCrupHandler(c echo.Context) error {
  var args alma.ProcessArgs
  args.Resource_id = c.Param("id")
  args.Repo_id = "2"
  args.Filename = file.Filename()
  var err error
  args.Session_id, err = utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(500, "Cannot retrieve session, try redoing login.") }

  //acquire aspace resource
  rjson, err := as.AcquireJson(args.Session_id, args.Repo_id, "resources/" + args.Resource_id)
  if err != nil { file.WriteReport(args.Filename, []string{ "Could not aquire JSON from aspace: " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}

  args.Oclc_id = as.GetOclcId(rjson)
  //try for mms_id and create based on presence in resource json
  args.Mms_id, args.Create = as.GetMmsId(rjson)
  //needed for holding record, appears as 099 in the aspace MARC but not OCLC's
  args.Id_0 = as.ExtractID0(rjson)

  //authenticate with OCLC
  args.Oclc_token, err = oclc.GetToken(c)
  if err != nil { file.WriteReport(args.Filename, []string{ "Could not authenticate with OCLC: " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}

  //get oclc marc
  oclc_marc, err := oclc.Record(args.Oclc_token, args.Oclc_id)
  if err != nil { file.WriteReport(args.Filename, []string{ "Could not acquire OCLC MARC " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}

  tcmap, errmsgs := as.ExtractTCData(args.Session_id, args.Repo_id, args.Resource_id)
  if len(errmsgs) != 0 { file.WriteReport(args.Filename, errmsgs); return c.String(http.StatusInternalServerError, "Error, please see report.") }
  //launch processing, starting with bib
  //eventually hand this off to a worker?
  fs := alma.FunMap{ BoundwithPF: alma.ProcessBoundwith, HoldingPF: alma.ProcessHolding, ItemsPF: alma.ProcessItems, ItemPF: alma.ProcessItem, AfterBib: as.AfterBibCreate }
  alma.ProcessBib(args, oclc_marc, rjson, tcmap, fs)

  base_url := os.Getenv("HOME_URL")
  return c.HTML(http.StatusOK, fmt.Sprintf("<p>Relevant updates will be written to <a href=\"%s/reports/%s\">%s</a></p>", base_url, args.Filename, args.Filename))
}
