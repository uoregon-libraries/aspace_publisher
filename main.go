package main

import (
  "aspace_publisher/handlers"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  slogecho "github.com/samber/slog-echo"
  "github.com/DeRuina/timberjack"
  "os"
  "log"
  "log/slog"
)

func main(){
  e := echo.New()
  logmode := os.Getenv("LOGMODE")
  var logging *slog.Logger
  if logmode == "file" {
    logr := &timberjack.Logger{
    Filename:   "logs/app.log", // path of log file
    MaxSize:    50, // file size in MB
    MaxBackups: 7, // number of files to retain
    MaxAge:     8, // how long (in days) to retain files
    Compression: "gzip", // archive files?
    LocalTime:  true, // re timestamps
    RotateAt: []string{"00:00"},
}
    logging = slog.New(slog.NewJSONHandler(logr, nil))
  } else { logging = slog.New(slog.NewJSONHandler(os.Stdout, nil)) }
  slog.SetDefault(logging)

  e.Use(slogecho.New(logging))
  e.Use(middleware.Recover())

  path := os.Getenv("HOME_DIR")
  e.GET("/version", handlers.VersionHandler)
  e.File("/as/login.html", path + "views/as/login.html") // as/login.html
  e.POST("login", handlers.AspaceLoginHandler)
  e.GET("/ead/validate/:id", handlers.ValidateEadHandler)
  e.GET("/ead/convert/:id", handlers.ConvertEadHandler)
  e.GET("/ead/upload/:id", handlers.UploadEadHandler)
  e.GET("/oclc/crup/:id", handlers.OclcCrupHandler)
  e.GET("/oclc/validate/:id", handlers.OclcValidateHandler)
  e.File("/as/do.html", path + "views/as/do.html") //urlpath,directorypath, uploads/do.html
  e.POST("/upload_do", handlers.UploadDigitalObjectsHandler)
  e.Static("/reports", "views/reports")
  e.GET("/alma/crup/:id", handlers.AlmaCrupHandler)
  e.GET("/test/status/:status", handlers.StatusHandler)
  log.Fatal(e.Start(os.Getenv("PORT")))

}

