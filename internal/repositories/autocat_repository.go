package repositories

import (
	"context"
	"fmt"
	"rudder/internal/database"
	"rudder/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type AutocatRepository struct {
	queries *sqlc.Queries
	db      *database.DBConnection
}

func NewAutocatRepository(db *database.DBConnection) *AutocatRepository {
	return &AutocatRepository{
		queries: sqlc.New(db.Pool),
		db:      db,
	}
}

func (r *AutocatRepository) MatchTransactions(
	ctx context.Context,
) ([]sqlc.MatchTransactionsRow, error) {
	matches, err := r.queries.MatchTransactions(ctx)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *AutocatRepository) UpdateTransactionCategories(
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

func (r *AutocatRepository) GetAutocatRules(
	ctx context.Context,
) ([]sqlc.GetAutocatRulesRow, error) {
	return r.queries.GetAutocatRules(ctx)
}
