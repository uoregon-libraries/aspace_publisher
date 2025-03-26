package handlers

import(
  "github.com/labstack/echo/v4"
  "log"
  "io"
  "aspace_publisher/as"
  "aspace_publisher/utils"
  "net/http"
//  "github.com/tidwall/gjson"
)

func UploadDigitalObjectsHandler(c echo.Context) error {
  //get file
  file, _ := c.FormFile("file")
  src, err := file.Open()
  if err != nil { log.Println(err); return c.String(http.StatusBadRequest, "Could not open file") }
  defer src.Close()
  body, err := io.ReadAll(src)
  if err != nil { log.Println(err); return c.String(http.StatusBadRequest, "Could not open file") }

  //get session id
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(520, "Aspace authorization is in progress, please wait a moment and try request again.") }

  //hand uploaded file to as
  result := as.CreateDigitalObjects(string(body), session_id)
  str_result := result.ResponsesToString()
  return c.String(http.StatusOK, str_result)
}
