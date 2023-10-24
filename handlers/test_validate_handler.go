package handlers

import(
  "log"
  "github.com/labstack/echo/v4"
  "os"
  "net/http"
  "aspace_publisher/utils"
  "aspace_publisher/aw"
  "aspace_publisher/as"
)


func TestValidateHandler(c echo.Context) error {
  ead_id := c.Param("id")
  repo_id := "2"
  verbose := os.Getenv("VERBOSE")
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(520, "Authorization is in progress, please wait a moment and try request again.") }

  ead_orig, err := as.AcquireEad(session_id, repo_id, ead_id, verbose)
  if err != nil { log.Println(err); return echo.NewHTTPError(400, ead_orig) }

  ead_prepped, _, ark, err := aw.PrepareEad(repo_id, ead_id, ead_orig)
  if err != nil { log.Println(err); return echo.NewHTTPError(400, "unable to prep ead") }

  ead_converted, err := aw.CallConversion(ead_prepped)
  if err != nil { log.Println(err); return echo.NewHTTPError(400, "unable to convert ead") }

  f, err := os.CreateTemp("", "ead-")
  if err != nil { log.Println(err); return echo.NewHTTPError(400, "unable to create temp dir") }
  defer f.Close()
  defer os.Remove(f.Name())
  _, err = f.Write([]byte(ead_converted))
  if err != nil { log.Println(err); return echo.NewHTTPError(400, "unable to write file") }

  //get session for aw
  aw_session, err := aw.GetSession(c, verbose)
  if err != nil { return echo.NewHTTPError(403, "Unable to complete ArchivesWest auth.") }
  vals, err := aw.MakeUploadMap(ark, "ead", f.Name())
  if err != nil { return echo.NewHTTPError(400, "Unable to create upload map.") }
  // create form
  form, boundary, err := utils.CreateMultipartFormData(vals)
  if err != nil { return echo.NewHTTPError(400, "Unable to create form.") }
  //validate
  response, err := aw.Validate(aw_session, boundary, verbose, form)

  if err != nil { return echo.NewHTTPError(400, "Unable to complete request.") }

  return c.String(http.StatusOK, response)
  //Use Inline or Attachment
  //return c.Inline(f.Name(), filename)
}
