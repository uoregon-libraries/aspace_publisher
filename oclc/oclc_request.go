package oclc

import(
  "fmt"
  "net/http"
  "io"
  "log"
  "os"
  "errors"
  "time"
  "strings"
  "net/http/httputil"
)


func Create(token string, marc string) (string, error){
  verbose := os.Getenv("VERBOSE")
  base_url := "https://metadata.api.oclc.org/worldcat"
  url := base_url + "/manage/bibs"
  data := strings.NewReader(marc)
  req, err := http.NewRequest("POST", url, data)
  if err != nil { log.Println(err); return "", errors.New("unable to create http request") }
  req.Header.Set("accept", "application/marcxml+xml")
  req.Header.Set("Content-Type", "application/marcxml+xml")
  req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { log.Println(err) } else {
      log.Printf("REQUEST:\n%s", string(reqdump)) }
  }

  client := &http.Client{
    Timeout: time.Second * 60,
  }
  response, err := client.Do(req)
  defer response.Body.Close()
  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else {
      log.Printf("RESPONSE:\n%s", string(respdump)) }
  }
  if err != nil { log.Println(err); return "", errors.New("unable to complete http request") }
  body, err := io.ReadAll(response.Body)
  if err != nil { log.Println(err); return "", errors.New("unable to read response from oclc") }
  if response.StatusCode != 200 {
    return "", errors.New(fmt.Sprintf("Unable to create worldcat record: %s", string(body)))
  }

  return string(body), nil

}

func Update(token string, marc string) {

}

func Validate(token string, marc string)(string, error){
  verbose := os.Getenv("VERBOSE")
  base_url := "https://metadata.api.oclc.org/worldcat"
  url := base_url + "/manage/bibs/validate/validateFull"
  data := strings.NewReader(marc)
  req, err := http.NewRequest("POST", url, data)
  if err != nil { log.Println(err); return "", errors.New("unable to create http request") }
  req.Header.Set("accept", "application/json")
  req.Header.Set("Content-Type", "application/marcxml+xml")
  req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { log.Println(err) } else {
      log.Printf("REQUEST:\n%s", string(reqdump)) }
  }
  client := &http.Client{
    Timeout: time.Second * 60,
  }
  response, err := client.Do(req)
  defer response.Body.Close()
  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else {
      log.Printf("RESPONSE:\n%s", string(respdump)) }
  }

  if err != nil { log.Println(err); return "", errors.New("unable to complete http request to oclc.") }
  body, err := io.ReadAll(response.Body)
  if err != nil { log.Println(err); return "", errors.New("unable to read response from oclc server.") }

  return string(body), nil
}

