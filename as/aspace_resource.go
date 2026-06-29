package as

import (
  "github.com/tidwall/gjson"
  "errors"
  "strings"
  "regexp"
)

func IsPublished(resource []byte)(string, error){
  result := gjson.GetBytes(resource, "publish")
  if !result.Exists() { return "", errors.New("unable to determine published?") }
  return result.String(), nil
}

func GetOclcId(resource []byte)(string, error){
  result := gjson.GetBytes(resource, "user_defined.string_1")
  err := ValidID(result.String())
  return result.String(), err
}

func GetMmsId(resource []byte)(string,bool,error){
  result := gjson.GetBytes(resource, "user_defined.string_2")
  err := ValidID(result.String())
  return result.String(), result.String() == "", err
  }

func ExtractID(_url string)string{
  parts := strings.Split(_url, "/")
  return parts[len(parts)-1]
}

func ExtractID0(resource []byte) string {
  result := gjson.GetBytes(resource, "id_0")
  return result.String()
}

// for use with mms_id, oclc
func ValidID(id string) error {
  re1 := regexp.MustCompile(`[0-9]+`)
  matched1 := re1.Find([]byte(id))
  if string(matched1) == id { return nil }
  return errors.New("not a valid ID")
}
