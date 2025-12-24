package as

import (
  "testing"
  "os"
)

func TestUpdateUserDefined2(t *testing.T){
  resource := `{ "user_defined": { "string_1": "98765" } }` 
  result,err := UpdateUserDefined2([]byte(resource), "123456")
  if err != nil { t.Errorf("incorrect response") }
  user_defined_str_2,_ := GetMmsId(result)
  if user_defined_str_2 != "123456" { t.Errorf("incorrect result") }
}

func TestAssemblePath(t *testing.T){
  parts := []string{"abc", "def", "", "ghi"}
  result := AssemblePath(parts)
  if result != "abc/def/ghi" { t.Errorf("incorrect result") }
}

func TestAssembleUrl(t *testing.T){
  os.Setenv("ASPACE_URL", "http://example.org/")
  parts :=  []string{"abc", "def", "", "ghi"}
  result,err := AssembleUrl(parts)
  if err != nil { t.Errorf("incorrect resposne") }
  if result != "http://example.org/abc/def/ghi" { t.Errorf("incorrect result") }
}
