package handlers

import(
  "fmt"
  "github.com/labstack/echo/v4"
  "aspace_publisher/echosession"
  "os"
  "aspace_publisher/utils"
  "aspace_publisher/aw"
  "aspace_publisher/as"
  "aspace_publisher/oclc"
)

func AsToOclcCreate(c echo.Context) error {

  ead_id := c.Param("id")
  repo_id := "2"
  store := echosession.FromContext(c)
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if err != nil { return echo.NewHTTPError(520, "Aspace authorization is in progress, please wait a moment and try request again.") }

  marc, err := as.AcquireMarc(session_id, repo_id, ead_id)
  if err != nil { fmt.Println(err); return err }

  token, err = oclc.OclcAuth()
  if err != nil { return echo.NewHTTPError(520, "Oclc authorization is in progress, please wait a moment and try request again.") }

  response, err = oclc.Create(token, marc)
  if err != nil { fmt.Println(err); return err }

  //response is xml, a marc record?
  //so unmarshal it into a marc record obj
  //extract oclc number
  //acquire aspace resource, which is in json
  //insert oclc
  //post resource json back to aspace

  var resp_marc MarcRecord
  resp_marc.Initialize(response)

  resource_json, err = as.AcquireJson(session_id, repo_id, ead_id)
  if err != nil { fmt.Println(err); return err }

  

  



  
}
