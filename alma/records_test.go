package alma

import (
  "testing"
  "fmt"
  "net/http/httptest"
  "net/http"
  "os"
  "encoding/json"
  "encoding/xml"
  xmlcmp "github.com/imflog/xmlcompare"

)
func TestItem( t *testing.T){
  itemjson := itemstring_fixture4
  var item Item
  json.Unmarshal([]byte(itemjson), &item)
  if item.Item_data.Item_pid != "23329204380001852" { t.Errorf("did not unmarshal correctly") }
}

func TestExtractBibID( t *testing.T){
  fstring := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><bib><mms_id>123456789111</mms_id></bib>"
  mms_id := ExtractBibID([]byte(fstring))
  if mms_id != "123456789111" { t.Errorf("incorrect mms_id") }
}
func TestGetHoldingId( t *testing.T){
  data := "{ \"holding\" : [{ \"holding_id\": \"123456789123\" }] }"
  bib_id := "98765432987"
  path := "/almaws/v1/bibs/" + bib_id + "/holdings"
   ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == path {
      fmt.Fprint(w, data)
    } else {
      t.Errorf("incorrect request url")
    }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  id := GetHoldingId(bib_id)
  if id != "123456789123" { t.Errorf("incorrect holding id") }
}

func TestExtractHoldingID( t *testing.T){
  data := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><holding><holding_id>12341234123</holding_id></holding>"
  holding_id := ExtractHoldingID([]byte(data))
  if holding_id != "12341234123" { t.Errorf("incorrect holding_id") }

}

func TestExtractItemID( t *testing.T){
  data := "{\"item_data\":{\"pid\": \"123456789\" }}"
  id := ExtractItemID([]byte(data))
  if id != "123456789" { t.Errorf("incorrect item id") }
}

func TestParseHoldingItem( t *testing.T){
  data := itemstring_fixture3
  holding,item := ParseHoldingItem([]byte(data))
  if holding != "98765432987" { t.Errorf("incorrect holding id") }
  if item != "456745674567" { t.Errorf("incorrect item id") }
}

func TestStringify(t *testing.T){
  fstring := bibstring_fixture4
  expected := bibstring_fixture5
  bib := ConstructBib("", fstring, "false")
  bib_string,err := bib.Stringify()
  if err != nil { t.Errorf("incorrect response") }
  if compareBibs([]byte(bib_string), []byte(expected)) != true { t.Errorf("incorrect rec") }
}
func TestGetSubfieldVal(t *testing.T){
  s1 := Subfield{ Code: "w", Value: "987654321" }
  s2 := Subfield{ Code: "t", Value: "Hither and yon" }
  d := Datafield{ Tag: "774", Ind1: "1", Ind2: " ", Subfield: []Subfield{ s1, s2 }}
  result := d.GetSubfieldVal("w")
  if result != "987654321" { t.Error("wrong result") }
}

func TestCheck774Id(t *testing.T){
  s1 := Subfield{ Code: "w", Value: "987654321" }
  s2 := Subfield{ Code: "t", Value: "Hither and yon" }
  d := Datafield{ Tag: "774", Ind1: "1", Ind2: " ", Subfield: []Subfield{ s1, s2 }}
  mms_ids := []string{"987654321", "876543219"}
  mms_ids2 := []string{"876543219", "789456123"}
  if d.Check774Id(mms_ids) != true { t.Error("incorrect response") }
  if d.Check774Id(mms_ids2) != false { t.Error("incorrect response") }
}

func TestCreate774(t *testing.T){
  coll := map[string]string{ "ref":"/repositories/2/resources/1", "identifier":"A 123","display_string":"Historic Histories", "mms_id": "987654321"}
  d := Create774(coll)
  expected := `<datafield tag="774" ind1="1" ind2=" "><subfield code="t">A 123 Historic Histories</subfield><subfield code="w">987654321</subfield></datafield>`
  //wanted to use reflect.DeepEqual, but could not get it to work here
  output, err := xml.Marshal(d)
  if err != nil { t.Error(err) }
  eq, err := xmlcmp.Equal(string(output), expected)
  if err != nil { t.Error(err) }
  if eq != true { t.Error("output != expected") }

}

func TestCreate774s(t *testing.T){
  var rec Record
  err := xml.Unmarshal([]byte(recfixture), &rec)
  if err != nil { t.Error(err) }

  tccolls := []map[string]string{
    map[string]string{ "ref":"/repositories/2/resources/1", "identifier":"A 044","display_string":"Lawrence Hall commonplace book", "mms_id":"99198735901852"},
    map[string]string{ "ref":"/repositories/2/resources/2", "identifier":"A 045","display_string":"James Watson Hamilton diary", "mms_id":"99900407345201852"},
map[string]string{ "ref":"/repositories/2/resources/3", "identifier":"A 046","display_string":"Edwin William Hampton correspondence", "mms_id":"99900407345101852"},
  }

  rec = Create774s(rec, tccolls)
  output, err := xml.Marshal(rec)
  if err != nil { t.Error(err) }
  eq, err := xmlcmp.Equal(string(output), recfixture2)
  if err != nil { t.Error(err) }
  if eq != true { t.Error("output != expected") }

}

func TestAdd774s(t *testing.T){
  tccolls := []map[string]string{
    map[string]string{ "ref":"/repositories/2/resources/1", "identifier":"A 044","display_string":"Lawrence Hall commonplace book", "mms_id":"99198735901852"},
    map[string]string{ "ref":"/repositories/2/resources/2", "identifier":"A 045","display_string":"James Watson Hamilton diary", "mms_id":"99900407345201852"},
map[string]string{ "ref":"/repositories/2/resources/3", "identifier":"A 046","display_string":"Edwin William Hampton correspondence", "mms_id":"99900407345101852"},
  }
  var rec Record
  err := xml.Unmarshal([]byte(recfixture1), &rec)
  if err != nil { t.Error(err) }

  rec = Add774s(rec, tccolls)

  output, err := xml.Marshal(rec)
  if err != nil { t.Error(err) }
  eq, err := xmlcmp.Equal(string(output), recfixture2)
  if err != nil { t.Error(err) }
  if eq != true { t.Error("output != expected") }
}

func TestRemove774s(t *testing.T){
  var rec Record
  err := xml.Unmarshal([]byte(recfixture2), &rec)
  if err != nil { t.Error(err) }

  mmsids := []string{"99198735901852"}
  rec = Remove774s(rec, mmsids)
  output, err := xml.Marshal(rec)
  if err != nil { t.Error(err) }
  eq, err := xmlcmp.Equal(string(output), recfixture1)
  if err != nil { t.Error(err) }
  if eq != true { t.Error("output != expected") }
}