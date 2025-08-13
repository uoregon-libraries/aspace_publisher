package as

import (
  "github.com/tidwall/sjson"
  "github.com/tidwall/gjson"
  "errors"
  "log"
  "os"
  "strings"
  "net/http"
  "net/http/httputil"
  "fmt"
  "time"
  "io"
)

func UpdateUserDefined1(record []byte, oclc string)([]byte, error){
  modified, err := sjson.SetBytes(record, "user_defined.string_1", oclc)
  if err != nil { log.Println(err); return nil, err }
  return modified, nil
}

// refactor
func UpdateUserDefined2(record []byte, mms_id string)([]byte, error){
  modified, err := sjson.SetBytes(record, "user_defined.string_2", mms_id)
  if err != nil { log.Println(err); return nil, err }
  return modified, nil
}

func UpdateResource(sessionid string, repo_id string, resource_id string, json_record string )(string, error){
  verbose := os.Getenv("VERBOSE")
  base_url := os.Getenv("ASPACE_URL")
  url := base_url + fmt.Sprintf("repositories/%s/resources/%s", repo_id, resource_id)
  data := strings.NewReader(json_record)
  req, err := http.NewRequest("POST", url, data)
  if err != nil { return "", errors.New("unable to create http request") }

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
  response, err := client.Do(req)
  if err != nil { return "", err }
  defer response.Body.Close()

  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else {
      log.Printf("RESPONSE:\n%s", string(respdump)) }
  }

  body, err := io.ReadAll(response.Body)
  if err != nil { log.Println(err); return "", errors.New("unable to read response") }

  if response.StatusCode != 200 {
    return "", errors.New(fmt.Sprintf("Unable to update aspace resource: %s", string(body)))
  }

  return string(body), nil
}

// takes AO json and inserts instance
func UpdateWithInstance(record []byte, instance string)([]byte, error){
  instance_json := gjson.Parse(instance)
  modified, err := sjson.SetBytes(record, "instances.-1", instance_json.Value())
  if err != nil { log.Println(err); return nil, err }
  return modified, nil
}

func Instance(path string) string {
  return fmt.Sprintf(`{"instance_type": "digital_object", "jsonmodel_type": "instance", "is_representative": false, "digital_object": { "ref": "%s"}`, path)
}
