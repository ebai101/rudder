package handlers

import (
	"rudder/internal/services"
	"rudder/internal/views"

	"github.com/labstack/echo/v4"
)

type InsightsHandlers struct {
	service *services.FinancialService
}

func NewInsightsHandlers(
	service *services.FinancialService,
) *InsightsHandlers {
	return &InsightsHandlers{
		service: service,
	}
}

func (ih *InsightsHandlers) insightsHandler(c echo.Context) error {
	c.Set("ISERROR", false)
	ctx := c.Request().Context()

	ins, err := ih.service.GetInsights(ctx)
	if err != nil {
		return err
	}

	component := views.Index(ins)
	return renderView(c, component)
}
