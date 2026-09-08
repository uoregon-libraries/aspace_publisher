package alma

import(
  "net/http"
  "log/slog"
  "time"
  "errors"
  "io"
  "strings"
  "aspace_publisher/connect"
  "github.com/tidwall/gjson"
  "encoding/xml"
)

//url /almaws/v1/conf/sets/<setid>/members
//params: limit=100, apikey=abcde12341234
//url /almaws/v1/bibs/<mms_id>/holdings/<holding_id>/items/<item_id>
//params: view=brief, apikey=abcde12341234

func Get(url string, params []string, accept string)([]byte, error){
  param_str := strings.Join(params[:], "&")
  final_url := url + "?" + param_str

  req, err := http.NewRequest("GET", final_url, nil)
  if err != nil { slog.Error(err.Error()); return nil, errors.New("unable to create http request") }
  req.Header.Set("accept", accept)
  connect.RequestDump(req)
  client := &http.Client{
    Timeout: time.Second * 60,
  }

  response, err := client.Do(req)
  connect.ResponseDump(response)
  if err != nil { slog.Error(err.Error()); return nil, errors.New("unable to complete http request") }
  defer response.Body.Close()
  body, err := io.ReadAll(response.Body)
  if err != nil { slog.Error(err.Error()); return nil, errors.New("unable to read response from alma") }
  if response.StatusCode != 200 {
    if accept == "application/json" {
      return body, errors.New(ExtractJsonError(body))
    } else {
      return body, errors.New(ExtractXmlError(body))
    }
  }
  return body, nil
}


func ExtractJsonError(body []byte)string{
  return gjson.GetBytes(body, "errorList.error.0.errorMessage").String()
}

func ExtractXmlError(body []byte)string{
  var wr WebServiceResult
  err := xml.Unmarshal(body, &wr)
  if err != nil { return err.Error() }
  return wr.ErrorList.AlmaError[0].ErrorMess
}

type WebServiceResult struct{
  XMLName xml.Name `xml:"web_service_result"`
  ErrorsExist bool `xml:"errorsExist"`
  ErrorList ErrorList `xml:"errorList"`
}
type ErrorList struct{
  XMLName xml.Name `xml:"errorList"`
  AlmaError []AlmaError `xml:"error"`
}
type AlmaError struct {
  XMLName xml.Name `xml:"error"`
  ErrorCode string `xml:"errorCode"`
  ErrorMess string `xml:"errorMessage"`
}
