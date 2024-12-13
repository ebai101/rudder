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

func (s *FinancialService) GetTransactions(
	ctx context.Context, limit int32,
) ([]models.Transaction, error) {
	rows, err := s.repo.GetTransactionRows(ctx, limit)
	if err != nil {
		return nil, err
	}

	var txns []models.Transaction
	for _, row := range rows {
		txns = append(txns, models.Transaction{
			TransactionID:   row.TransactionID,
			PostedDate:      row.PostedDate,
			Description:     row.Description.String,
			Category:        row.Category.String,
			Amount:          row.Amount,
			AccountID:       row.AccountID,
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
