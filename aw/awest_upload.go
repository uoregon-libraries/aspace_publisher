package aw

import(
  "bytes"
  "io"
  "fmt"
  "errors"
  "net/http"
  "net/http/httputil"
  "log"
  "os"
  "strings"
  "time"
  "github.com/PuerkitoBio/goquery"
)

// an ark_url looks like https://archiveswest.orbiscascade.org/ark:80444/xv205342
func MakeUploadMap(ark, filekey, filepath string)(map[string]string, error){
  vals := make(map[string]string)
  vals["filekey"] = filekey
  vals["filepath"] = filepath
  if ark != "" {
    vals["ark"] = ark
    exists, err := TestArk(ark)
    if err != nil { return nil, err }
    if exists == true {
      vals["replace"] = "1"
    }
  }
  return vals, nil
}

// tries to upload an ead in two steps
// first post the upload
// second get confirmation upload was successful
// return confirmation page
func Upload(sessionid string, boundary string, verbose string, form *bytes.Buffer)(string, error){
  url := os.Getenv("AWEST_URL") + "upload-process.php"
  req, err := http.NewRequest("POST", url, form)
  if err != nil { return "", errors.New("unable to create http request") }

  req.Header.Set("cookie", "PHPSESSID=" + sessionid)
  req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", boundary))

  client := &http.Client{
    Timeout: time.Second * 30,
  }
  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { log.Println(err) } else { log.Printf("REQUEST:\n%s", string(reqdump)) }
  }
  response, err := client.Do(req); if err != nil { return "", err }
  defer response.Body.Close()
  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else {  log.Printf("RESPONSE:\n%s", string(respdump)) }
  }
  body, err := io.ReadAll(response.Body); if err != nil { return "", err }

  return string(body), nil
}

// returns true if there a record published at the ark url
func TestArk(ark string)(bool, error){
  ark_url := os.Getenv("ARK_URL_BASE") + ark
  req, err := http.NewRequest("GET", ark_url, nil)
  if err != nil { return false, errors.New("unable to create http request") }
  client := &http.Client{
    Timeout: time.Second * 30,
  }
  response, err := client.Do(req); if err != nil { return false, err }
  body, err := io.ReadAll(response.Body); if err != nil { return false, err }
  defer response.Body.Close()
  doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
  if err != nil { return false, errors.New("unable to read response") }
  tag := doc.Find("#toc")
  if len(tag.Nodes) == 0 {
    return false, nil
  }
  return true, nil
}

func Validate(sessionid string, boundary string, verbose string, form *bytes.Buffer)(string, error){
  url := os.Getenv("AWEST_URL") + "validation-process.php"
  req, err := http.NewRequest("POST", url, form)
  if err != nil { return "", errors.New("unable to create http request") }

  req.Header.Set("cookie", "PHPSESSID=" + sessionid)
  req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", boundary))
  client := &http.Client{
    Timeout: time.Second * 30,
  }
  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { log.Println(err) } else { log.Printf("REQUEST:\n%s", string(reqdump)) }
  }
  response, err := client.Do(req); if err != nil { return "", err }
  defer response.Body.Close()
  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else { log.Printf("RESPONSE:\n%s", string(respdump)) }
  }
  body, err := io.ReadAll(response.Body); if err != nil { return "", err }
  return string(body), nil
}
