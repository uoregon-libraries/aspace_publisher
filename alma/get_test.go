package alma

import (
  "testing"
)

func TestExtractJsonError(t *testing.T){
  json_str := `{"errorsExist":true,"errorList":{"error":[{"errorCode":"123","errorMessage":"Item not found"}]},"result":null}`
  err_mess := ExtractJsonError([]byte(json_str))
  if err_mess != "Item not found" { t.Errorf("failed to extract json error message") }
}

func TestExtractXmlError(t *testing.T){
  xml_str := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><web_service_result xmlns="http://com/exlibris/urm/general/xmlbeans"><errorsExist>true</errorsExist><errorList><error><errorCode>123</errorCode><errorMessage>Item not found</errorMessage></error></errorList></web_service_result>`
  err_mess := ExtractXmlError([]byte(xml_str))
  if err_mess != "Item not found" { t.Errorf("failed to extract xml error message") }
}
