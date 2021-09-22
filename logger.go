package main

import (
  "log"
  "os"
  "sync"
)
var logger *log.Logger
var once sync.Once

func GetLogger() *log.Logger {
    once.Do(func() { logger = createLogger() })
  return logger
}

func createLogger() *log.Logger{
  file, _ := os.OpenFile("aspace-export-log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  return log.New(file, "aspace-export", log.Lshortfile)
}
