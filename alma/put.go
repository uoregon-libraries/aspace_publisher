package alma

import (

  "os"
  "net/http"
  "log"
  "time"
  "errors"
  "io"
  "strings"
  "aspace_publisher/connect"
)

//url /almaws/v1/bibs/<mms_id>/holdings/<holding_id>/items/<item_id>
//params: apikey=abcde12341234
func Put(_url string, params []string, data string, format string)([]byte, error){
  body, err := push("PUT", _url, params, data, format)
  return body, err
}

func Post(_url string, params []string, data string, format string)([]byte, error){
  body, err := push("POST", _url, params, data, format)
  return body, err
}

func push(method string, _url string, params []string, _data string, format string)([]byte, error){
  debug := os.Getenv("DEBUG")
  verbose := os.Getenv("VERBOSE")
  param_str := strings.Join(params[:], "&")
  final_url := _url + "?" + param_str
  if debug == "true" {
    log.Println("Swapping " + final_url + "for test url")
    final_url = os.Getenv("TEST_URL")
  }
  data := strings.NewReader(_data)
  req, err := http.NewRequest(method, final_url, data)
  if err != nil { log.Println(err); return nil, errors.New("unable to create http request") }
  req.Header.Set("Content-Type", "application/" + format)
  req.Header.Set("accept", "application/" + format)
  connect.RequestDump(verbose, req)
  client := &http.Client{
    Timeout: time.Second * 60,
  }
  response, err := client.Do(req)
  connect.ResponseDump(verbose, response)
  if err != nil { log.Println(err); return nil, errors.New("unable to complete http request") }
  defer response.Body.Close()
  body, err := io.ReadAll(response.Body)
  if err != nil { log.Println(err); return nil, errors.New("unable to read response from alma") }
  if response.StatusCode != 200 { return body, errors.New("alma errors") }

  return body, nil
}
