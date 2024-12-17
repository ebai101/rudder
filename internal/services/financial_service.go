package services

import (
	"context"
	"rudder/internal/models"
	"rudder/internal/repositories"
	"rudder/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type FinancialService struct {
	repo *repositories.FinancialRepository
}

func NewFinancialService(
	repo *repositories.FinancialRepository,
) *FinancialService {
	return &FinancialService{
		repo: repo,
	}
}

func (s *FinancialService) GetTransactionRows(
	ctx context.Context, limit, offset int32, desc string,
) ([]models.Transaction, error) {
	if desc == "" {
		desc = "%"
	}

	params := sqlc.GetTransactionRowsParams{
		Description: pgtype.Text{String: desc, Valid: true},
		Limit:       limit,
		Offset:      offset,
	}
	rows, err := s.repo.GetTransactionRows(ctx, params)
	if err != nil {
		return nil, err
	}

	var txns []models.Transaction
	for _, row := range rows {
		txns = append(txns, models.Transaction{
			ID:              row.ID,
			TransactionID:   row.TransactionID,
			PostedDate:      row.PostedDate,
			Description:     row.Description.String,
			Category:        row.Category.String,
			Amount:          row.Amount,
			AccountName:     row.AccountName,
			InstName:        row.InstName,
			FullDescription: row.FullDescription,
			AddedDate:       row.AddedDate,
			CategorizedDate: row.CategorizedDate.Time,
			Note:            row.Note.String,
			CheckNum:        row.CheckNum.String,
		})
	}

	return txns, nil
}

func (s *FinancialService) GetTransaction(
	ctx context.Context,
	id int64,
) (models.Transaction, error) {
	row, err := s.repo.GetTransaction(ctx, id)
	if err != nil {
		return models.Transaction{}, err
	}

	txn := models.Transaction{
		ID:              row.ID,
		TransactionID:   row.TransactionID,
		PostedDate:      row.PostedDate,
		Description:     row.Description.String,
		Category:        row.Category.String,
		Amount:          row.Amount,
		AccountName:     row.AccountName,
		InstName:        row.InstName,
		FullDescription: row.FullDescription,
		AddedDate:       row.AddedDate,
		CategorizedDate: row.CategorizedDate.Time,
		Note:            row.Note.String,
		CheckNum:        row.CheckNum.String,
	}

	return txn, nil

}

func (s *FinancialService) GetAccountRows(ctx context.Context) ([]models.Account, error) {
	rows, err := s.repo.GetAccountRows(ctx)
	if err != nil {
		return nil, err
	}

	var accs []models.Account
	for _, row := range rows {
		accs = append(accs, models.Account{
			ID:               row.ID,
			AccountID:        row.AccountID,
			AccountName:      row.AccountName,
			InstName:         row.InstName,
			AccountType:      string(row.AccountType.AccountTypeT),
			AccountClass:     string(row.AccountClass.AccountClassT),
			Currency:         row.Currency,
			Active:           row.Active,
			Balance:          row.Balance,
			BalanceDate:      row.BalanceDate,
			BalanceAddedDate: row.AddedDate,
		})
	}

	return accs, nil
}
