package handlers

import (
	"rudder/internal/services"
	"rudder/internal/views"
	"strconv"

	"github.com/labstack/echo/v4"
)

const pageSize int64 = 20

type TransactionsHandlers struct {
	service *services.FinancialService
}

func NewTransactionsHandlers(
	service *services.FinancialService,
) *TransactionsHandlers {
	return &TransactionsHandlers{
		service: service,
	}
}

func (th *TransactionsHandlers) txnsMainHandler(c echo.Context) error {
	c.Set("ISERROR", false)
	ctx := c.Request().Context()

	txns, err := th.service.GetTransactionRows(ctx, 21, 0, "")
	if err != nil {
		return err
	}

	return renderView(c, views.Transactions(txns))
}

func (th *TransactionsHandlers) txnsScrollHandler(c echo.Context) error {
	c.Set("ISERROR", false)
	ctx := c.Request().Context()

	page, err := strconv.ParseInt(c.Param("page"), 10, 32)
	if err != nil {
		page = 0
	}
	nextPage := page + pageSize

	txns, err := th.service.GetTransactionRows(ctx, 21, int32(page), "")
	if err != nil {
		return err
	}

	return renderView(c, views.TransactionsList(txns, "/transactions/", nextPage))
}

func (th *TransactionsHandlers) txnDetailsHandler(c echo.Context) error {
	c.Set("ISERROR", false)
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return nil
	}

	txn, err := th.service.GetTransaction(ctx, id)
	if err != nil {
		return err
	}

	return renderView(c, views.TransactionDetail(txn))
}

func (th *TransactionsHandlers) txnDetailsEditHandler(c echo.Context) error {
	c.Set("ISERROR", false)
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return nil
	}

	txn, err := th.service.GetTransaction(ctx, id)
	if err != nil {
		return err
	}

	return renderView(c, views.TransactionDetailEdit(txn))
}
