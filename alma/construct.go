package alma

import (
  "aspace_publisher/marc"
  "log"
)

func ConstructBib(marc_string string)string{
  marc_stripped, err := marc.StripOuterTags(marc_string)
  if err != nil { log.Println(err); return "" }
  return "<bib>" + marc_stripped + "</bib>"
}
