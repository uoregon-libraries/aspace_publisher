package alma

import (
  "strings"
  "time"
  "errors"
  "fmt"
  "log/slog"
  "github.com/tidwall/sjson"
  "encoding/xml"
  "aspace_publisher/file"
  "aspace_publisher/as"
)

type BWFunMap struct {
  BibPF ProcessBWBibFun
  HoldingPF ProcessBWHoldingFun
  ItemPF ProcessBWItemFun
  Time TimeFun
  UpdateTC as.UpdateTCFun
}

type ProcessBWBibFun func(ProcessArgs, map[string]string, []map[string]string, BWFunMap)
//following the api bib update ex, i.e. just bib tags wrapping the record, no other bib data
func ProcessBWBib(args ProcessArgs, tcmap map[string]string, tccolls []map[string]string, fs BWFunMap){
  slog.Info(fmt.Sprintf("Creating bib %+v", args))
  var bib Bib
  var rec Record
  var _url string
  if args.Create {
    rec = ConstructBWMarc(tcmap, fs.Time())
    rec = Create774s(rec, tccolls)
    _url = BuildUrl([]string{"bibs"})
  } else {
    bibxml, err := FetchByHolding(tcmap["ils_holding"])
    if err != nil { file.WriteReport(args.Filename, []string{ err.Error() }); return }
    err = xml.Unmarshal(bibxml, &bib)
    if err != nil { file.WriteReport(args.Filename, []string{ err.Error() }); return }
    rec = UpdateBWBibRec(bib.Rec, args, tccolls)
    _url = BuildUrl([]string{"bibs", bib.Mms_id })
    args.Mms_id = bib.Mms_id
  }
  bib = ConstructBWBib(rec)
  params := []string{ ApiKey() }
  var result []byte
  strbib, err := bib.Stringify()
  if err != nil { file.WriteReport(args.Filename, []string{ err.Error() }); return }
  if args.Create {
    result, err = Post(_url, params, strbib, "xml")
  } else { result, err = Put(_url, params, strbib, "xml") }
  if err != nil { file.WriteReport(args.Filename, []string{ err.Error() }); return }
  if args.Mms_id == "" { args.Mms_id = ExtractBibID(result) }
  fs.HoldingPF(args, rec, tcmap, fs)
}

func ConstructBWBib(rec Record) Bib {
  bib := ConstructBib("", "<record></record>", "true")
  bib.Rec = rec
  return bib
}

func ConstructBWMarc(tc_data map[string]string, t time.Time) Record{
  date1 := t.Format("20060102150405")
  date2 := t.Format("20060102")

  var rec = Record{}
  rec.Leader = "00000npcaa2200000 i 4500"
  c005 := Controlfield{ Tag:"005", Value: date1 + ".0"}
  c008 := Controlfield{ Tag:"008", Value: fmt.Sprintf("%si19001920oru                 eng d", date2) }
  rec.Controlfield = []Controlfield{c005, c008}
  d245 := Datafield{ Tag:"245", Ind1: "0", Ind2: "0" }
  s245a := Subfield{ Code:"a", Value: "Multiple archival collections in a shared box with barcode " + tc_data["barcode"]}
  d245.Subfield = []Subfield{ s245a }
  d962 := Datafield{ Tag: "962", Ind1: " ", Ind2: " " }
  s962a := Subfield{ Code: "a", Value: "BoundwithRecord" }
  s9629 := Subfield{ Code: "9", Value: "local"}
  d962.Subfield = []Subfield{ s962a, s9629 }
  rec.Datafield = []Datafield{ d245, d962 }
  return rec
}

func extractMmsids(tccolls []map[string]string)[]string{
  var mmsids []string
  for _,coll := range tccolls{
    if coll["mms_id"] != "" { mmsids = append(mmsids, coll["mms_id"]) }
  }
  return mmsids
}
// requires the tccolls with mmsid
func UpdateBWBibRec(rec Record, args ProcessArgs, tccolls []map[string]string) Record{
  rec = Add774s(rec, tccolls)
  mmsids := extractMmsids(tccolls)
  rec = Remove774s(rec, mmsids)
  return rec
}

type TimeFun func()time.Time
func GetTime()time.Time{
  return time.Now()
}

type ProcessBWHoldingFun func(ProcessArgs, Record, map[string]string, BWFunMap)
func ProcessBWHolding(args ProcessArgs, rec Record, tcmap map[string]string, fs BWFunMap){
  slog.Info(fmt.Sprintf("Creating holding %+v", args))
  path := []string{"bibs", args.Mms_id, "holdings", args.Holding_id}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var holding string
  var err error
  if args.Create != true {
    holdxml, err := Get(_url, params, "application/xml")
    if err != nil { file.WriteReport(args.Filename, []string{"Unable to obtain current holding: " + err.Error()}); return }
    holding, err = UpdateBWHolding(string(holdxml), tcmap) } else {
    holding, err = ConstructBWHolding(rec, tcmap ) }
  if err != nil { file.WriteReport(args.Filename, []string{"Did not construct holding: " + err.Error()}); return }
  var result []byte
  if args.Create {
    result, err = Post(_url, params, holding, "xml") } else {
    result, err = Put(_url, params, holding, "xml") }
  if err != nil { file.WriteReport(args.Filename, []string{"Did not push to alma: " + err.Error()}); return }
  args.Holding_id = ExtractHoldingID(result)
  fs.ItemPF(args, tcmap, fs)
}
// update barcode if needed
func UpdateBWHolding(hold string, tcmap map[string]string)(string, error){
  holding, err := ParseXML(hold, "holding")
  if err != nil { return "", err }
  sfh := holding.FindElement("//subfield[@code='h']")
  if strings.Contains(sfh.Text(), tcmap["barcode"]) {
    return "", errors.New("skip update, no change to barcode")
  }
  sfh.SetText(fmt.Sprintf("[%s]", tcmap["barcode"]))
  // acc to the API docs, remove the holdingId
  id_ptr := holding.FindElement("//holding_id")
  parent := id_ptr.Parent()
  parent.RemoveChild(id_ptr)
  str, err := holding.WriteToString()
  return str, err
}

