package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"rudder/backend/config"
	"rudder/backend/models"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Config config.AppConfig
	Pool   *pgxpool.Pool
}

func upsertOrganizations(conn *pgxpool.Conn, orgs []models.OrganizationRow) (int64, error) {
	batch := &pgx.Batch{}
	for _, row := range orgs {
		batch.Queue(`
			insert into organizations (inst_name, sfin_url, domain_name, url)
			values ($1, $2, $3, $4)
			on conflict do nothing
		`, row.InstName, row.SfinUrl, row.DomainName, row.URL)
	}
	br := conn.SendBatch(context.Background(), batch)
	defer br.Close()

	var totalInserted int64
	for range orgs {
		cmdTag, err := br.Exec()
		if err != nil {
			return 0, err
		}
		if cmdTag.RowsAffected() > 0 {
			totalInserted++
		}
	}

	return totalInserted, nil
}

func upsertAccounts(conn *pgxpool.Conn, accs []models.AccountRow) (int64, error) {
	batch := &pgx.Batch{}
	for _, row := range accs {
		batch.Queue(`
			insert into accounts (account_id, account_name, inst_name, currency)
			values ($1, $2, $3, $4)
			on conflict do nothing
		`, row.AccountID, row.AccountName, row.InstName, row.Currency)
	}
	br := conn.SendBatch(context.Background(), batch)
	defer br.Close()

	var totalInserted int64
	for range accs {
		cmdTag, err := br.Exec()
		if err != nil {
			return 0, err
		}
		if cmdTag.RowsAffected() > 0 {
			totalInserted++
		}
	}

	return totalInserted, nil
}

func upsertBalances(conn *pgxpool.Conn, bals []models.BalanceRow) (int64, error) {
	batch := &pgx.Batch{}
	for _, row := range bals {
		batch.Queue(`
			insert into balances (balance_id, balance_date, balance, account_id)
			values ($1, $2, $3, $4)
			on conflict do nothing
		`, row.BalanceID, row.BalanceDate, row.Balance, row.AccountID)
	}
	br := conn.SendBatch(context.Background(), batch)
	defer br.Close()

	var totalInserted int64
	for range bals {
		cmdTag, err := br.Exec()
		if err != nil {
			return 0, err
		}
		if cmdTag.RowsAffected() > 0 {
			totalInserted++
		}
	}

	return totalInserted, nil
}

func upsertTransactions(conn *pgxpool.Conn, txns []models.TransactionRow) (int64, error) {
	batch := &pgx.Batch{}
	for _, row := range txns {
		batch.Queue(`
			insert into transactions (transaction_id, posted_date, description, amount, account_id, inst_name, full_description)
			values ($1, $2, $3, $4, $5, $6, $7)
			on conflict do nothing
			`,
			row.TransactionID,
			row.PostedDate,
			row.Description,
			row.Amount,
			row.AccountID,
			row.InstName,
			row.FullDescription,
		)
	}
	br := conn.SendBatch(context.Background(), batch)
	defer br.Close()

	var totalInserted int64
	for range txns {
		cmdTag, err := br.Exec()
		if err != nil {
			return 0, err
		}
		if cmdTag.RowsAffected() > 0 {
			totalInserted++
		}
	}

	return totalInserted, nil
}

func (db Database) UpsertAll(model models.RowModel) error {
	conn, err := db.Pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	orgsInserted, err := upsertOrganizations(conn, model.Organizations)
	if err != nil {
		return err
	}
	if orgsInserted > 0 {
		log.Printf("Inserted %d organizations\n", orgsInserted)
	}

	accsInserted, err := upsertAccounts(conn, model.Accounts)
	if err != nil {
		return err
	}
	if accsInserted > 0 {
		log.Printf("Inserted %d accounts\n", accsInserted)
	}

	balsInserted, err := upsertBalances(conn, model.Balances)
	if err != nil {
		return err
	}
	if balsInserted > 0 {
		log.Printf("Inserted %d balances\n", balsInserted)
	}

	txnsInserted, err := upsertTransactions(conn, model.Transactions)
	if err != nil {
		return err
	}
	if txnsInserted > 0 {
		log.Printf("Inserted %d transactions\n", txnsInserted)
	}

	if orgsInserted == 0 && accsInserted == 0 && balsInserted == 0 && txnsInserted == 0 {
		log.Printf("No updates at this time.")
	}

	return nil
}

func (db Database) GetAutocatRules() ([]models.AutoCatRule, error) {
	var rules []models.AutoCatRule

	conn, err := db.Pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `
		select 
			(
				select json_agg(
					json_build_object(
						'column_name', c.column_name,
						'operator', c."operator",
						'filter_value', c.filter_value,
						'criteria_order', c.criteria_order
					)
				) 
				from autocat_criteria c 
				where c.rule_id = r.id
			) as criteria,
			(
				select json_agg(
					json_build_object(
						'column_name', o.column_name,
						'override_value', o.override_value,
						'override_order', o.override_order
					)
				) 
				from autocat_overrides o 
				where o.rule_id = r.id
			) as overrides
		from 
			autocat_rules r;	
	`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var rawCriteria string
		var rawOverrides string
		var rule models.AutoCatRule

		if err := rows.Scan(&rawCriteria, &rawOverrides); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(rawCriteria), &rule.Criteria); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(rawOverrides), &rule.Overrides); err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

func (db Database) GetTransactionRows() ([]models.TransactionRow, error) {
	var txns []models.TransactionRow

	conn, err := db.Pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `
		select
			transaction_id, posted_date, description, category,
			amount, account_id, inst_name, full_description,
			added_date, categorized_date, note, check_num
		from transactions`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var txn models.TransactionRow

		err := rows.Scan(
			&txn.TransactionID,
			&txn.PostedDate,
			&txn.Description,
			&txn.Category,
			&txn.Amount,
			&txn.AccountID,
			&txn.InstName,
			&txn.FullDescription,
			&txn.AddedDate,
			&txn.CategorizedDate,
			&txn.Note,
			&txn.CheckNum,
		)
		if err != nil {
			return nil, err
		}

		txns = append(txns, txn)
	}

	return txns, nil
}

func (db Database) Close() {
	db.Pool.Close()
}

func OpenDatabase(appConfig *config.AppConfig) (Database, error) {
	db := Database{}

	poolConfig, err := pgxpool.ParseConfig(appConfig.DatabaseUrl)
	if err != nil {
		return Database{}, fmt.Errorf("unable to parse db config: %v", err)
	}
	poolConfig.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		pgxdecimal.Register(c.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return Database{}, fmt.Errorf("unable to connect to database: %v", err)
	}
	db.Pool = pool

	return db, nil
}
