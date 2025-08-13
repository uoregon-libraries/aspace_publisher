package alma

import (
  "testing"
  "os"
  "io/ioutil"
  "encoding/xml"
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
