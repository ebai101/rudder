package services

import (
	"context"
	"encoding/json"
	"log"
	"rudder/internal/models"
	"rudder/internal/repositories"
	"rudder/sqlc"
)

type AutocatService struct {
	repo        *repositories.AutocatRepository
	sfinService *SimpleFINService
}

func NewAutocatService(
	acRepo *repositories.AutocatRepository,
	sfinService *SimpleFINService,
) *AutocatService {
	return &AutocatService{
		repo:        acRepo,
		sfinService: sfinService,
	}
}

func (s *AutocatService) CategorizeTransactions(
	ctx context.Context,
) (*sqlc.UpdateTransactionCategoriesBatchResults, error) {
	log.Printf("Categorizing transactions...")
	matches, err := s.repo.MatchTransactions(ctx)
	if err != nil {
		return nil, err
	}

	return s.repo.UpdateTransactionCategories(ctx, matches), nil
}

func (s *AutocatService) GetAutocatRules(ctx context.Context) ([]models.AutocatRule, error) {
	rows, err := s.repo.GetAutocatRules(ctx)
	if err != nil {
		return nil, err
	}

	var rules []models.AutocatRule
	for _, row := range rows {
		var cta []models.AutocatCriterion
		var ovr []models.AutocatOverride
		if err := json.Unmarshal(row.Criteria, &cta); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(row.Overrides, &ovr); err != nil {
			return nil, err
		}

		rules = append(rules, models.AutocatRule{
			RuleID:    row.ID,
			Criteria:  cta,
			Overrides: ovr,
		})
	}

	return rules, nil
}
