package aw

import(
  "os"
  "strings"
  "net/http"
  "errors"
  "log"
  "net/url"
  "strconv"
  "net/http/cookiejar"
  "net/http/httputil"
  "aspace_publisher/utils"
  "github.com/labstack/echo/v4"
)

func GetSession(c echo.Context, verbose bool) (string, error){  
  session_id, err := utils.FetchExpirableCookie(c, "aw_session")
  if session_id == "" || err != nil {
    session_id, err = authenticate(verbose)
    if err != nil { return "", err }
    utils.WriteCookie(c, 1, "aw_session", session_id)
  }
  return session_id, nil
}

func authenticate(verbose bool)(string, error){
  authurl, _ := url.Parse(os.Getenv("AWEST_URL") + "login.php")
  data := url.Values{}
  data.Set("username", os.Getenv("AWEST_NAME"))
  data.Set("password", os.Getenv("AWEST_PASS"))
  jar, err := cookiejar.New(nil)
  if err != nil { log.Println(err); return "", err }
  client := &http.Client{
    Jar: jar,
  }

  request, err := http.NewRequest("POST", authurl.String(), strings.NewReader(data.Encode()))
  request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

  if verbose {
    reqdump, err := httputil.DumpRequest(request, true)
    if err != nil { log.Println(err); return "", err }
    log.Printf("REQUEST:\n%s", string(reqdump))
  }
  if err != nil { log.Println(err); return "", errors.New("Unable to create login request") }
  response, err := client.Do(request)
  if err != nil { log.Println(err); return "", errors.New("Unable to complete login to awest") }

  if verbose {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err); return "", err }
    log.Printf("RESPONSE:\n%s", string(respdump))
  }
  defer response.Body.Close()

  session, err := parse_session(jar, authurl)
  if err != nil { return "", err }

  return session, nil
}

func parse_session(jar *cookiejar.Jar, url *url.URL) (string, error){
  for _, cookie := range jar.Cookies(url) {
    if cookie.Name == "PHPSESSID"{
      return cookie.Value, nil
    }
  }
  return "", errors.New("could not find session in cookies")
}
