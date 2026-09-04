package aw

import(
  "os"
  "strings"
  "net/http"
  "errors"
  "log/slog"
  "net/url"
  "strconv"
  "net/http/cookiejar"
  "aspace_publisher/utils"
  "aspace_publisher/connect"
  "github.com/labstack/echo/v4"
)

func GetSession(c echo.Context) (string, error){
  session_id, err := utils.FetchCookieVal(c, "aw_session")
  if session_id == "" || err != nil {
    session_id, err = authenticate()
    if err != nil { return "", err }
    utils.WriteCookie(c, 60, "aw_session", session_id)
  }
  return session_id, nil
}

func authenticate()(string, error){
  authurl, _ := url.Parse(os.Getenv("AWEST_URL") + "login.php")
  data := url.Values{}
  data.Set("username", os.Getenv("AWEST_NAME"))
  data.Set("password", os.Getenv("AWEST_PASS"))
  jar, err := cookiejar.New(nil)
  if err != nil { slog.Error(err.Error()); return "", errors.New("could not create cookie jar") }
  client := &http.Client{
    Jar: jar,
  }

  request, err := http.NewRequest("POST", authurl.String(), strings.NewReader(data.Encode()))
  request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

  connect.RequestDump(request, "DEBUG")
  if err != nil { slog.Error(err.Error()); return "", errors.New("Unable to create login request") }
  response, err := client.Do(request)
  connect.ResponseDump(response, "DEBUG")
  if err != nil { slog.Error(err.Error()); return "", errors.New("Unable to complete login to awest") }
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
  slog.Error("awest_auth error: could not find session")
  return "", errors.New("could not find session in cookies")
}
