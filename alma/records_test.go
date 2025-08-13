package alma

import (
  "testing"
)

func TestExtractBibID( t *testing.T){
  fstring := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><bib><mms_id>123456789111</mms_id></bib>"
  mms_id := ExtractBibID([]byte(fstring))
  if mms_id != "123456789111" { t.Errorf("incorrect mms_id") }
}
