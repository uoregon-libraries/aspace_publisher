package main

import (
  "aspace_publisher/handlers"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "os"
)

func main(){
  e := echo.New()
  // Middleware
  e.Use(middleware.Logger())
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

  e.Logger.Fatal(e.Start(os.Getenv("PORT")))

}

