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
	limit int32,
) ([]sqlc.GetTransactionRowsRow, error) {
	return r.queries.GetTransactionRows(ctx, limit)
}

func (r *FinancialRepository) GetAccountRows(
	ctx context.Context,
) ([]sqlc.GetAccountRowsRow, error) {
	return r.queries.GetAccountRows(ctx)
}

func (r *FinancialRepository) GetBalanceRows(
	ctx context.Context,
	limit int32,
) ([]sqlc.GetBalanceRowsRow, error) {
	return r.queries.GetBalanceRows(ctx, limit)
}

func (r *FinancialRepository) GetOrganizationRows(
	ctx context.Context,
	limit int32,
) ([]sqlc.GetOrganizationRowsRow, error) {
	return r.queries.GetOrganizationRows(ctx)
}
