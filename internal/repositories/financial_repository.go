package repositories

import (
	"context"
	"rudder/internal/database"
	"rudder/internal/models"
	"rudder/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type FinancialRepository struct {
	queries *sqlc.Queries
	db      *database.DBConnection
}

func NewFinancialRepository(db *database.DBConnection) *FinancialRepository {
	return &FinancialRepository{
		queries: sqlc.New(db.Pool),
		db:      db,
	}
}

func (r *FinancialRepository) GetTransactionRows(
	ctx context.Context,
	limit, offset int32,
	desc string,
) ([]sqlc.TransactionsView, error) {
	params := sqlc.GetTransactionRowsParams{
		Description: pgtype.Text{String: desc, Valid: true},
		Limit:       limit,
		Offset:      offset,
	}
	return r.queries.GetTransactionRows(ctx, params)
}

func (r *FinancialRepository) GetTransaction(
	ctx context.Context,
	id int64,
) (sqlc.TransactionsView, error) {
	return r.queries.GetTransaction(ctx, id)
}

func (r *FinancialRepository) GetAccountRows(
	ctx context.Context,
) ([]sqlc.AccountsView, error) {
	return r.queries.GetAccountRows(ctx)
}

func (r *FinancialRepository) GetAccount(
	ctx context.Context,
	id int64,
) (sqlc.AccountsView, error) {
	return r.queries.GetAccount(ctx, id)
}

func (r *FinancialRepository) GetAccountTransactions(
	ctx context.Context,
	limit, offset int32,
	id int64,
) ([]sqlc.TransactionsView, error) {
	params := sqlc.GetAccountTransactionsParams{
		ID:     id,
		Limit:  limit,
		Offset: offset,
	}
	return r.queries.GetAccountTransactions(ctx, params)
}

func (r *FinancialRepository) GetAccountBalances(
	ctx context.Context,
) ([]sqlc.GetAccountBalancesRow, error) {
	return r.queries.GetAccountBalances(ctx)
}

func (r *FinancialRepository) GetInsights(
	ctx context.Context,
	interval models.IntervalPair,
) (sqlc.GetInsightsRow, error) {
	args := sqlc.GetInsightsParams{
		PostedDate:   interval.Start,
		PostedDate_2: interval.End,
	}
	return r.queries.GetInsights(ctx, args)
}
