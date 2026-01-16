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

type XMLObject map[string]any

func compareXML(stringA string, stringB string) bool{
  var objA XMLObject
  var objB XMLObject
  xml.Unmarshal([]byte(stringA), &objA)
  xml.Unmarshal([]byte(stringB), &objB)
  return reflect.DeepEqual(objA, objB)
}

type JSONObject map[string]any
func compareJSON(stringA string, stringB string) bool{
  var objA JSONObject
  var objB JSONObject
  json.Unmarshal([]byte(stringA), &objA)
  json.Unmarshal([]byte(stringB), &objB)
  return reflect.DeepEqual(objA, objB)
}

func TestConstructBib( t *testing.T){
  fstring := bibstring_fixture4
  expected := bibstring_fixture5
  bib := ConstructBib(fstring, false)
  bibstr,_ := bib.Stringify()
  if compareXML(bibstr, expected) != true { t.Errorf("incorrect bib rec") }
}

func TestConstructBoundwith(t *testing.T){
  bwmarc := bibstring_fixture1
  bwmmsid := "9999123456456"
  expected := bibstring_fixture2
  bibmarc := bibstring_fixture3
  bibmmsid := "999954325432"
  tcmap :=  map[string]string{ "mms_id": bwmmsid }
  bib, err := ConstructBoundwith([]byte(bwmarc),bibmarc,bibmmsid,tcmap)
  bibstr, err := bib.Stringify()
  if err != nil { t.Errorf("error in stringify") }
  if compareXML(bibstr, expected) != true { t.Errorf("incorrect boundwith rec") }
}

// case where df774exists == true
func TestConstructBoundwith2(t *testing.T){
  bwmarc := bibstring_fixture2
  bwmmsid := "9999123456456"
  expected := bibstring_fixture2
  bibmarc := bibstring_fixture3
  bibmmsid := "999954325432"
  tcmap :=  map[string]string{ "mms_id": bwmmsid }
  bib, err := ConstructBoundwith([]byte(bwmarc),bibmarc,bibmmsid,tcmap)
  bibstr, err := bib.Stringify()
  if err != nil { t.Errorf("error in stringify") }
  if compareXML(bibstr, expected) != true { t.Errorf("incorrect boundwith rec") }
}

func TestConstructHolding(t *testing.T){
  home := os.Getenv("HOME_DIR")
  hold, err := ioutil.ReadFile(home + "fixtures/marc_3464.xml")
  if err != nil { t.Errorf("error reading file") }
  expected := holdingstring_fixture1
  if err != nil { t.Errorf("error reading file") }
  var h = Holding{}
  result, _ := ConstructHolding(string(hold), h, "Coll 408")
  holdstr, err := result.Stringify()
  if err != nil { t.Errorf("stringify error") }
  if compareXML(holdstr, expected) != true { t.Errorf("incorrect holding rec") }
}

func TestConstructItem(t *testing.T){
  home := os.Getenv("HOME_DIR")
  tcdata, err := ioutil.ReadFile(home + "fixtures/top_container.json")
  if err != nil { t.Errorf("error reading file") }
  var tc as.TopContainer
  err = json.Unmarshal(tcdata, &tc)
  if err != nil { t.Errorf("error unmarshalling tc data") } 
  expected := itemstring_fixture1
  if err != nil { t.Errorf("error reading file") }
  item := Item{}
  result, _ := ConstructItem("98765432987",item, tc.Mapify())
  itemstr, err := result.Stringify()
  if err != nil { t.Errorf("error in stringify") }
  if compareJSON(itemstr, expected) != true { t.Errorf("incorrect item rec") }
}

func TestDf774Exists(t *testing.T){
  bwbibstr := bibstring_fixture2
  bwbib_xml,_ := ParseMarc(bwbibstr)
  res := df774Exists(bwbib_xml, "999954325432")
  if res != true { t.Errorf("incorrect result") }
  res2 := df774Exists(bwbib_xml, "999912341234")
  if res2 != false { t.Errorf("incorrect result") }
}
