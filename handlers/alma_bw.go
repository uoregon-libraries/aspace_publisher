package handlers

import (
  "aspace_publisher/as"
  "aspace_publisher/alma"
  "aspace_publisher/file"
  "aspace_publisher/utils"
)

func AlmaBW(tc_id, filename, session string){
  var args alma.ProcessArgs
  args.Repo_id = "2"
  args.Filename = filename
  var err error
  d_session, err := utils.Decrypt(session)
  if err != nil { file.WriteReport(args.Filename, []string{ "could not decrypt session" + err.Error() }); return }
  args.Session_id = d_session
  
  tc, err := as.GetTC(args.Session_id, args.Repo_id, tc_id)
  if err != nil { file.WriteReport(args.Filename, []string{ "could not obtain top container" + err.Error() }); return}
  tcmap := tc.Mapify()
  if tcmap["ils_holding"] == "" { args.Create = true } else { args.Create = false; args.Holding_id = tcmap["ils_holding"] }
  tccolls, err := tc.CollsWithMmsids(args.Repo_id, args.Session_id)
  if err != nil { file.WriteReport(args.Filename, []string{ "problem obtaining coll mms_ids" + err.Error() }); return}
  fs := alma.BWFunMap{ BibPF: alma.ProcessBWBib, HoldingPF: alma.ProcessBWHolding, ItemPF: alma.ProcessBWItem, Time: alma.GetTime, UpdateTC: as.UpdateTC}
  alma.ProcessBWBib(args, tcmap, tccolls, fs)
  return
}