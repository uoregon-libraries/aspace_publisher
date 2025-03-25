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


func UploadEadHandler(c echo.Context) error {
  ead_id := c.Param("id")
  repo_id := "2"
  verbose := os.Getenv("VERBOSE")
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(520, "Cannot retrieve session, try redoing login.") }

  ead_orig, err := as.AcquireEad(session_id, repo_id, ead_id, verbose)
  if err != nil { log.Println(err); return echo.NewHTTPError(400, ead_orig) }

  ead_prepped, filename, ark, err := aw.PrepareEad(repo_id, ead_id, ead_orig)
  if err != nil { log.Println(err); return echo.NewHTTPError(400, "unable to prep ead") }

  ead_converted, err := aw.CallConversion(ead_prepped)
  if err != nil { log.Println(err); return echo.NewHTTPError(400, "unable to convert ead") }

  f, err := os.Create(filename)
  if err != nil { log.Println(err); return echo.NewHTTPError(400, "unable to create temp file") }
  defer f.Close()
  defer os.Remove(f.Name())
  _, err = f.Write([]byte(ead_converted))
  if err != nil { log.Println(err); return echo.NewHTTPError(400, "unable to write file") }

  //get session for aw
  aw_session, err := aw.GetSession(c, verbose)
  if err != nil { return echo.NewHTTPError(403, "Unable to complete ArchivesWest auth.") }
  vals, err := aw.MakeUploadMap(ark, "ead", f.Name())
  if err != nil { return echo.NewHTTPError(400, "Unable to create map object.") }
  // create form
  form, boundary, err := utils.CreateMultipartFormData(vals)
  if err != nil { return echo.NewHTTPError(400, "Unable to create upload form.") }
  //upload
  response, err := aw.Upload(aw_session, boundary, verbose, form)
  if err != nil { return echo.NewHTTPError(400, "Unable to complete request.") }

  parsed, err := aw.ParseResult(response)
  if err != nil { return echo.NewHTTPError(400, "Unable to parse response.") }
  return c.HTML(http.StatusOK, parsed)
}
