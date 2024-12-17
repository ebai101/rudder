package repositories

import (
	"context"
	"fmt"
	"rudder/internal/database"
	"rudder/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CategoriesRepository struct {
	queries *sqlc.Queries
	db      *database.DBConnection
}

func NewCategoriesRepository(db *database.DBConnection) *CategoriesRepository {
	return &CategoriesRepository{
		queries: sqlc.New(db.Pool),
		db:      db,
	}
}

func (r *CategoriesRepository) MatchTransactions(
	ctx context.Context,
) ([]sqlc.MatchTransactionsRow, error) {
	matches, err := r.queries.MatchTransactions(ctx)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *CategoriesRepository) UpdateTransactionCategories(
	ctx context.Context,
	matches []sqlc.MatchTransactionsRow,
) *sqlc.UpdateTransactionCategoriesBatchResults {
	var params []sqlc.UpdateTransactionCategoriesParams
	now := time.Now().UTC()

	for _, match := range matches {
		category := fmt.Sprintf("%v", match.NewCategory)

		params = append(params, sqlc.UpdateTransactionCategoriesParams{
			TransactionID:   match.TransactionID,
			Category:        pgtype.Text{String: category, Valid: true},
			CategorizedDate: pgtype.Timestamptz{Time: now, Valid: true},
		})
	}

	return r.queries.UpdateTransactionCategories(ctx, params)
}

func (r *CategoriesRepository) GetAutocatRules(
	ctx context.Context,
) ([]sqlc.GetAutocatRulesRow, error) {
	return r.queries.GetAutocatRules(ctx)
}
