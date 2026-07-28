package handlers

import (
  "github.com/labstack/echo/v4"
  "aspace_publisher/utils"
  "aspace_publisher/as"
  "aspace_publisher/alma"
  "aspace_publisher/file"
  "net/http"
  "os"
  "fmt"
)

// Iterates through the tcmap, fetches item by barcode, 
func UpdateAspaceTCHandler(c echo.Context) error{
  var args alma.ProcessArgs
  args.Resource_id = c.Param("id")
  args.Repo_id = "2"
  args.Filename = file.Filename()
  var err error
  args.Session_id, err = utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(500, "Cannot retrieve session, try redoing login.") }

  fs := alma.FunMap{ UpdateTC: as.UpdateTC }
  tcmap, errmsgs := as.ExtractTCData(args.Session_id, args.Repo_id, args.Resource_id)
  if len(errmsgs) != 0 { file.WriteReport(args.Filename, errmsgs); return c.String(http.StatusInternalServerError, "Error, please see report.") }

  for _,tc := range tcmap{
    //don't overwrite. future to do: use param overwrite=true/false
    if tc["ils_item"] != "" && tc["ils_holding"] != "" { continue }
    item, err := alma.FetchByBarcode(tc["barcode"])
    if err != nil { errmsgs = append(errmsgs, "Unable to fetch barcode " + tc["barcode"]); continue }
    holding_id, item_id := alma.ParseHoldingItem(item)
    if item_id == "" || holding_id == "" { errmsgs = append(errmsgs, "Unable to extract ids for item with barcode " + tc["barcode"]); continue }
    err = fs.UpdateTC(args.Repo_id, holding_id, item_id, args.Session_id, tc)
    if err != nil { errmsgs = append(errmsgs, "TC update error " + tc["barcode"]); continue }
  }
  file.WriteReport(args.Filename, errmsgs)
  base_url := os.Getenv("HOME_URL")
  return c.HTML(http.StatusOK, fmt.Sprintf("<p>Relevant updates will be written to <a href=\"%s/reports/%s\">%s</a></p>", base_url, args.Filename, args.Filename))
}
