package handlers

import (
  "github.com/labstack/echo/v4"
  "net/http"
  "aspace_publisher/alma"
  "aspace_publisher/utils"
)

func StatusHandler(c echo.Context) error {
  e_other, err := utils.Encrypt("abcdefghijklmnopqrstuvwxyz0123456789")
  if err != nil { return c.HTML(400, "could not encrypt other") }
  alma.CallWorker("startStatusJob", map[string]string{ "status": c.Param("status"), "other": e_other })
  return c.HTML(http.StatusOK, "test output available in log")
}
