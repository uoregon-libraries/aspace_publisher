package as

import (
  "testing"
  "os"
  "github.com/tidwall/gjson"
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

func TestUpdateWithInstance(t *testing.T){
  ao_string := `{"jsonmodel_type":"archival_object","instances":[{"instance_type":"mixed_materials","jsonmodel_type":"instance"}]}`
  ao_modified,_ := UpdateWithInstance([]byte(ao_string), Instance("/repositories/2/digital_objects/123"))
  do_instances := gjson.Get(string(ao_modified), "instances")
  arr := do_instances.Array()
  if len(arr) != 2 { t.Fatalf("adding to instances fail") }
  inst_type := gjson.Get(arr[1].String(), "instance_type")
  if inst_type.String() != "digital_object" { t.Fatalf("adding to instances fail") }
}
