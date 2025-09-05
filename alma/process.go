package alma

import(
  "strings"
  "os"
  "net/url"
  "log"
  "time"
  "slices"
  "errors"
  "aspace_publisher/file"
)

func ProcessBib(mms_id string, marc_string string, create bool)(string, error){
  bib, err := ConstructBib(marc_string)
  if err != nil { return "", errors.New("unable to construct bib: " + err.Error()) }
  path := []string{"bibs", mms_id}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var result []byte
  if create { 
    result, err = Post(_url, params, bib, "xml") } else {
    result, err = Put(_url, params, bib, "xml")
  }
  if err != nil { log.Println(err); return "", err }
  if create {
    mms_id = ExtractBibID(result)
  }
  return mms_id, err
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

func ProcessHolding(mms_id string, holding_id string, marc_string string, id_0 string, create bool)(string, error){
  holding, err := ConstructHolding(marc_string, id_0)
  if err != nil { return "", errors.New("Unable to construct holding: " + err.Error()) }
  path := []string{"bibs", mms_id, "holdings", holding_id}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var result []byte
  if create {
    result, err = Post(_url, params, holding, "xml") } else {
    result, err = Put(_url, params, holding, "xml")
  }
  if err != nil { return "", err }
  if create {
    holding_id = ExtractHoldingID(result)
  }
  return holding_id, err
}

func ProcessItem(mms_id string, holding_id string, item_id string, container_data map[string]string, create bool)(string, error){
  item, err := ConstructItem(item_id, holding_id, container_data)
  if err != nil { return "", errors.New("Unable to construct item" + err.Error()) }
  path := []string{"bibs", mms_id, "holdings", holding_id, "items", item_id}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var result []byte
  if create {
    result, err = Post(_url, params, item, "json") } else {
    result, err = Put(_url, params, item, "json")
  }
  if err != nil { return "", errors.New("problem posting to alma: " + err.Error()) }
  if create {
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
