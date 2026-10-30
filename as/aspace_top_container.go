package as

import (
  "github.com/tidwall/gjson"
  "github.com/tidwall/sjson"
  "slices"
  "encoding/json"
  "log/slog"
  "strconv"
  "regexp"
  "fmt"
  "errors"
)

type TopContainer struct{
  Uri string `json:"uri"`
  Barcode string `json:"barcode"`
  Indicator string `json:"indicator"`
  Type string `json:"type"`
  Ils_holding string `json:"ils_holding_id"`
  Ils_item string `json:"ils_item_id"`
  Boundwith bool
  UpdateRefs bool
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
  tc_map["update_refs"] = strconv.FormatBool(t.UpdateRefs)
  return tc_map
}

//returns array of top container ids (paths)
//json_list example: [{"ref":"/repositories/2/top_containers/59527"}]
func TCListPublished(session_id, repo_id, id string)([]string, error){
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

func TCList(session_id, repo_id, id string)([]string, error){
  json_list, err := AcquireJson(session_id, repo_id, "top_containers/search?type[]=resource&q=/resources/" + id)
  if err != nil { return nil, err }
  r_uri := fmt.Sprintf("/repositories/%s/resources/%s", repo_id, id)
  tclist := gjson.GetBytes(json_list, fmt.Sprintf("response.docs.#(collection_uri_u_sstr.#(== %s))#.id", r_uri))
  cleanlist := []string{}
  // ensure list content unique
  for _, tc := range tclist.Array() {
    if slices.Contains(cleanlist, tc.String()) == false {
      cleanlist = append(cleanlist, tc.String())
    }
  }
  return cleanlist, nil
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
  if err != nil { slog.Error(err.Error()); return nil, err }
  modified2, err := sjson.SetBytes(modified1, "ils_item_id", item_id)
  if err != nil { slog.Error(err.Error()); return nil, err }
  return modified2, nil
}

type UpdateTCFun func(string, string, string, string, string)(string, error)
// adds ils ids to top container
// requires the ids be passed in
func UpdateTC(repo_id, holding_id, item_id, session_id, tc_id string)(string, error){
  jsonTC, err := AcquireJson(session_id, repo_id, "top_containers/" + tc_id)
  if err != nil { return "", err }
  modified, err := UpdateIlsIds(jsonTC, holding_id, item_id)
  if err != nil { return "", err }
  path := []string{ "repositories", repo_id, "top_containers", tc_id }
  _url, err := AssembleUrl(path)
  if err != nil { return "", err }
  result,err := Update(session_id, _url, string(modified))
  return result, err
}

// gets tcmaps for all tcs associated with a given resource
func ExtractTCData(session, repo_id, resource_id string)([]map[string]string, []string){
  msgs := []string{}
  tclist,err := TCList(session, repo_id, resource_id)
  if err != nil { msgs = append(msgs, "Unable to acquire TC list: " + err.Error()); return nil, msgs }
  var top_containers []map[string]string
  for _,tc_path := range tclist{
    tc_id := ExtractID(tc_path)
    tc,err := ExtractSingleTC(session, repo_id, tc_id)
    if err != nil { msgs = append(msgs, fmt.Sprintf("Unable to acquire %s: %s",tc_id, err.Error())) }
    top_containers = append(top_containers, tc.Mapify())
  }
  return top_containers, msgs
}

// get tc for a given tc_id
func ExtractSingleTC(session, repo_id, tc_id string)(TopContainer, error){
  var tc TopContainer
  jsonTC, err := AcquireJson(session, repo_id, "top_containers/" + tc_id)
  if err != nil { return tc, err }
  err = json.Unmarshal(jsonTC, &tc)
  if err != nil { return tc, err }
  tc.Boundwith = IsBoundwith(jsonTC)
  if ValidateTCData(tc) != true { return tc, errors.New("TopContainer contains invalid value: " + tc_id) }
  return tc, nil
}

func IsBoundwith(jsontc []byte)bool{
  result := gjson.GetBytes(jsontc, "collection")
  if len(result.Array()) > 1 { return true }
  return false
}

// for now, the only field that gets validation is the barcode
func ValidateTCData(tc TopContainer) bool{
  if !validateBarcode(tc.Barcode) { return false }
  //more validations here
  return true
}

func validateBarcode(barcode string) bool{
  re1 := regexp.MustCompile(`[A-Za-z0-9]+`)
  matched1 := re1.Find([]byte(barcode))
  return string(matched1) == barcode
}
