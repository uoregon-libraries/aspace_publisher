package handlers

import(
  "fmt"
  "github.com/labstack/echo/v4"
  "os"
  "aspace_publisher/utils"
  "aspace_publisher/aw"
  "aspace_publisher/as"
)


func ConvertEadHandler(c echo.Context) error {
    ead_id := c.Param("id")
    repo_id := "2"

    session_id, err := utils.FetchCookieVal(c, "as_session")
    if err!= nil { fmt.Println(err); return err }

    ead_orig, err := as.AcquireEad(session_id, repo_id, ead_id)
    if err != nil { fmt.Println(err); return err }

    ead_prepped, filename, err := aw.PrepareEad(repo_id, ead_id, ead_orig)
    if err != nil { fmt.Println(err); return err }
    fmt.Println(filename)

    ead_converted, err := aw.CallConversion(ead_prepped)
    if err != nil { fmt.Println(err); return err }

    f, err := os.CreateTemp("", "ead-")
    if err != nil { fmt.Println(err); return err }
    defer f.Close()
    defer os.Remove(f.Name())
    _, err = f.Write([]byte(ead_converted))
    if err != nil { fmt.Println(err); return err }
    
    return c.Attachment(f.Name(), filename)

  return c.String(200, "OK")
  }


