package as

import (
  "github.com/tidwall/gjson"
  "errors"
)

func IsPublished(record []byte)(string, error){
  result := gjson.GetBytes(record, "publish")
  if !result.Exists() { return "", errors.New("unable to determine published?") }
  return result.String(), nil
}

func GetOclcId(record []byte)(string){
  result := gjson.GetBytes(record, "user_defined.string_1")
  return result.String()
}
