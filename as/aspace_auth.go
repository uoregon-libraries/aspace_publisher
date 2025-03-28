package as

import(
  "fmt"
  "net/http"
  "net/http/httputil"
  "net/url"
  "io"
  "os"
  "log"
  "time"
  "strings"
  "strconv"
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
    session_id, err = AuthenticateAS(username, password)
    if err != nil { return false, err }
    utils.WriteCookie(c, 120, "as_session", session_id)
  }
  return true, nil
}

//Note: this will work on the server. Or from a local machine using VPN
func AuthenticateAS(uname string, pass string) (string, error){
  var authresp AuthResp
  debug := os.Getenv("DEBUG")
  authurl := os.Getenv("ASPACE_URL") + fmt.Sprintf("users/%s/login", uname)
  data := url.Values{}
  data.Set("password", pass)
  request, err := http.NewRequest("POST", authurl, strings.NewReader(data.Encode()))
  request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
  request.Header.Set("Accept", "*/*")
  request.Header.Set("User-Agent", "curl/7.61.1")

  if debug == "true" {
    reqdump, err := httputil.DumpRequestOut(request, true)
    if err != nil { log.Println(err) } else {
      log.Printf("REQUEST:\n%s", string(reqdump)) }
  }
  if err != nil { log.Println(err); return "", errors.New("Unable to create login request") }
  client := http.Client{
	 Timeout: 60 * time.Second,
  }
  response, err := client.Do(request)
  if debug == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else {
      log.Printf("RESPONSE:\n%s", string(respdump)) }
  }
  if err != nil { log.Println(err); return "", errors.New("Unable to complete login to aspace") }
  if response.StatusCode != 200 { log.Println("unable to log into Aspace"); return "", errors.New("Unable to complete login to aspace") }
  defer response.Body.Close()
  byteVal, _ := io.ReadAll(response.Body)
  err = json.Unmarshal(byteVal, &authresp)
  if err != nil { return "", errors.New("Unable to extract session id") }

  return authresp.Session, nil
}
