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
  bib := ConstructBib("", fstring, "false")
  bibstr,_ := bib.Stringify()
  if compareBibs([]byte(bibstr), []byte(expected)) != true { t.Errorf("incorrect bib rec") }
}

func TestUpdateBoundwith(t *testing.T){
  bwmarc := bibstring_fixture1
  bwmmsid := "9999123456456"
  expected := bibstring_fixture2
  bibmarc := bibstring_fixture3
  bibmmsid := "345634563456"
  tcmap :=  map[string]string{ "mms_id": bwmmsid }
  bib, err := UpdateBoundwith([]byte(bwmarc),bibmarc,bibmmsid,tcmap)
  bibstr, err := bib.Stringify()
  if err != nil { t.Errorf("error in stringify") }
  if compareBibs([]byte(bibstr), []byte(expected)) != true { t.Errorf("incorrect boundwith rec") }
}

// case where df774exists == true
func TestUpdateBoundwith2(t *testing.T){
  bwmarc := bibstring_fixture2
  bwmmsid := "9999123456456"
  expected := bibstring_fixture2
  bibmarc := bibstring_fixture3
  bibmmsid := "345634563456"
  tcmap :=  map[string]string{ "mms_id": bwmmsid }
  bib, err := UpdateBoundwith([]byte(bwmarc),bibmarc,bibmmsid,tcmap)
  bibstr, err := bib.Stringify()
  if err != nil { t.Errorf("error in stringify") }
  if compareBibs([]byte(bibstr), []byte(expected)) != true { t.Errorf("incorrect boundwith rec") }
}

func TestConstructHolding(t *testing.T){
  home := os.Getenv("HOME_DIR")
  hold, err := ioutil.ReadFile(home + "fixtures/marc_3464.xml")
  if err != nil { t.Errorf("error reading file") }
  expected := holdingstring_fixture1
  if err != nil { t.Errorf("error reading file") }
  result, _ := ConstructHolding(string(hold), "Coll 408")
  holdstr, err := result.Stringify()
  if err != nil { t.Errorf("stringify error") }
  if compareHolds([]byte(holdstr), []byte(expected)) != true { t.Errorf("incorrect holding rec") }
}

func TestUpdateHolding(t *testing.T){
  home := os.Getenv("HOME_DIR")
  marc, err := ioutil.ReadFile(home + "fixtures/marc_3464b.xml")
  if err != nil { t.Errorf("error reading file") }
  hold := holdingstring_fixture2
  expected := holdingstring_fixture3
  holdstr, _ := UpdateHolding(string(marc), hold)
    if compareHolds([]byte(holdstr), []byte(expected)) != true { t.Errorf("incorrect holding rec") }

  hold2 := holdingstring_fixture3
  holdstr,err = UpdateHolding(string(marc), hold2)
  if err.Error() != "skip update" { t.Errorf("error should be skip") }
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
  var bwbib Bib
  bwbibstr := bibstring_fixture2
  xml.Unmarshal([]byte(bwbibstr), &bwbib)
  res := df774Exists(bwbib, "9999123456123")
  if res != true { t.Errorf("incorrect result") }
  res2 := df774Exists(bwbib, "999912341234")
  if res2 != false { t.Errorf("incorrect result") }
}
