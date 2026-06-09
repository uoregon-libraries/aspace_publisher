package as

import(
  "log"
  "strings"
  "github.com/tidwall/gjson"
  "github.com/tidwall/sjson"
  "regexp"
)
// 2 step process: first create the digital object, then pull the archival object and update with the instance
func CreateDigitalObjects(digital_obj_list string, sessionid string) (Responses){
  var r Responses
  items := gjson.Get(digital_obj_list, "digital_objects")
  // returning false from inside the ForEach will stop the loop
  // value is the raw digital object json
  // modify value for final version to submit to aspace
  items.ForEach(func(key, value gjson.Result) bool {
    doident := value.Get("digital_object_id")
    digital_obj, err := modify(value)
    if err != nil {
      r.responses = append(r.responses, Response{ doident.String(), BuildErrorMessage(err.Error()) } )
      return true
    }

    short_ref := validateRefId(extractRefPathFromString(digital_obj))
    if short_ref == "" {
      r.responses = append(r.responses, Response{ doident.String(), BuildErrorMessage("no valid reference id") } )
      return true
    }
    // if ref does not exist, request will 404, skip and go on
    json, err := AcquireJson(sessionid, "2", short_ref)
    if err != nil {
      log.Println(err)
      r.responses = append(r.responses, Response{ short_ref, BuildErrorMessage(err.Error()) } )
      return true
    }

    result := Post(sessionid, doident.String(), "2", "digital_objects", digital_obj)
    r.responses = append(r.responses, result)
    if strings.Contains(result.ResponseToString(), "error"){ return true }
    doid := extractIdFromResponse(result.ResponseToString())
    if doid == "" {
      log.Println("failed to extract doid from aspace response for " + short_ref)
      return true }
    inst := Instance("/repositories/2/digital_objects/" + doid)
    modified, err := UpdateWithInstance(json, inst)
    result = Post(sessionid, short_ref, "2", short_ref, string(modified))
    r.responses = append(r.responses, result)
    return true
  })
  return r
}

// currently the only mod needed is the ref url
func modify(val gjson.Result)(string, error){
  temp_ref := extractRefPathFromResult(val)
  final_ref := strings.Replace(temp_ref, "x", "2", 1)
  dig_obj, err := sjson.Set(val.String(),"linked_instances.0.ref", final_ref)
  return dig_obj, err
}

func extractRefPathFromResult(val gjson.Result) string {
  ref_path := val.Get("linked_instances.0.ref")
  return ref_path.String()
}

func extractRefPathFromString(value string) string{
  ref_path := gjson.Get(value, "linked_instances.0.ref")
  return ref_path.String()
}

func validateRefId(refid string) string {
  re1 := regexp.MustCompile(`archival_objects/[0-9]+`)
  matched1 := re1.Find([]byte(refid))
  if matched1 != nil { return string(matched1) }
  re2 := regexp.MustCompile(`resources/[0-9]+`)
  matched2 := re2.Find([]byte(refid))
  if matched2 != nil { return string(matched2) }
  return ""
}

func extractIdFromPath(ref string) string {
  arr := strings.Split(ref, "/")
  return arr[len(arr) - 1]
}

// extract DO id from aspace response, returns "" if not found
func extractIdFromResponse(body string) string {
  id := gjson.Get(body, "message.id")
  return id.String()
}
