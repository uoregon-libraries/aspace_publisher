package as

import (
  "github.com/tidwall/gjson"
  "errors"
  "strings"
  "regexp"
  "fmt"
  "net/url"
  "encoding/json"

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

// for the initial launch confirmation step
type QueryResponse struct {
  Results []ShortResource `json:"results"`
}

type ShortResource struct{
  Id string `json:"id"`
  Title string `json:"title"`
  Identifier string `json:"identifier"`
}

func (sr ShortResource) Mapify() map[string]string{
  return map[string]string{ "id": sr.Id, "title": sr.Title, "identifier": sr.Identifier }
}

func GetResourceShort(resource_id, repo_id, session_id string)(map[string]string, error){
  query := fmt.Sprintf(`id:"/repositories/%s/resources/%s"`, repo_id, resource_id)
  escaped := url.QueryEscape(query)
  search_string := fmt.Sprintf(`search?type[]=resource&q=%s&fields[]=id,identifier,title&page=1&page_size=1`, escaped)
  json_list, err := AcquireJson(session_id, repo_id, search_string)
  var srmap map[string]string
  if err != nil { return srmap, err }
  var qr QueryResponse
  err = json.Unmarshal(json_list, &qr)
  if err != nil { return srmap, err }
  qr.Results[0].Id = ExtractID(qr.Results[0].Id) //just want the end bits
  if qr.Results[0].Id != resource_id { return srmap, errors.New("Incorrect response for resource lookup") }
  return qr.Results[0].Mapify(), nil
}