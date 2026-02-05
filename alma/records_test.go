package alma

import (
  "testing"
  "fmt"
  "net/http/httptest"
  "net/http"
  "os"
  "strings"
)

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

func TestFetchBibID(t *testing.T){
  data := `{ "bib_data": {"mms_id":"123456789123"} }`
  barcode := "123123123123123"
  path := "/almaws/v1/items"
   ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if strings.Contains(r.URL.String(), "item_barcode=" + barcode) != true { t.Errorf("incorrect params") }
    if r.URL.Path == path {
      fmt.Fprint(w, data)
    } else {
      t.Errorf("incorrect request url")
    }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "abcdeabcdeabcde")
  id := FetchBibID(barcode)
  if id != "123456789123" { t.Errorf("incorrect mms id") }

}

func TestStringify(t *testing.T){
  fstring := bibstring_fixture4
  expected := bibstring_fixture5
  bib := ConstructBib(fstring, false)
  bib_string,err := bib.Stringify()
  if err != nil { t.Errorf("incorrect response") }
  if compareXML(bib_string, expected) != true { t.Errorf("incorrect rec") }
}
