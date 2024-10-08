package aw

import (
  "fmt"
  "strings"
  "errors"
  "os/exec"
  "os"
  "path/filepath"
  "github.com/beevik/etree"
  "log"
)

// needs refactoring, but for now...
func PrepareEad(repo_id string, resource_id string, xml string)(string, string, string, error){
  aw_xml := etree.NewElement("ead")
  as_xml := etree.NewDocument()
  err := as_xml.ReadFromString(xml)
  if err != nil { log.Println(err); return "","", "", errors.New("could not read EAD") }

  eadheader_copy := as_xml.FindElement("//eadheader").Copy()
  eadheader_copy.RemoveAttr("findaidstatus")
  eadid := eadheader_copy.FindElement("//eadid")
  eadid.CreateAttr("encodinganalog", "identifier")
  ark := eadid.SelectAttrValue("url","")
  split_ark := strings.Split(ark,"ark:")[1]
  ark_id := strings.TrimPrefix(split_ark, "/")
  eadid.CreateAttr("identifier", ark_id)
  extptr := eadheader_copy.FindElement("//extptr")
  addressline := extptr.Parent()
  url := extptr.SelectAttrValue("xlink:href", "")
  addressline.SetText(url)
  addressline.RemoveChild(extptr)
  aw_xml.AddChild(eadheader_copy)

  control := as_xml.FindElement("//control")
  if control != nil {
    control_copy := control.Copy()
    aw_xml.AddChild(control_copy)
  }

  as_archdesc := as_xml.FindElement("//archdesc")
  aw_archdesc := aw_xml.CreateElement("archdesc")
  for _, attr := range as_archdesc.Attr {
    aw_archdesc.CreateAttr(attr.Key, attr.Value)
  }
  children := []string{"did", "accessrestrict", "controlaccess", "otherfindaid", "bioghist", "scopecontent"}
  for _, child := range children {
    ch := as_xml.FindElement("//archdesc/" + child)
    if ch != nil {
      aw_archdesc.AddChild(ch.Copy())
    }
  }

  archdesc_title := aw_xml.FindElement("//archdesc/did/unittitle").Text()
  resource_uri := fmt.Sprintf("https://scua.uoregon.edu/repositories/%s/resources/%s",repo_id, resource_id)
  unittitle := etree.NewElement("unittitle")
  extref := unittitle.CreateElement("extref")
  extref.SetText(archdesc_title)
  cleantitle := strings.Replace(archdesc_title, " ", "-", -1)
  attrs := map[string]string{"title": cleantitle, "show": "new", "href": resource_uri, "actuate": "onrequest"}
  for key, value := range attrs {
    extref.CreateAttr(key, value)
  }
  i := aw_xml.FindElement("//archdesc/did/unittitle").Index()
  did := aw_xml.FindElement("//archdesc/did")
  did.RemoveChildAt(i)
  did.InsertChildAt(i, unittitle)
  unitid := aw_xml.FindElement("//archdesc/did/unitid[@type='aspace_uri']")
  if unitid != nil {
    i := unitid.Index()
    did.RemoveChildAt(i)
  }

  filedesc_title := aw_xml.FindElement("//eadheader/filedesc/titlestmt/titleproper").Text()
  dsc := aw_archdesc.CreateElement("dsc")
  dsc.CreateAttr("type", "analyticover")
  c01 := dsc.CreateElement("c01")
  c01.CreateAttr("level", "otherlevel")
  c01.CreateAttr("otherlevel", "Heading")
  did = c01.CreateElement("did")
  unittitle = did.CreateElement("unittitle")
  extref = unittitle.CreateElement("extref")
  extref.SetText(filedesc_title)
  cleantitle = strings.Replace(filedesc_title, " ", "-", -1)
  attrs["title"] = cleantitle
  for key, value := range attrs {
    extref.CreateAttr(key, value)
  }

  otherfindaid := aw_archdesc.CreateElement("otherfindaid")
  p := otherfindaid.CreateElement("p")
  extref = p.CreateElement("extref")
  extref.SetText("See the Current Collection Guide for detailed description and requesting options.")
  cleantitle = "see-current-collection-guide-and-requesting-options"
  attrs = map[string]string{"title": cleantitle, "show": "new", "href": resource_uri, "actuate": "onrequest"}
  for key, value := range attrs {
    extref.CreateAttr(key, value)
  }

  d := etree.NewDocumentWithRoot(aw_xml)
  s, err := d.WriteToString()
  if err != nil { log.Println(err); return "", "", "", errors.New("could not write EAD to string") }
return s, eadid.Text(), ark_id, nil
}

func CallConversion(xml string)(string, error){
  home_dir := os.Getenv("HOME_DIR")
  converter_location := filepath.Join(home_dir, "aw/converter.php")
  cmd := exec.Command("php", converter_location, xml)
  var out strings.Builder
  cmd.Stdout = &out
  cmd.Run()
  if strings.Contains(out.String(), "error"){
    log.Println(out.String()); return "", errors.New("unable to complete the awest conversion")
  }
  return out.String(), nil
}
