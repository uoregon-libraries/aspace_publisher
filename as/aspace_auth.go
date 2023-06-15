package as

import(
  "fmt"
  "net/http"
  "net/url"
  "io"
  "time"
  "encoding/json"
  "errors"
  "aspace_publisher/utils"
  "github.com/labstack/echo/v4"
)
type AuthResp struct {
  Session string
}

func As_basic(username, password string, c echo.Context) (bool, error){
  session_id, err := utils.FetchCookieVal(c, "as_session")
  if session_id == "" || err != nil {
    session_id, err := AuthenticateAS(username, password)
    if err != nil { return false, err }
    utils.WriteCookie(c, "as_session", session_id)
    time.Sleep(3 * time.Second)
  }
  return true, nil
}

//Note: this will work on the server. Or from a local machine using VPN
func AuthenticateAS(uname string, pass string) (string, error){
  var authresp AuthResp 
  authurl := fmt.Sprintf("https://archives.uoregon.edu/api/users/%s/login", uname)
  response, err := http.PostForm(authurl, url.Values{"password": {pass}})
  if err != nil { return "", errors.New("Unable to complete login") }

  defer response.Body.Close()
  byteVal, _ := io.ReadAll(response.Body)
  err = json.Unmarshal(byteVal, &authresp)
  if err != nil { return "", errors.New("Unable to extract session id") }

  return authresp.Session, nil
}
