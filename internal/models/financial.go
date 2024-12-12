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

type Organization struct {
	InstName   string `validate:"required"`
	SfinUrl    string `validate:"required,url"`
	DomainName string
	Url        string
}

type Account struct {
	AccountID    string `validate:"required"`
	AccountName  string `validate:"required"`
	InstName     string `validate:"required"`
	AccountType  string
	AccountClass string
	Currency     string `validate:"required"`
	Active       bool   `validate:"required"`
}

type Balance struct {
	BalanceID   string          `validate:"required"`
	BalanceDate time.Time       `validate:"required"`
	Balance     decimal.Decimal `validate:"required"`
	AccountID   string          `validate:"required"`
	AddedDate   time.Time       `validate:"required"`
}

type Transaction struct {
	TransactionID   string    `validate:"required"`
	PostedDate      time.Time `validate:"required"`
	Description     string
	Category        string
	Amount          decimal.Decimal `validate:"required"`
	AccountID       string          `validate:"required"`
	InstName        string          `validate:"required"`
	FullDescription string          `validate:"required"`
	AddedDate       time.Time       `validate:"required"`
	CategorizedDate time.Time
	Note            string
	CheckNum        string
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

func (acc SimpleFINAccount) GenerateBalanceID() string {
	id := acc.AccountId
	date := time.Unix(acc.BalanceDate, 0).UTC().Format("%-m/%-d/%y %I:%M %p")
	bal := acc.Balance.String()

	plaintext := []byte(fmt.Sprintf("%s%s%s", id, date, bal))
	hash := md5.Sum(plaintext)
	return hex.EncodeToString(hash[:])
}
