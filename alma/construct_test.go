package alma

import (
  "testing"
  "os"
  "io/ioutil"
  "encoding/xml"
  "encoding/json"
  "aspace_publisher/as"
  "reflect"
  "strings"
)

func TestConstructBib( t *testing.T){
  fstring := bibstring_fixture4
  expected := bibstring_fixture5
  bib := ConstructBib(fstring, false)
  bibstr,_ := bib.Stringify()
  if strings.Compare(bibstr, expected) != 1 { t.Errorf("incorrect bib rec") }
}

func TestConstructBoundwith(t *testing.T){
  fstring := bibstring_fixture1
  expected := bibstring_fixture2
  bibstring := bibstring_fixture3
  tcmap :=  map[string]string{ "mms_id": "9999123456456" }
  bib, err := ConstructBoundwith(fstring,bibstring,tcmap)
  bibstr, err := bib.Stringify()
  if err != nil { t.Errorf("error in stringify") }
  if strings.Compare(bibstr, expected) != 1 { t.Errorf("incorrect boundwith rec") }
}

func TestConstructHolding(t *testing.T){
  home := os.Getenv("HOME_DIR")
  hold, err := ioutil.ReadFile(home + "fixtures/marc_3464.xml")
  if err != nil { t.Errorf("error reading file") }
  expected, err := ioutil.ReadFile(home + "/fixtures/holding_alma.xml")
  if err != nil { t.Errorf("error reading file") }
  var h = Holding{}
  result, _ := ConstructHolding(string(hold), h, "Coll 408")
  holdstr, err := result.Stringify()
  if err != nil { t.Errorf("stringify error") }
  var itemA Holding
  var itemB Holding
  xml.Unmarshal([]byte(holdstr), &itemA)
  xml.Unmarshal(expected, &itemB)
  if reflect.DeepEqual(itemA, itemB) != true { t.Errorf("incorrect holding rec") }
}

func TestConstructItem(t *testing.T){
  home := os.Getenv("HOME_DIR")
  tcdata, err := ioutil.ReadFile(home + "fixtures/top_container.json")
  if err != nil { t.Errorf("error reading file") }
  var tc as.TopContainer
  err = json.Unmarshal(tcdata, &tc)
  if err != nil { t.Errorf("error unmarshalling tc data") } 
  expected, err := ioutil.ReadFile(home + "/fixtures/item_alma.json")
  if err != nil { t.Errorf("error reading file") }
  item := Item{}
  result, _ := ConstructItem("98765432987",item, tc.Mapify())
  itemstr, err := result.Stringify()
  if err != nil { t.Errorf("error in stringify") }
  if strings.Compare(itemstr, string(expected)) != 1 { t.Errorf("inccorect item rec") }
}
