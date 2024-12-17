package services

import (
	"context"
	"rudder/internal/models"
	"rudder/internal/repositories"
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

	rows, err := s.repo.GetTransactionRows(ctx, limit, offset, desc)
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

func (s *FinancialService) GetAccountBalances(
	ctx context.Context,
) ([]models.AccountBalance, error) {
	rows, err := s.repo.GetAccountBalances(ctx)
	if err != nil {
		return nil, err
	}

	var bals []models.AccountBalance
	for _, row := range rows {
		bals = append(bals, models.AccountBalance{
			ID:          row.ID,
			AccountName: row.AccountName,
			Balance:     row.Balance,
		})
	}

	return bals, nil
}

func (s *FinancialService) GetAccount(ctx context.Context, id int64) (models.Account, error) {
	row, err := s.repo.GetAccount(ctx, id)
	if err != nil {
		return models.Account{}, err
	}

	acc := models.Account{
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
	}

	return acc, nil
}

func (s *FinancialService) GetAccountTransactions(
	ctx context.Context,
	limit, offset int32,
	id int64,
) ([]models.Transaction, error) {
	rows, err := s.repo.GetAccountTransactions(ctx, limit, offset, id)
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
		})
	}

	return txns, nil
}
