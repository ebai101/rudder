package handlers

import (
	"net/http"
	"rudder/internal/views"
	"rudder/util/template"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func renderView(c echo.Context, component templ.Component) error {
	isHtmxRequest := c.Request().Header.Get("HX-Request") == "true"

	if isHtmxRequest {
		return template.AssertRender(c, http.StatusOK, component)
	}
	return template.AssertRender(c, http.StatusOK, views.FullPage(component))
}

func RegisterRoutes(
	e *echo.Echo,
	th *TransactionsHandlers,
	ah *AccountsHandlers,
	ach *AutocatHandlers,
) {
	e.GET("/", func(c echo.Context) error {
		component := views.Index()
		return renderView(c, component)
	})
	e.GET("/transactions", th.txnsMainHandler)
	e.GET("/transactions/:page", th.txnsScrollHandler)
	e.GET("/transactions/detail/:id", th.txnDetailsHandler)
	e.GET("/accounts", ah.accsListHandler)
	e.GET("/autocat", ach.acatListHandler)
}
