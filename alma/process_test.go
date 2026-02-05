package alma

import (
  "testing"
  "os"
  "net/http"
  "net/http/httptest"
  "io/ioutil"
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
  fs := FunMap{ BoundwithPF: DummyBoundwithPF, NZPF: DummyNZPF, AfterBib: DummyAfterBib, SetHolding: DummySetHolding }
  path := "/almaws/v1/bibs" //test post
  rjson := []byte{}

  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil { t.Errorf("error reading request body") }
    if compareXML(string(body), expected) != true { t.Errorf("incorrect record posted") }
    if r.URL.Path != path {
      t.Errorf("incorrect request url")
    }
    fmt.Fprint(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?><bib><mms_id>123456789111</mms_id></bib>")
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
      if compareXML(string(body), expected) != true { t.Errorf("incorrect record posted") }
      fmt.Fprint(w, "No good deed goes unpunished")
    } else if r.Method == "GET" {
      fmt.Fprint(w, fstring)
    } else { t.Errorf("incorrect http method") }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  ProcessBoundwith(args, bibstring, tcmap, fs)
}

func TestProcessHolding(t *testing.T){
  path := "/almaws/v1/bibs/345634563456/holdings" //test Get/Put
  args2 := ProcessArgs{ Mms_id: "345634563456", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: true, Id_0: "Coll 408"}
  fs := FunMap{ ItemsPF: DummyItemsPF }
tcmap := []map[string]string{ map[string]string{ "boundwith": "false", "ils_holding": "", "ils_item": "" } }
  home := os.Getenv("HOME_DIR")
  hold, err := ioutil.ReadFile(home + "fixtures/marc_3464.xml")
  expected := holdingstring_fixture1
  if err != nil { t.Errorf("error reading file") }
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != path { t.Errorf("incorrect alma path") }
    if r.Method == "POST" {
      body, err := ioutil.ReadAll(r.Body)
      if err != nil { t.Errorf("error reading request body") }
      if compareXML(string(body), expected) != true { t.Errorf("incorrect record posted") }
      fmt.Fprint(w, "why we can't have nice things")
    } else if r.Method == "GET" {
      fmt.Fprint(w, string(hold))
    } else { t.Errorf("incorrect http method") }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  ProcessHolding(args2, string(hold), tcmap, fs)
}

func TestProcessItems(t *testing.T){
  tcmap := []map[string]string{ map[string]string{ "boundwith": "false", "ils_holding": "234523452345", "ils_item": "", "mms_id": "345634563456" } }
  args := ProcessArgs{ Mms_id: "345634563456", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: true }
  fs := FunMap{ UpdateTC: DummyUpdateTC, ItemPF: DummyItemPF }
  item := ""
  path := "/almaws/v1/bibs/345634563456/holdings/234523452345/items/"
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != path { t.Errorf("incorrect alma path") }
    fmt.Fprint(w, string(item))
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  ProcessItems(args, tcmap, fs)
}

func TestProcessItem(t *testing.T){
  tcmap := map[string]string{ "boundwith": "false", "ils_holding": "234523452345", "ils_item": "", "mms_id": "345634563456" }
  args := ProcessArgs{ Mms_id: "345634563456", Holding_id: "234523452345", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: true }
  path := "/almaws/v1/bibs/345634563456/holdings/234523452345/items"
  item := Item{}
  expected := itemstring_fixture2
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != path { t.Errorf("incorrect alma path") }
    if r.Method == "POST" {
      body, err := ioutil.ReadAll(r.Body)
      if err != nil { t.Errorf("error reading request body") }
      if compareJSON(string(body), expected) != true { t.Errorf("incorrect record posted") }
 }
    fmt.Fprint(w, itemstring_fixture2)
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  id,_ := ProcessItem(args, item, tcmap)
  if id != "456745674567" { t.Errorf("incorrect id returned") }
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
func DummyUpdateTC(repo_id string, holding_id string, item_id string, session_id string, tcmap map[string]string)error{ return nil }
func DummySetHolding(oclc_id string, token string)(string, error){ return fmt.Sprintf("holding %s is set", oclc_id), nil }

