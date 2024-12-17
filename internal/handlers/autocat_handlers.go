package handlers

import (
	"rudder/internal/services"
	"rudder/internal/views"

	"github.com/labstack/echo/v4"
)

type AutocatHandlers struct {
	service *services.AutocatService
}

func NewAutocatHandlers(
	service *services.AutocatService,
) *AutocatHandlers {
	return &AutocatHandlers{
		service: service,
	}
}

func (ah *AutocatHandlers) acatListHandler(c echo.Context) error {
	c.Set("ISERROR", false)
	ctx := c.Request().Context()

	accs, err := ah.service.GetAutocatRules(ctx)
	if err != nil {
		return err
	}

	component := views.Autocat(accs)
	return renderView(c, component)
}
