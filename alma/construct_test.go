package alma

import (
  "testing"
)

func TestConstructBib( t *testing.T){
  fstring := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><record><blah>banana</blah></record>"
  bib_string := ConstructBib(fstring)
  if bib_string != "<bib><record><blah>banana</blah></record></bib>" {
    t.Errorf("incorrect bib record")
  }
}
