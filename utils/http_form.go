package utils

import (
    "bytes"
    "io"
    "mime/multipart"
    "os"
    "errors"
    "path/filepath"
    "log"
)

// expects a map of string keys and vals
// this method meant to be used for uploading a file
// map should include filekey => fieldname for the file
// map should include filepath => path to file being uploaded
func CreateMultipartFormData(vals map[string]string)(*bytes.Buffer, string, error){
  form := new(bytes.Buffer)
  writer := multipart.NewWriter(form)
  err := AddFile(writer, vals)
  if err != nil { return form, "", err }
  for key, val := range vals {
    formField, err := writer.CreateFormField(key)
    if err != nil { log.Println(err); return form, "", errors.New("MultipartForm/writer error") }
    _, err = formField.Write([]byte(val))
    if err != nil { log.Println(err); return form, "", errors.New("MultipartForm/writer error") }
  }
  boundary := writer.Boundary()
  writer.Close()
  return form, boundary, nil
}

func AddFile(writer *multipart.Writer, vals map[string]string)(error){
  fw, err := writer.CreateFormFile(vals["filekey"], filepath.Base(vals["filepath"]))
  if err != nil { log.Println(err); return errors.New("AddFile/writer error") }
  fd, err := os.Open(vals["filepath"])
  if err != nil { log.Println(err); return errors.New("AddFile/open file error") }
  defer fd.Close()
  _, err = io.Copy(fw, fd)
  if err != nil { log.Println(err); return errors.New("AddFile/copy file error") }
  delete(vals, "filekey")
  delete(vals, "filepath")
  return nil
}

// mainly here for testing
func ReadMultipartFormData(form *bytes.Buffer, boundary string){
  reader := multipart.NewReader(form, boundary)
  for {
    p, err := reader.NextPart()
    if err == io.EOF {
      log.Print("EOF")
      return
    }
    if err != nil {
      log.Fatal(err)
    }
    slurp, err := io.ReadAll(p)
    if err != nil {
      log.Fatal(err)
    }
    log.Println(string(slurp))
  }
}

