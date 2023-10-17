package handlers

import(
  "github.com/labstack/echo/v4"
  "os"
  "log"
  "aspace_publisher/utils"
  "aspace_publisher/aw"
  "aspace_publisher/as"
)


func ConvertEadHandler(c echo.Context) error {
    ead_id := c.Param("id")
    repo_id := "2"

    session_id, err := utils.FetchCookieVal(c, "as_session")
    if err != nil { return echo.NewHTTPError(520, "Authorization is in progress, please wait a moment and try request again.") }

    ead_orig, err := as.AcquireEad(session_id, repo_id, ead_id)
    if err != nil { log.Println(err); return err }

    ead_prepped, filename, _, err := aw.PrepareEad(repo_id, ead_id, ead_orig)
    if err != nil { log.Println(err); return err }

    ead_converted, err := aw.CallConversion(ead_prepped)
    if err != nil { log.Println(err); return err }

    f, err := os.CreateTemp("", "ead-")
    if err != nil { log.Println(err); return err }
    defer f.Close()
    defer os.Remove(f.Name())
    _, err = f.Write([]byte(ead_converted))
    if err != nil { log.Println(err); return err }
    //Use Inline or Attachment
    return c.Inline(f.Name(), filename)
  }


