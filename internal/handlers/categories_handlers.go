package handlers

import (
	"rudder/internal/services"
	"rudder/internal/views"

	"github.com/labstack/echo/v4"
)

type CategoriesHandlers struct {
	service *services.CategoriesService
}

func NewCategoriesHandlers(
	service *services.CategoriesService,
) *CategoriesHandlers {
	return &CategoriesHandlers{
		service: service,
	}
}

func (ah *CategoriesHandlers) acatRuleListHandler(c echo.Context) error {
	c.Set("ISERROR", false)
	ctx := c.Request().Context()

	accs, err := ah.service.GetAutocatRules(ctx)
	if err != nil {
		return err
	}

	component := views.Autocat(accs)
	return renderView(c, component)
}
