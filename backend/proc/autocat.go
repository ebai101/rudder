package proc

import (
	"log"
	"rudder/backend/models"
	"time"
)

func CategorizeTransactions(txns []models.TransactionRow, rules []models.AutoCatRule) ([]models.TransactionRow, error) {
	if len(txns) == 0 {
		return txns, nil
	}

	newTxns := txns

	rulesApplied := 0
	for idx, txn := range newTxns {
		for _, rule := range rules {
			if rule.MatchesRow(txn) {
				newTxn, err := rule.ApplyOverrides(txn)
				if err != nil {
					return txns, err
				}
				newTxn.CategorizedDate = time.Now().UTC()
				newTxns[idx] = newTxn

				log.Printf("Rule criteria %v matches transaction %v\n", rule.Criteria, txn.TransactionID)
				log.Println("New transaction:")
				log.Printf("%+v\n", newTxn)
				rulesApplied++
			}
		}
	}
	log.Printf("Applied %d categories to fetched transactions\n", rulesApplied)

	return newTxns, nil
}
