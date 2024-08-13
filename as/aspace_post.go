package as

import (
  "log"
  "os"
  "strings"
  "net/http"
  "net/http/httputil"
  "github.com/tidwall/sjson"
  "fmt"
  "time"
  "io"
)
type Responses struct {
  responses []Response
}

type Response struct {
  id string
  response string
}

func (r Response) ResponseToString() string {
  return fmt.Sprintf(`{"id":"%s", "response":"%s"}`, r.id, r.response)
}

func (r Responses) ResponsesToString() string {
  all_resp := `{"responses":[]}`
  for _, elt := range r.responses {
    temp, _ := sjson.Set(`{"id":"", "response":""}`, "id", elt.id)
    temp2, _ := sjson.Set(temp, "response", elt.response)
    all_resp, _ = sjson.Set(all_resp, "responses.-1", temp2)
  }
  return all_resp
}

func Post(sessionid string, identifier string, repo_id string, record_id string, json_record string ) Response {
  verbose := os.Getenv("VERBOSE")
  test := os.Getenv("TEST")
  base_url := os.Getenv("ASPACE_URL")
  url := base_url + fmt.Sprintf("repositories/%s/%s", repo_id, record_id)
  data := strings.NewReader(json_record)
  req, err := http.NewRequest("POST", url, data)
if err != nil { log.Println(err); return Response{identifier, "unable to create http request"} }

  req.Header.Set("X-ArchivesSpace-Session", sessionid)
  req.Header.Set("Accept", "*/*")
  req.Header.Set("User-Agent", "curl/7.61.1")

  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { log.Println(err) } else {
      log.Printf("REQUEST:\n%s", string(reqdump)) }
  }
  if test == "true" { return Response { identifier, "test mode" } }

  client := &http.Client{
    Timeout: time.Second * 60,
  }
  response, err := client.Do(req)
  if err != nil { log.Println(err); return Response{ identifier, "unable to make request to aspace" } }
  defer response.Body.Close()

  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else {
      log.Printf("RESPONSE:\n%s", string(respdump)) }
  }

  body, err := io.ReadAll(response.Body)
  if err != nil { log.Println(err); return Response{ identifier, "unable to read response" } }

  return Response{ identifier, string(body) }
}
