package as

import (
  "testing"
  "os"
  "io/ioutil"
  "log"
)

func loadfile(filename string)[]byte{
  home := os.Getenv("HOME_DIR")
  f, err := ioutil.ReadFile(home + "fixtures/" + filename)
  if err != nil { log.Println("error reading file") }
  return f
}

func TestIsPublished(t *testing.T){
  resource := loadfile("9730.json")
  is_pub,_ := IsPublished(resource)
  if is_pub != "true" { t.Errorf("incorrect result") }
}

func TestGetOclcId(t *testing.T){
  resource := loadfile("9730.json")
  oclc_id := GetOclcId(resource)
  if oclc_id != "1535209600" { t.Errorf("incorrect result") }
}

func TestGetMmsId(t *testing.T){
  resource := loadfile("9730.json")
  id,is_empty := GetMmsId(resource)
  if  id != "" { t.Errorf("incorrect result")}
  if is_empty != true { t.Errorf("incorrect result") }
  resource = loadfile("2023.json")
  id, is_empty = GetMmsId(resource)
  if id != "99107164901852" { t.Errorf("incorrect result") }
  if is_empty != false { t.Errorf("incorrect result") }
}

func TestExtractID(t *testing.T){
  url := "http://example.org/blah/blah/123456"
  result := ExtractID(url)
  if result != "123456" { t.Errorf("incorrect result") }
}

func TestExtractID0(t *testing.T){
  resource := loadfile("2023.json")
  result := ExtractID0(resource)
  if result != "Coll 100" { t.Errorf("incorrect result") }
}
