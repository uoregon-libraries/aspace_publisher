package as

import (
  "github.com/tidwall/sjson"
  "github.com/tidwall/gjson"
  "log"
  "fmt"
)

// takes AO json and inserts instance
func UpdateWithInstance(record []byte, instance string)([]byte, error){
  instance_json := gjson.Parse(instance)
  modified, err := sjson.SetBytes(record, "instances.-1", instance_json.Value())
  if err != nil { log.Println(err); return nil, err }
  return modified, nil
}

func Instance(path string) string {
  return fmt.Sprintf(`{"instance_type": "digital_object", "jsonmodel_type": "instance", "is_representative": false, "digital_object": { "ref": "%s"}`, path)
}
