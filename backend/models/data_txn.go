package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/shopspring/decimal"
	"gopkg.in/validator.v2"
)

type OrganizationRow struct {
	InstName   string `validate:"nonnil"`
	SfinUrl    string `validate:"nonnil"`
	DomainName string `validate:"nonnil"`
	URL        string `validate:"nonnil"`
}

type AccountRow struct {
	AccountID   string `validate:"nonnil"`
	AccountName string `validate:"nonnil"`
	InstName    string `validate:"nonnil"`
	Currency    string `validate:"nonnil"`
}

type BalanceRow struct {
	BalanceID   string          `validate:"nonnil"`
	BalanceDate time.Time       `validate:"nonnil,timestamp"`
	Balance     decimal.Decimal `validate:"nonnil,amount"`
	AccountID   string          `validate:"nonnil"`
	AddedDate   time.Time       `validate:"nonnil,timestamp"`
}

type TransactionRow struct {
	TransactionID   string    `validate:"nonnil"`
	PostedDate      time.Time `validate:"nonnil,timestamp"`
	Description     string    `validate:"nonnil"`
	Category        string
	Amount          decimal.Decimal `validate:"amount"`
	AccountID       string          `validate:"nonnil"`
	InstName        string          `validate:"nonnil"`
	FullDescription string          `validate:"nonnil"`
	AddedDate       time.Time       `validate:"nonnil,timestamp"`
	CategorizedDate time.Time
	Note            string
	CheckNum        string
}

type RowModel struct {
	Organizations []OrganizationRow
	Accounts      []AccountRow
	Balances      []BalanceRow
	Transactions  []TransactionRow
}

func validateTimestamp(v any, param string) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.TypeFor[time.Time]().Kind() {
		return validator.ErrUnsupported
	}

	timestamp, _ := val.Interface().(time.Time)
	if timestamp.IsZero() {
		return nil
	}

	if timestamp.Compare(time.Now().AddDate(0, 0, 1)) > 0 {
		return fmt.Errorf("timestamp %v cannot be in the future", timestamp)
	}
	earlyDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	if timestamp.Compare(earlyDate) < 0 {
		return fmt.Errorf("timestamp %v cannot be before the year 2000", timestamp)
	}

	return nil
}

func validateAmount(v any, param string) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.TypeFor[decimal.Decimal]().Kind() {
		return validator.ErrUnsupported
	}

	amount, _ := val.Interface().(decimal.Decimal)
	if amount.Compare(decimal.NewFromFloat(99999999.99)) > 0 {
		return errors.New("amount exceeds maximum allowed value")
	}
	return nil
}

func getBalanceID(acc Account) string {
	id := acc.AccountId
	date := time.Unix(acc.BalanceDate, 0).UTC().Format("%-m/%-d/%y %I:%M %p")
	bal := acc.Balance.String()

	plaintext := []byte(fmt.Sprintf("%s%s%s", id, date, bal))
	hash := md5.Sum(plaintext)
	return hex.EncodeToString(hash[:])
}

func parseOrganizations(accs []Account) ([]OrganizationRow, error) {
	var instRows []OrganizationRow

	for _, acc := range accs {
		row := OrganizationRow{
			InstName:   acc.Org.Name,
			SfinUrl:    acc.Org.SfinUrl,
			DomainName: acc.Org.Domain,
			URL:        acc.Org.Url,
		}
		if err := validator.Validate(row); err != nil {
			return nil, fmt.Errorf("error validating OrganizationRow: %v", err)
		}
		instRows = append(instRows, row)
	}

	return instRows, nil
}

func parseAccounts(accs []Account) ([]AccountRow, error) {
	var accRows []AccountRow

	for _, acc := range accs {
		row := AccountRow{
			AccountID:   acc.AccountId,
			AccountName: acc.AccountName,
			InstName:    acc.Org.Name,
			Currency:    acc.Currency,
		}
		if err := validator.Validate(row); err != nil {
			return nil, fmt.Errorf("error validating AccountRow: %v", err)
		}
		accRows = append(accRows, row)
	}

	return accRows, nil
}

func parseBalances(accs []Account) ([]BalanceRow, error) {
	var balRows []BalanceRow

	for _, acc := range accs {
		balID := getBalanceID(acc)
		balDate := time.Unix(acc.BalanceDate, 0).UTC()
		row := BalanceRow{
			BalanceID:   balID,
			BalanceDate: balDate,
			Balance:     acc.Balance,
			AccountID:   acc.AccountId,
		}
		if err := validator.Validate(row); err != nil {
			return nil, fmt.Errorf("error validating BalanceRow: %v", err)
		}
		balRows = append(balRows, row)
	}

	return balRows, nil
}

func parseTransactions(accs []Account) ([]TransactionRow, error) {
	var txnRows []TransactionRow

	for _, acc := range accs {
		for _, txn := range acc.Transactions {
			postedDate := time.Unix(txn.PostedDate, 0).UTC()
			row := TransactionRow{
				TransactionID:   txn.TransactionId,
				PostedDate:      postedDate,
				Description:     txn.Payee,
				Amount:          txn.Amount,
				AccountID:       acc.AccountId,
				InstName:        acc.Org.Name,
				FullDescription: txn.Description,
			}
			if err := validator.Validate(row); err != nil {
				return nil, fmt.Errorf("error validating TransactionRow: %v", err)
			}
			txnRows = append(txnRows, row)
		}
	}

	return txnRows, nil
}

func ParseSimpleFINResponse(sfinResponse *SimpleFINResponse) (RowModel, error) {
	validator.SetValidationFunc("timestamp", validateTimestamp)
	validator.SetValidationFunc("amount", validateAmount)

	rowModel := RowModel{}

	orgs, err := parseOrganizations(sfinResponse.Accounts)
	if err != nil {
		return RowModel{}, err
	}
	rowModel.Organizations = orgs

	accs, err := parseAccounts(sfinResponse.Accounts)
	if err != nil {
		return RowModel{}, err
	}
	rowModel.Accounts = accs

	bals, err := parseBalances(sfinResponse.Accounts)
	if err != nil {
		return RowModel{}, err
	}
	rowModel.Balances = bals

	txns, err := parseTransactions(sfinResponse.Accounts)
	if err != nil {
		return RowModel{}, err
	}
	rowModel.Transactions = txns

	return rowModel, nil
}
