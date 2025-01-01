package handlers

import (
	"rudder/internal/models"
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

	ins, err := ih.service.GetInsights(ctx, models.PAST_30_DAYS)
	if err != nil {
		return err
	}

	chartData, err := ih.service.GetInsightsChartData(ctx)
	if err != nil {
		return err
	}

	component := views.Index(ins, chartData)
	return renderView(c, component)
}
