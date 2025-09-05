package alma

import(
  "github.com/tidwall/gjson"
  "time"
  "encoding/json"
  "log"
  "strings"
  "os"
  "strconv"
  "net/url"
  "fmt"
  "aspace_publisher/file"
)

// filename, list
type ProcessFunc func(string, []string)

func CheckJob(joblink string, nextFun ProcessFunc, filename string, list []string){
  MAX, _ := strconv.Atoi(os.Getenv("JOB_MAX_TRIES"))
  span,_ := time.ParseDuration(os.Getenv("JOB_WAIT_TIME"))
  i := 0
  params := []string{ ApiKey() }
  var result map[string]string
  for i < MAX {
    resp,err := Get(joblink, params)
    if err != nil { 
      log.Println(err)
      /*** count this as one try ****/
      i += 1
      time.Sleep(span)
      continue
    }
    result = ExtractJobResults(resp)
    if !strings.Contains(result["status"], "COMPLETED") {
      i += 1
      time.Sleep(span)
      continue
    }
    file.WriteReport(filename, []string{ fmt.Sprintf("%s: %s, %s", result["jobname"], result["status"], joblink) } )
    if result["status"] == "COMPLETED_SUCCESS" {
      if nextFun != nil {
        nextFun(filename, list)
      }
    }
    return
  }
  file.WriteReport(filename, []string{ fmt.Sprintf("See %s re: %s", joblink, result["alert"]) } )
}

func DummyFunc(word string, list map[string][]bool){
  log.Println("word is " + word)
  for k, v := range list{
    log.Println(fmt.Sprintf("%s, %t", k, v[0]))
  }
}

func ExtractJobResults(resp []byte)map[string]string{
  //the docs on using progress are unclear
  result := map[string]string{}
  jobname := gjson.GetBytes(resp, "job_info.name")
  result["jobname"] = jobname.String()
  status := gjson.GetBytes(resp, "status.value")
  result["status"] = status.String()
  alert := gjson.GetBytes(resp, "alert.value")
  result["alert"] = alert.String()
  return result
}

func SubmitJob(jobid string, job_params []Param)(string, error){
  _url,_ := url.Parse(BaseUrl())
  _url = _url.JoinPath("conf", "jobs", jobid)
  params := []string{ "op=run", ApiKey() }
  job := JobInit(job_params)
  json,_ := json.Marshal(job)
  resp,err := Post(_url.String(), params, string(json), "json")
  if err != nil { 
    log.Println(err)
    return "", err
  }
  link := ExtractJobInstance(resp)
  return link, nil
}

func ExtractJobInstance(resp []byte)(string){
  instance := gjson.GetBytes(resp, "additional_info.link")
  return instance.String()
}

// boilerplate
// finish this when the actual jobs are available
func JobInit(params []Param)Job{
  job := Job{ Parameter: params}
  return job
}

type Job struct{
  Parameter []Param `json:"parameter"`
}

type Param struct{
  Name Val     `json:"name"`
  Value string `json:"value"`
}

// Val declaration in update_set
