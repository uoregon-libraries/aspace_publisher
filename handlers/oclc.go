package handlers

import (
  "aspace_publisher/as"
  "aspace_publisher/file"
  "aspace_publisher/oclc"
  "aspace_publisher/marc"
)

func validateMarc(resource_id string, repo_id, session_id string, oclc_token string, fm oclc.FunMap) (string, error){
  fname := file.Filename()
  published, err := as.CheckIsPublished(resource_id, repo_id, session_id)
  if err != nil { return published, err}

  marc_rec, err := as.AcquireMarc(session_id, repo_id, resource_id, published)
  if err != nil { return marc_rec, err  }

  //strip outer tag
  marc_stripped,err := marc.StripOuterTags(marc_rec)

  if err != nil {
    file.WriteReport(fname, []string{ marc_stripped, err.Error() })
    return fname, err
  }
  errstr := ""
  oclc_resp, err := oclc.Request(oclc_token, "POST", marc_stripped, "manage/bibs/validate/validateFull", "","json")
  if err != nil { errstr = err.Error() }
  file.WriteReport(fname, []string{ oclc_resp, errstr })
  return fname, err
}

func oclcCrup(resource_id string, repo_id string, session_id string, oclc_token string, fm oclc.FunMap) (string, error){
  fname := file.Filename()
  json, err := as.AcquireJson(session_id, repo_id, "resources/" + resource_id)
  if err != nil {
    file.WriteReport(fname, []string{ string(json), err.Error() })
    return fname, err
  }

  published, err := as.IsPublished(json)
  if err != nil {
    file.WriteReport(fname, []string{ "Problem checking as.IsPublished", err.Error() })
    return fname, err
  }

  marc_rec, err := as.AcquireMarc(session_id, repo_id, resource_id, published)
  if err != nil {
    file.WriteReport(fname, []string{ marc_rec, err.Error() })
    return fname, err
  }

  //strip outer tag
  marc_stripped,err := marc.StripOuterTags(marc_rec)

  if err != nil {
    file.WriteReport(fname, []string{ marc_stripped, err.Error() })
    return fname, err
  }

  //is it a new record?
  oclc_id, err := as.GetOclcId(resource_id, repo_id, session_id)
  if err != nil { 
    file.WriteReport(fname, []string{ oclc_id, err.Error() })
    return fname, err
  }
  var oclc_resp string
  //Update
  if oclc_id != "" {
    oclc_marc, err_ := oclc.Record(oclc_token, oclc_id)
    if err_ != nil{
      file.WriteReport(fname, []string{ oclc_marc, err.Error() })
      return fname, err
    }
    edited_marc, err_ := marc.EditMarcForOCLC(oclc_marc, marc_stripped)
    if err_ != nil{
      file.WriteReport(fname, []string{ edited_marc, err.Error() })
      return fname, err
    }

    oclc_resp, err = oclc.Request(oclc_token, "PUT", edited_marc, "manage/bibs", oclc_id, "marcxml+xml")
  } else {
    oclc_resp, err = oclc.Request(oclc_token, "POST", marc_stripped, "manage/bibs", "", "marcxml+xml")
  }
  if err != nil {
    file.WriteReport(fname, []string{ oclc_resp, err.Error() })
    return fname, err
  }

  //if updating, done
  if oclc_id != "" {
    file.WriteReport(fname, []string{ oclc_resp })
    return fname, nil
  }

  oclc_id, err = marc.ExtractOclc(string(oclc_resp))
  if err != nil {
    file.WriteReport(fname, []string{ "Could not extract oclc id", err.Error() })
    return fname, err
  }
  //insert oclc
  modified, err := as.UpdateUserDefined1(json, oclc_id)
  if err != nil { 
    file.WriteReport(fname, []string{ string(modified), err.Error() })
    return fname, err
  }
  //post resource json back to aspace
  as_resp := as.Post(session_id, resource_id, repo_id, "resources/" + resource_id, string(modified))
  file.WriteReport(fname, []string{ as_resp.ResponseToString() })
  return fname, nil
}
