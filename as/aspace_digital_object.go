package as

import(
  "log"
  "strings"
  "github.com/tidwall/gjson"
)
// 2 step process: first create the digital object, then pull the archival object and update with the instance
func CreateDigitalObjects(digital_obj_list string, sessionid string) (Responses){
  var r Responses
  items := gjson.Get(digital_obj_list, "digital_objects")
  // returning false from inside the ForEach will stop the loop
  items.ForEach(func(key, value gjson.Result) bool {
    doident := gjson.Get(value.String(), "digital_object_id")
    aoid := extractIdFromInstance(value)
    result := Post(sessionid, doident.String(), "2", "digital_objects", value.String())
    r.responses = append(r.responses, result)
    if strings.Contains(result.response, "error"){ return false }
    doid := extractIdFromResponse(result.response)

    json, err := AcquireJson(sessionid, "2", "archival_objects/" + aoid)
    if err != nil {
      log.Println(err)
      r.responses = append(r.responses, Response{ aoid, err.Error() } )
      return false
    }
    inst := Instance("/repositories/2/digital_objects/" + doid)
    modified, err := UpdateWithInstance(json, inst)
    result = Post(sessionid, aoid, "2", "archival_objects/" + aoid, string(modified))
    r.responses = append(r.responses, result)
    return true
  })
  return r
}

func extractIdFromInstance(val gjson.Result) string {
  arr := gjson.Get(val.String(), "linked_instances")
  long_id := arr.Array()[0].Get("ref")
  short_id := extractIdFromPath(long_id.String())
  return short_id
}

func extractIdFromPath(ref string) string {
  arr := strings.Split(ref, "/")
  return arr[len(arr) - 1]
}

// extract DO id from aspace response, returns "" if not found
func extractIdFromResponse(body string) string {
  id := gjson.Get(body, "id")
  return id.String()
}


