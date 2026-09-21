package as

import (
  "encoding/json"
  "fmt"
  "aspace_publisher/file"
  "strings"
)
//event_type: EAD published; EAD revised, MARC published, MARC updated
func ProcessEvent(session, agent, resource_id, repo_id, event_type string) Response{
  resource := fmt.Sprintf("/repositories/%s/resources/%s", repo_id, resource_id)
  event := ConstructEvent(agent, resource, event_type)
  body, err := json.Marshal(event)
  if err != nil {}
  response := Post(session, resource_id + "-event", repo_id, "events", string(body))
  return response
}

func ConstructEvent(agent, resource, event_type string) Event{
  var event Event
  event.ModelType = "event"
  event.Agents = []Rec{ Rec{ Ref: agent, Role: role(event_type) } }
  event.Records = []Rec{ Rec{ Ref: resource, Role: "Source" } }
  event.Date = ConstructDate()
  event.Outcome = "Pass"
  event.EventType = event_type
  return event
}

func role(event_type string) string{
  if strings.Contains(event_type, "EAD") { return "uploaded_by" }
  return "cataloger"
}

func ConstructDate()Date{
  var date Date
  date.ModelType = "date"
  date.DateType = "Single"
  date.Label = "Agent Relation"
  date.Begin = ConstructTime()
  return date
}

func ConstructTime()string{
  t := file.TimeNow()
  y := t.Format("2006")
  m := t.Format("01")
  d := t.Format("02")
  return fmt.Sprintf("%s-%s-%s", y,m,d)
}
type Event struct{
  ModelType string `json:"jsonmodel_type"`
  Agents []Rec `json:"linked_agents"`
  Records []Rec `json:"linked_records"`
  Date Date `json:"date"`
  EventType string `json:"event_type"`
  Outcome string `json:"outcome"`
}

type Date struct{
  ModelType string `json:"jsonmodel_type"`
  DateType string `json:"date_type"`
  Label string `json:"label"`
  Begin string `json:"begin"`
}

type Rec struct{
  Ref string `json:"ref"`
  Role string`json:"role"`
}

