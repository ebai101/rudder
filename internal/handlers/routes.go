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
		ctx := c.Request().Context()

		txns, err := app.F.GetTransactions(ctx, 10)
		if err != nil {
			c.Logger().Error(err)
		}

		component := views.Transactions(txns)
		return template.AssertRender(c, http.StatusOK, component)
	})
	app.E.GET("/click-me", func(c echo.Context) error {
		component := views.ClickMeBody()
		return template.AssertRender(c, http.StatusOK, component)
	})
}
