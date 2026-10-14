package file

import(
  "bufio"
  "math/rand"
  "os"
  "fmt"
  "log/slog"
  "strings"
  "time"
)

func buildFromFile(filename string)[]string{
  str_arr := []string{}
  f,_ := os.Open(filename)
  defer f.Close()
  scanner := bufio.NewScanner(f)
  for scanner.Scan(){
    arr := strings.Split(scanner.Text(), ",")
    str_arr = append(str_arr, arr...)
  }
  return str_arr
}

func Filename()string{
  homedir := os.Getenv("HOME_DIR")
  nouns := buildFromFile(homedir + "/file/nouns.txt")
  mods := buildFromFile(homedir + "/file/modifiers.txt")
  t := TimeNow()
  date1 := t.Format("2006-01-02-15-04")

  s1 := mods[rand.Intn(len(mods))]
  s2 := nouns[rand.Intn(len(nouns))]
  return fmt.Sprintf("%s_%s-%s", date1, s1, s2)
}

func WriteReport(filename string, messages []string){
  report_dir := os.Getenv("REPORT_DIR")
  path := report_dir + "/" + filename
  f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  defer f.Close()
  if err != nil { slog.Error(err.Error()) }
  for _,v := range messages{
    _, err = fmt.Fprintln(f, v)
    if err != nil { slog.Error(err.Error()) }
  }
}
func TimeNow()time.Time{
  loc, _ := time.LoadLocation("America/Los_Angeles")
  t := time.Now().In(loc)
  return t
}


