package connect

import(
  "net/http"
  "net/http/httputil"
  "log/slog"
  "fmt"
  "os"
)

func RequestDump(req *http.Request, gate ...string){
  verbose := os.Getenv("VERBOSE")
  if len(gate) > 0 { verbose = os.Getenv(gate[0]) }
  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { slog.Error(err.Error()) } else {
      slog.Info(fmt.Sprintf("REQUEST:\n%s", string(reqdump))) }
  }
}

func ResponseDump(response *http.Response, gate ...string){
  verbose := os.Getenv("VERBOSE")
  if len(gate) > 0 { verbose = os.Getenv(gate[0]) }
  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { slog.Error(err.Error()) } else {
      slog.Info(fmt.Sprintf("RESPONSE:\n%s", string(respdump))) }
  }
}
