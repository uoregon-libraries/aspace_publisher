package marc

import (
  "github.com/beevik/etree"
  "log"
  "errors"
  "strings"
)

func ExtractOclc(marc string)(string, error){
  marc_xml := etree.NewDocument()
  err := marc_xml.ReadFromString(marc)
  if err != nil { log.Println(err); return "", errors.New("Unable to read XML response from OCLC.") }
  oclc_val := marc_xml.FindElement("//controlfield[@tag='001']").Text()
  oclc := cleanOclc(oclc_val)
  return oclc, nil
}

func cleanOclc(oclc_id string)string {
  return strings.TrimFunc(oclc_id, func(r rune) bool{
    return !unicode.IsNumber(r)
  })
}

func StripOuterTags(marc string)(string, error){
  marc_xml := etree.NewDocument()
  err := marc_xml.ReadFromString(marc)
  if err != nil { log.Println(err); return "", errors.New("Unable to read XML response from archivesspace.") }
  record := marc_xml.FindElement("//record")
  stripped := etree.NewDocumentWithRoot(record)
  s, err := stripped.WriteToString()
  if err != nil { log.Println(err); return "", errors.New("Unable to process marc xml") }
  return s, nil
}
