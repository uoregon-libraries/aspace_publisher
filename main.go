package main

import (
  "aspace_publisher/as"
  "aspace_publisher/handlers"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
)


func main() {
  e := echo.New()
  // Middleware
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  e.Use(middleware.BasicAuth(as.As_basic))

  e.GET("/ead/validate/:id", handlers.TestValidateHandler)
  e.GET("/ead/convert/:id", handlers.ConvertEadHandler)
  e.GET("/ead/upload/:id", handlers.UploadEadHandler)
  e.Logger.Fatal(e.Start(":3000"))
}

