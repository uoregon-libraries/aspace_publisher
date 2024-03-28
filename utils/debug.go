package utils

import(
  "os"
)

func WriteFile(name string, to_write string) error {
  f, err := os.CreateTemp("", name)
  if err != nil { return err }
  _, err = f.Write([]byte(to_write))
  f.Close()
  return err
}
