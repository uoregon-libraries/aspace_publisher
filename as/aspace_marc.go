package as

import (
  "aspace_publisher/marc"
)

type AspaceMarc struct{
  MarcRec MarcRecord
  AspaceId string
}

func (am *AspaceMarc) Initialize(marc string) string,error{
  am.MarcRec.Initialize(marc)
  err = am.set_AspaceId(); if err != nil { return "", err }

}

func (am *AspaceMarc) set_AspaceId() error{
  datafield := am.MarcRec.get_datafield("856")
  if datafield == nil { return errors.New("Cannot set ApaceId, datafield not found") }
  subfield := marc.Get_subfield("u", datafield)
  if subfield == nil { return errors.New("Cannot set AspaceId, subfield not found") }
  if subfield.Value == "" { return errors.New("Cannot set AspaceId, value empty") }
  am.AspaceId = subfield.Value
  return nil
}
