package utils

import(
  "net/http"
  "time"
  "github.com/labstack/echo/v4"
)

func WriteCookie(c echo.Context, expires int, name string, value string) {
  cookie := new(http.Cookie)
  cookie.Name = name
  cookie.Value = value
  cookie.Expires = time.Now().Add(time.Duration(expires) * time.Hour)
  c.SetCookie(cookie)
}

func FetchExpirableCookie(c echo.Context, name string) (string, error){
  if Expired(c, name) { return "", nil }
  cookie, err := c.Cookie(name)
  if err != nil { return "", err }

  return cookie.Value, nil
}

func FetchCookieVal(c echo.Context, name string) (string, error) {
  cookie, err := c.Cookie(name)
  if err != nil { return "", err }

  return cookie.Value, nil
}

func Expired(c echo.Context, name string) (bool) {
  cookie, err := c.Cookie(name)
  if err != nil { return true }
  if cookie.Expires.After(time.Now()) {
    return false
  } else { return true }
}
