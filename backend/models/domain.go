package models

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/shopspring/decimal"
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
	Amount        decimal.Decimal `json:"amount,string"`
	Description   string          `json:"description"`
	Payee         string          `json:"payee"`
	TransactedAt  int64           `json:"transacted_at"`
}

type Account struct {
	Org          Organization    `json:"org"`
	AccountId    string          `json:"id"`
	AccountName  string          `json:"name"`
	Currency     string          `json:"currency"`
	Balance      decimal.Decimal `json:"balance,string"`
	BalanceAvail decimal.Decimal `json:"available-balance,string"`
	BalanceDate  int64           `json:"balance-date"`
	Transactions []Transaction   `json:"transactions"`
}

type SimpleFINResponse struct {
	Errors   []string  `json:"errors"`
	Accounts []Account `json:"accounts"`
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
