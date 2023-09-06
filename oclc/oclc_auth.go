package oclc

import(
  "net/http"
  "encoding/json"
  "os"
  "fmt"
  "time"
  "io/ioutil"
  "github.com/labstack/echo/v4"
  "aspace_publisher/echosession"
)

type OclcToken struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
}

func oclcAuth() string, error {

  base := os.Getenv("OCLC_TOKEN_URL")
  key := os.Getenv("OCLC_KEY")
  secret := os.Getenv("OCLC_SECRET")

  url := base + "?grant_type=client_credentials&scope=WorldCatMetadataAPI"
  req, err := http.NewRequest("POST", url, nil); if err != nil { return "", err }
  req.SetBasicAuth(key, secret)
  client := &http.Client{
    Timeout: time.Second * 10,
  }
  var ot OclcToken
  response, err := client.Do(req); if err != nil { return "", err }
  byteVal, _ := ioutil.ReadAll(response.Body)
  err = json.Unmarshal(byteVal, &ot); if err != nil { return "", err }
  return ot.AccessToken, nil
}


func GetToken(c echo.Context) string, error{
  store := echosession.FromContext(c)
  t, err := store.Get("oclc_token")
  if t == "" || err != nil {
    t, err = oclcAuth()
    if err != nil { return "", err }
    store.Set("oclc_token", t)
  }
  return t, nil
}

