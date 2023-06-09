package handlers

import(
  "aspace_publisher/as"
  "aspace_publisher/aw"
  "aspace_publisher/utils"
  "io/ioutil"
  "fmt"
  "os"
  "github.com/labstack/echo/v4"
)

func Get_ead_handler(c echo.Context) error {
    ead_id := c.Param("id")
    repo_id := "2"

    session_id, err := utils.FetchCookieVal(c, "as_session")
    if err!= nil { fmt.Println(err); return err }

    ead_orig, err := as.AcquireEad(session_id, repo_id, ead_id)
    if err != nil { fmt.Println(err); return err }

    f, err := os.CreateTemp("", "orig-")
    if err != nil { fmt.Println(err); return err }
    defer f.Close()
    defer os.Remove(f.Name())
    _, err = f.Write([]byte(ead_orig))
    if err != nil { fmt.Println(err); return err }
    filename := "acquired.xml"
    return c.Attachment(f.Name(), filename)

}

func Prep_ead_handler(c echo.Context) error {
    ead_id := c.Param("id")
    repo_id := "2"

    ead_orig, err := ioutil.ReadFile("ead_orig.txt")
    if err != nil { fmt.Println(err); return err }
    
    ead_prepped, filename, err := aw.PrepareEad(repo_id, ead_id, string(ead_orig))
    if err != nil { fmt.Println(err); return err }
    fmt.Println(filename)

    f, err := os.CreateTemp("", "prepped-")
    if err != nil { fmt.Println(err); return err }
    defer f.Close()
    defer os.Remove(f.Name())
    _, err = f.Write([]byte(ead_prepped))
    if err != nil { fmt.Println(err); return err }
    
    return c.Attachment(f.Name(), filename)

}

func Php_ead_handler(c echo.Context) error {

    ead_prepped, err := ioutil.ReadFile("ead_prepped.txt")
    if err != nil { fmt.Println("failed read"); return err }

    ead_converted, err := aw.CallConversion(string(ead_prepped))
    if err != nil { fmt.Println("failed convert"); return err }

    f, err := os.CreateTemp("", "converted-")
    if err != nil { fmt.Println(err); return err }
    defer f.Close()
    defer os.Remove(f.Name())
    _, err = f.Write([]byte(ead_converted))
    if err != nil { fmt.Println(err); return err }
    filename := "converted.xml"
    return c.Attachment(f.Name(), filename)
}


