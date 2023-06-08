package utils

import(
  "net/http"
  "time"
  "github.com/labstack/echo/v4"
)

func WriteCookie(c echo.Context, name string, value string) {
  cookie := new(http.Cookie)
  cookie.Name = name
  cookie.Value = value
  cookie.Expires = time.Now().Add(5 * time.Hour)
  c.SetCookie(cookie)
}

func FetchCookieVal(c echo.Context, name string) (string, error) {
  cookie, err := c.Cookie(name)
  if err != nil { return "", err }

  return cookie.Value, nil
}
