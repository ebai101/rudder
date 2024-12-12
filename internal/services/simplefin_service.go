package services

import (
	"context"
	"log"
	"rudder/internal/clients"
	"rudder/internal/config"
	"rudder/internal/models"
	"rudder/internal/repositories"
	"rudder/sqlc"
	"time"
)

type SimpleFINService struct {
	config         *config.AppConfig
	apiClient      *clients.SimpleFINClient
	repo           *repositories.SimpleFINRepository
	lastSyncTime   time.Time
	syncInProgress bool
}

func NewSimpleFINService(
	config *config.AppConfig,
	apiClient *clients.SimpleFINClient,
	repo *repositories.SimpleFINRepository,
) *SimpleFINService {
	return &SimpleFINService{
		config:    config,
		apiClient: apiClient,
		repo:      repo,
	}
}

func (s *SimpleFINService) InsertOrganizations(
	ctx context.Context,
	accs []models.SimpleFINAccount,
) *sqlc.InsertOrganizationsBatchResults {
	return s.repo.InsertOrganizations(ctx, accs)
}

func (s *SimpleFINService) InsertAccounts(
	ctx context.Context,
	accs []models.SimpleFINAccount,
) *sqlc.InsertAccountsBatchResults {
	return s.repo.InsertAccounts(ctx, accs)
}

func (s *SimpleFINService) InsertBalances(
	ctx context.Context,
	accs []models.SimpleFINAccount,
) *sqlc.InsertBalancesBatchResults {
	return s.repo.InsertBalances(ctx, accs)
}

func (s *SimpleFINService) InsertTransactions(
	ctx context.Context,
	accs []models.SimpleFINAccount,
) *sqlc.InsertTransactionsBatchResults {
	return s.repo.InsertTransactions(ctx, accs)
}

func (s *SimpleFINService) SyncSimpleFIN(
	ctx context.Context,
	useCached bool,
	saveCached bool,
	numDays int,
) error {
	var sfinResp models.SimpleFINResponse
	respFilename := "sfin_last.json"

	if useCached {
		if err := models.LoadResponseJSON(respFilename, &sfinResp); err != nil {
			return err
		}
	} else {
		log.Printf("Fetching %d days...\n", numDays)
		if err := s.apiClient.GetAccounts(numDays, &sfinResp); err != nil {
			return err
		}
		if saveCached {
			log.Printf("Saving response to %v...\n", respFilename)
			err := sfinResp.SaveResponseJSON(respFilename)
			if err != nil {
				return err
			}
		}
	}

	s.InsertOrganizations(ctx, sfinResp.Accounts)
	s.InsertAccounts(ctx, sfinResp.Accounts)
	s.InsertBalances(ctx, sfinResp.Accounts)
	s.InsertTransactions(ctx, sfinResp.Accounts)

	// log.Println("Fetching AutoCat rules...")
	// acRules, err := s.repo.
	// if err != nil {
	// 	return err
	// }

	// log.Println("Categorizing new transactions...")
	// newTxns, err := CategorizeTransactions(rowModel.Transactions, acRules)
	// if err != nil {
	// 	return err
	// }
	// rowModel.Transactions = newTxns

	// log.Println("Updating transaction categories...")
	// if err := db.UpdateTransactionCategories(rowModel.Transactions); err != nil {
	// 	return err
	// }

	return nil
}
