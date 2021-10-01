package aspace_export

import(
  "testing"
  "fmt"
)

var ae *AspaceExport

func TestNoFile(t *testing.T) {
  ae = &AspaceExport{ FileName: "../files/mark.txt" }
  _, e := ae.Initialize()
  if e == nil {
    t.Errorf("error not returned")
  } else { fmt.Printf("Nofile error expected: %s\n", e.Error()) }
}

func TestBadFile(t *testing.T){
  ae = &AspaceExport{ FileName: "../files/badmarc.txt" }
  _, e := ae.Initialize()
  if e == nil {
    t.Errorf("error not returned")
  } else { fmt.Printf("Bad XML eror expected: %s\n", e.Error()) }
}

//This has to be set. Error if not found
func TestAspaceID(t *testing.T) {
  ae = &AspaceExport{ FileName: "../files/marc.txt" }
  _, e := ae.Initialize()
  if e != nil { fmt.Printf("Unexpected error: %s\n", e.Error()) }
  if ae.AspaceId != "https://scua.uoregon.edu/repositories/2/resources/2023" {
       t.Errorf("incorrect id set")
  } else { fmt.Println("AspaceId ok") }
}

//Ok for this not to be found and set, as data may be missing
//In this case, data is NOT missing
func TestOclcId(t *testing.T) {
  if ae.OclcId != "ocm18421997" {
    t.Errorf("incorrect id set")
  } else { fmt.Println("OclcId ok") }
  if ae.OclcDate != "20210720192458.8"{
    t.Errorf("incorrect date set")
  } else { fmt.Println("OclcDate ok") }
}

//Example has an Oclc ID so this should be PUT
func TestProtocol(t *testing.T) {
  if ae.Protocol != "PUT" {
    t.Errorf("incorrect protocol set")
  } else { fmt.Println("protocol ok") }
}



