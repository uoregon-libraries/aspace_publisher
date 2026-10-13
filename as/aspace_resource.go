package as

import (
  "github.com/tidwall/gjson"
  "errors"
  "strings"
  "regexp"
)

func CheckIsPublished(resource_id string, repo_id string, session_id string)(string, error){
  json, err := AcquireJson(session_id, repo_id, "resources/" + resource_id)
  if err != nil { return string(json), err }

  //is it published?
  return IsPublished(json)
}

func IsPublished(resource []byte)(string, error){
  result := gjson.GetBytes(resource, "publish")
  if !result.Exists() { return "", errors.New("unable to determine published?") }
  return result.String(), nil
}

func GetOclcId(resource_id string, repo_id string, session_id string)(string, error){
  json, err := AcquireJson(session_id, repo_id, "resources/" + resource_id)
  if err != nil { return string(json), err }

  return ExtractOclcId(json)
}
func ExtractOclcId(resource []byte)(string, error){
  result := gjson.GetBytes(resource, "user_defined.string_1")
  err := ValidID(result.String())
  return result.String(), err
}

// returns string ID, is-empty?, error
func ExtractMmsId(resource []byte)(string,bool,error){
  result := gjson.GetBytes(resource, "user_defined.string_2")
  err := ValidID(result.String())
  return result.String(), result.String() == "", err
  }

// returns string ID, is-empty?, error
// special handling if 404
func GetMmsId(repo_id, resource_id, session_id string)(string, bool, error){
  json, err := AcquireJson(session_id, repo_id, "resources/" + resource_id)
  if err != nil {
    if strings.Contains(string(json), "404 Not Found") == true {
      return "", true, nil } else {
      return string(json), true, err
    }
  }
  return ExtractMmsId(json)
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
