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

// func (s *AutocatService) GetAutocatRules(ctx context.Context) ([]models.AutocatRule, error) {
// 	rows, err := s.acRepo.GetAutocatRules(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting autocat rules: %v", err)
// 	}

// 	var rules []models.AutocatRule
// 	for _, row := range rows {
// 		var rule models.AutocatRule

// 		if err := json.Unmarshal(row.Criteria, &rule.Criteria); err != nil {
// 			return nil, fmt.Errorf("error unmarshalling JSON response: %v", err)
// 		}
// 		if err := json.Unmarshal(row.Overrides, &rule.Overrides); err != nil {
// 			return nil, fmt.Errorf("error unmarshalling JSON response: %v", err)
// 		}

// 		rules = append(rules, rule)
// 	}
// 	return rules, nil
// }

func (s *AutocatService) CategorizeTransactions(
	ctx context.Context,
) (*sqlc.UpdateTransactionCategoriesBatchResults, error) {
	matches, err := s.acRepo.MatchTransactions(ctx)
	if err != nil {
		return nil, err
	}

	return s.acRepo.UpdateTransactionCategories(ctx, matches), nil
}

// func CategorizeTransactions(
// 	txns []models.Transaction,
// 	rules []models.AutocatRule,
// ) ([]models.Transaction, error) {
// 	if len(txns) == 0 {
// 		return txns, nil
// 	}
// 	newTxns := txns

// 	rulesApplied := 0
// 	for idx, txn := range newTxns {
// 		for _, rule := range rules {
// 			if rule.MatchesRow(txn) {
// 				newTxn, err := rule.ApplyOverrides(txn)
// 				if err != nil {
// 					return txns, err
// 				}
// 				newTxn.CategorizedDate = time.Now().UTC()
// 				newTxns[idx] = newTxn

// 				log.Printf(
// 					"Rule criteria %v matches transaction %v\n",
// 					rule.Criteria,
// 					txn.TransactionId,
// 				)
// 				log.Println("New transaction:")
// 				log.Printf("%+v\n", newTxn)
// 				rulesApplied++
// 			}
// 		}
// 	}
// 	log.Printf("Applied %d categories to fetched transactions\n", rulesApplied)

// 	return newTxns, nil
// }
