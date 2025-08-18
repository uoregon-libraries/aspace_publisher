package alma

import (
  "testing"
  "fmt"
  "net/http/httptest"
  "net/http"
  "os"
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
      fmt.Fprintf(w, data)
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
