package repositories

import (
	"context"
	"rudder/internal/database"
	"rudder/sqlc"
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
	params sqlc.GetTransactionRowsParams,
) ([]sqlc.GetTransactionRowsRow, error) {
	return r.queries.GetTransactionRows(ctx, params)
}

func (r *FinancialRepository) GetTransaction(
	ctx context.Context,
	id int64,
) (sqlc.GetTransactionRow, error) {
	return r.queries.GetTransaction(ctx, id)
}

func (r *FinancialRepository) GetAccountRows(
	ctx context.Context,
) ([]sqlc.GetAccountRowsRow, error) {
	return r.queries.GetAccountRows(ctx)
}

func (r *FinancialRepository) GetAccount(
	ctx context.Context,
	id int64,
) (sqlc.GetAccountRow, error) {
	return r.queries.GetAccount(ctx, id)
}
