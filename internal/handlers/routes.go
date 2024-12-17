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
	ih *InsightsHandlers,
	th *TransactionsHandlers,
	ah *AccountsHandlers,
	ch *CategoriesHandlers,
) {
	e.GET("/", ih.insightsHandler)
	e.GET("/transactions", th.txnsMainHandler)
	e.GET("/transactions/:page", th.txnsScrollHandler)
	e.GET("/transactions/detail/:id", th.txnDetailsHandler)
	e.GET("/transactions/detail/:id/edit", th.txnDetailsEditHandler)
	e.GET("/accounts", ah.accsListHandler)
	e.GET("/accounts/:id", ah.accsDetailHandler)
	e.GET("/accounts/:id/transactions", ah.accsTransactionsHandler)
	e.GET("/accounts/navbar", ah.accsNavbarHandler)
	e.GET("/categories", ch.acatRuleListHandler)
}
