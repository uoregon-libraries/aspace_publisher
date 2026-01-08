package alma

import (
  "testing"
  "os"
  "net/http"
  "net/http/httptest"
  "io/ioutil"
  "strings"
  "fmt"
  "log"
)

func TestBuildUrl(t *testing.T){
  path := []string{"one", "two", "three", ""}
  os.Setenv("ALMA_URL", "http://blah.org")
  url := BuildUrl(path)
  if url != "http://blah.org/one/two/three" { t.Errorf("incorrect url") }

  path = []string{"one", "two", "", "three"}
  url = BuildUrl(path)
  if url != "http://blah.org/one/two/three" { t.Errorf("incorrect url") }
}

//will call ConstructBib, Post/Put, ExtractBibID
func TestProcessBib(t *testing.T){
  args := ProcessArgs{ Mms_id: "", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: true }
  fstring := bibstring_fixture4
  expected := bibstring_fixture5

  tcmap := []map[string]string{ map[string]string{} }
  fs := FunMap{ BoundwithPF: DummyBoundwithPF, NZPF: DummyNZPF, AfterBib: DummyAfterBib }
  path := "/almaws/v1/bibs" //test post
  rjson := []byte{}

  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil { t.Errorf("error reading request body") }
    if strings.Compare(string(body), expected) != 1 { t.Errorf("incorrect record posted") }
    if r.URL.Path != path {
      log.Println("Path: " + r.URL.Path)
      t.Errorf("incorrect request url")
    }
    fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?><bib><mms_id>123456789111</mms_id></bib>")
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  ProcessBib(args, fstring, rjson, tcmap, fs)
}

// tests ConstructBoundwith, calls Get/Put
func TestProcessBoundwith(t *testing.T){
  args := ProcessArgs{ Mms_id: "345634563456", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: true }
  tcmap := []map[string]string{
    map[string]string{ "boundwith": "true", "barcode":"123412341234", "ils_holding": "", "ils_item": "" },
    map[string]string{ "boundwith": "false", "ils_holding": "", "ils_item": ""},
  }
  fs := FunMap{ FetchBib: DummyFetchBibID, HoldingPF: DummyHoldingPF, ItemsPF: DummyItemsPF }
  path := "/almaws/v1/bibs/234523452345" //test Get/Put
  fstring := bibstring_fixture1
  expected := bibstring_fixture2
  bibstring := bibstring_fixture3
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != path { t.Errorf("incorrect alma path") }
    if r.Method == "PUT" {
      body, err := ioutil.ReadAll(r.Body)
      if err != nil { t.Errorf("error reading request body") }
      if strings.Compare(string(body), expected) != 1 { t.Errorf("incorrect record posted") }
      fmt.Fprintf(w, "No good deed goes unpunished")
    } else if r.Method == "GET" {
      fmt.Fprintf(w, fstring)
    } else { t.Errorf("incorrect http method") }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  ProcessBoundwith(args, bibstring, tcmap, fs)
}

func DummyBoundwithPF(args ProcessArgs, marc_string string, tcmap []map[string]string, fs FunMap){ return }
func DummyHoldingPF(args ProcessArgs, marc_string string, tcmap []map[string]string, fs FunMap){
  if tcmap[0]["mms_id"] != "234523452345" { log.Println("incorrect update to tcmap") }
  return
}
func DummyItemsPF(args ProcessArgs, tcmap []map[string]string, fs FunMap){ return }
func DummyItemPF(args ProcessArgs, item Item, tcmap map[string]string)(string, error){ return "", nil }
func DummyNZPF(list []string, filename string){ return }
func DummyAfterBib(rjson []byte, args_map map[string]string)error{
  if args_map["mms_id"] != "123456789111" { log.Println("incorrect mms_id") }
  return nil
 }
func DummyFetchBibID(barcode string)string{
  if barcode != "123412341234" { log.Println("incorrect barcode sent") }
  return "234523452345"
}
