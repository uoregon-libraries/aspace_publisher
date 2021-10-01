package main

import(
  "golang.org/x/oauth2/clientcredentials"
  "golang.org/x/net/context"
  "net/http"
  "os"
)

func oclc_conf() clientcredentials.Config {
  secret := os.Getenv("OCLC_SECRET")
  key := os.Getenv("OCLC_KEY")
  oclc_token_url := os.Getenv("OCLC_TOKEN_URL")

  config := clientcredentials.Config{
    ClientID: key,
    ClientSecret: secret,
    TokenURL: oclc_token_url,
    Scopes: []string{"WorldCatMetadataAPI"},
  }
  return config
}

func authenticated_client() *http.Client {
  config := oclc_conf()
  client := config.Client(context.Background())
  return client
}

