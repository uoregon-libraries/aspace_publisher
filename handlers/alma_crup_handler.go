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
  "regexp"
)

func AlmaCrupHandler(c echo.Context) error {
  var args alma.ProcessArgs
  args.Resource_id = c.Param("id")
  args.Holding_id = validHolding(c)
  args.Repo_id = "2"
  args.Filename = file.Filename()
  var err error
  args.Session_id, err = utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(500, "Cannot retrieve session, try redoing login.") }

  //acquire aspace resource
  rjson, err := as.AcquireJson(args.Session_id, args.Repo_id, "resources/" + args.Resource_id)
  if err != nil { file.WriteReport(args.Filename, []string{ "Could not aquire JSON from aspace: " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}

  args.Oclc_id, err = as.GetOclcId(rjson)
  if err != nil { file.WriteReport(args.Filename, []string{ "Problem retrieving OCLC id: " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}
  if args.Oclc_id == "" { file.WriteReport(args.Filename, []string{ "Could not acquire OCLC id from resource" }); return c.String(http.StatusInternalServerError, "Error, please see report.")}
  //try for mms_id and create based on presence in resource json
  args.Mms_id, args.Create, err = as.GetMmsId(rjson)
  if err != nil { file.WriteReport(args.Filename, []string{ "Problem retrieving mms id: " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}
  //needed for holding record, appears as 099 in the aspace MARC but not OCLC's
  args.Id_0 = as.ExtractID0(rjson)

  //authenticate with OCLC
  args.Oclc_token, err = oclc.GetToken(c)
  if err != nil { file.WriteReport(args.Filename, []string{ "Could not authenticate with OCLC: " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}

  //get oclc marc
  oclc_marc, err := oclc.Record(args.Oclc_token, args.Oclc_id)
  if err != nil { file.WriteReport(args.Filename, []string{ "Could not acquire OCLC MARC " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.")}
  marc_clean := oclc.UnformatXML(oclc_marc)
  tcmap, errmsgs := as.ExtractTCData(args.Session_id, args.Repo_id, args.Resource_id)
  if len(errmsgs) != 0 { file.WriteReport(args.Filename, errmsgs); return c.String(http.StatusInternalServerError, "Error, please see report.") }

  tcmap, err = alma.CheckTCMap(tcmap)
  if err != nil { file.WriteReport(args.Filename, []string{ "Error attempting to acquire ids: " + err.Error() }); return c.String(http.StatusInternalServerError, "Error, please see report.") }
  //launch processing, starting with bib
  //eventually hand this off to a worker?
  fs := alma.FunMap{ BoundwithPF: alma.ProcessBoundwith, HoldingAPF: alma.ProcessHoldingA, HoldingBPF: alma.ProcessHoldingB, ItemsPF: alma.ProcessItems, ItemPF: alma.ProcessItem, AfterBib: as.AfterBibCreate, SetHolding: oclc.SetHolding, NZPF: alma.LinkToNetwork, UpdateTC: as.UpdateTC, CheckForMissingPF: alma.CheckItemsForMissing }
  alma.ProcessBib(args, marc_clean, rjson, tcmap, fs)

  base_url := os.Getenv("HOME_URL")
  return c.HTML(http.StatusOK, fmt.Sprintf("<p>Relevant updates will be written to <a href=\"%s/reports/%s\">%s</a></p>", base_url, args.Filename, args.Filename))
}

// hopefully unlikely case for this to be useful:
// processing has borked after the holding was created but NO items were created
// thus the holding id has not been added to any top containers
func validHolding(c echo.Context) string{
  h := c.QueryParam("holding")
  if h == "" { return "" }
  re1 := regexp.MustCompile(`[0-9]+`)
  matched1 := re1.Find([]byte(h))
  if string(matched1) == h {
    return h
  } else {
    return ""
  }
}
