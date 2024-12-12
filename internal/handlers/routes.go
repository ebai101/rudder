package handlers

import (
	"net/http"
	"rudder/internal"
	"rudder/internal/views"
	"rudder/util/template"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(app *internal.Application) {
	app.E.GET("/", func(c echo.Context) error {
		component := views.Index("Rudder")
		return template.AssertRender(c, http.StatusOK, component)
	})
	app.E.GET("/click-me", func(c echo.Context) error {
		component := views.ClickMeBody()
		return template.AssertRender(c, http.StatusOK, component)
	})
}
