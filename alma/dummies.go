package alma

import(
  "fmt"
  "log"
)
func DummyLinkToNetwork(list []string, filename string){
  fmt.Println("from dummy link")
  fmt.Println("ID: " + list[0])
  fmt.Println("Filename: " + filename)
}

func DummyBoundwithPF(args ProcessArgs, marc_string string, tcmap []map[string]string, fs FunMap){ return }
func DummyHoldingPF(args ProcessArgs, marc_string string, tcmap []map[string]string, fs FunMap){
  if args.Holding_id != "" { log.Fatal("incorrect holding set") }
  return
}
func DummyItemsPF(args ProcessArgs, tcmap []map[string]string, fs FunMap){ return }
func DummyItemPF(args ProcessArgs, item Item, tcmap map[string]string)(string, error){
  return "456745674567", nil
}
func DummyNZPF(list []string, filename string){ return }
func DummyAfterBib(rjson []byte, args_map map[string]string)error{
  if args_map["mms_id"] != "654365436543" { log.Fatal("incorrect mms_id") }
  return nil
 }
func DummyFetchBibID(barcode string)string{
  if barcode != "123412341234" { log.Fatal("incorrect barcode sent") }
  return "234523452345"
}
func DummyUpdateTC(repo_id string, holding_id string, item_id string, session_id string, tcmap map[string]string)error{

  if item_id != "456745674567" { log.Fatal("incorrect value sent to DummyUpdateTC") }
  return nil
}
func DummySetHolding(oclc_id string, token string)(string, error){ return fmt.Sprintf("holding %s is set", oclc_id), nil }

