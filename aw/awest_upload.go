package aw

import(
  "bytes"
  "io"
  "fmt"
  "errors"
  "net/http"
  "golang.org/x/net/html"
  "log/slog"
  "os"
  "strings"
  "time"
  "regexp"
  "github.com/PuerkitoBio/goquery"
  "github.com/beevik/etree"
  "aspace_publisher/connect"
)

// an ark_url looks like https://archiveswest.orbiscascade.org/ark:80444/xv205342
func MakeUploadMap(ark string, exists bool, filepath string)(map[string]string){
  vals := make(map[string]string)
  vals["filekey"] = "ead"
  vals["filepath"] = filepath
  vals["ark"] = ark
  if exists { vals["replace"] = "1" }
  return vals
}

func Request(sessionid string, boundary string, form *bytes.Buffer, operation string)(io.Reader, error){
  url := ""
  if operation == "upload" {
    url = os.Getenv("AWEST_URL") + "upload-process.php"
  } else {
    url = os.Getenv("AWEST_URL") + "validation-process.php"
  }
  req, err := http.NewRequest("POST", url, form)
  if err != nil { slog.Error(err.Error()); return nil, errors.New("unable to create http request") }

  req.Header.Set("cookie", "PHPSESSID=" + sessionid)
  req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", boundary))

  connect.RequestDump(req)

  client := &http.Client{
    Timeout: time.Second * 30,
  }
  response, err := client.Do(req); if err != nil { slog.Error(err.Error()); return nil, err }
  defer response.Body.Close()

  connect.ResponseDump(response) //check response only after err check

  return response.Body, nil
}

func ark_url(ark string)(string){
  base_url := os.Getenv("ARK_URL_BASE")
  url := fmt.Sprintf("%s%s", base_url, ark)
  return url
}

//80444/xv181026
func ValidArk(arkId string)(bool){
  re1 := regexp.MustCompile(`/*80444/xv[0-9]+`)
  matched1 := re1.Find([]byte(arkId))
  if string(matched1) == arkId { return true }
  return false
}

func ExtractArk(ead string)(string,error){
  et,err := ParseXML(ead)
  if err != nil{ return "", err}
  eadid := et.FindElement("//eadid")
  ark := eadid.SelectAttrValue("identifier","")
  return ark, nil
}

func ParseXML(xml_string string)(*etree.Document, error){
  xml_doc := etree.NewDocument()
  err := xml_doc.ReadFromString(xml_string)
  if err != nil { slog.Error(err.Error()); return xml_doc, errors.New("Unable to read XML") }
  return xml_doc, nil
}

// returns true if there a record published at the ark url
func CheckArk(ark string)(bool, error){
  req, err := http.NewRequest("GET", ark_url(ark), nil)
  if err != nil { slog.Error(err.Error()); return false, errors.New("unable to create http request") }
  client := &http.Client{
    Timeout: time.Second * 30,
  }
  connect.RequestDump(req)
  response, err := client.Do(req); if err != nil { slog.Error(err.Error()); return false, err }
  body, err := io.ReadAll(response.Body); if err != nil { slog.Error(err.Error()); return false, err }
  defer response.Body.Close()
  connect.ResponseDump(response) //check response only after err check
  doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
  if err != nil { slog.Error(err.Error()); return false, errors.New("unable to read response") }
  tag := doc.Find("#toc")
  if len(tag.Nodes) == 0 {
    return false, nil
  }
  return true, nil
}

func ParseResult(r io.Reader)(string, error){
  var b bytes.Buffer
  doc, err := goquery.NewDocumentFromReader(r)
  if err != nil { slog.Error(err.Error()); return "", err }
  success := doc.Find(".success")
  if success.Length() > 0 {
    err := html.Render(&b, success.Nodes[0])
    if err != nil { slog.Error(err.Error()); return "", err }
    return b.String(), nil
  } else {
    err_nodes := doc.Find(".errors")
    if err_nodes.Length() > 0 {
      err := html.Render(&b, err_nodes.Nodes[0])
      if err != nil { slog.Error(err.Error()); return "", err }
      return b.String(), nil
    }
  }
  return "", nil
}
