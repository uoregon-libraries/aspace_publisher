package oclc

import(
  "net/http"
  "encoding/json"
  "os"
  "fmt"
  "time"
  "io/ioutil"
)

type OclcToken struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
}

func (ot *OclcToken) GetToken() error{
  base := os.Getenv("OCLC_TOKEN_URL")
  key := os.Getenv("OCLC_KEY")
  secret := os.Getenv("OCLC_SECRET")

  url := base + "?grant_type=client_credentials&scope=WorldCatMetadataAPI"
  req, err := http.NewRequest("POST", url, nil); if err != nil { return err }
  req.SetBasicAuth(key, secret)
  client := &http.Client{
    Timeout: time.Second * 10,
  }
  response, err := client.Do(req); if err != nil { return err }
  byteVal, _ := ioutil.ReadAll(response.Body)
  err = json.Unmarshal(byteVal, &ot); if err != nil { return err }
  return nil
}


