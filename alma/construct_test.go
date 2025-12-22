package alma

import (
  "testing"
  "os"
  "io/ioutil"
  "encoding/xml"
  "encoding/json"
  "aspace_publisher/as"
  "reflect"
  "fmt"
)

func TestConstructBib( t *testing.T){
  fstring := `<?xml version=\"1.0\" encoding=\"UTF-8\"?><record><leader>123</leader><controlfield tag="001">on1097882240</controlfield><controlfield tag="003">OCoLC</controlfield><datafield ind1="0" ind2=" " tag="041"><subfield code="a">eng</subfield></datafield><datafield ind1=" " ind2=" " tag="099"><subfield code="a">Coll 408</subfield></datafield></record>`
  expected := `<bib><suppress_from_publishing>false</suppress_from_publishing><suppress_from_external_search>true</suppress_from_external_search><record><leader>123</leader><controlfield tag="001">on1097882240</controlfield><controlfield tag="003">OCoLC</controlfield><datafield ind1="0" ind2=" " tag="041"><subfield code="a">eng</subfield></datafield><datafield ind1=" " ind2=" " tag="099"><subfield code="a">Coll 408</subfield></datafield></record></bib>`
  bib_string, _ := ConstructBib(fstring)
  var itemA Bib
  var itemB Bib
  xml.Unmarshal([]byte(bib_string), &itemA)
  xml.Unmarshal([]byte(expected), &itemB)
  if reflect.DeepEqual(itemA, itemB) != true { t.Errorf("incorrect bib rec") }
}

func TestConstructHolding(t *testing.T){
  home := os.Getenv("HOME_DIR")
  hold, err := ioutil.ReadFile(home + "fixtures/marc_3464.xml")
  if err != nil { t.Errorf("error reading file") }
  expected, err := ioutil.ReadFile(home + "/fixtures/holding_alma.xml")
  if err != nil { t.Errorf("error reading file") }
  result, _ := ConstructHolding(string(hold), "Coll 408")
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
  result, _ := ConstructItem("", "98765432987", tc.Mapify())
  var itemA Item
  var itemB Item
  json.Unmarshal([]byte(result), &itemA)
  json.Unmarshal(expected, &itemB)
  if itemA != itemB { t.Errorf("incorrect item record") }
}
