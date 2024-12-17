package handlers

import (
	"rudder/internal/services"
	"rudder/internal/views"

	"github.com/labstack/echo/v4"
)

type AccountsHandlers struct {
	service *services.FinancialService
}

func NewAccountsHandlers(
	service *services.FinancialService,
) *AccountsHandlers {
	return &AccountsHandlers{
		service: service,
	}
}

func (ah *AccountsHandlers) accsListHandler(c echo.Context) error {
	c.Set("ISERROR", false)
	ctx := c.Request().Context()

	accs, err := ah.service.GetAccountRows(ctx)
	if err != nil {
		return err
	}

	component := views.Accounts(accs)
	return renderView(c, component)
}
