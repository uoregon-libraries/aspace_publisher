package alma

import (
  "testing"
  "reflect"
  "time"
  "slices"
//  "fmt"
)

func dummyTime()time.Time{
  t, _:= time.Parse("20060102150405", "20070515123442")
  return t
}
func TestConstructBWMarc(t *testing.T){
  tcmap := map[string]string{ "barcode":"123456789"}
  r := ConstructBWMarc(tcmap, dummyTime())

  expected := []Controlfield{ 
    Controlfield{ Tag:"005", Value: "20070515123442.0" },
    Controlfield{ Tag:"008", Value: "20070515i19001920oru                 eng d" },
  }
  if reflect.DeepEqual(r.Controlfield, expected) != true { t.Error("controlfields failed") }
  val := "Multiple archival collections in a shared box with barcode 123456789"
  if CheckSubfieldValue( r.Datafield, val) != true { t.Error("245 datafield has failed") }
  if CheckSubfieldValue( r.Datafield, "BoundwithRecord") != true{ t.Error("962a subfield has failed") }
  if CheckSubfieldValue( r.Datafield, "local") != true{ t.Error("9629 subfield has failed") }
}

func CheckSubfieldValue(dfs []Datafield, val string)bool{
  return slices.ContainsFunc(dfs, func(d Datafield)bool{
    return slices.ContainsFunc(d.Subfield, func(s Subfield) bool { 
      return s.Value == val
    })
  })
}

func TestConstructBWBib(t *testing.T){
  tcmap := map[string]string{ "barcode":"123456789"}
  r := ConstructBWMarc(tcmap, dummyTime())
  bib := ConstructBWBib(r)
  bstr, err := bib.Stringify()
  if err != nil { t.Error(err) }
  if compareBibs([]byte(bstr), []byte(bwbibfixture1)) != true { t.Error("construct bib failed") }  
}

func TestConstructBWHolding(t *testing.T){
  tcmap := map[string]string{ "barcode":"123456789"}
  r := ConstructBWMarc(tcmap, dummyTime())
  h, err := ConstructBWHolding(r, tcmap)
  if err != nil { t.Error(err) }
  if compareHolds([]byte(h), []byte(bwholdfixture1)) != true { t.Error("construct hold failed") }
}

func TestConstructBWItem(t *testing.T){
  tcmap := map[string]string{ "barcode":"123456789", "type":"Multiple Collection Box", "indicator":"[123456789]"}
  istr, err := ConstructBWItem("234567891", tcmap)
  if err != nil { t.Error(err) }
  if compareJSON(istr, bwitemfixture1) != true { t.Error("construct item failed") }
}

var bwbibfixture1 = `<bib><suppress_from_publishing>true</suppress_from_publishing><suppress_from_external_search>true</suppress_from_external_search><record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20070515123442.0</controlfield><controlfield tag="008">20070515i19001920oru                 eng d</controlfield><datafield tag="245" ind1="0" ind2="0"><subfield code="a">Multiple archival collections in a shared box with barcode 123456789</subfield></datafield><datafield tag="962" ind1=" " ind2=" "><subfield code="a">BoundwithRecord</subfield><subfield code="9">local</subfield></datafield></record></bib>`

var bwholdfixture1 = `<holding><suppress_from_publishing>false</suppress_from_publishing><record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20070515123442.0</controlfield><controlfield tag="008">20070515i19001920oru                 eng d</controlfield><datafield tag="852" ind1="8" ind2=" "><subfield code="b">SpecColl</subfield><subfield code="c">spmanus</subfield><subfield code="h">[123456789]</subfield></datafield></record></holding>`

var bwitemfixture1 = `{"holding_data":{"holding_id":"234567891","copy_id":"1"},"item_data":{"barcode":"123456789","policy":{"value":""},"description":"Multiple Collection Box [123456789]","base_status":{"value":"1"},"library":{"value":"SpecColl"},"location":{"value":"spmanus"},"physical_material_type":{"value":"MANUSCRIPT"}}}`