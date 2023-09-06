package oclc

import(
  "net/http"
  "net/http/httptest"
  "testing"
  "os"
  "fmt"
)

func TestTokenRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/token?grant_type=client_credentials&scope=WorldCatMetadataAPI" {
			t.Errorf("authenticate client request URL = %q; want %q", r.URL, "/token?grant_type=client_credentials&scope=WorldCatMetadataAPI")
            fmt.Println(r.URL.String())
		}
        //r should have basic auth
		headerAuth := r.Header.Get("Authorization")
		if headerAuth != "Basic Q0xJRU5UX0lEOkNMSUVOVF9TRUNSRVQ=" {
			t.Errorf("Unexpected authorization header, %v is found.", headerAuth)
		}
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
        json := `{"access_token":"90d64460d14870c08c81352a05dedd3465940a7c","token_type":"bearer"}`
		w.Write([]byte(json))
	}))

	defer ts.Close()
    os.Setenv("OCLC_TOKEN_URL", ts.URL + "/token")
    os.Setenv("OCLC_KEY", "CLIENT_ID")
    os.Setenv("OCLC_SECRET", "CLIENT_SECRET")
    var ot OclcToken
    _ = ot.GetToken()
	if ot.AccessToken != "90d64460d14870c08c81352a05dedd3465940a7c" {
		t.Errorf("Access token = %q; want %q", ot.AccessToken, "90d64460d14870c08c81352a05dedd3465940a7c")
	}
	if ot.TokenType != "bearer" {
		t.Errorf("token type = %q; want %q", ot.TokenType, "bearer")
	}
    fmt.Println("oclc_auth testing")
}
