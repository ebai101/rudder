package services

import (
	"context"
	"rudder/internal/repositories"
	"rudder/sqlc"
)

type AutocatService struct {
	acRepo      *repositories.AutocatRepository
	sfinService *SimpleFINService
}

func NewAutocatService(
	acRepo *repositories.AutocatRepository,
	sfinService *SimpleFINService,
) *AutocatService {
	return &AutocatService{
		acRepo:      acRepo,
		sfinService: sfinService,
	}
}

func (s *AutocatService) CategorizeTransactions(
	ctx context.Context,
) (*sqlc.UpdateTransactionCategoriesBatchResults, error) {
	matches, err := s.acRepo.MatchTransactions(ctx)
	if err != nil {
		return nil, err
	}

	return s.acRepo.UpdateTransactionCategories(ctx, matches), nil
}
