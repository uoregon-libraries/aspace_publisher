package handlers

import(
  "github.com/labstack/echo/v4"
  "net/http"
  "aspace_publisher/as"
  "aspace_publisher/utils"
)

func AspaceLoginHandler(c echo.Context) error {
  name := c.FormValue("name")
  password := c.FormValue("password")
  session_id, err := as.AuthenticateAS(name, password)
  if err != nil { return echo.NewHTTPError(400, err.Error()) }
  utils.WriteCookie(c, 60, "as_session", session_id)
  return c.String(http.StatusOK, "You have logged in.")
}

