package repositories

import (
	"context"
	"rudder/internal/database"
	"rudder/internal/models"
	"rudder/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type SimpleFINRepository struct {
	queries *sqlc.Queries
	db      *database.DBConnection
}

func NewSimpleFINRepository(db *database.DBConnection) *SimpleFINRepository {
	return &SimpleFINRepository{
		queries: sqlc.New(db.Pool),
		db:      db,
	}
}

func (r *SimpleFINRepository) InsertAccounts(
	ctx context.Context,
	accs []models.SimpleFINAccount,
) *sqlc.InsertAccountsBatchResults {
	var params []sqlc.InsertAccountsParams
	for _, acc := range accs {
		params = append(params, sqlc.InsertAccountsParams{
			AccountID:   acc.AccountId,
			AccountName: acc.AccountName,
			InstName:    acc.Org.Name,
			Currency:    acc.Currency,
		})
	}
	return r.queries.InsertAccounts(context.Background(), params)
}

func (r *SimpleFINRepository) InsertOrganizations(
	ctx context.Context,
	accs []models.SimpleFINAccount,
) *sqlc.InsertOrganizationsBatchResults {
	var params []sqlc.InsertOrganizationsParams
	for _, acc := range accs {
		params = append(params, sqlc.InsertOrganizationsParams{
			InstName:   acc.Org.Name,
			SfinUrl:    acc.Org.SfinUrl,
			DomainName: pgtype.Text{String: acc.Org.Domain, Valid: true},
			Url:        pgtype.Text{String: acc.Org.Url, Valid: true},
		})
	}
	return r.queries.InsertOrganizations(context.Background(), params)
}

func (r *SimpleFINRepository) InsertBalances(
	ctx context.Context,
	accs []models.SimpleFINAccount,
) *sqlc.InsertBalancesBatchResults {
	var params []sqlc.InsertBalancesParams
	for _, acc := range accs {
		balID := acc.GenerateBalanceID()
		balDate := time.Unix(acc.BalanceDate, 0).UTC()

		params = append(params, sqlc.InsertBalancesParams{
			BalanceID:   balID,
			BalanceDate: balDate,
			Balance:     decimal.Decimal{},
			AccountID:   acc.AccountId,
		})
	}
	return r.queries.InsertBalances(context.Background(), params)
}

func (r *SimpleFINRepository) InsertTransactions(
	ctx context.Context,
	accs []models.SimpleFINAccount,
) *sqlc.InsertTransactionsBatchResults {
	var params []sqlc.InsertTransactionsParams

	for _, acc := range accs {
		for _, txn := range acc.Transactions {
			postedDate := time.Unix(txn.PostedDate, 0).UTC()

			params = append(params, sqlc.InsertTransactionsParams{
				TransactionID:   txn.TransactionId,
				PostedDate:      postedDate,
				Description:     pgtype.Text{String: txn.Payee, Valid: true},
				Amount:          txn.Amount,
				AccountID:       acc.AccountId,
				InstName:        acc.Org.Name,
				FullDescription: txn.Description,
			})
		}
	}
	return r.queries.InsertTransactions(ctx, params)
}
