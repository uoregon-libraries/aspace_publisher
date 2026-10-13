package as

import (
  "testing"
  "os"
  "net/http"
  "net/http/httptest"
  "fmt"
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

func TestExtractOclcId(t *testing.T){
  resource := loadfile("9730.json")
  oclc_id,_ := ExtractOclcId(resource)
  if oclc_id != "1535209600" { t.Errorf("incorrect result") }
}

func TestExtractMmsId(t *testing.T){
  resource := loadfile("9730.json")
  id,is_empty,_ := ExtractMmsId(resource)
  if  id != "" { t.Errorf("incorrect result")}
  if is_empty != true { t.Errorf("incorrect result") }
  resource = loadfile("2023.json")
  id, is_empty,_ = ExtractMmsId(resource)
  if id != "99107164901852" { t.Errorf("incorrect result") }
  if is_empty != false { t.Errorf("incorrect result") }
}
func TestGetMmsId(t *testing.T){

  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    path1 := "/api/repositories/2/resources/123"
    if r.URL.Path == path1 {
      fmt.Fprint(w, `{ "user_defined": { "string_2": "987654321" } }`)
    } else {
      fmt.Fprint(w, `<center><h1>404 Not Found</h1></center>`)
    }
  }))
  defer ts.Close()
  os.Setenv("ASPACE_URL", ts.URL + "/api/")
  id, _, _ := GetMmsId("2", "123", "2345678901")
  if id != "987654321" { t.Error("incorrect response") }
  _, is_empty, _ := GetMmsId("2", "234", "2345678901")
  if is_empty != true { t.Error("incorrect response to 404") }

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

func TestValidID(t *testing.T){
  if ValidID(" 12345") == nil { t.Errorf("incorrect validation") }
  if ValidID("on12345") == nil { t.Errorf("incorrect validation") }
  if ValidID("12345") != nil { t.Errorf("incorrect validation") }
}
