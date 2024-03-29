package aw

import (
  "testing"
  "net/http"
  "net/http/httptest"
  "os"
  "path/filepath"
  "bytes"
  "fmt"
  "mime/multipart"
)

func TestPrepare(t *testing.T) {
  home_dir := os.Getenv("HOME_DIR")
  //fmt.Println(filepath.Join(home_dir, "fixtures/ead_orig.xml"))
  orig_xml, err := os.ReadFile(filepath.Join(home_dir, "fixtures/ead_orig.xml"))
  if err != nil { t.Errorf(err.Error()) } 
  repo_id := "2"
  resource_id := "123"
  _, filename, arkid, err := PrepareEad(repo_id, resource_id, string(orig_xml))
  if filename != "ORU_Coll746.xml" { t.Errorf("wrong filename: %s", filename) }
  if arkid != "80444/xv804414" { t.Errorf("wrong arkid: %s", arkid) }
  if err != nil { t.Errorf(err.Error()) }
}

func TestConvert(t *testing.T) {
  home_dir := os.Getenv("HOME_DIR")
  prepped_xml, err := os.ReadFile(filepath.Join(home_dir, "fixtures/ead_prepped.xml")) 
  if err != nil { t.Errorf(err.Error()) } 
  _, err = CallConversion(string(prepped_xml))
  if err != nil { t.Errorf(err.Error()) }
}

func TestAuthenticate(t *testing.T) {
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.PostFormValue("password") != "mypass" { t.Errorf("request bungled data") }
    if r.PostFormValue("username") != "myname" { t.Errorf("request bungled data") }
    cookie := http.Cookie{
      Name: "PHPSESSID",
      Value: "abcde1234yaddayadda",
    }
    http.SetCookie(w, &cookie)
  }))
  defer ts.Close()
  os.Setenv("AWEST_URL", ts.URL + "/")
  os.Setenv("AWEST_NAME", "myname")
  os.Setenv("AWEST_PASS", "mypass")
  session, err := authenticate("false")
  if session != "abcde1234yaddayadda" { t.Errorf("fail to parse session id") }
  if err != nil { t.Errorf(err.Error()) }
}

func dummyform()(*bytes.Buffer, string){
  form := new(bytes.Buffer)
  writer := multipart.NewWriter(form)
  boundary := writer.Boundary()
  writer.Close()
  return form, boundary
}

// Upload is essentially the same method
func TestValidate(t *testing.T) {
  form, boundary := dummyform()
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("cookie") != "PHPSESSID=abcde1234yaddayadda" { t.Errorf("request bungled adding session id to header") }
    if r.Header.Get("Content-Type") != fmt.Sprintf("multipart/form-data; boundary=%s", boundary) { t.Errorf("request bungled adding content type to header") }
    fmt.Fprintf(w, "success")
  }))
  defer ts.Close()
  os.Setenv("AWEST_URL", ts.URL + "/")
  response, err := Validate("abcde1234yaddayadda", boundary, "false", form)
  if response != "success" { t.Errorf("fail validate") }
  if err != nil { t.Errorf(err.Error()) }
}

func TestMap(t *testing.T){
  good_ark := "80444/xv205342"
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprintf(w, "<div id=\"toc\"><p>Kevin Bacon</p></div>")
  }))
  defer ts.Close()
  os.Setenv("ARK_URL_BASE", ts.URL + "/ark:")
  vals, err := MakeUploadMap(good_ark, "ead", "/tmp/abcde.xml")
  if err != nil { t.Errorf(err.Error()) }
  if vals["ark"] != good_ark { t.Errorf("did not map ark") }
  if vals["filekey"] != "ead" { t.Errorf("did not map filekey") }
  if vals["filepath"] != "/tmp/abcde.xml" { t.Errorf("did not map filepath") }
  if vals["replace"] != "1" { t.Errorf("did not correctly complete arkid test") }
}

func TestTestArk(t *testing.T){
  good_ark := "80444/xv205342"
  bad_ark := "80444/xv205343"
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/ark:" + good_ark {
      fmt.Fprintf(w, "<div id=\"toc\"><p>Kevin Bacon</p></div>")
      t.Logf("path= %s", r.URL.Path)
    } else { fmt.Fprintf(w, "<p></p>") }
  }))
  defer ts.Close()
  os.Setenv("ARK_URL_BASE", ts.URL + "/ark:")
  response, err := TestArk(good_ark)
  if response != true { t.Errorf("fail good ark") }
  if err != nil { t.Errorf(err.Error()) }
  response, err = TestArk(ts.URL + "/ark:" + bad_ark)
  if response != false { t.Errorf("fail bad ark") }
  if err != nil { t.Errorf(err.Error()) }
}
