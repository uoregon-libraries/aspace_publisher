package handlers

import (
  "aspace_publisher/file"
  "aspace_publisher/as"
  "aspace_publisher/aw"
)

//resource_id, session_id, repo_id, aw_session, operation
func processEad(resource_id string, session_id string, repo_id string, aw_session string, operation string) (string, error){
  fname := file.Filename()
  ead_orig, err := as.AcquireEad(session_id, repo_id, resource_id)
  if err != nil {
    file.WriteReport(fname, []string{ "Could not aquire JSON from aspace: ", err.Error() })
    return fname, err
  }
  result,err := aw.ProcessEad(repo_id, resource_id, ead_orig, aw_session, operation)
  if err != nil {
    file.WriteReport(fname, []string{ "Error processing EAD", err.Error() })
    return fname, err
  }
  file.WriteReport(fname, []string{ result })
  return fname, nil
}
