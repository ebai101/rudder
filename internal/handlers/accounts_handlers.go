package handlers

import (
	"fmt"
	"rudder/internal/services"
	"rudder/internal/views"
	"strconv"

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

func (ah *AccountsHandlers) accsNavbarHandler(c echo.Context) error {
	c.Set("ISERROR", false)
	ctx := c.Request().Context()

	accs, err := ah.service.GetAccountBalances(ctx)
	if err != nil {
		return err
	}

	component := views.AccountsNavbar(accs)
	return renderView(c, component)
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

func (ah *AccountsHandlers) accsDetailHandler(c echo.Context) error {
	c.Set("ISERROR", false)
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return nil
	}

	acc, err := ah.service.GetAccount(ctx, id)
	if err != nil {
		return err
	}

	component := views.AccountsDetail(acc)
	return renderView(c, component)
}

func (ah *AccountsHandlers) accsTransactionsHandler(c echo.Context) error {
	c.Set("ISERROR", false)
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return nil
	}

	txns, err := ah.service.GetAccountTransactions(ctx, int32(PAGE_SIZE)+1, 0, id)
	if err != nil {
		return err
	}

	txnListTarget := fmt.Sprintf("/accounts/%d/transactions/", id)
	component := views.TransactionsList(txns, txnListTarget, 20)
	return renderView(c, component)
}

func (ah *AccountsHandlers) accsTransactionsScrollHandler(c echo.Context) error {
	c.Set("ISERROR", false)
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return nil
	}

	page, err := strconv.ParseInt(c.Param("page"), 10, 32)
	if err != nil {
		page = 0
	}
	nextPage := page + PAGE_SIZE

	txns, err := ah.service.GetAccountTransactions(ctx, int32(PAGE_SIZE)+1, int32(page), id)
	if err != nil {
		return err
	}

	txnListTarget := fmt.Sprintf("/accounts/%d/transactions/", id)
	component := views.TransactionsList(txns, txnListTarget, nextPage)
	return renderView(c, component)
}
