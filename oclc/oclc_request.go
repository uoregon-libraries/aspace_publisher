package oclc

import(
  "fmt"
  "net/http"
  "io"
//  "os"

)


func Create(token string, marc string) string, error{
  //base_url := os.Getenv("OCLC_URL")
  base_url := "https://metadata.api.oclc.org/worldcat"
  url := base_url + "/manage/bibs"
  req, err := http.NewRequest("POST", url, nil)
  if err != nil { return "", errors.New("unable to create http request") }
  req.Header.Set("accept", "application/marcxml+xml")
  req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
  client := &http.Client{
    Timeout: time.Second * 30,
  }
  response, err := client.Do(req); if err != nil { return "", err }
  body, err := io.ReadAll(response.Body); if err != nil { return "", err }
  response.Body.Close()
  return string(body), nil

}

func Update(token string, marc string) {

}


