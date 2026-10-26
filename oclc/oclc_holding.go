package oclc

import(
  "fmt"
)

type SetHoldingFun func(string, string)(string, string, error)
func SetHolding(oclc_num string, token string)(string, string, error){
  url := fmt.Sprintf("manage/institution/holdings/%s/set", oclc_num)
  oclc_resp, contype, err := Request(token, "POST", "", url, "", "json")
  return oclc_resp, contype, err
}
