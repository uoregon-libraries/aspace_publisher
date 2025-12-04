package oclc

import(
  "fmt"
)

func SetHolding(oclc_num string, token string)(string, error){
  url := fmt.Sprintf("manage/institution/holdings/%s/set", oclc_num)
  oclc_resp, err := Request(token, "POST", "", url, "", "json")
  return oclc_resp, err
}
