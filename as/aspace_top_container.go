package as

import (
  "github.com/tidwall/gjson"
  "github.com/tidwall/sjson"
  "slices"
  "encoding/json"
  "log"
  "strconv"
)

type TopContainer struct{
  Uri string `json:"uri"`
  Barcode string `json:"barcode"`
  Indicator string `json:"indicator"`
  Type string `json:"type"`
  Ils_holding string `json:"ils_holding_id"`
  Ils_item string `json:"ils_item_id"`
  Boundwith bool
}

// helper to avoid including the as package in alma construct
func (t TopContainer)Mapify()map[string]string{
  tc_map := map[string]string{}
  tc_map["type"] = t.Type
  tc_map["indicator"] = t.Indicator
  tc_map["barcode"] = t.Barcode
  tc_map["boundwith"] = strconv.FormatBool(t.Boundwith)
  tc_map["uri"] = t.Uri
  tc_map["ils_holding"] = t.Ils_holding
  tc_map["ils_item"] = t.Ils_item
  return tc_map
}

//returns array of top container ids (paths)
//json_list example: [{"ref":"/repositories/2/top_containers/59527"}]
func TCList(session_id, repo_id, id string)([]string, error){
  json_list, err := AcquireJson(session_id, repo_id, "resources/" + id + "/top_containers")
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

// puts modified top container json to aspace
func UpdateTC(repo_id string, holding_id string, item_id string, session_id string, tcmap map[string]string)error{
  tc_id := ExtractID(tcmap["uri"])
  jsonTC, err := AcquireJson(session_id, repo_id, "top_containers/" + tc_id)
  if err != nil { return err }
  modified, err := UpdateIlsIds(jsonTC, holding_id, item_id)
  if err != nil { return err }
  path := []string{ "repositories", repo_id, "top_containers", tc_id }
  _url, err := AssembleUrl(path)
  if err != nil { return err }
  _,err = Update(session_id, _url, string(modified))
  if err != nil { return err }
  return nil
}

func ExtractTCData(session string, repo_id string, resource_id string)([]map[string]string, []string){
  msgs := []string{}
  tclist,err := TCList(session, repo_id, resource_id)
  if err != nil { msgs = append(msgs, "Unable to acquire TC list: " + err.Error()); return nil, msgs}
  var top_containers []map[string]string
  for _,tc_path := range tclist{
    tc_id := ExtractID(tc_path)
    jsonTC, err := AcquireJson(session, repo_id, "top_containers/" + tc_id)
    if err != nil { msgs = append(msgs, "Unable to acquire TC json " + err.Error()); continue }

    var tc TopContainer
    err = json.Unmarshal(jsonTC, &tc)
    if err != nil { msgs = append(msgs, "Unable to process TC json: " + err.Error()); continue }
    tc.Boundwith = IsBoundwith(jsonTC)
    top_containers = append(top_containers, tc.Mapify())
  }
  return top_containers, msgs
}

func IsBoundwith(jsontc []byte)bool{
  result := gjson.GetBytes(jsontc, "collection")
  if len(result.Array()) > 1 { return true }
  return false
}
