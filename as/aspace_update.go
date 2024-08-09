package as

import (
  "github.com/tidwall/sjson"
  "log"
)

func UpdateUserDefined1(record []byte, oclc string)([]byte, error){
  modified, err := sjson.SetBytes(record, "user_defined.string_1", oclc)
  if err != nil { log.Println(err); return nil, err }
  return modified, nil
}