// keeps leader and controlfield, overwrites datafield
func ConstructBWHolding(rec Record, tcmap map[string]string)(string, error){
  var h = Holding{}
  h.Suppress = false
  h.Rec.Leader = rec.Leader
  h.Rec.Controlfield = rec.Controlfield
  sfb := Subfield{Code:"b", Value:"SpecColl"}
  sfc := Subfield{Code:"c", Value: "spmanus"}
  sfh := Subfield{Code:"h", Value: fmt.Sprintf("[%s]", tcmap["barcode"])}
  df852 := Datafield{Ind1:"8", Ind2:" ", Tag:"852"}
  df852.Subfield = []Subfield{sfb, sfc, sfh}
  h.Rec.Datafield = []Datafield{ df852 }
  str, err := h.Stringify()
  return str, err
}
type ProcessBWItemFun func(ProcessArgs, map[string]string, BWFunMap)
//must wrap up at the end!
func ProcessBWItem(args ProcessArgs, tcmap map[string]string, fs BWFunMap){
  slog.Info(fmt.Sprintf("Starting item...args %+v, tcmap %+v", args, tcmap))
  var itembyte []byte
  var err error
  if tcmap["ils_item"] != "" {
    path := []string{"bibs", args.Mms_id, "holdings", tcmap["ils_holding"], "items", tcmap["ils_item"]}
    _url := BuildUrl(path)
    params := []string{ ApiKey() }
    itembyte, err = Get(_url, params, "application/json")
    if err != nil { file.WriteReport(args.Filename, []string{ "Error processing item: " + err.Error()}); return }
  }

  var itemstr string
  if itembyte != nil { itemstr, err = UpdateBWItem(args.Holding_id, string(itembyte), tcmap)
  } else { itemstr, err = ConstructBWItem(args.Holding_id, tcmap) }
  if err != nil { file.WriteReport(args.Filename, []string{ "Error processing item: " + err.Error()}); return }

  path := []string{ "bibs", args.Mms_id, "holdings", args.Holding_id, "items", tcmap["ils_item"]}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var result []byte
  //push record to alma
  if tcmap["ils_item"] == "" {
    result, err = Post(_url, params, itemstr, "json") } else {
    result, err = Put(_url, params, itemstr, "json")
  }
  if err != nil { file.WriteReport(args.Filename, []string{ "Error processing item: " + err.Error()}); return }
  var item_id string
  if tcmap["ils_item"] == "" {
    item_id = ExtractItemID(result)
    err = fs.UpdateTC(args.Repo_id, args.Holding_id, item_id, args.Session_id, tcmap)
    if err != nil { file.WriteReport(args.Filename, []string{ "error updating tc: " + err.Error() } ) }
    } else {
    item_id = tcmap["ils_item"]
  }
  file.WriteReport(args.Filename, []string{ "item processed: " + item_id })
}

// no policy
func ConstructBWItem(holding_id string, tc_data map[string]string)(string, error){
  var item Item
  item.Holding_data.Holding_id = holding_id
  item.Holding_data.Copy_id = "1"
  item.Item_data.Library = Value{ Val: "SpecColl"}
  item.Item_data.Location = Value{ Val: "spmanus"}
  item.Item_data.Base_status = Value{ Val: "1" }
  item.Item_data.Physical_material_type = Value{ Val: "MANUSCRIPT" }
  item.Item_data.Barcode = tc_data["barcode"] //may change, eg falls off
  item.Item_data.Description = fmt.Sprintf("%s %s", tc_data["type"], tc_data["indicator"])
  istring, err := item.Stringify()
  if err != nil { return "", err}
  return istring, nil
}

//if check barcode at the holding phase is unchanged, this will not be called
// holding id is not in docs for update, but *is* in the create, leaving it in for now
func UpdateBWItem(holding_id string, item_json string, tc_data map[string]string)(string, error){
  i2, err := sjson.Set(item_json, "item_data.description", fmt.Sprintf("%s %s", tc_data["type"], tc_data["indicator"]))
  i3, err := sjson.Set(i2, "item_data.barcode", tc_data["barcode"])
  if err != nil { return "", err }
  i4, err := sjson.Delete(i3, "item_data.modification_date")
  if err != nil { return "", err }
  i5, err := sjson.Delete(i4, "item_data.creation_date")
  if err != nil { return "", err }
  i6, err := sjson.Delete(i5, "item_data.pid")
  return i6, err
}