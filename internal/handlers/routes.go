package handlers

import (
	"net/http"
	"rudder/internal"
	"rudder/internal/views"
	"rudder/util/template"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderHTMX(c echo.Context, component templ.Component) error {
	isHtmxRequest := c.Request().Header.Get("HX-Request") == "true"

	if isHtmxRequest {
		return template.AssertRender(c, http.StatusOK, component)
	}
	return template.AssertRender(c, http.StatusOK, views.FullPage(component))
}

func RegisterRoutes(app *internal.Application) {
	app.E.GET("/", func(c echo.Context) error {
		component := views.Index(template.ChartComponent(BarChart()))
		return RenderHTMX(c, component)
	})

	app.E.GET("/transactions", func(c echo.Context) error {
		ctx := c.Request().Context()

		txns, err := app.F.GetTransactions(ctx, 50)
		if err != nil {
			c.Logger().Error(err)
		}

		component := views.Transactions(txns)
		return RenderHTMX(c, component)
	})

	app.E.GET("/accounts", func(c echo.Context) error {
		ctx := c.Request().Context()

		accs, err := app.F.GetAccounts(ctx)
		if err != nil {
			c.Logger().Error(err)
		}

		component := views.Accounts(accs)
		return RenderHTMX(c, component)
	})

	app.E.GET("/autocat", func(c echo.Context) error {
		ctx := c.Request().Context()

		rules, err := app.AC.GetAutocatRules(ctx)
		if err != nil {
			c.Logger().Error(err)
		}

		component := views.Autocat(rules)
		return RenderHTMX(c, component)
	})
}
