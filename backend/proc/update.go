package proc

import (
	"log"
	"rudder/backend/config"
	"rudder/backend/models"
	"rudder/backend/resource"
)

func Update(appConfig *config.AppConfig, db *resource.Database, sfinAPI *resource.SimpleFINAPI, args *config.Args, numDays int) error {
	var sfinResp models.SimpleFINResponse
	respFilename := "sfin_last.json"

	if args.UseCached {
		if err := models.LoadResponseJSON(respFilename, &sfinResp); err != nil {
			return err
		}
	} else {
		log.Printf("Fetching %d days...\n", numDays)
		if err := sfinAPI.GetAccounts(numDays, &sfinResp); err != nil {
			return err
		}
		if args.SaveCached {
			log.Printf("Saving response to %v...\n", respFilename)
			err := sfinResp.SaveResponseJSON(respFilename)
			if err != nil {
				return err
			}
		}
	}

	log.Println("Parsing response...")
	rowModel, err := models.ParseSimpleFINResponse(&sfinResp)
	if err != nil {
		return err
	}

	log.Println("Fetching AutoCat rules...")
	acRules, err := db.GetAutocatRules()
	if err != nil {
		return err
	}

	log.Println("Updating records...")
	if err := db.InsertAll(rowModel); err != nil {
		return err
	}

	log.Println("Categorizing new transactions...")
	newTxns, err := CategorizeTransactions(rowModel.Transactions, acRules)
	if err != nil {
		return err
	}
	rowModel.Transactions = newTxns

	log.Println("Updating transaction categories...")
	if err := db.UpdateTransactionCategories(rowModel.Transactions); err != nil {
		return err
	}

	return nil
}
