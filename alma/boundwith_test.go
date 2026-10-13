package alma

import (
  "testing"
  "reflect"
  "time"
//  "slices"
  "net/http/httptest"
  "net/http"
  "os"
  "io/ioutil"
  "github.com/google/go-cmp/cmp"
  "encoding/xml"
  //xmlcmp "github.com/imflog/xmlcompare"
  "fmt"
)

func dummyTime()time.Time{
  t, _:= time.Parse("20060102150405", "20070515123442")
  return t
}
//    body, err := ioutil.ReadAll(r.Body)
//  r.Method == "POST"
func DummyProcessBWHolding(args ProcessArgs, rec Record, tcmap map[string]string, fs BWFunMap){}
func DummyProcessBWItem(args ProcessArgs, tcmap map[string]string, fs BWFunMap){}

func TestProcessBWBib(t *testing.T){
  //args: args, tcmap, tccolls
  args1 := ProcessArgs{ Create: true, Filename: "test" }
  args2 := ProcessArgs{ Create: false, Filename: "test"}
  tcmap := map[string]string{ "barcode":"123456789", "type":"Multiple Collection Box", "indicator":"[123456789]" }
  tccolls1 := []map[string]string{}
  tccolls2 := []map[string]string{ 
    map[string]string{ "ref":"/repositories/2/resources/1", "identifier":"A 044","display_string":"Lawrence Hall commonplace book", "mms_id":"99198735901852"},
    map[string]string{ "ref":"/repositories/2/resources/2", "identifier":"A 045","display_string":"James Watson Hamilton diary", "mms_id":"99900407345201852"}}
  fs := BWFunMap{ BibPF: ProcessBWBib, HoldingPF: DummyProcessBWHolding, ItemPF: DummyProcessBWItem, Time: dummyTime}
  //mock get/post/put, check if the submitted content is correct
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
      body, err := ioutil.ReadAll(r.Body) //wil be []byte
      if err != nil { t.Error(err) }

      if diff := cmp.Diff(string(body), bwbibfixture1); diff != "" { t.Error("post incorrect: " + diff) }
      fmt.Fprint(w, "ok")
    } else if r.Method == "PUT" {
      body, err := ioutil.ReadAll(r.Body) //wil be []byte
      if err != nil { t.Error(err) }
      if diff := cmp.Diff(string(body), bwbibfixture4); diff != "" { t.Error("put incorrect: " + diff) }
      fmt.Fprint(w, "ok")
    } else if r.Method == "GET" {
      fmt.Print(w, bwbibfixture3)
    }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")

  ProcessBWBib(args1, tcmap, tccolls1, fs)
  ProcessBWBib(args2, tcmap, tccolls2, fs)
}

func TestProcessBWHolding(t *testing.T){
  //args: args, rec, tcmap,
  args1 := ProcessArgs{ Create: true, Filename: "test", Mms_id: "987654321" }
  args2 := ProcessArgs{ Create: false, Filename: "test", Mms_id: "987654321", Holding_id: "456789123" }
  var rec Record
  err := xml.Unmarshal([]byte(recfixture3), &rec)
  if err != nil { t.Error(err) }
  tcmap1 := map[string]string{ "barcode":"3502501231242352345214", "type":"Multiple Collection Box", "indicator":"[3502501231242352345214]" }
  tcmap2 := map[string]string{ "barcode":"123456789", "type":"Multiple Collection Box", "indicator":"[123456789]" }

  fs := BWFunMap{ HoldingPF: ProcessBWHolding, ItemPF: DummyProcessBWItem, Time: dummyTime}
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
      body, err := ioutil.ReadAll(r.Body) //wil be []byte
      if err != nil { t.Error(err) }
      if diff := cmp.Diff(string(body), bwholdfixture); diff != "" { t.Error("post incorrect: " + diff) }
      fmt.Fprint(w, "ok")
    } else if r.Method == "PUT" {
      body, err := ioutil.ReadAll(r.Body) //wil be []byte
      if err != nil { t.Error(err) }
      if diff := cmp.Diff(string(body), bwholdfixture1); diff != "" { t.Error("put incorrect: " + diff) }
      fmt.Fprint(w, "ok")
    } else if r.Method == "GET" {
      fmt.Fprint(w, bwholdfixturea)
    }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")

  ProcessBWHolding(args1, rec, tcmap1, fs) //create holding
  ProcessBWHolding(args2, rec, tcmap2, fs) //update barcode
}

