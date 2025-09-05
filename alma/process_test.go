package alma

import (
  "testing"
  "os"
)

func TestBuildUrl(t *testing.T){
  path := []string{"one", "two", "three", ""}
  os.Setenv("ALMA_URL", "http://blah.org")
  url := BuildUrl(path)
  if url != "http://blah.org/one/two/three" { t.Errorf("incorrect url") }

  path = []string{"one", "two", "", "three"}
  url = BuildUrl(path)
  if url != "http://blah.org/one/two/three" { t.Errorf("incorrect url") }

}
