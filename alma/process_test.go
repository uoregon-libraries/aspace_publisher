package alma

import (
  "testing"
  "os"
  "net/http"
  "net/http/httptest"
  "io/ioutil"
  "fmt"
  "log"
  "strings"
  "reflect"
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
    map[string]string{ "mms_id": "561235612355", "boundwith": "true", "barcode":"123412341234", "ils_holding": "", "ils_item": "" },
    map[string]string{ "mms_id": "345634563456","boundwith": "false", "barcode":"234562345623", "ils_holding": "", "ils_item": ""},
  }
  fs := FunMap{ HoldingPF: DummyHoldingPF, ItemsPF: DummyItemsPF }
  path := "/almaws/v1/bibs/561235612355" //test Get/Put
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
      if r.URL.Path != path { t.Errorf("incorrect alma path") }
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

//CheckTCMap calls FetchByBarcode, ParseHoldingItem
  //cases: barcode == ""
  //boundwith is true, fetch the id
  //boundwith is true, unsuccessful fetch
  //boundwith is false, create is true
  //boundwith is false, create is false

func TestCheckTCMap(t *testing.T){
  barcodes := []string{ "35025042674552","35025042674553" }
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if strings.Contains(r.URL.String(), "item_barcode=" + barcodes[0]) == true {
      fmt.Fprint(w, itemstring_fixture3)
    }
    if strings.Contains(r.URL.String(), "item_barcode=" + barcodes[1]) == true {
      w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
      w.WriteHeader(http.StatusNotFound)
      fmt.Fprint(w, "Item not found")
    }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  tcmap0 := map[string]string{ "barcode":"", "boundwith": "true", "ils_holding": "", "ils_item": "", "mms_id": "" }
  tcmap := []map[string]string{tcmap0}

  // case no barcode
  tcmapR1, err1 := CheckTCMap(tcmap)
  if err1 == nil { t.Errorf("empty barcode should error") }
  if !reflect.DeepEqual(tcmapR1,tcmap) { t.Errorf("there should be no change in returned map") }

  // case barcode, successful fetch
  tcmap[0]["barcode"] = barcodes[0]
  tcmapR2, err2 := CheckTCMap(tcmap)
  if err2 != nil { t.Errorf("there should be no errors raised") }
  if tcmapR2[0]["mms_id"] != "1231231234" { t.Errorf("mms_id is incorrect")}
  if tcmapR2[0]["ils_holding"] != "98765432987" { t.Errorf("holding is incorrect")}
  if tcmapR2[0]["ils_item"] != "456745674567" { t.Errorf("pid is incorrect")}

  //case barcode, boundwith true¸ unsuccessful fetch
  tcmap[0] = map[string]string{ "barcode":barcodes[1], "boundwith": "true", "ils_holding": "", "ils_item": "", "mms_id": "" }
  tcmapR3, err3 := CheckTCMap(tcmap)
  if err3 == nil { t.Errorf("error should be populated for failure to fetch") }
  if !reflect.DeepEqual(tcmapR3,tcmap) { t.Errorf("there should be no change in returned map") }

  // case barcode boundwith false, unsuccessful fetch
  tcmap[0]["boundwith"] = "false"
  tcmapR6, err6 := CheckTCMap(tcmap)
  if err6 != nil { t.Errorf("there should not be any errors") }
  if !reflect.DeepEqual(tcmapR6,tcmap) { t.Errorf("the map should not have changed") }
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

