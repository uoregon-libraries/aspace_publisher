package handlers

import(
  "github.com/labstack/echo/v4"
  "aspace_publisher/as"
)

func LauncherHandler(c echo.Context) error {
  resource_id := c.FormValue("resource_id")
  err := as.ValidID(resource_id)
  if err != nil { return c.String(400, "bad resource id") }
  workflow := c.FormValue("workflow")
  switch workflow {
  case "validate_ead":
    return c.Redirect(302, "/ead/validate/" + resource_id)
  case "upload_ead":
    return c.Redirect(302, "/ead/upload/" + resource_id)
  case "validate_marc":
    return c.Redirect(302, "/oclc/validate/" + resource_id)
  case "publish_marc":
    return c.Redirect(302, "/oclc/crup/" + resource_id)
  case "publish_alma":
    return c.Redirect(302, "/alma/crup/" + resource_id)
  default:
    return c.String(400, "No workflow submitted")
  }
}
