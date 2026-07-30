package connect

import(
  "net/http"
  "net/http/httputil"
  "log/slog"
  "fmt"
)

func RequestDump(verbose string, req *http.Request){
  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { slog.Error(err.Error()) } else {
      slog.Info(fmt.Sprintf("REQUEST:\n%s", string(reqdump))) }
  }
}

func ResponseDump(verbose string, response *http.Response){
  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { slog.Error(err.Error()) } else {
      slog.Info(fmt.Sprintf("RESPONSE:\n%s", string(respdump))) }
  }
}
