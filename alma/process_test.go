package alma

import (
  "testing"
  "os"
  "net/http"
  "net/http/httptest"
  "io/ioutil"
  "fmt"
  "strings"
  "reflect"
  "encoding/json"
  "encoding/xml"
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

func compareBibs(b1 []byte, b2 []byte)bool{
  var bib1 Bib
  var bib2 Bib
  xml.Unmarshal(b1, &bib1)
  xml.Unmarshal(b2, &bib2)
  return reflect.DeepEqual(bib1, bib2)
}

func compareHolds(h1 []byte, h2 []byte)bool{
  var hold1 Holding
  var hold2 Holding
  xml.Unmarshal(h1, &hold1)
  xml.Unmarshal(h2, &hold2)
  return reflect.DeepEqual(hold1, hold2)
}

//will call ConstructBib, Post/Put, ExtractBibID
func TestProcessBib(t *testing.T){
  args1 := ProcessArgs{ Mms_id: "", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: true }
  args2 := ProcessArgs{ Mms_id: "654365436543", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: false }
  marc := bibstring_fixture4
  expected1 := bibstring_fixture5
  expected2 := bibstring_fixture6
  tcmap := []map[string]string{ map[string]string{} }
  fs := FunMap{ BoundwithPF: DummyBoundwithPF, CallWorker: DummyCallWorker, AfterBib: DummyAfterBib, SetHolding: DummySetHolding }
  path1 := "/almaws/v1/bibs" //test post
  path2 := "/almaws/v1/bibs/654365436543"
  rjson := []byte{}

  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil { t.Errorf("error reading request body") }

    if r.Method == "POST" {
      if r.URL.Path != path1 { t.Errorf("incorrect request url") }
      if compareBibs(body, []byte(expected1)) != true { t.Errorf("incorrect record posted") }
      fmt.Fprint(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?><bib><mms_id>654365436543</mms_id></bib>")
    }
    if r.Method == "PUT" {
      if r.URL.Path != path2 { t.Errorf("incorrect request url") }
      if compareBibs(body, []byte(expected2)) != true { t.Errorf("incorrect record posted") }
      fmt.Fprint(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?><bib><mms_id>654365436543</mms_id></bib>")
    }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  ProcessBib(args1, marc, rjson, tcmap, fs)
  ProcessBib(args2, marc, rjson, tcmap, fs)
}

// tests ConstructBoundwith, calls Get/Put
func TestProcessBoundwith(t *testing.T){
  args := ProcessArgs{ Mms_id: "345634563456", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: true }
  tcmap := []map[string]string{
    map[string]string{ "mms_id": "561235612355", "boundwith": "true", "barcode":"123412341234", "ils_holding": "234567234567", "ils_item": "765476547654" },
    map[string]string{ "mms_id": "345634563456","boundwith": "false", "barcode":"234562345623", "ils_holding": "", "ils_item": ""},
  }
  fs := FunMap{ HoldingAPF: DummyHoldingAPF, ItemsPF: DummyItemsPF, CheckForMissingPF: DummyCheckForMissing }
  path := "/almaws/v1/bibs/561235612355" //test Get/Put
  initialbw := bibstring_fixture1
  expected := bibstring_fixture2
  marc := bibstring_fixture3
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != path { t.Errorf("incorrect alma path") }
    if r.Method == "PUT" {
      body, err := ioutil.ReadAll(r.Body)
      if err != nil { t.Errorf("error reading request body") }
      if compareBibs(body, []byte(expected)) != true { t.Errorf("incorrect record posted") }
      fmt.Fprint(w, expected)
    } else if r.Method == "GET" {
      if r.URL.Path != path { t.Errorf("incorrect alma path") }
      fmt.Fprint(w, initialbw)
    } else { t.Errorf("incorrect http method") }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  ProcessBoundwith(args, marc, tcmap, fs)
}

func TestProcessHoldingA(t *testing.T){
  path1 := "/almaws/v1/bibs/345634563456/holdings"
  args1 := ProcessArgs{ Mms_id: "345634563456", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: true, Id_0: "Coll 408"}
  fs := FunMap{ HoldingBPF: DummyHoldingBPF, ItemsPF: DummyItemsPF }
  tcmap1 := []map[string]string{ map[string]string{ "boundwith": "false", "ils_holding": "", "ils_item": "" } }
  home := os.Getenv("HOME_DIR")
  marc1, err := ioutil.ReadFile(home + "fixtures/marc_3464.xml")
  if err != nil { t.Errorf("error reading file") }
  expected1 := holdingstring_fixture1

  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil { t.Errorf("error reading request body") }
    if r.Method == "POST" {
      if r.URL.Path != path1 { t.Errorf("incorrect alma path") }
      if compareHolds(body, []byte(expected1)) != true { t.Errorf("incorrect record posted for POST") }
      fmt.Fprint(w, "fiddledeedee")
    }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  ProcessHoldingA(args1, string(marc1), tcmap1, fs)
}

func TestProcessHoldingB(t *testing.T){
  args2 := ProcessArgs{ Mms_id: "345634563456", Holding_id:"456745674567", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: false, Id_0: "Coll 408"}
  path2 := "/almaws/v1/bibs/345634563456/holdings/456745674567"
  fs := FunMap{ ItemsPF: DummyItemsPF }
  tcmap2 := []map[string]string{ map[string]string{ "boundwith": "false", "ils_holding": "456745674567", "ils_item": "234523452345" } }
  home := os.Getenv("HOME_DIR")
  marc2, err := ioutil.ReadFile(home + "fixtures/marc_3464b.xml")
  if err != nil { t.Errorf("error reading file") }
  expected2 := holdingstring_fixture3
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil { t.Errorf("error reading request body") }
    if r.Method == "PUT" {
      if r.URL.Path != path2 { t.Errorf("incorrect alma path") }
      if compareHolds(body, []byte(expected2)) != true { t.Errorf("incorrect record posted for PUT") }
      fmt.Fprint(w, "arglebarglesnickersnack")
    } else { fmt.Fprint(w, holdingstring_fixture2) }// only happens on an update
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  ProcessHoldingB(args2, string(marc2), tcmap2, fs)
}

func TestProcessItems(t *testing.T){
  tcmap1 := []map[string]string{ map[string]string{ "boundwith": "false", "ils_holding": "98765432987", "ils_item": "", "mms_id": "345634563456" } }
  tcmap2 := []map[string]string{ map[string]string{ "boundwith": "false", "ils_holding": "98765432987", "ils_item": "456745674567", "mms_id": "345634563456" } }
  tcmap3 := []map[string]string{ map[string]string{ "boundwith": "true", "ils_holding": "98765432988", "ils_item": "456745674568", "mms_id": "345634563457" } }

  args1 := ProcessArgs{ Mms_id: "345634563456", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: true }
  args2 := ProcessArgs{ Mms_id: "345634563456", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: false }
  args3 := ProcessArgs{ Mms_id: "345634563456", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: false }

  fs := FunMap{ UpdateTC: DummyUpdateTC, ItemPF: DummyItemPF }
  path := "/almaws/v1/bibs/345634563456/holdings/98765432987/items/456745674567"
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != path { t.Errorf("incorrect alma path") }
    fmt.Fprint(w, itemstring_fixture2)
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  ProcessItems(args1, tcmap1, fs)
  ProcessItems(args2, tcmap2, fs)
  ProcessItems(args3, tcmap3, fs)
}

func TestProcessItem(t *testing.T){
  tcmap1 := map[string]string{ "boundwith": "false", "ils_holding": "98765432987", "ils_item": "", "mms_id": "345634563456", "barcode": "35025042674552", "type":"Box", "indicator":"1" }
  tcmap2 :=  map[string]string{ "boundwith": "false", "ils_holding": "98765432987", "ils_item": "456745674567", "mms_id": "345634563456", "barcode": "35025042674552", "type":"Box", "indicator":"1" }
  args1 := ProcessArgs{ Mms_id: "345634563456", Holding_id: "98765432987", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: true }
  args2 := ProcessArgs{ Mms_id: "345634563456", Holding_id: "98765432987", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: false }

  path1 := "/almaws/v1/bibs/345634563456/holdings/98765432987/items"
  path2 := "/almaws/v1/bibs/345634563456/holdings/98765432987/items/456745674567"
  expected1 := itemstring_fixture1 //create, no pid
  expected2 := itemstring_fixture2 //update, pid, plus changes to policy, description
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil { t.Errorf("error reading request body") }
    if r.Method == "POST" {
      if r.URL.Path != path1 { t.Errorf("incorrect alma path") }

      if compareJSON(string(body), expected1) != true { t.Errorf("incorrect record posted") }
      fmt.Fprint(w, itemstring_fixture2)
    } else {
      if r.URL.Path != path2 { t.Errorf("incorrect alma path") }
      if compareJSON(string(body), expected2) != true { t.Errorf("incorrect record posted") }
      fmt.Fprint(w, itemstring_fixture2)
    }
    }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  item := Item{}
  id,_ := ProcessItem(args1, item, tcmap1)
  if id != "456745674567" { t.Errorf("incorrect id returned") }
  json.Unmarshal([]byte(itemstring_fixture5), &item)
  id,_ = ProcessItem(args2, item, tcmap2)
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

  //case barcode, boundwith true¸ unsuccessful fetch
  tcmap[0] = map[string]string{ "barcode":barcodes[1], "boundwith": "true", "ils_holding": "", "ils_item": "", "mms_id": "" }
  tcmapR3, err3 := CheckTCMap(tcmap)
  if err3 == nil { t.Errorf("error should be populated for failure to fetch") }
  if !reflect.DeepEqual(tcmapR3,tcmap) { t.Errorf("there should be no change in returned map") }
}

func TestGetBarcodes(t *testing.T){
  tcdata0 := map[string]string{ "barcode":"ronco3000", "boundwith": "true", "ils_holding": "", "ils_item": "", "mms_id": "" }
  tcdata1 := map[string]string{ "barcode":"ronco3001", "boundwith": "false", "ils_holding": "", "ils_item": "", "mms_id": "" }
  tcdata2 := map[string]string{ "barcode":"ronco3002", "boundwith": "false", "ils_holding": "", "ils_item": "", "mms_id": "" }

  tcmap := []map[string]string{tcdata0, tcdata1, tcdata2}
  bc := GetBarcodes(tcmap)
  if !reflect.DeepEqual(bc, []string{"ronco3001", "ronco3002"}) { t.Errorf("incorrect result") }
}

func TestParseShortItem(t *testing.T){
  item := `{"bib_data":{"title":"Rotten banana"},"holding_data":{"call_number":"Fruit1223"},"item_data":{"barcode":"alma3000","description":"unarranged basket"}}`
  si := ParseShortItem(item)
  expected := ShortItem{Barcode:"alma3000", Description: "unarranged basket",  CallNumber: "Fruit1223", Title: "Rotten banana"}
  if !reflect.DeepEqual(si, expected) { t.Errorf("incorrect result") }
  fmt.Println(si)
  fmt.Println(expected)
}

func TestSIStringify(t *testing.T){
  si := ShortItem{Barcode:"alma3000", Description: "unarranged basket",  CallNumber: "Fruit1223", Title: "Rotten banana"}
  if si.Stringify() != "alma3000,Fruit1223,Rotten banana,unarranged basket" { t.Errorf("incorrect reponse") }
}

func TestBuildMessageFromMissing(t *testing.T){
  si1 := ShortItem{Barcode:"alma3000", Description: "unarranged basket",  CallNumber: "Fruit1223", Title: "Rotten banana"}
  si2 := ShortItem{Barcode:"alma3001", Description: "unarranged basket",  CallNumber: "Fruit1223", Title: "Rotten apple"}
  items := []ShortItem{si1, si2}
  msg_arr := BuildMessageFromMissing(items)
  expected := []string{"alma3000,Fruit1223,Rotten banana,unarranged basket","alma3001,Fruit1223,Rotten apple,unarranged basket"}
  if !reflect.DeepEqual(msg_arr, expected) { t.Errorf("incorrect result") }
}

func TestCompileMissing(t *testing.T){
  jsonstr := `{"item":[{"bib_data":{"title":"Rotten banana"},"holding_data":{"call_number":"Fruit1223"},"item_data":{"barcode":"alma3000","description":"unarranged basket"}},{"bib_data":{"title":"Rotten apple"},"holding_data":{"call_number":"Fruit1223"},"item_data":{"barcode":"alma3001","description":"unarranged basket"}}]`
  si1 := ShortItem{Barcode:"alma3000", Description: "unarranged basket",  CallNumber: "Fruit1223", Title: "Rotten banana"}
  si2 := ShortItem{Barcode:"alma3001", Description: "unarranged basket",  CallNumber: "Fruit1223", Title: "Rotten apple"}

  barcodes := []string{"ronco3100", "ronco3101"}
  expected := []ShortItem{si1, si2}
  response := CompileMissing([]byte(jsonstr), barcodes)
  if !reflect.DeepEqual(response, expected) { t.Errorf("incorrect response") }
}

func TestCheckItemsForMissing(t *testing.T){
  path1 := "/almaws/v1/bibs/345634563456/holdings/98765432987/items"
  jsonstr := `{"item":[{"bib_data":{"title":"Rotten banana"},"holding_data":{"call_number":"Fruit1223"},"item_data":{"barcode":"alma3000","description":"unarranged basket"}},{"bib_data":{"title":"Rotten apple"},"holding_data":{"call_number":"Fruit1223"},"item_data":{"barcode":"alma3001","description":"unarranged basket"}}]`

  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != path1 { t.Errorf("incorrect request url") }
    fmt.Fprint(w, jsonstr)
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")

  args := ProcessArgs{ Mms_id: "345634563456", Holding_id: "98765432987", Filename: "test", Session_id: "123123123", Repo_id: "2", Resource_id: "1234", Create: true }
  tcdata0 := map[string]string{ "barcode":"ronco3100", "boundwith": "true", "ils_holding": "", "ils_item": "", "mms_id": "" }
  tcdata1 := map[string]string{ "barcode":"ronco3101", "boundwith": "true", "ils_holding": "", "ils_item": "", "mms_id": "" }
  tcmap := []map[string]string{tcdata0, tcdata1}
  CheckItemsForMissing(args, tcmap)
}

func TestCallWorker(t *testing.T){
  worker_path := "startJob"
  args := map[string]string{ "id": "banana", "value": "yellow" }
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" + worker_path { t.Errorf("incorrect request path: %v", r.URL.Path) }
    if r.URL.RawQuery != "id=banana&value=yellow" { t.Errorf("incorrect param: ") }
    fmt.Fprint(w, "ok")
  }))
  defer ts.Close()
  os.Setenv("WORKER_URL", ts.URL)
  CallWorker(worker_path, args) 
}
func TestBuildWorkerUrl( t *testing.T){
  worker_path := "startJob"
  os.Setenv("WORKER_URL", "http://riverservice.org")
  args := map[string]string{ "id": "banana", "value": "yellow" }
  expected := "http://riverservice.org/startJob?id=banana&value=yellow"
  _url := BuildWorkerUrl(worker_path, args)
  if _url != expected { t.Errorf("incorrect response") }
}
