package aw

import(
  "bytes"
  "io"
  "fmt"
  "errors"
  "net/http"
  "net/http/httputil"
  "golang.org/x/net/html"
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

func Upload(sessionid string, boundary string, verbose string, form *bytes.Buffer)(io.Reader, error){
  url := os.Getenv("AWEST_URL") + "upload-process.php"
  req, err := http.NewRequest("POST", url, form)
  if err != nil { return nil, errors.New("unable to create http request") }

  req.Header.Set("cookie", "PHPSESSID=" + sessionid)
  req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", boundary))

  client := &http.Client{
    Timeout: time.Second * 30,
  }
  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { log.Println(err) } else { log.Printf("REQUEST:\n%s", string(reqdump)) }
  }
  response, err := client.Do(req); if err != nil { return nil, err }
  defer response.Body.Close()
  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else {  log.Printf("RESPONSE:\n%s", string(respdump)) }
  }
  return response.Body, nil
}

func ark_url(ark string)(string){
  url := fmt.Sprintf("https://archiveswest.orbiscascade.org/ark:%s", ark)
  return url
}

// returns true if there a record published at the ark url
func TestArk(ark string)(bool, error){
  req, err := http.NewRequest("GET", ark_url(ark), nil)
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

func Validate(sessionid string, boundary string, verbose string, form *bytes.Buffer)(io.Reader, error){
  url := os.Getenv("AWEST_URL") + "validation-process.php"
  req, err := http.NewRequest("POST", url, form)
  if err != nil { return nil, errors.New("unable to create http request") }

  req.Header.Set("cookie", "PHPSESSID=" + sessionid)
  req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", boundary))
  client := &http.Client{
    Timeout: time.Second * 30,
  }
  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { log.Println(err) } else { log.Printf("REQUEST:\n%s", string(reqdump)) }
  }
  response, err := client.Do(req); if err != nil { return nil, err }
  defer response.Body.Close()
  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else { log.Printf("RESPONSE:\n%s", string(respdump)) }
  }
  return response.Body, nil
}

func ParseResult(r io.Reader)(string, error){
  var b bytes.Buffer
  doc, err := goquery.NewDocumentFromReader(r)
  if err != nil { log.Println(err); return "", err }
  success := doc.Find(".success")
  if success.Length() > 0 {
    err := html.Render(&b, success.Nodes[0])
    if err != nil { log.Println(err); return "", err }
    return b.String(), nil
  } else {
    err_nodes := doc.Find(".errors")
    if err_nodes.Length() > 0 {
      err := html.Render(&b, err_nodes.Nodes[0])
      if err != nil { log.Println(err); return "", err }
      return b.String(), nil
    }
  }
  return "", nil
}
