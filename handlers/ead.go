package handlers

import (
  "aspace_publisher/file"
  "aspace_publisher/as"
  "aspace_publisher/aw"
  "errors"
  "fmt"
)

//resource_id, session_id, agent, repo_id, aw_session, operation
//assumes that ark id has been generated and set
func processEad(resource_id, session_id, agent, repo_id, aw_session, operation string) (string, error){
  fname := file.Filename()
  ead_orig, err := as.AcquireEad(session_id, repo_id, resource_id)
  if err != nil {
    file.WriteReport(fname, []string{ "Could not aquire JSON from aspace: ", err.Error() })
    return fname, err
  }

  ark, err := aw.ExtractArk(ead_orig)
  if err != nil { return "", err }
  if aw.ValidArk(ark) != true { return "", errors.New("arkid not valid") }
  exists, err := aw.CheckArk(ark)
  if err != nil { return "", errors.New(fmt.Sprintf("could not load arkid: %v", err.Error())) }

  result,err := aw.ProcessEad(repo_id, resource_id, ead_orig, exists, aw_session, operation)
  if err != nil {
    file.WriteReport(fname, []string{ "Error processing EAD", err.Error() })
    return fname, err
  }
  file.WriteReport(fname, []string{ result })

  var event_type string
  if exists { event_type = "EAD revised" } else { event_type = "EAD published" }
  response := as.ProcessEvent(session_id, agent, resource_id, repo_id, event_type)
  file.WriteReport(fname, []string{ fmt.Sprintf("event posted: %v", response.ResponseToString()) })
  return fname, nil
}
