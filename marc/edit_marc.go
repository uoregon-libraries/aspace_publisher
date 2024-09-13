package marc

import (
  "errors"
  "log"
  "github.com/beevik/etree"
)

func EditMarcForOCLC(oclc_marc string, as_marc string)(string,error){
  oclc := etree.NewDocument()
  err := oclc.ReadFromString(oclc_marc)
  if err != nil { log.Println(err); return "", errors.New("unable to read oclc marc") }

  as := etree.NewDocument()
  err = as.ReadFromString(as_marc)
  if err != nil { log.Println(err); return "", errors.New("unable to read as marc") }

  //fix 040
  as_040 := as.FindElement("//datafield[@tag='040']")
  d := oclc.FindElements("//datafield[@tag='040']/subfield[@code='d']")
  for _, sub := range d {
    as_040.AddChild(sub.Copy())
  }

  //fix 000
  as_000 := as.FindElement("//leader")
  as_l := as_000.Text()
  oclc_l := oclc.FindElement("//leader").Text()
  new_l := edit_status(as_l, oclc_l)
  as_000.SetText(new_l)

  //fix 008
  as_008 := as.FindElement("//controlfield[@tag='008']")
  as_008_t := as_008.Text()
  oclc_008_t := oclc.FindElement("//controlfield[@tag='008']").Text()
  new_008 := edit_008(as_008_t, oclc_008_t)
  as_008.SetText(new_008)

  s, err := as.WriteToString()
  if err != nil { log.Println(err); return "", errors.New("unable to write marc to string") }
  return s, nil
}

func edit_status(as_l string, oclc_l string) string{
  as_r := []rune(as_l)
  oclc_r := []rune(oclc_l)
  as_r[5] = oclc_r[5]
  return string(as_r)
}

func edit_008(as_008_t string, oclc_008_t string) string {
  as_r := []rune(as_008_t)
  oclc_r := []rune(oclc_008_t)
  new_008 := append(oclc_r[0:6], as_r[6:len(as_r)]...)
  return string(new_008)
}
