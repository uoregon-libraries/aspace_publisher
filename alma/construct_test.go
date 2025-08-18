package alma

import (
  "testing"
  "os"
  "io/ioutil"
  "encoding/xml"
  "encoding/json"
  "aspace_publisher/as"
  "reflect"
)

func TestConstructBib( t *testing.T){
  fstring := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><record><blah>banana</blah></record>"
  bib_string := ConstructBib(fstring)
  if bib_string != "<bib><record><blah>banana</blah></record></bib>" {
    t.Errorf("incorrect bib record")
  }
}

func TestConstructHolding(t *testing.T){
  home := os.Getenv("HOME_DIR")
  hold, err := ioutil.ReadFile(home + "fixtures/marc_3464.xml")
  if err != nil { t.Errorf("error reading file") }
  expected, err := ioutil.ReadFile(home + "/fixtures/holding_alma.xml")
  if err != nil { t.Errorf("error reading file") }
  result := ConstructHolding(string(hold))
  var itemA Holding
  var itemB Holding
  xml.Unmarshal([]byte(result), &itemA)
  xml.Unmarshal(expected, &itemB)
  if reflect.DeepEqual(itemA, itemB) != true { t.Errorf("incorrect holding rec") }
  strA,_ := xml.Marshal(itemA)
  strB,_ := xml.Marshal(itemB)
  if string(strA) != string(strB) { t.Errorf("incorrect holding record") }
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
  result := ConstructItem("", "98765432987", "true", tc.Mapify())
  var itemA Item
  var itemB Item
  json.Unmarshal([]byte(result), &itemA)
  json.Unmarshal(expected, &itemB)
  if itemA != itemB { t.Errorf("incorrect item record") }
}
