package oclc

import(
  "net/http"
  "encoding/json"
  "os"
  "log/slog"
  "time"
  "io/ioutil"
  "github.com/labstack/echo/v4"
  "aspace_publisher/utils"
  "aspace_publisher/connect"
  "errors"
)

type OclcToken struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
}

func OclcAuth() (string, error) {
  url := os.Getenv("OCLC_AUTH_URL")
  name := os.Getenv("OCLC_NAME")
  pass := os.Getenv("OCLC_PASS")

  req, err := http.NewRequest("POST", url, nil)
  if err != nil { slog.Error(err.Error()); return "", errors.New("unable to create request") }
  req.SetBasicAuth(name, pass)

  connect.RequestDump(req, "DEBUG")

  client := &http.Client{
    Timeout: time.Second * 10,
  }
  response, err := client.Do(req)
  if err != nil { slog.Error(err.Error()); return "", errors.New("unable to complete http request") }
  defer response.Body.Close()

  connect.ResponseDump(response, "DEBUG") //check response only after err check

  byteVal, _ := ioutil.ReadAll(response.Body)
  var ot OclcToken
  err = json.Unmarshal(byteVal, &ot)
  if err != nil { slog.Error(err.Error()); return "", errors.New("unable to extract token") }
  return ot.AccessToken, nil
}


func GetToken(c echo.Context) (string, error){
  token, err := utils.FetchCookieVal(c, "oclc_token")
  if token == "" || err != nil {
    token, err = OclcAuth()
    if err != nil { slog.Error(err.Error()); return "", err }
    utils.WriteCookie(c, 20, "oclc_token", token)
  }
  return token, nil
}

