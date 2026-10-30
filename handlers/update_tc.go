package handlers

import (
  "aspace_publisher/as"
  "aspace_publisher/alma"
  "aspace_publisher/file"
)

//standalone, to be called from launcher
func UpdateTopContainers(repo_id, resource_id, session_id string)string{
  fname := file.Filename()
  tcmaps, msgs := as.ExtractTCData(session_id, repo_id, resource_id)
  for _,tcmap := range tcmaps {
    if tcmap["boundwith"] != "true" {
      resp, err := UpdateTopContainer(tcmap, session_id, repo_id)
      msgs = append(msgs, resp)
      if err != nil { msgs = append(msgs, err.Error()) }
    }
  }
  file.WriteReport(fname, msgs)
  return fname
}

//standalone, to be called from launcher
func UpdateTopContainerSA(tc_id, session_id, repo_id string)string{
  fname := file.Filename()
  tc, err := as.ExtractSingleTC(session_id, repo_id, tc_id)
  if err != nil { file.WriteReport(fname, []string{ err.Error() })}
  response, err := UpdateTopContainer(tc.Mapify(), session_id, repo_id)
  if err != nil { file.WriteReport(fname, []string{ response, err.Error() })
  } else { file.WriteReport(fname, []string{ response }) }
  return fname
}

func UpdateTopContainer(tcmap map[string]string, session_id, repo_id string)(string,error){
  item, err := alma.FetchByBarcode(tcmap["barcode"])
  if err != nil { return "",err }
  item_id, holding_id := alma.ParseHoldingItem(item)
  tc_id := as.ExtractID(tcmap["uri"])
  response, err := as.UpdateTC(repo_id,holding_id, item_id, session_id, tc_id)
  return response, err
}