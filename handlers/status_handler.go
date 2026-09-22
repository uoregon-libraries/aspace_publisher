package handlers

import (
  "github.com/labstack/echo/v4"
  "net/http"
  "aspace_publisher/alma"
  "aspace_publisher/utils"
  "log"
)

func StatusHandler(c echo.Context) error {
  log.Println(c.Param("status"))
  e_other, err := utils.Encrypt("abcdefghijklmnopqrstuvwxyz0123456789")
  if err != nil { return c.HTML(400, "could not encrypt other") }
  err = alma.CallWorker("startStatusJob", map[string]string{ "status": c.Param("status"), "other": e_other })
  if err != nil { return c.HTML(400, err.Error()) }
  return c.HTML(http.StatusOK, c.Param("status"))
}
