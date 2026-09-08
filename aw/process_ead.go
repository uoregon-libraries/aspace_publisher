package aw

import (
  "os"
  "aspace_publisher/utils"
)

func ProcessEad(repo_id, resource_id, ead_orig string, aw_session string, operation string) (string, error) {

  ead_prepped, eadname, ark, err := PrepareEad(repo_id, resource_id, ead_orig)
  if err != nil{
    return "", err
  }

  ead_converted, err := CallConversion(ead_prepped)
  if err != nil {
    return "", err
  }

 if operation == "convert" {
    return ead_converted, nil
  }

  f, err := os.Create(eadname)
  if err != nil {
    return "", err
  }
  defer f.Close()
  defer os.Remove(f.Name())
  _, err = f.Write([]byte(ead_converted))
  if err != nil {
    return "", err
  }

   vals, err := MakeUploadMap(ark, "ead", f.Name())
  if err != nil {
    return "", err
  }
  // create form
  form, boundary, err := utils.CreateMultipartFormData(vals)
  if err != nil {
    return "", err
  }
  response, err := Request(aw_session, boundary, form, operation)
  if err != nil {
    return "", err
  }

  parsed, err := ParseResult(response)
  if err != nil {
    return "", err
  }
  return parsed, nil
}
