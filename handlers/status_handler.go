package handlers

import (
  "github.com/labstack/echo/v4"
  "net/http"
  "aspace_publisher/alma"
  "log"
)

func StatusHandler(c echo.Context) error {
  log.Println(c.Param("status"))
  err := alma.CallWorker("startStatusJob", map[string]string{ "status": c.Param("status"), "other": "hippo" })
  if err != nil { return c.HTML(400, err.Error()) }
  return c.HTML(http.StatusOK, c.Param("status"))
}
