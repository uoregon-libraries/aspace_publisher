package as

import (
  "github.com/tidwall/gjson"
  "errors"
  "strings"
)

func IsPublished(resource []byte)(string, error){
  result := gjson.GetBytes(resource, "publish")
  if !result.Exists() { return "", errors.New("unable to determine published?") }
  return result.String(), nil
}

func GetOclcId(resource []byte)(string){
  result := gjson.GetBytes(resource, "user_defined.string_1")
  return result.String()
}

func GetMmsId(resource []byte)string{
  result := gjson.GetBytes(resource, "user_defined.string_2")
  return result.String()
  }

func ExtractID(_url string)string{
  parts := strings.Split(_url, "/")
  return parts[len(parts)-1]
}