func TestProcessBWItem(t *testing.T){
  args := ProcessArgs{ Create: true, Filename: "test", Mms_id: "987654321", Holding_id: "234567891" }
  tcmap1 := map[string]string{ "barcode":"123456789", "type":"Multiple Collection Box", "indicator":"[123456789]" }
  tcmap2 := map[string]string{ "barcode":"456789123", "type":"Multiple Collection Box", "indicator":"[456789123]", "ils_holding":"234567891", "ils_item":"345678912" }

  fs := BWFunMap{ HoldingPF: DummyProcessBWHolding, ItemPF: ProcessBWItem, Time: dummyTime, UpdateTC: DummyUpdateTC2}
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
      body, err := ioutil.ReadAll(r.Body) //wil be []byte
      if err != nil { t.Error(err) }
      if diff := cmp.Diff(string(body), bwitemfixture1); diff != "" { t.Error("post incorrect: " + diff) }
      fmt.Fprint(w, "ok")
    } else if r.Method == "PUT" {
      body, err := ioutil.ReadAll(r.Body) //wil be []byte
      if err != nil { t.Error(err) }
      if diff := cmp.Diff(string(body), bwitemfixture3); diff != "" { t.Error("put incorrect: " + diff) }
      fmt.Fprint(w, "ok")
    } else if r.Method == "GET" {
      fmt.Fprint(w, bwitemfixture1a)
    }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")

  ProcessBWItem(args, tcmap1, fs) //create holding
  ProcessBWItem(args, tcmap2, fs) //update barcode}
}

func TestConstructBWMarc(t *testing.T){
  tcmap := map[string]string{ "barcode":"123456789"}
  r := ConstructBWMarc(tcmap, dummyTime())

  expected := []Controlfield{ 
    Controlfield{ Tag:"005", Value: "20070515123442.0" },
    Controlfield{ Tag:"008", Value: "20070515i19001920oru                 eng d" },
  }
  if reflect.DeepEqual(r.Controlfield, expected) != true { t.Error("controlfields failed")     }

  expectedDF := `<datafield tag="245" ind1="0" ind2="0"><subfield code="a">Multiple archival collections in a shared box with barcode 123456789</subfield></datafield><datafield tag="962" ind1=" " ind2=" "><subfield code="a">BoundwithRecord</subfield><subfield code="9">local</subfield></datafield>`

  dfs, err := xml.Marshal(r.Datafield)
  if err != nil { t.Error(err) }
  if cmp.Diff(string(dfs), expectedDF) != "" { t.Error("datafields not equal") }
/*
  eq, err := xmlcmp.Equal(string(dfs), expectedDF)
  if err != nil { t.Error(err) }
  if eq != false { t.Error("output does not equal expected") }
*/
}

