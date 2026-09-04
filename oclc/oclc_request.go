package oclc

import(
  "fmt"
  "net/http"
  "io"
  "log/slog"
  "os"
  "errors"
  "time"
  "strings"
  "slices"
  "regexp"
  "aspace_publisher/connect"
)

type RequestFun func(string, string, string, string, string, string)(string, error)
func Request(token string, method string, marc string, path string, id string, accept string) (string, error){
  base_url := os.Getenv("OCLC_URL")
  test := os.Getenv("TEST")
  url := assembleUrl([]string{base_url,path,id})
  data := strings.NewReader(marc)
  req, err := http.NewRequest(method, url, data)
  if err != nil { slog.Error(err.Error()); return "", errors.New("unable to create http request") }
  req.Header.Set("accept", "application/" + accept)
  if marc != "" {
    req.Header.Set("Content-Type", "application/marcxml+xml")
  }
  req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

  connect.RequestDump(req)

  client := &http.Client{
    Timeout: time.Second * 60,
  }

  if test == "true" { return `<record></record>`, nil }

  response, err := client.Do(req)
  if err != nil { slog.Error(err.Error()); return "", errors.New("unable to complete http request") }
  defer response.Body.Close()

  connect.ResponseDump(response) //check response only after err check

  body, err := io.ReadAll(response.Body)
  if err != nil { slog.Error(err.Error()); return "", errors.New("unable to read response from oclc") }
  if response.StatusCode != 200 { return string(body), errors.New("oclc errors") }

  return string(body), nil

}

func assembleUrl(parts []string) string{
  parts = slices.DeleteFunc(parts, func(str string) bool{
    return str == "" } )
  return strings.Join(parts, "/")
}

type RecordFun func(string, string)(string, error)
func Record(token string, id string)(string, error){
  base_url := os.Getenv("OCLC_URL")
  url := assembleUrl([]string{base_url,"manage/bibs", id})
  req, err := http.NewRequest("GET", url, nil)
  if err != nil { slog.Error(err.Error()); return "", errors.New("unable to create http request") }
  req.Header.Set("accept", "application/marcxml+xml")
  req.Header.Set("Content-Type", "application/marcxml+xml")
  req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

  connect.RequestDump(req)

  client := &http.Client{
    Timeout: time.Second * 60,
  }
  response, err := client.Do(req)
  if err != nil { slog.Error(err.Error()); return "", errors.New("unable to complete http request") }
  defer response.Body.Close()

  connect.ResponseDump(response) //check response only after err check

  body, err := io.ReadAll(response.Body)
  if err != nil { slog.Error(err.Error()); return "", errors.New("unable to read response from oclc") }
  if response.StatusCode != 200 { return string(body), errors.New("oclc errors") }

  return string(body), nil
}
func UnformatXML(xmlString string) string {
  var unformatXMLRegEx = regexp.MustCompile(`>\s+<`)
  unformatBetweenTags := unformatXMLRegEx.ReplaceAllString(xmlString, "><")
  return strings.TrimSpace(unformatBetweenTags)
}
