package alma

import(
  "net/http"
  "log"
  "time"
  "errors"
  "io"
  "strings"
)

//url /almaws/v1/conf/sets/<setid>/members
//params: limit=100, apikey=abcde12341234
//url /almaws/v1/bibs/<mms_id>/holdings/<holding_id>/items/<item_id>
//params: view=brief, apikey=abcde12341234

func Get(url string, params []string)([]byte, error){
  param_str := strings.Join(params[:], "&")
  final_url := url + "?" + param_str

  req, err := http.NewRequest("GET", final_url, nil)
  if err != nil { log.Println(err); return nil, errors.New("unable to create http request") }
  req.Header.Set("accept", "application/json")

  client := &http.Client{
    Timeout: time.Second * 60,
  }

  response, err := client.Do(req)

  if err != nil { log.Println(err); return nil, errors.New("unable to complete http request") }
  defer response.Body.Close()
  body, err := io.ReadAll(response.Body)
  if err != nil { log.Println(err); return nil, errors.New("unable to read response from alma") }
  if response.StatusCode != 200 { return body, errors.New("alma errors") }

  return body, nil
}
