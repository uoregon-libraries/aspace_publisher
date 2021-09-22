package aspace_export

import(
  "testing"
  "fmt"
)

var ae *AspaceExport

func TestAspaceID(t *testing.T) {
  ae = &AspaceExport{ FileName: "/home/lsato/aspace_publisher/mark.txt" }
  str,e := ae.set_xml()
  if e != nil { t.Errorf("error not returned") }
  if str != ae.FileName { t.Errorf("filename not reported") }
}

func TestAspaceID(t *testing.T) {
  ae = &AspaceExport{ FileName: "/home/lsato/aspace_publisher/marc.txt" }
  ae.set_xml()
  ae.set_AspaceId()
  if ae.AspaceId != "https://scua.uoregon.edu/repositories/2/resources/2023" {
    t.Errorf("incorrect id set")
  }
}

func TestOclcId(t *testing.T) {
  ae.set_oclc_fields()
  if ae.OclcId != "ocm18421997" {
    t.Errorf("incorrect id set")
  }
  if ae.OclcDate != "20210720192458.8"{
    t.Errorf("incorrect date set")
  }
}

func TestProtocol(t *testing.T) {
  ae.set_protocol()
  if ae.Protocol != "POST" {
    t.Errorf("incorrect protocol set")
  }
}



