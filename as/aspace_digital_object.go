package as

import(
  "net/http"
  "log"
  "os"
  "fmt"
  "io"
  "time"
  "strings"
  "net/http/httputil"
  "github.com/tidwall/gjson"
  "github.com/tidwall/sjson"
)

type Responses struct {
  responses []Response
}

type Response struct {
  id string
  response string
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

func CreateDigitalObjects(digital_obj_list string, sessionid string) (Responses){
  var r Responses
  items := gjson.Get(digital_obj_list, "digital_objects")
  items.ForEach(func(key, value gjson.Result) bool {
    aoid := gjson.Get(value.String(), "digital_object_id")
    result := CreateDigitalObject(aoid.String(), value.String(), sessionid)
    r.responses = append(r.responses, result)
    return true
  })
  return r
}

func CreateDigitalObject(identifier string, digital_obj string, sessionid string) (Response){
  repo_id := "2"
  verbose := os.Getenv("VERBOSE")
  base_url := os.Getenv("ASPACE_URL")
  url := base_url + fmt.Sprintf("repositories/%s/digital_objects", repo_id)
  data := strings.NewReader(digital_obj)
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

  client := &http.Client{
    Timeout: time.Second * 60,
  }
  resp, err := client.Do(req)
  if err != nil { log.Println(err); return Response{ identifier, "unable to make request to aspace" } }
  defer resp.Body.Close()

  if verbose == "true" {
    respdump, err := httputil.DumpResponse(resp, true)
    if err != nil { log.Println(err) } else {
      log.Printf("RESPONSE:\n%s", string(respdump)) }
  }

  body, err := io.ReadAll(resp.Body)
  if err != nil { log.Println(err); return Response{ identifier, "unable to read response" } }

  return Response{identifier, string(body)}
}
