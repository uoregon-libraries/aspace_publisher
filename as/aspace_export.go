package as

import(
  "fmt"
  "time"
  "net/http"
  "errors"
  "io"
  "os"
)

func AcquireMarc(sessionid string, repo_id string, resource_id string) (string, error){
  base_url := os.Getenv( "ASPACE_URL")
  url := base_url + fmt.Sprintf("repositories/%s/resources/marc21/%s.xml?include_unpublished_marc=%s", repo_id, resource_id, "False")
  req, err := http.NewRequest("GET", url, nil)
  if err != nil { return "", errors.New("unable to create http request") }

  req.Header.Set("X-ArchivesSpace-Session", sessionid)
  client := &http.Client{
    Timeout: time.Second * 30,
  }
  response, err := client.Do(req); if err != nil { return "", err }
  body, err := io.ReadAll(response.Body); if err != nil { return "", err }
  response.Body.Close()
  return string(body), nil
}

func AcquireJson(sessionid string, repo_id string, resource_id string) (string, error){
  base_url := os.Getenv( "ASPACE_URL")
  url := base_url + fmt.Sprintf("/repositories/%s/resources/%s", repo_id, resource_id)
  req, err := http.NewRequest("GET", url, nil)
  if err != nil { return "", errors.New("unable to create http request") }

  req.Header.Set("X-ArchivesSpace-Session", sessionid)
  client := &http.Client{
    Timeout: time.Second * 30,
  }
  response, err := client.Do(req); if err != nil { return "", err }
  body, err := io.ReadAll(response.Body); if err != nil { return "", err }
  response.Body.Close()
  return string(body), nil
}

func SetUserDefined(field string, value string, json ){


}