func TestExtractMmsids(t *testing.T){
  tccolls := []map[string]string{
    map[string]string{ "ref":"/repositories/2/resources/1", "identifier":"A 044","display_string":"Lawrence Hall commonplace book", "mms_id":"99198735901852"},
    map[string]string{ "ref":"/repositories/2/resources/2", "identifier":"A 045","display_string":"James Watson Hamilton diary", "mms_id":"99900407345201852"}}
  expected := []string{"99198735901852", "99900407345201852"}
  if diff := cmp.Diff(extractMmsids(tccolls), expected); diff != "" { t.Error("error extracting mmsids: " + diff) }

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

func TestUpdateBWHolding(t *testing.T){
  tcdata := map[string]string{ "barcode":"345678912"}
  holdingxml, err := UpdateBWHolding(bwholdfixture1a, tcdata)
  if err != nil { t.Error(err) }
/*
  eq, err := xmlcmp.Equal(holdingxml, bwholdfixture2)
  if err != nil { t.Error(err) }
  if eq != false { t.Error("output does not equal expected") }
  if compareHolds([]byte(holdingxml), []byte(bwholdfixture2)) { t.Error("incorrect holding update") }
*/
  if diff := cmp.Diff(holdingxml, bwholdfixture2); diff != "" { t.Error("second try failed") }
}

func TestConstructBWItem(t *testing.T){
  tcmap := map[string]string{ "barcode":"123456789", "type":"Multiple Collection Box", "indicator":"[123456789]"}
  istr, err := ConstructBWItem("234567891", tcmap)
  if err != nil { t.Error(err) }
  if compareJSON(istr, bwitemfixture1) != true { t.Error("construct item failed") }
}

func TestUpdateBWItem(t *testing.T){
  tcdata := map[string]string{ "barcode": "345678912", "type": "Multiple Collection Box", "indicator":"[345678912]"}
  itemjson, err := UpdateBWItem("234567891", bwitemfixture1a, tcdata)
  if err != nil { t.Error(err) }
  if diff := cmp.Diff(itemjson, bwitemfixture2); diff != "" { t.Error("incorrect item update") }
}

var recfixture = `<record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20231122114938.0</controlfield><controlfield tag="008">130912i19001920oru                 eng d</controlfield><datafield ind1="0" ind2="0" tag="245"><subfield code="a">Multiple archival collections in a shared box with barcode 3502501231242352345214</subfield></datafield><datafield ind1=" " ind2=" " tag="962"><subfield code="a">BoundWithRecord</subfield><subfield code="9">local</subfield></datafield></record>`

//contains Lawrence Hall
var recfixture1 = `<record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20231122114938.0</controlfield><controlfield tag="008">130912i19001920oru                 eng d</controlfield><datafield ind1="0" ind2="0" tag="245"><subfield code="a">Multiple archival collections in a shared box with barcode 3502501231242352345214</subfield></datafield><datafield ind1="1" ind2=" " tag="774"><subfield code="t">A 044 Lawrence Hall commonplace book</subfield><subfield code="w">99198735901852</subfield></datafield><datafield ind1=" " ind2=" " tag="962"><subfield code="a">BoundWithRecord</subfield><subfield code="9">local</subfield></datafield></record>`

//contains Lawrence Hall, James Watson Hamilton, Edwin William Hampton
var recfixture2 = `<record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20231122114938.0</controlfield><controlfield tag="008">130912i19001920oru                 eng d</controlfield><datafield ind1="0" ind2="0" tag="245"><subfield code="a">Multiple archival collections in a shared box with barcode 3502501231242352345214</subfield></datafield><datafield ind1=" " ind2=" " tag="962"><subfield code="a">BoundWithRecord</subfield><subfield code="9">local</subfield></datafield><datafield ind1="1" ind2=" " tag="774"><subfield code="t">A 044 Lawrence Hall commonplace book</subfield><subfield code="w">99198735901852</subfield></datafield><datafield ind1="1" ind2=" " tag="774"><subfield code="t">A 045 James Watson Hamilton diary</subfield><subfield code="w">99900407345201852</subfield></datafield><datafield ind1="1" ind2=" " tag="774"><subfield code="t">A 046 Edwin William Hampton correspondence</subfield><subfield code="w">99900407345101852</subfield></datafield></record>`

var recfixture3 = `<record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20070515123442.0</controlfield><controlfield tag="008">20070515i19001920oru                 eng d</controlfield><datafield tag="245" ind1="0" ind2="0"><subfield code="a">Multiple archival collections in a shared box with barcode 3502501231242352345214</subfield></datafield><datafield tag="962" ind1=" " ind2=" "><subfield code="a">BoundwithRecord</subfield><subfield code="9">local</subfield></datafield></record>`
//contains Lawrence Hall
var tcfixture1 = `{"lock_version":25,"barcode":"35025045445273","ils_holding_id":"22497408420001852","ils_item_id":"23496874430001852","indicator":"[35025045445273]","created_by":"admin","last_modified_by":"scuaapi","create_time":"2018-02-08T05:10:49Z","system_mtime":"2026-04-06T21:13:33Z","user_mtime":"2026-04-06T21:13:33Z","type":"Multiple Collection Box","jsonmodel_type":"top_container","active_restrictions":[],"container_locations":[{"jsonmodel_type":"container_location","status":"current","start_date":"2017-12-28","system_mtime":"2026-04-06 21:13:33 UTC","user_mtime":"2026-04-06 21:13:33 UTC","ref":"/locations/2826"}],"series":[],"collection":[{"ref":"/repositories/2/resources/72","identifier":"A 044","display_string":"Lawrence Hall commonplace book"}],"uri":"/repositories/2/top_containers/4001","repository":{"ref":"/repositories/2"},"restricted":false,"is_linked_to_published_record":true,"display_string":"Multiple collection box [35025045445273]: [Barcode: 35025045445273]","long_display_string":"Multiple collection box [35025045445273] [Barcode: 35025045445273], Knight, M180_Deck 3 [4SCUA0005520, Range:  A, Bay: 1, Shelf: 2], A 044, Lawrence Hall commonplace book"}`

var tcfixture2 = `{"lock_version":25,"barcode":"35025045445273","ils_holding_id":"22497408420001852","ils_item_id":"23496874430001852","indicator":"[35025045445273]","created_by":"admin","last_modified_by":"scuaapi","create_time":"2018-02-08T05:10:49Z","system_mtime":"2026-04-06T21:13:33Z","user_mtime":"2026-04-06T21:13:33Z","type":"Multiple Collection Box","jsonmodel_type":"top_container","active_restrictions":[],"container_locations":[{"jsonmodel_type":"container_location","status":"current","start_date":"2017-12-28","system_mtime":"2026-04-06 21:13:33 UTC","user_mtime":"2026-04-06 21:13:33 UTC","ref":"/locations/2826"}],"series":[],"collection":[{"ref":"/repositories/2/resources/72","identifier":"A 044","display_string":"Lawrence Hall commonplace book"},{"ref":"/repositories/2/resources/73","identifier":"A 045","display_string":"James Watson Hamilton diary"}],"uri":"/repositories/2/top_containers/4001","repository":{"ref":"/repositories/2"},"restricted":false,"is_linked_to_published_record":true,"display_string":"Multiple collection box [35025045445273]: [Barcode: 35025045445273]","long_display_string":"Multiple collection box [35025045445273] [Barcode: 35025045445273], Knight, M180_Deck 3 [4SCUA0005520, Range:  A, Bay: 1, Shelf: 2], A 044, Lawrence Hall commonplace book"}`
var bwbibfixture1 = `<bib><suppress_from_publishing>true</suppress_from_publishing><suppress_from_external_search>true</suppress_from_external_search><record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20070515123442.0</controlfield><controlfield tag="008">20070515i19001920oru                 eng d</controlfield><datafield tag="245" ind1="0" ind2="0"><subfield code="a">Multiple archival collections in a shared box with barcode 123456789</subfield></datafield><datafield tag="962" ind1=" " ind2=" "><subfield code="a">BoundwithRecord</subfield><subfield code="9">local</subfield></datafield></record></bib>`
var bwbibfixture1a =
`<bib><suppress_from_publishing>true</suppress_from_publishing><suppress_from_external_search>true</suppress_from_external_search><mms_id></mms_id><record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20070515123442.0</controlfield><controlfield tag="008">20070515i19001920oru                 eng d</controlfield><datafield tag="245" ind1="0" ind2="0"><subfield code="a">Multiple archival collections in a shared box with barcode 123456789</subfield></datafield><datafield tag="962" ind1=" " ind2=" "><subfield code="a">BoundwithRecord</subfield><subfield code="9">local</subfield></datafield></record></bib>`
//has lawrence
var bwbibfixture2 = `<bib><suppress_from_publishing>true</suppress_from_publishing><suppress_from_external_search>true</suppress_from_external_search><record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20070515123442.0</controlfield><controlfield tag="008">20070515i19001920oru                 eng d</controlfield><datafield tag="245" ind1="0" ind2="0"><subfield code="a">Multiple archival collections in a shared box with barcode 123456789</subfield></datafield><datafield ind1="1" ind2=" " tag="774"><subfield code="t">A 044 Lawrence Hall commonplace book</subfield><subfield code="w">99198735901852</subfield></datafield><datafield tag="962" ind1=" " ind2=" "><subfield code="a">BoundwithRecord</subfield><subfield code="9">local</subfield></datafield></record></bib>`
//has lawrence and edwin
var bwbibfixture3 = `<?xml version="1.0" encoding="UTF-8"?><bib><suppress_from_publishing>true</suppress_from_publishing><suppress_from_external_search>true</suppress_from_external_search><record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20070515123442.0</controlfield><controlfield tag="008">20070515i19001920oru                 eng d</controlfield><datafield tag="245" ind1="0" ind2="0"><subfield code="a">Multiple archival collections in a shared box with barcode 123456789</subfield></datafield><datafield ind1="1" ind2=" " tag="774"><subfield code="t">A 044 Lawrence Hall commonplace book</subfield><subfield code="w">99198735901852</subfield></datafield><datafield ind1="1" ind2=" " tag="774"><subfield code="t">A 046 Edwin William Hampton correspondence</subfield><subfield code="w">99900407345101852</subfield></datafield><datafield tag="962" ind1=" " ind2=" "><subfield code="a">BoundwithRecord</subfield><subfield code="9">local</subfield></datafield></record></bib>`
//has lawrence and james
var bwbibfixture4 = `<bib><suppress_from_publishing>true</suppress_from_publishing><suppress_from_external_search>true</suppress_from_external_search><record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20070515123442.0</controlfield><controlfield tag="008">20070515i19001920oru                 eng d</controlfield><datafield tag="245" ind1="0" ind2="0"><subfield code="a">Multiple archival collections in a shared box with barcode 123456789</subfield></datafield><datafield ind1="1" ind2=" " tag="774"><subfield code="t">A 044 Lawrence Hall commonplace book</subfield><subfield code="w">99198735901852</subfield></datafield><datafield ind1="1" ind2=" " tag="774"><subfield code="t">A 045 James Watson Hamilton diary</subfield><subfield code="w">99900407345201852</subfield></datafield><datafield tag="962" ind1=" " ind2=" "><subfield code="a">BoundwithRecord</subfield><subfield code="9">local</subfield></datafield></record></bib>`

var bwholdfixture = `<holding><suppress_from_publishing>false</suppress_from_publishing><record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20070515123442.0</controlfield><controlfield tag="008">20070515i19001920oru                 eng d</controlfield><datafield tag="852" ind1="8" ind2=" "><subfield code="b">SpecColl</subfield><subfield code="c">spmanus</subfield><subfield code="h">[3502501231242352345214]</subfield></datafield></record></holding>`

var bwholdfixturea = `<?xml version="1.0" encoding="UTF-8"?><holding><holding_id>234567891</holding_id><suppress_from_publishing>false</suppress_from_publishing><record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20070515123442.0</controlfield><controlfield tag="008">20070515i19001920oru                 eng d</controlfield><datafield tag="852" ind1="8" ind2=" "><subfield code="b">SpecColl</subfield><subfield code="c">spmanus</subfield><subfield code="h">[3502501231242352345214]</subfield></datafield></record></holding>`
var bwholdfixture1 = `<holding><suppress_from_publishing>false</suppress_from_publishing><record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20070515123442.0</controlfield><controlfield tag="008">20070515i19001920oru                 eng d</controlfield><datafield tag="852" ind1="8" ind2=" "><subfield code="b">SpecColl</subfield><subfield code="c">spmanus</subfield><subfield code="h">[123456789]</subfield></datafield></record></holding>`

var bwholdfixture1a = `<holding><holding_id>234567891</holding_id><suppress_from_publishing>false</suppress_from_publishing><record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20070515123442.0</controlfield><controlfield tag="008">20070515i19001920oru                 eng d</controlfield><datafield tag="852" ind1="8" ind2=" "><subfield code="b">SpecColl</subfield><subfield code="c">spmanus</subfield><subfield code="h">[123456789]</subfield></datafield></record></holding>`

var bwholdfixture2 = `<holding><suppress_from_publishing>false</suppress_from_publishing><record><leader>00000npcaa2200000 i 4500</leader><controlfield tag="005">20070515123442.0</controlfield><controlfield tag="008">20070515i19001920oru                 eng d</controlfield><datafield tag="852" ind1="8" ind2=" "><subfield code="b">SpecColl</subfield><subfield code="c">spmanus</subfield><subfield code="h">[345678912]</subfield></datafield></record></holding>`

var bwitemfixture1 = `{"holding_data":{"holding_id":"234567891","copy_id":"1"},"item_data":{"barcode":"123456789","policy":{"value":""},"description":"Multiple Collection Box [123456789]","base_status":{"value":"1"},"library":{"value":"SpecColl"},"location":{"value":"spmanus"},"physical_material_type":{"value":"MANUSCRIPT"}}}`
var bwitemfixture1a = `{"holding_data":{"holding_id":"234567891","copy_id":"1"},"item_data":{"pid":"345678912","barcode":"123456789","policy":{"value":""},"description":"Multiple Collection Box [123456789]","base_status":{"value":"1"},"library":{"value":"SpecColl"},"location":{"value":"spmanus"},"physical_material_type":{"value":"MANUSCRIPT"},"creation_date": "2026-04-06Z","modification_date": "2026-04-06Z"}}`

var bwitemfixture2 = `{"holding_data":{"holding_id":"234567891","copy_id":"1"},"item_data":{"barcode":"345678912","policy":{"value":""},"description":"Multiple Collection Box [345678912]","base_status":{"value":"1"},"library":{"value":"SpecColl"},"location":{"value":"spmanus"},"physical_material_type":{"value":"MANUSCRIPT"}}}`
var bwitemfixture3 = `{"holding_data":{"holding_id":"234567891","copy_id":"1"},"item_data":{"barcode":"456789123","policy":{"value":""},"description":"Multiple Collection Box [456789123]","base_status":{"value":"1"},"library":{"value":"SpecColl"},"location":{"value":"spmanus"},"physical_material_type":{"value":"MANUSCRIPT"}}}`
