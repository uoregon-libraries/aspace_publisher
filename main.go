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
  e.GET("/version", handlers.VersionHandler)
  e.Static("/as", "views/as") // as/login.html
  e.POST("login", handlers.AspaceLoginHandler)
  e.GET("/ead/validate/:id", handlers.ValidateEadHandler)
  e.GET("/ead/convert/:id", handlers.ConvertEadHandler)
  e.GET("/ead/upload/:id", handlers.UploadEadHandler)
  e.GET("/oclc/crup/:id", handlers.OclcCrupHandler)
  e.GET("/oclc/validate/:id", handlers.OclcValidateHandler)

  e.Logger.Fatal(e.Start(os.Getenv("PORT")))

}

