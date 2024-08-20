package as

import (
  "log"
  "os"
  "strings"
  "net/http"
  "net/http/httputil"
  "fmt"
  "time"
  "io"
  "encoding/json"
)
type Responses struct {
  responses []Response
}

type Response struct {
  id string
  message Message
}

type Message map[string]any

func (r Response) ResponseToString() string{
  var output []byte
  var err error
  output, err = json.Marshal(r.message)
  if err != nil { log.Println(err); return `{"id":` + r.id + `", "error": "unable to marshal message" }` }
  return `{"id":"` + r.id + `", "message":` + string(output) + "}"
}

func (r Responses) ResponsesToString() string {
  all_resp := ""
  for _, elt := range r.responses {
    all_resp += elt.ResponseToString()
  }
  return all_resp
}

func BuildMessage(message string) Message{
  var m Message
  err := json.Unmarshal([]byte(message), &m)
  if err != nil { log.Println(err); return BuildErrorMessage("unable to unmarshal message") }
  return m
}

func BuildErrorMessage(message string) Message{
  var m Message
  e := `{"error":"` + message + `"}`
  _ = json.Unmarshal([]byte(e), &m)
  return m
}

func Post(sessionid string, identifier string, repo_id string, record_id string, json_record string ) Response {
  verbose := os.Getenv("VERBOSE")
  test := os.Getenv("TEST")
  base_url := os.Getenv("ASPACE_URL")
  url := base_url + fmt.Sprintf("repositories/%s/%s", repo_id, record_id)
  data := strings.NewReader(json_record)
  req, err := http.NewRequest("POST", url, data)
if err != nil { log.Println(err); return Response{identifier, BuildErrorMessage("unable to create http request")} }

  req.Header.Set("X-ArchivesSpace-Session", sessionid)
  req.Header.Set("Accept", "*/*")
  req.Header.Set("User-Agent", "curl/7.61.1")

  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { log.Println(err) } else {
      log.Printf("REQUEST:\n%s", string(reqdump)) }
  }
  if test == "true" { return Response { identifier, BuildErrorMessage("test mode") } }

  client := &http.Client{
    Timeout: time.Second * 60,
  }
  response, err := client.Do(req)
  if err != nil { log.Println(err); return Response{ identifier, BuildErrorMessage("unable to make request to aspace") } }
  defer response.Body.Close()

  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else {
      log.Printf("RESPONSE:\n%s", string(respdump)) }
  }

  body, err := io.ReadAll(response.Body)
  if err != nil { log.Println(err); return Response{ identifier, BuildErrorMessage("unable to read response") } }

  var r Response
  err = json.Unmarshal(body, &r)
  if err != nil { log.Println(err); return Response{ identifier, BuildErrorMessage("unable to unmarshal response") } }

  return Response{ identifier, BuildMessage(string(body)) }
}
