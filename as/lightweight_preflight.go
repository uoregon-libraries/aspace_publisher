package as

import (
  "aspace_publisher/file"
  "encoding/json"
)

type Preflight struct {
  CreateBib bool
  CreateHolding bool
  BoundwithBarcodes []string
  CreateItemsBarcodes []string
  UpdateItemsBarcodesPids []string
}

func LightweightPreflight(session_id string, resource_id string, repo_id string)string{
  //setup
  fname := file.Filename()
  rjson, err := AcquireJson(session_id, repo_id, "resources/" + resource_id)
  if err != nil { file.WriteReport(fname, []string{ "Could not aquire JSON from aspace: " + err.Error() }); return fname }
  _, create, err := GetMmsId(rjson)
  if err != nil { file.WriteReport(fname, []string{ "Problem retrieving mms id: " + err.Error() }); return fname }
  tcmap, errmsgs := ExtractTCData(session_id, repo_id, resource_id)
  if len(errmsgs) != 0 { file.WriteReport(fname, errmsgs); return fname }

  p := Preflight{}
  p.CreateBib = create
  //boundwith:
  bw := []string{}
  for _, tc := range tcmap{
    if tc["boundwith"] == "true" { bw = append(bw, tc["barcode"]) }
  }
  p.BoundwithBarcodes = bw
  //holding
  holding := false
  holding_id := ""
  for _, tc := range tcmap{
    if tc["boundwith"] == "false" { 
      holding = true
      if holding_id == "" { holding_id =  tc["ils_holding"]}
    }
  }
  if holding {
    p.CreateHolding = (holding_id == "")
  }

  //items
  creates := []string{}
  updates := []string{}
  for _, tc := range tcmap{
    if tc["boundwith"] == "false" { 
      if tc["ils_item"] == "" { creates = append(creates, tc["barcode"])
      } else { updates = append(updates, tc["barcode"] + "/" + tc["ils_item"]) }
    }
  }
  if len(creates) > 0 { p.CreateItemsBarcodes = creates }
  if len(updates) > 0 { p.UpdateItemsBarcodesPids = updates }

  jsonbytes,err := json.Marshal(p)
  if err != nil {}
  file.WriteReport(fname, []string{ string(jsonbytes) })
  return fname
}