package alma

import(
  "strings"
  "os"
  "net/url"
  "log"
  "time"
  "slices"
)

func ProcessBib(mms_id string, marc_string string, create bool)(string, error){
  bib := ConstructBib(marc_string)
  path := []string{"bibs", mms_id}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var result []byte
  var err error
  if create { 
    result, err = Post(_url, params, bib) } else {
    result, err = Put(_url, params, bib)
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

func ProcessHolding(mms_id string, holding_id string, marc_string string, create bool)(string, error){
  holding := ConstructHolding(marc_string)
  path := []string{"bibs", mms_id, "holdings", holding_id}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var result []byte
  var err error
  if create {
    result, err = Post(_url, params, holding) } else {
    result, err = Put(_url, params, holding)
  }
  if err != nil { return "", err }
  if create {
    holding_id = ExtractHoldingID(result)
  }
  return holding_id, err
}

func ProcessItem(mms_id string, holding_id string, item_id string, container_data map[string]string, create bool)(string, error){
  item := ConstructItem(item_id, holding_id, container_data)
  path := []string{"bibs", mms_id, "holdings", holding_id, "items", item_id}
  _url := BuildUrl(path)
  params := []string{ ApiKey() }
  var result []byte
  var err error
  if create {
    result, err = Post(_url, params, item) } else {
    result, err = Put(_url, params, item)
  }
  if err != nil { return "", err }
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
  err := UpdateSet(setid, "BIB_MMS", list)
  if err != nil {}
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
  if err != nil {}
  span,_ := time.ParseDuration(os.Getenv("JOB_WAIT_TIME"))
  time.Sleep(span)
  CheckJob(instance, nil, filename, nil)
}

func BaseUrl()string{
  return os.Getenv("ALMA_URL")
}
