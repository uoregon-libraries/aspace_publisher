package as

import (
  "testing"
  "net/http/httptest"
  "os"
  "net/http"
  "fmt"
)

func TestAuthenticateAS(t *testing.T){
  resp := `{"session":"abcdefghijk123456789", "agent_record":{"ref":"/agents/people/123"}}`
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.FormValue("password") != "green" { t.Errorf("pass incorrect") }
    if r.URL.Path != "/api/users/kermit/login" { t.Errorf("incorrect path: %v", r.URL.Path) }
    fmt.Fprint(w, resp)
  }))
  defer ts.Close()
  os.Setenv("ASPACE_URL", ts.URL + "/api/")
  session, agent, err := AuthenticateAS("kermit", "green")
  if err != nil { t.Errorf("auth test failed") }
  if session != "abcdefghijk123456789" { t.Errorf("session incorrect") }
  if agent != "/agents/people/123" { t.Errorf("agent incorrect") }

}
