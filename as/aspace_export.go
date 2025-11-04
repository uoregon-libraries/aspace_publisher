package as

import(
  "fmt"
  "time"
  "net/http"
  "net/http/httputil"
  "log"
  "errors"
  "io"
  "os"
)

func AcquireMarc(sessionid string, repo_id string, resource_id string, published string) (string, error){
  include := "false"
  if published == "false" { include = "true" }
  verbose := os.Getenv("VERBOSE")
  base_url := os.Getenv("ASPACE_URL")
  url := base_url + fmt.Sprintf("repositories/%s/resources/marc21/%s.xml?include_unpublished_marc=%s", repo_id, resource_id, include)
  req, err := http.NewRequest("GET", url, nil)
  if err != nil { log.Println(err); return "", errors.New("unable to create http request") }

  req.Header.Set("X-ArchivesSpace-Session", sessionid)
  req.Header.Set("Accept", "*/*")
  req.Header.Set("User-Agent", "curl/7.61.1")

  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { log.Println(err) } else { log.Printf("REQUEST:\n%s", string(reqdump)) }
  }

  client := &http.Client{
    Timeout: time.Second * 60,
  }
  response, err := client.Do(req)
  if err != nil { log.Println(err); return "", errors.New("unable to complete request to archivesspace") }
  defer response.Body.Close()
  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else { log.Printf("RESPONSE:\n%s", string(respdump)) }
  }

  body, err := io.ReadAll(response.Body)
  if err != nil { log.Println(err); return "", errors.New("unable to read response from archivesspace") }
  if response.StatusCode != 200 { return body, errors.New("aspace error exporting MARC") }
  return string(body), nil
}

func AcquireJson(sessionid string, repo_id string, record_id string) ([]byte, error){
  base_url := os.Getenv( "ASPACE_URL")
  url := base_url + fmt.Sprintf("/repositories/%s/%s", repo_id, record_id)
  req, err := http.NewRequest("GET", url, nil)
  if err != nil { log.Println(err); return nil, errors.New("unable to create http request") }

  req.Header.Set("X-ArchivesSpace-Session", sessionid)
  req.Header.Set("Accept", "*/*")
  req.Header.Set("User-Agent", "curl/7.61.1")

  client := &http.Client{
    Timeout: time.Second * 60,
  }
  response, err := client.Do(req)
  if err != nil { log.Println(err); return nil, errors.New("unable to complete request to archivesspace.") }
  defer response.Body.Close()
  body, err := io.ReadAll(response.Body)
  if err != nil { log.Println(err); return nil, errors.New("unable to read response from archivesspace") }
  if response.StatusCode != 200 { return body, errors.New("aspace error exporting record") }
  return body, nil
}


