package handlers

import (
  "aspace_publisher/as"
  "aspace_publisher/oclc"
  "aspace_publisher/alma"
  "aspace_publisher/file"
)

func AlmaCrup(resource_id string, filename string, session_id string, oclc_token string){
  var args alma.ProcessArgs
  args.Resource_id = resource_id
  args.Repo_id = "2"
  args.Filename = filename
  var err error
  args.Session_id = session_id
  args.Oclc_token = oclc_token

  //acquire aspace resource
  rjson, err := as.AcquireJson(args.Session_id, args.Repo_id, "resources/" + args.Resource_id)
  if err != nil { file.WriteReport(args.Filename, []string{ "Could not aquire JSON from aspace: " + err.Error() }); return }

  args.Oclc_id, err = as.ExtractOclcId(rjson)
  if err != nil { file.WriteReport(args.Filename, []string{ "Problem retrieving OCLC id: " + err.Error() }); return }
  if args.Oclc_id == "" { file.WriteReport(args.Filename, []string{ "Could not acquire OCLC id from resource" }); return }
  //try for mms_id and create based on presence in resource json
  args.Mms_id, args.Create, err = as.GetMmsId(rjson)
  if err != nil { file.WriteReport(args.Filename, []string{ "Problem retrieving mms id: " + err.Error() }); return }
  //needed for holding record, appears as 099 in the aspace MARC but not OCLC's
  args.Id_0 = as.ExtractID0(rjson)


  //get oclc marc
  oclc_marc, err := oclc.Record(args.Oclc_token, args.Oclc_id)
  if err != nil { file.WriteReport(args.Filename, []string{ "Could not acquire OCLC MARC " + err.Error() }); return }
  marc_clean := oclc.UnformatXML(oclc_marc)
  tcmap, errmsgs := as.ExtractTCData(args.Session_id, args.Repo_id, args.Resource_id)
  if len(errmsgs) != 0 { file.WriteReport(args.Filename, errmsgs); return }

  tcmap, err = alma.CheckTCMap(tcmap)
  if err != nil { file.WriteReport(args.Filename, []string{ "Error attempting to acquire ids: " + err.Error() }); return }
  //launch processing, starting with bib
  //eventually hand this off to a worker?

  fs := alma.FunMap{ BoundwithPF: alma.ProcessBoundwith, HoldingAPF: alma.ProcessHoldingA, HoldingBPF: alma.ProcessHoldingB, ItemsPF: alma.ProcessItems, ItemPF: alma.ProcessItem, AfterBib: as.AfterBibCreate, SetHolding: oclc.SetHolding, UpdateTC: as.UpdateTC, CheckForMissingPF: alma.CheckItemsForMissing, CallWorker: alma.CallWorker }
  alma.ProcessBib(args, marc_clean, rjson, tcmap, fs)
  return
}
