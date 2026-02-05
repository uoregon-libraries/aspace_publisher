package alma

import(
  "strings"
  "os"
  "net/url"
  "time"
  "slices"
  "errors"
  "encoding/json"
  "encoding/xml"
  "aspace_publisher/file"
  "aspace_publisher/as"
  "aspace_publisher/oclc"
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

func (p ProcessArgs)Mapify()map[string]string{
  pa_map := map[string]string{}
  pa_map["repo_id"] = p.Repo_id
  pa_map["resource_id"] = p.Resource_id
  pa_map["mms_id"] = p.Mms_id
  pa_map["holding_id"] = p.Holding_id
  pa_map["session_id"] = p.Session_id
  pa_map["oclc_token"] = p.Oclc_token
  pa_map["oclc_id"] = p.Oclc_id
  pa_map["id_0"] = p.Id_0
  pa_map["filename"] = p.Filename
  return pa_map
}

type FunMap struct {
  BoundwithPF ProcessBoundwithFun
  HoldingPF ProcessHoldingFun
  ItemsPF ProcessItemsFun
  ItemPF ProcessItemFun
  NZPF LinkToNetworkFun
  FetchBib FetchBibIDFun
  AfterBib as.AfterBibFun
  UpdateTC as.UpdateTCFun
  SetHolding oclc.SetHoldingFun
}

func ProcessBib(args ProcessArgs, marc_string string, rjson []byte, tcmap []map[string]string, fs FunMap){
  // assemble record
  bib := ConstructBib(marc_string, false)
  bib_str, err := bib.Stringify()
  if err != nil { file.WriteReport(args.Filename, []string{ "Unable to construct bib: " + err.Error() }); return }
  path := []string{"bibs", args.Mms_id}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var result []byte
  // push to alma
  if args.Create {
    result, err = Post(_url, params, bib_str, "xml") } else {
    result, err = Put(_url, params, bib_str, "xml")
  }
  if err != nil { file.WriteReport(args.Filename, []string{"Unable to publish bib" + err.Error()}); return }
  if args.Create {
    args.Mms_id = ExtractBibID(result)
    //update the aspace resource
    err = fs.AfterBib(rjson, args.Mapify())
    if err != nil { file.WriteReport(args.Filename, []string{ err.Error() }) }
    //todo: switch to worker.
    fs.NZPF([]string{ args.Mms_id }, args.Filename)
    res,err := fs.SetHolding(args.Oclc_id, args.Oclc_token)
    if err != nil {
      file.WriteReport(args.Filename, []string{ err.Error() }) } else {
      file.WriteReport(args.Filename, []string{ res }) 
    }
  }
  fs.BoundwithPF(args, marc_string, tcmap, fs)
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

type ProcessBoundwithFun func(ProcessArgs, string, []map[string]string, FunMap)
// if a given top container is not a multi-collection box, 
// the only action sets the tc['mms_id'] to the bib being created/updated
// if boundwith true and error occurs, write report and stop once loop is complete
func ProcessBoundwith(args ProcessArgs,marc_string string, tcmap []map[string]string, fs FunMap){
  var process_holding = false
  msgs := []string{}
  for _,tc := range tcmap{
    if tc["boundwith"] == "true" {
      mms_id := fs.FetchBib(tc["barcode"])//get boundwith mms_id
      // currently only handling bwbib exists case
      tc["mms_id"] = mms_id
      path := []string{"bibs", mms_id}
      _url := BuildUrl(path)
      params := []string{ ApiKey() }
      bwbib_byte, err := Get(_url, params, "application/xml")
      if err != nil { msgs = append(msgs, err.Error()); continue }
      bwbib, err := ConstructBoundwith(bwbib_byte, marc_string, args.Mms_id, tc)
      if err != nil { msgs = append(msgs, err.Error()); continue }
      bwbib_str, err := bwbib.Stringify()
      if err != nil { msgs = append(msgs, err.Error()); continue }
      _, err = Put(_url, params, bwbib_str, "xml")
      if err != nil { msgs = append(msgs, err.Error()); continue }
    } else { tc["mms_id"] = args.Mms_id; args.Holding_id = tc["ils_holding"]; process_holding = true }
  }
  if len(msgs) != 0 { file.WriteReport(args.Filename, msgs); return }
  if process_holding { fs.HoldingPF(args, marc_string, tcmap, fs) } else {
    fs.ItemsPF(args, tcmap, fs)
  }
}

type ProcessHoldingFun func(ProcessArgs, string, []map[string]string, FunMap)
// does not need tcmap, passes it to items processing which does
func ProcessHolding(args ProcessArgs, marc_string string, tcmap []map[string]string, fs FunMap){
  //assemble holding record
    path := []string{"bibs", args.Mms_id, "holdings", args.Holding_id}
    _url := BuildUrl(path)
    params := []string{ ApiKey() }

  var holding = Holding{}
  if args.Holding_id != "" {
    holdxml, err := Get(_url, params, "application/xml")
    if err != nil { file.WriteReport(args.Filename, []string{"Unable to obstain current holding: " + err.Error()}); return }
    xml.Unmarshal(holdxml, &holding)
  }
  holding, err := ConstructHolding(marc_string, holding, args.Id_0)
  if err != nil { file.WriteReport(args.Filename, []string{"Unable to construct holding: " + err.Error()}); return }
  holdingstr, err := holding.Stringify()
  if err != nil {}
  var result []byte
  // push record to alma
  if args.Create {
    result, err = Post(_url, params, holdingstr, "xml") } else {
    result, err = Put(_url, params, holdingstr, "xml")
  }
  if err != nil { file.WriteReport(args.Filename, []string{"Unable to push to alma: " + err.Error()}); return }
  if args.Create {
    args.Holding_id = ExtractHoldingID(result)
  }
  fs.ItemsPF(args, tcmap, fs)
}

type ProcessItemsFun func(ProcessArgs, []map[string]string, FunMap)
func ProcessItems(args ProcessArgs, tcmap []map[string]string, fs FunMap){
  itemlist := []string{}
  msgs := []string{}

  // iterate through the top containers
  // if an error occurs during the loop, report and continue
  for _,tc := range tcmap{
    var item = Item{}
    if tc["ils_item"] != "" { //this is an update. mms_id is hopefully set in ProcessBoundwith
      path := []string{"bibs", tc["mms_id"], "holdings", tc["ils_holding"], "items", tc["ils_item"]}
      _url := BuildUrl(path)
      params := []string{ ApiKey() }
      itemjson, err := Get(_url, params, "application/json")
      if err != nil { msgs = append(msgs, "Unable to request Alma item: " + err.Error()); continue }
      json.Unmarshal(itemjson, &item)
    }
    item_id, err := fs.ItemPF(args, item, tc)
    if err != nil { msgs = append(msgs, "Unable to process Alma item: " + err.Error()); continue }
    itemlist = append(itemlist, item_id)
    if tc["ils_item"] == "" {
      err = fs.UpdateTC(args.Repo_id, args.Holding_id, item_id, args.Session_id, tc)
      if err != nil { msgs = append(msgs, "Unable to update TC in aspace: " + err.Error()); continue }
    }
  }
  msgs = append(msgs, "items created: " + strings.Join(itemlist, ", "))
  file.WriteReport(args.Filename, msgs)
}

type ProcessItemFun func(ProcessArgs, Item, map[string]string)(string, error)
// does not log or write reports
func ProcessItem(args ProcessArgs, item Item, tcmap map[string]string)(string, error){
  //assemble item record
  item, err := ConstructItem(args.Holding_id, item, tcmap)
  if err != nil { return "", errors.New("Unable to construct item" + err.Error()) }
  itemstr, err := item.Stringify()
  if err != nil { return "", errors.New("Unable to construct item" + err.Error()) }
  path := []string{"bibs", args.Mms_id, "holdings", args.Holding_id, "items", tcmap["ils_item"]}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var result []byte
  //push record to alma
  if tcmap["ils_item"] != "" {
    result, err = Post(_url, params, itemstr, "json") } else {
    result, err = Put(_url, params, itemstr, "json")
  }
  if err != nil { return "", errors.New("problem posting to alma: " + err.Error()) }
  var item_id string
  if args.Create {
    item_id = ExtractItemID(result) } else {
    item_id = tcmap["Ils_item_id"]
  }
  return item_id, err
}

func ApiKey()string{
  key := os.Getenv("ALMA_KEY")
  return "apikey=" + key
}

type LinkToNetworkFun func([]string, string)
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
