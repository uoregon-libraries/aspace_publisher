package as

import (
  "github.com/tidwall/gjson"
  "github.com/tidwall/sjson"
  "slices"
  "encoding/json"
  "log"
)

type TopContainer struct{
  Barcode string `json:"barcode"`
  Indicator string `json:"indicator"`
  Type string `json:"type"`
  Ils_holding string `json:"ils_holding_id"`
  Ils_item string `json:"ils_item_id"`
}

// helper to avoid including the as package in alma construct
func (t TopContainer)Mapify()map[string]string{
  tc_map := map[string]string{}
  tc_map["type"] = t.Type
  tc_map["indicator"] = t.Indicator
  tc_map["barcode"] = t.Barcode
  return tc_map
}

//returns array of top container ids (paths)
//json_list example: [{"ref":"/repositories/2/top_containers/59527"}]
func TCList(session_id, repo_id, mms_id string)([]string, error){
  json_list, err := AcquireJson(session_id, repo_id, "resources/" + mms_id + "/top_containers")
  if err != nil { return nil, err }
  pathlist := gjson.GetBytes(json_list, "#.ref")
  idlist := []string{}
  for _, tc_id := range pathlist.Array() {
    if slices.Contains(idlist, tc_id.String()) == false {
      idlist = append(idlist, tc_id.String())
    }
  }
  return idlist, nil
}

//pulls ils ids out of top container json
func GetTCRefs(record []byte)(string, string){
  var tc TopContainer
  json.Unmarshal(record, &tc)
  return tc.Ils_item, tc.Ils_holding
}

//modifies top container json to add ils ids
func UpdateIlsIds(record []byte, holding_id, item_id string)([]byte, error){
  modified1, err := sjson.SetBytes(record, "ils_holding_id", holding_id)
  if err != nil { log.Println(err); return nil, err }
  modified2, err := sjson.SetBytes(modified1, "ils_item_id", item_id)
  if err != nil { log.Println(err); return nil, err }
  return modified2, nil
}

// puts modified top container json to asapce
func UpdateTC(repo_id string, tc_id string, jsonTC []byte, holding_id string, item_id string, session_id string)error{
  path := []string{ "repositories", repo_id, "top_containers", tc_id }
  modified, err := UpdateIlsIds(jsonTC, holding_id, item_id)
  if err != nil { return err }
  _url, err := AssembleUrl(path)
  if err != nil { return err }
  _,err = Update(session_id, _url, string(modified))
  if err != nil { return err }
  return nil
}
