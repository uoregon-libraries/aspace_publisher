package handlers

import (
  "aspace_publisher/file"
)

func validateEad(resource_id string, session_id string) (string, error){
  fname := file.Filename()
  return fname, nil
}

func uploadEad(resource_id string, session_id string) (string, error){
  fname := file.Filename()
  return fname, nil
}
