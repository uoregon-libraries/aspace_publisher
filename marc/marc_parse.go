package marc

import (
  "github.com/beevik/etree"
  "log"
  "errors"
)

func ExtractOclc(marc string)(string, error){
  marc_xml := etree.NewDocument()
  err := marc_xml.ReadFromString(marc)
  if err != nil { log.Println(err); return "", errors.New("Unable to read XML response from OCLC.") }
  oclc := marc_xml.FindElement("//controlfield[@tag='001']").Text()
  return oclc, nil
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
