package oclc

import (
  "fmt"
  "aspace_publisher/as"
)

type FunMap struct {
  AsCheckIsPublished as.CheckIsPublishedFun
  AsAcquireJson as.AcquireJsonFun
  AsAcquireMarc as.AcquireMarcFun
  AsGetOclcId as.GetOclcIdFun
  OclcRequest RequestFun
  OclcRecord RecordFun
  AsPost as.PostFun
}

func DummyRequest(token string, method string, marc string, path string, id string, accept string) (string, error){
  fmt.Sprintf("token: %v, method: %v, marc: %v, path: %v, id: %v, accept: %v", token, method, marc, path, id, accept)
  return `<record></record>`, nil
}

func DummyReport(token string, id string)(string, error){
  fmt.Sprintf("token: %v, id: %v",token, id)
  return `<record></record>`, nil
}

func DummyAsCheckIsPublished(resource_id string, repo_id string, session_id string)(string, error){
  fmt.Sprintf("resource_id: %v, repo_id: %v, session_id: %v", resource_id, repo_id, session_id)
  return "true", nil
}

func DummyAsAcquireJson(sessionid string, repo_id string, resource_id string){}
//could also return fixtures/marc_1867.xml
func DummyAsAcquireMarc(sessionid string, repo_id string, resource_id string, published string) (string, error){
  fmt.Sprintf("sessionid: %v, repo_id: %v, resource_id: %v, published: %v", sessionid, repo_id, resource_id, published)
  return `<?xml version='1.0' encoding='UTF-8'?><record></record>`, nil
}

func DummyAsGetOclcId(resource_id string, repo_id string, session_id string)(string, error){
  fmt.Sprintf("resource_id: %v, repo_id: %v, session_id: %v", resource_id, repo_id, session_id)
  return "123456789", nil
}

func DummyAsPost(sessionid string, identifier string, repo_id string, record_id string, json_record string ) as.Response {
  fmt.Sprintf("sessionid: %v, identifier: %v, repo_id: %v, record_id: %v, json_record: %v", sessionid, identifier, repo_id, record_id, json_record)
  return as.Response{"1234", as.BuildMessage(`{"status": "success"}`)}
}
