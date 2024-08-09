package as

import (
  "log"
  "os"
  "strings"
  "net/http"
  "net/http/httputil"
<<<<<<< HEAD
  "fmt"
  "time"
  "io"
  "encoding/json"
=======
  "github.com/tidwall/sjson"
  "fmt"
  "time"
  "io"
>>>>>>> 2 step process
)
type Responses struct {
  responses []Response
}

type Response struct {
  id string
<<<<<<< HEAD
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
=======
  response string
}

func (r Responses) ResponsesToString() string {
  all_resp := `{"responses":[]}`
  for _, elt := range r.responses {
    temp, _ := sjson.Set(`{"id":"", "response":""}`, "id", elt.id)
    temp2, _ := sjson.Set(temp, "response", elt.response)
    all_resp, _ = sjson.Set(all_resp, "responses.-1", temp2)
>>>>>>> 2 step process
  }
  return all_resp
}

<<<<<<< HEAD
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

=======
>>>>>>> 2 step process
func Post(sessionid string, identifier string, repo_id string, record_id string, json_record string ) Response {
  verbose := os.Getenv("VERBOSE")
  test := os.Getenv("TEST")
  base_url := os.Getenv("ASPACE_URL")
  url := base_url + fmt.Sprintf("repositories/%s/%s", repo_id, record_id)
  data := strings.NewReader(json_record)
  req, err := http.NewRequest("POST", url, data)
<<<<<<< HEAD
if err != nil { log.Println(err); return Response{identifier, BuildErrorMessage("unable to create http request")} }
=======
if err != nil { log.Println(err); return Response{identifier, "unable to create http request"} }
>>>>>>> 2 step process

  req.Header.Set("X-ArchivesSpace-Session", sessionid)
  req.Header.Set("Accept", "*/*")
  req.Header.Set("User-Agent", "curl/7.61.1")

  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { log.Println(err) } else {
      log.Printf("REQUEST:\n%s", string(reqdump)) }
  }
<<<<<<< HEAD
  if test == "true" { return Response { identifier, BuildErrorMessage("test mode") } }
=======
  if test == "true" { return Response { identifier, "test mode" } }
>>>>>>> 2 step process

  client := &http.Client{
    Timeout: time.Second * 60,
  }
  response, err := client.Do(req)
<<<<<<< HEAD
  if err != nil { log.Println(err); return Response{ identifier, BuildErrorMessage("unable to make request to aspace") } }
=======
  if err != nil { log.Println(err); return Response{ identifier, "unable to make request to aspace" } }
>>>>>>> 2 step process
  defer response.Body.Close()

  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else {
      log.Printf("RESPONSE:\n%s", string(respdump)) }
  }

  body, err := io.ReadAll(response.Body)
<<<<<<< HEAD
  if err != nil { log.Println(err); return Response{ identifier, BuildErrorMessage("unable to read response") } }

  var r Response
  err = json.Unmarshal(body, &r)
  if err != nil { log.Println(err); return Response{ identifier, BuildErrorMessage("unable to unmarshal response") } }

  return Response{ identifier, BuildMessage(string(body)) }
=======
  if err != nil { log.Println(err); return Response{ identifier, "unable to read response" } }

  return Response{ identifier, string(body) }
>>>>>>> 2 step process
}
