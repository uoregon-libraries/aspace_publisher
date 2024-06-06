package as

import(
  "fmt"
  "time"
  "net/http"
  "net/http/httputil"
  "errors"
  "io"
  "os"
  "log"
)

func AcquireEad(sessionid string, repo_id string, resource_id string, verbose string) (string, error){
  base_url := os.Getenv( "ASPACE_URL")
  url := base_url + fmt.Sprintf("repositories/%s/resource_descriptions/%s.xml?include_unpublished=%s&include_daos=%s&numbered_cs=%s&ead3=%s", repo_id, resource_id, "False", "True", "True", "False")
  req, err := http.NewRequest("GET", url, nil)
  if err != nil { log.Println(err); return "", errors.New("unable to create http request") }

  req.Header.Set("X-ArchivesSpace-Session", sessionid)
  req.Header.Set("Accept", "*/*")
  req.Header.Set("User-Agent", "curl/7.61.1")

  client := &http.Client{
    Timeout: time.Second * 90,
  }
  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { log.Println(err) } else { log.Printf("REQUEST:\n%s", string(reqdump)) }
  }
  response, err := client.Do(req); if err != nil { log.Println(err); return "", err }
  defer response.Body.Close()
  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else { log.Printf("RESPONSE:\n%s", string(respdump)) }
  }

  body, err := io.ReadAll(response.Body); if err != nil { log.Println(err); return "", err }
  if response.StatusCode != 200 {
    return string(body), errors.New("unable to retrieve ead")
  } else { return string(body), nil }
}
