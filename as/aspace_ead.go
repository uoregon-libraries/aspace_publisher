package as

import(
  "fmt"
  "time"
  "net/http"
  "errors"
  "io"
  "os"
  "log/slog"
  "aspace_publisher/connect"
)

func AcquireEad(sessionid string, repo_id string, resource_id string) (string, error){
  base_url := os.Getenv( "ASPACE_URL")
  url := base_url + fmt.Sprintf("repositories/%s/resource_descriptions/%s.xml?include_unpublished=%s&include_daos=%s&numbered_cs=%s&ead3=%s", repo_id, resource_id, "False", "True", "True", "False")
  req, err := http.NewRequest("GET", url, nil)
  if err != nil { slog.Error(err.Error()); return "", errors.New("unable to create http request") }

  req.Header.Set("X-ArchivesSpace-Session", sessionid)
  req.Header.Set("Accept", "*/*")
  req.Header.Set("User-Agent", "curl/7.61.1")

  connect.RequestDump(req)

  client := &http.Client{
    Timeout: time.Second * 90,
  }

  response, err := client.Do(req); if err != nil { slog.Error(err.Error()); return "", err }
  defer response.Body.Close()
  connect.ResponseDump(response) //check response only after err check

  body, err := io.ReadAll(response.Body); if err != nil { slog.Error(err.Error()); return "", err }
  if response.StatusCode != 200 {
    return string(body), errors.New("problem retrieving ead")
  } else { return string(body), nil }
}
