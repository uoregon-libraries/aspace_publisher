package marc

import(

  "testing"
  "github.com/beevik/etree"
  "os"
  "path/filepath"
)

func TestExtractOclc(t *testing.T){
  home_dir := os.Getenv("HOME_DIR")
  oclc_marc,_ := os.ReadFile(filepath.Join(home_dir, "fixtures/oclc_marc_3464.xml"))
  id,_ := ExtractOclc(string(oclc_marc))
  if id != "1097882240" { t.Fatalf("id is not correct") }
}

func TestStripOuterTags(t *testing.T){
  home_dir := os.Getenv("HOME_DIR")
  as_marc,_ := os.ReadFile(filepath.Join(home_dir, "fixtures/marc_3464.xml"))
  stripped,_ := StripOuterTags(string(as_marc))
  marc_tree := etree.NewDocument()
  marc_tree.ReadFromString(stripped)
  if marc_tree.Root().FullTag() != "record" { t.Fatalf("did not strip outer tag") }
}

func TestEdit000(t *testing.T){
  as_l := "00000npcaa2200000 i 4500"
  oclc_l := "00000cpc a22000007i 4500"
  edited := edit_status(as_l, oclc_l)
  if edited != "00000cpcaa2200000 i 4500" { t.Fatalf("leader is not correct") }
}

func TestEdit008(t *testing.T){
  as_008_t := "240910i19752006xxu                 eng d"
  oclc_008_t := "190419i19752006oru                 eng d"
  edited := edit_008(as_008_t, oclc_008_t)
  if edited != "190419i19752006xxu                 eng d" {t.Fatalf("008 is not correct") }
}

func TestEditMarcForOCLC(t *testing.T){
  home_dir := os.Getenv("HOME_DIR")
  oclc_marc,_ := os.ReadFile(filepath.Join(home_dir, "fixtures/oclc_marc_3464.xml"))
  as_marc,_ := os.ReadFile(filepath.Join(home_dir, "fixtures/marc_3464.xml"))
  edited, _ := EditMarcForOCLC(string(oclc_marc), string(as_marc))
  edit_tree := etree.NewDocument()
  _ = edit_tree.ReadFromString(edited)
  elts := edit_tree.FindElements("//datafield[@tag='040']/subfield[@code='d']")
  if len(elts) != 3 { t.Fatalf("did not insert all of the subfields") }
  if elts[0].Text() != "OCLCF" { t.Fatalf("did not insert the subfields in order") }
  leader := edit_tree.FindElement("//leader")
  if leader.Text() != "00000cpcaa2200000 i 4500" { t.Fatalf("leader is not correct") }
  oo8 := edit_tree.FindElement("//controlfield[@tag='008']")
  if oo8.Text() != "190419i19752006xxu                 eng d" {t.Fatalf("008 is not correct") }
}
