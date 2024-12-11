package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"

	"github.com/shopspring/decimal"
	"gopkg.in/validator.v2"
)

type Organization struct {
	Domain  string `json:"domain"`
	SfinUrl string `json:"sfin-url"`
	Name    string `json:"name"`
	Url     string `json:"url"`
}

type Transaction struct {
	TransactionId string          `json:"id"`
	PostedDate    int64           `json:"posted"`
	Amount        decimal.Decimal `json:"amount"`
	Description   string          `json:"description"`
	Payee         string          `json:"payee"`
	TransactedAt  int64           `json:"transacted_at"`
}

type Account struct {
	Org          Organization    `json:"org"`
	AccountId    string          `json:"id"`
	AccountName  string          `json:"name"`
	Currency     string          `json:"currency"`
	Balance      decimal.Decimal `json:"balance"`
	BalanceAvail decimal.Decimal `json:"available-balance"`
	BalanceDate  int64           `json:"balance-date"`
	Transactions []Transaction   `json:"transactions"`
}

type SimpleFINResponse struct {
	Errors   []string  `json:"errors"`
	Accounts []Account `json:"accounts"`
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

func (acc Account) GenerateBalanceID() string {
	id := acc.AccountId
	date := time.Unix(acc.BalanceDate, 0).UTC().Format("%-m/%-d/%y %I:%M %p")
	bal := acc.Balance.String()

	plaintext := []byte(fmt.Sprintf("%s%s%s", id, date, bal))
	hash := md5.Sum(plaintext)
	return hex.EncodeToString(hash[:])
}

func (resp SimpleFINResponse) SaveResponseJSON(filename string) error {
	respJSON, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("error formatting SimpleFINResponse json: %v", err)
	}

	if err := os.WriteFile(filename, respJSON, 0644); err != nil {
		return fmt.Errorf("error writing SimpleFINResponse file to %v: %v", filename, err)
	}

	return nil
}

func LoadResponseJSON(filename string, sfinResponse *SimpleFINResponse) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error loading response json from %v: %v", filename, err)
	}
	defer file.Close()

	fileContents, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading response json: %v", err)
	}

	if err := json.Unmarshal(fileContents, &sfinResponse); err != nil {
		return err
	}

	return nil
}
