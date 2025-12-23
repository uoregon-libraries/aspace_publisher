package alma

import(
  "strings"
  "os"
  "net/url"
  "time"
  "slices"
  "errors"
  "encoding/json"
  "aspace_publisher/file"
  "aspace_publisher/as"
)

type ProcessArgs struct {
  Repo_id string
  Resource_id string
  Mms_id string
  Holding_id string
  Session_id string
  Oclc_token string
  Oclc_id string
  Create bool
  Id_0 string
  Filename string
}

func ProcessBib(args ProcessArgs, marc_string string, rjson []byte){
  // assemble record
  bib, err := ConstructBib(marc_string)
  if err != nil { file.WriteReport(args.Filename, []string{ "Unable to construct bib: " + err.Error() }); return }
  path := []string{"bibs", args.Mms_id}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var result []byte
  // push to alma
  if args.Create {
    result, err = Post(_url, params, bib, "xml") } else {
    result, err = Put(_url, params, bib, "xml")
  }
  if err != nil { file.WriteReport(args.Filename, []string{"Unable to publish bib" + err.Error()}); return }
  if args.Create {
    args.Mms_id = ExtractBibID(result)
    //update the aspace resource
    modified, err := as.UpdateUserDefined2(rjson, args.Mms_id)
    if err != nil { file.WriteReport(args.Filename, []string{ err.Error() }); return }
    as.UpdateResource(args.Session_id, args.Repo_id, args.Resource_id, string(modified))
    //todo: switch to worker.
    LinkToNetwork([]string{ args.Mms_id }, args.Filename)
  }
  ProcessHolding(args, marc_string)
}

func BuildUrl(path []string)string{
  _url,_ := url.Parse(os.Getenv("ALMA_URL"))
  path = slices.DeleteFunc(path, func(str string) bool {
    return str == ""
  })
  string_path := strings.Join(path, "/")
  _url = _url.JoinPath(string_path)
  return _url.String()
}

func ProcessHolding(args ProcessArgs, marc_string string){
  //assemble holding record
  holding, err := ConstructHolding(marc_string, args.Id_0)
  if err != nil { file.WriteReport(args.Filename, []string{"Unable to construct holding: " + err.Error()}); return }
  path := []string{"bibs", args.Mms_id, "holdings", args.Holding_id}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var result []byte
  // push record to alma
  if args.Create {
    result, err = Post(_url, params, holding, "xml") } else {
    result, err = Put(_url, params, holding, "xml")
  }
  if err != nil { file.WriteReport(args.Filename, []string{"Unable to push to alma: " + err.Error()}); return }
  if args.Create {
    args.Holding_id = ExtractHoldingID(result)
    // call oclc holding set here
  }
  ProcessItems(args)
}

func ProcessItems(args ProcessArgs){
  itemlist := []string{}
  msgs := []string{}
  tclist,err := as.TCList(args.Session_id, args.Repo_id, args.Resource_id)
  if err != nil { file.WriteReport(args.Filename, []string{ "Unable to acquire TC list: " + err.Error() }); return }

  // iterate through the top containers
  // if an error occurs during the loop, report and continue
  for _,tc_path := range tclist{
    tc_id := as.ExtractID(tc_path)
    jsonTC, err := as.AcquireJson(args.Session_id, args.Repo_id, "top_containers/" + tc_id)
    if err != nil { msgs = append(msgs, "Unable to acquire TC json " + err.Error()); continue }
    item_id, _ := as.GetTCRefs(jsonTC)
    var tc as.TopContainer
    err = json.Unmarshal(jsonTC, &tc)
    if err != nil { msgs = append(msgs, "Unable to process TC json: " + err.Error()); continue }
    item_id, err = ProcessItem(args, item_id, tc.Mapify())
    if err != nil { msgs = append(msgs, "Unable to process Alma item: " + err.Error()); continue }
    itemlist = append(itemlist, item_id)
    if args.Create {
      err = as.UpdateTC(args.Repo_id, tc_id, jsonTC, args.Holding_id, item_id, args.Session_id)
      if err != nil { msgs = append(msgs, "Unable to update TC in aspace: " + err.Error()); continue }
    }
  }
  msgs = append(msgs, "items created: " + strings.Join(itemlist, ", "))
  file.WriteReport(args.Filename, msgs)
}

// does not log or write reports
func ProcessItem(args ProcessArgs, item_id string, container_data map[string]string)(string, error){
  //assemble item record
  item, err := ConstructItem(item_id, args.Holding_id, container_data)
  if err != nil { return "", errors.New("Unable to construct item" + err.Error()) }
  path := []string{"bibs", args.Mms_id, "holdings", args.Holding_id, "items", item_id}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var result []byte
  //push record to alma
  if args.Create {
    result, err = Post(_url, params, item, "json") } else {
    result, err = Put(_url, params, item, "json")
  }
  if err != nil { return "", errors.New("problem posting to alma: " + err.Error()) }
  if args.Create {
    item_id = ExtractItemID(result)
  }
  return item_id, err
}

func ApiKey()string{
  key := os.Getenv("ALMA_KEY")
  return "apikey=" + key
}

func LinkToNetwork(list []string, filename string){
  setid := os.Getenv("LINK_TO_NETWORK_SET")
  jobid := os.Getenv("LINK_TO_NETWORK_JOB")
  err := UpdateSet("LINK_TO_NETWORK_SET", "BIB_MMS", list)
  if err != nil { file.WriteReport(filename, []string{ "problem updating alma set: " + err.Error()}); return }
  var params = []Param{
    Param{ Name: Val{ Value: "set_id" }, Value: setid },
    Param{ Name: Val{ Value: "contribute_nz" }, Value: "true" },
    Param{ Name: Val{ Value: "non_serial_match_profile" }, Value: "com.exlibris.repository.mms.match.uniqueOCLC" },
    Param{ Name: Val{ Value: "non_serial_match_prefix" }, Value: "" },
    Param{ Name: Val{ Value: "serial_match_profile" }, Value: "com.exlibris.repository.mms.match.uniqueOCLC" },
    Param{ Name: Val{ Value: "serial_match_prefix" }, Value: "" },
    Param{ Name: Val{ Value: "ignoreResourceType" }, Value: "false" },
  }
  instance,err := SubmitJob(jobid, params)
  if err != nil { file.WriteReport(filename, []string{ "problem submitting alma job: " + err.Error() } ); return }
  span,_ := time.ParseDuration(os.Getenv("JOB_WAIT_TIME"))
  time.Sleep(span)
  CheckJob(instance, nil, filename, nil)
}

func BaseUrl()string{
  return os.Getenv("ALMA_URL")
}
