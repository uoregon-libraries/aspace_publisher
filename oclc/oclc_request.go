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
  "slices"
  "net/http/httputil"
)

func Request(token string, method string, marc string, path string, id string, accept string) (string, error){
  verbose := os.Getenv("VERBOSE")
  base_url := os.Getenv("OCLC_URL")
  test := os.Getenv("TEST")
  url := assembleUrl([]string{base_url,path,id})
  data := strings.NewReader(marc)
  req, err := http.NewRequest(method, url, data)
  if err != nil { log.Println(err); return "", errors.New("unable to create http request") }
  req.Header.Set("accept", "application/" + accept)
  if marc != "" {
    req.Header.Set("Content-Type", "application/marcxml+xml")
  }
  req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { log.Println(err) } else {
      log.Printf("REQUEST:\n%s", string(reqdump)) }
  }

  client := &http.Client{
    Timeout: time.Second * 60,
  }

  if test == "true" { return `<record></record>`, nil }

  response, err := client.Do(req)
  if err != nil { log.Println(err); return "", errors.New("unable to complete http request") }
  defer response.Body.Close()
  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else {
      log.Printf("RESPONSE:\n%s", string(respdump)) }
  }
  body, err := io.ReadAll(response.Body)
  if err != nil { log.Println(err); return "", errors.New("unable to read response from oclc") }
  if response.StatusCode != 200 { return string(body), errors.New("oclc errors") }

  return string(body), nil

}

func assembleUrl(parts []string) string{
  parts = slices.DeleteFunc(parts, func(str string) bool{
    return str == "" } )
  return strings.Join(parts, "/")
}

func Record(token string, id string)(string, error){
  verbose := os.Getenv("VERBOSE")
  base_url := os.Getenv("OCLC_URL")
  url := assembleUrl([]string{base_url,"manage/bibs", id})
  req, err := http.NewRequest("GET", url, nil)
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
  if err != nil { log.Println(err); return "", errors.New("unable to complete http request") }
  defer response.Body.Close()
  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else {
      log.Printf("RESPONSE:\n%s", string(respdump)) }
  }
  body, err := io.ReadAll(response.Body)
  if err != nil { log.Println(err); return "", errors.New("unable to read response from oclc") }
  if response.StatusCode != 200 { return string(body), errors.New("oclc errors") }

  return string(body), nil
}
