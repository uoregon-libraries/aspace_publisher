package handlers

import (
  "github.com/labstack/echo/v4"
  "net/http"
  "aspace_publisher/alma"
)

func StatusHandler(c echo.Context) error {
  alma.CallWorker("startStatusJob", map[string]string{ "status": c.Param("status"), "other": "hippo" })
  return c.HTML(http.StatusOK, "test output available in log")
}
