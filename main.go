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

  e.GET("/as_ead/:id", handlers.Get_ead_handler)
  e.GET("/prep_ead/:id", handlers.Prep_ead_handler)
  e.GET("/convert_ead/:id", handlers.Php_ead_handler)

  e.Logger.Fatal(e.Start(":3000"))
}

