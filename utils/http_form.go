package utils

import (
    "bytes"
    "io"
    "mime/multipart"
    "os"
    "errors"
    "path/filepath"
    "log/slog"
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
    if err != nil { slog.Error(err.Error()); return form, "", errors.New("MultipartForm/writer error") }
    _, err = formField.Write([]byte(val))
    if err != nil { slog.Error(err.Error()); return form, "", errors.New("MultipartForm/writer error") }
  }
  boundary := writer.Boundary()
  writer.Close()
  return form, boundary, nil
}

func AddFile(writer *multipart.Writer, vals map[string]string)(error){
  fw, err := writer.CreateFormFile(vals["filekey"], filepath.Base(vals["filepath"]))
  if err != nil { slog.Error(err.Error()); return errors.New("AddFile/writer error") }
  fd, err := os.Open(vals["filepath"])
  if err != nil { slog.Error(err.Error()); return errors.New("AddFile/open file error") }
  defer fd.Close()
  _, err = io.Copy(fw, fd)
  if err != nil { slog.Error(err.Error()); return errors.New("AddFile/copy file error") }
  delete(vals, "filekey")
  delete(vals, "filepath")
  return nil
}

// for testing, not called elsewhere
func ReadMultipartFormData(form *bytes.Buffer, boundary string){
  reader := multipart.NewReader(form, boundary)
  for {
    p, err := reader.NextPart()
    if err == io.EOF {
      slog.Error("EOF")
      return
    }
    if err != nil {
      slog.Error(err.Error())
    }
    slurp, err := io.ReadAll(p)
    if err != nil {
      slog.Error(err.Error())
    }
    slog.Info(string(slurp))
  }
}

