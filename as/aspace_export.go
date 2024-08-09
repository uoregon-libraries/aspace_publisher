package as

import(
  "fmt"
  "time"
  "net/http"
  "log"
  "errors"
  "io"
  "os"
)

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

  return body, nil
}

