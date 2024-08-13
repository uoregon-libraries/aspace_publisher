package oclc

import(
  "testing"
)

func TestAssembleUrl(t *testing.T){
  parts_no_id := []string{"https://blah.org","path",""}
  parts_id := []string{"https://blah.org","path","abcd"}
  url_no_id := assembleUrl(parts_no_id)
  url_id := assembleUrl(parts_id)
  if url_no_id != "https://blah.org/path" {t.Fatalf("assembled url is incorrect")}
  if url_id != "https://blah.org/path/abcd" {t.Fatalf("assembled url is incorrect")}
}
