package alma

import(
  "strings"
  "os"
  "net/url"
  "time"
  "slices"
  "errors"
  "fmt"
  "encoding/json"
  "aspace_publisher/file"
  "aspace_publisher/as"
  "aspace_publisher/oclc"
  "github.com/tidwall/gjson"
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
  HoldingAPF ProcessHoldingAFun
  HoldingBPF ProcessHoldingBFun
  ItemsPF ProcessItemsFun
  ItemPF ProcessItemFun
  NZPF LinkToNetworkFun
  AfterBib as.AfterBibFun
  UpdateTC as.UpdateTCFun
  SetHolding oclc.SetHoldingFun
  CheckForMissingPF CheckItemsForMissingFun
}

// only run creates
// updates are handled magically via oclc
func ProcessBib(args ProcessArgs, marc_string string, rjson []byte, tcmap []map[string]string, fs FunMap){
  // assemble record
  if args.Create {
    bib := ConstructBib(args.Mms_id, marc_string, "false")
    bib_str, err := bib.Stringify()
    if err != nil { file.WriteReport(args.Filename, []string{ "Unable to construct bib: " + err.Error() }); return }
    path := []string{"bibs", args.Mms_id}
    _url := BuildUrl(path)
    params := []string{ ApiKey() }

    // push to alma
    result, err := Post(_url, params, bib_str, "xml")
    if err != nil { file.WriteReport(args.Filename, []string{"Unable to publish bib" + err.Error()}); return }
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

// only handles updates
type ProcessBoundwithFun func(ProcessArgs, string, []map[string]string, FunMap)
// if not boundwith, set the args.Holding_id to tc["ils_holding"]
// will only call ProcessHolding if finds non-boundwith tc
// if boundwith true and error occurs, write report and stop once loop is complete
func ProcessBoundwith(args ProcessArgs,marc_string string, tcmap []map[string]string, fs FunMap){
  var process_holding = false
  msgs := []string{}
  for _,tc := range tcmap{
    if tc["boundwith"] == "true" {
      path := []string{"bibs", tc["mms_id"] }
      _url := BuildUrl(path)
      params := []string{ ApiKey() }
      bwbib_byte, err := Get(_url, params, "application/xml")
      if err != nil { msgs = append(msgs, err.Error()); continue }
      bwbib, err := UpdateBoundwith(bwbib_byte, marc_string, args.Mms_id, tc)
      if err != nil { msgs = append(msgs, err.Error()); continue }
      bwbib_str, err := bwbib.Stringify()
      if err != nil { msgs = append(msgs, err.Error()); continue }
      _, err = Put(_url, params, bwbib_str, "xml")
      if err != nil { msgs = append(msgs, err.Error()); continue }
    } else {
      if args.Holding_id == "" { args.Holding_id = tc["ils_holding"] }
      process_holding = true
    }
  }
  //check for potential alma items that no longer exist in aspace
  //do this here bc we need the holding id which has just been set
  if !args.Create { fs.CheckForMissingPF(args, tcmap) }

  if !process_holding { msgs = append(msgs, "no items to update") }
  if len(msgs) != 0 { file.WriteReport(args.Filename, msgs) }
  if process_holding { fs.HoldingAPF(args, marc_string, tcmap, fs) }
}

type ProcessHoldingAFun func(ProcessArgs, string, []map[string]string, FunMap)
// does not use tcmap, passes it to items processing which does
// "A" method handles creates
func ProcessHoldingA(args ProcessArgs, marc_string string, tcmap []map[string]string, fs FunMap){
  if args.Holding_id != "" {
    fs.HoldingBPF(args, marc_string, tcmap, fs)
    return
  }
  path := []string{"bibs", args.Mms_id, "holdings", args.Holding_id}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var result []byte
  holding, err := ConstructHolding(marc_string, args.Id_0)
  if err != nil { file.WriteReport(args.Filename, []string{"Unable to construct holding: " + err.Error()}); return }
  holdingstr, err := holding.Stringify()
  if err != nil { file.WriteReport(args.Filename, []string{"Unable to construct holding: " + err.Error()}); return }
  result, err = Post(_url, params, holdingstr, "xml")
  if err != nil { file.WriteReport(args.Filename, []string{"Unable to push to alma: " + err.Error()}); return }
  args.Holding_id = ExtractHoldingID(result)
  fs.ItemsPF(args, tcmap, fs)
}

type ProcessHoldingBFun func(ProcessArgs, string, []map[string]string, FunMap)
// "B" method handles updates
func ProcessHoldingB(args ProcessArgs, marc_string string, tcmap []map[string]string, fs FunMap){
  path := []string{"bibs", args.Mms_id, "holdings", args.Holding_id}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  holdxml, err := Get(_url, params, "application/xml")
  if err != nil { file.WriteReport(args.Filename, []string{"Unable to obstain current holding: " + err.Error()}); return }
  holdingstr, err := UpdateHolding(marc_string, string(holdxml))
  if err != nil {
    if err.Error() == "skip update" {
      fs.ItemsPF(args, tcmap, fs)
      return
    } else {
      file.WriteReport(args.Filename, []string{"Unable to construct holding: " + err.Error()}); return
    }
  }
  _, err = Put(_url, params, holdingstr, "xml")
  if err != nil { file.WriteReport(args.Filename, []string{"Unable to push to alma: " + err.Error()}); return }
  fs.ItemsPF(args, tcmap, fs)
}

type ProcessItemsFun func(ProcessArgs, []map[string]string, FunMap)
func ProcessItems(args ProcessArgs, tcmap []map[string]string, fs FunMap){
  itemlist := []string{} //for reporting
  msgs := []string{} //also for reporting

  // iterate through the top containers
  // if an error occurs during the loop, report and continue
  for _,tc := range tcmap{
    var item = Item{}
    if tc["boundwith"] == "true" { continue } // skip boundwith containers
    if tc["ils_item"] != "" { //this is an update
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
  msgs = append(msgs, "items processed: " + strings.Join(itemlist, ", "))
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

  path := []string{ "bibs", args.Mms_id, "holdings", args.Holding_id, "items", tcmap["ils_item"]}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var result []byte
  //push record to alma
  if tcmap["ils_item"] == "" {
    result, err = Post(_url, params, itemstr, "json") } else {
    result, err = Put(_url, params, itemstr, "json")
  }
  if err != nil { return "", errors.New("problem posting to alma: " + err.Error()) }
  var item_id string
  if tcmap["ils_item"] == "" {
    item_id = ExtractItemID(result) } else {
    item_id = tcmap["ils_item"]
  }
  return item_id, err
}
// to be run at the start of the alma workflow.
// necessary to acquire boundwith bib ids
// currently assumes that if boundwith is true, bwbib exists
// if barcode on boundwith has changed, failed lookup is a fatal error here
// barcode must be fixed using the boundwith endpoint

func CheckTCMap(tcmap []map[string]string)([]map[string]string, error){
  for _,tc := range tcmap{
    if tc["barcode"] == "" {
      return tcmap, errors.New(fmt.Sprintf("error: item does not have barcode"))
    }
    if tc["boundwith"] != "true" { // skip unless bw
      continue
    }
    item, err := FetchByBarcode(tc["barcode"])
    if err != nil { return tcmap, err }
    tc["mms_id"] = ParseBibID(item)
  }
  return tcmap, nil
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

func CompileMissing(itemlist []byte, tc_barcodes []string) []ShortItem{
  items := gjson.GetBytes(itemlist, "item")
  result := []ShortItem{}
  for _, item :=  range items.Array(){
    p := ParseShortItem(item.String())
    if InRange(p.Barcode, tc_barcodes) == false { result = append(result, p) }
  }
  return result
}

func InRange(str string, arr []string) bool{
  for _, val := range arr {
    if str == val { return true }
  }
  return false
}

func BuildMessageFromMissing(missing []ShortItem) []string{
  messagelist := []string{}
  for _, item := range missing {
    messagelist = append(messagelist, item.Stringify())
  }
  return messagelist
}

func (s ShortItem) Stringify() string{
  return fmt.Sprintf("%s,%s,%s,%s", s.Barcode, s.CallNumber, s.Title, s.Description)
}

func ParseShortItem(item string) ShortItem{
  b := gjson.Get(item, "item_data.barcode").String()
  d := gjson.Get(item, "item_data.description").String()
  c := gjson.Get(item, "holding_data.call_number").String()
  t := gjson.Get(item, "bib_data.title").String()
  return ShortItem{ Barcode: b, Description: d, CallNumber: c, Title: t }
}

type ShortItem struct{
  Barcode string
  Description string
  CallNumber string
  Title string
}

// exclude barcdoes that are for boundwith items
func GetBarcodes(tcmap []map[string]string) []string{
  barcodes := []string{}
  for _,item := range tcmap{
    if item["boundwith"] == "true" { continue }
    barcodes = append(barcodes, item["barcode"])
  }
  return barcodes
}

type CheckItemsForMissingFun func(ProcessArgs, []map[string]string)
func CheckItemsForMissing(args ProcessArgs, tcmap []map[string]string){
  tc_barcodes := GetBarcodes(tcmap)
  if args.Holding_id == "" { file.WriteReport(args.Filename, []string{"Skipping barcode comparison, no holding available for lookup."}); return }
  path := []string{"bibs", args.Mms_id, "holdings", args.Holding_id, "items"}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  itemsjson, err := Get(_url, params, "application/json")
  if err != nil { file.WriteReport(args.Filename, []string{"Unable to request Alma item: " + err.Error()}); return }
  missing := CompileMissing(itemsjson, tc_barcodes)
  if len(missing) != 0 {
    export_name := fmt.Sprintf("Resource%s_ItemsForDelete.csv", args.Resource_id)
    messages := []string{ "Barcode,CallNumber,Title,Description" }
    messages = append(messages, BuildMessageFromMissing(missing)...)
    file.WriteReport(export_name, messages)
    base_url := os.Getenv("HOME_URL")
    file.WriteReport(args.Filename, []string{fmt.Sprintf("see items for removal at <a href=\"%s/reports/%s\">%s</a></p>", base_url, export_name, export_name)})
  }
}
