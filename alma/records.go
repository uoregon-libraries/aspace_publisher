package alma

import (
  "encoding/xml"
  "log"
)

func ExtractBibID(data []byte)string{
  var b Bib
  err := xml.Unmarshal(data, &b)
  if err != nil { log.Println(err); return "" }
  return b.Mms_id
}

// GET returns json
// POST requires XML in and out
type Bib struct {
  XMLName xml.Name 	`xml:"bib"`	
  Mms_id string 	`xml:"mms_id"`
}

type Value struct {
  Val string `json:"value"`
}
