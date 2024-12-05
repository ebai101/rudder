package resource

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"rudder/backend/config"
	"strconv"
	"strings"
	"time"

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

func (resp SimpleFINResponse) SaveResponse(filename string) error {
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

type SimpleFINAPI struct {
	Config *config.AppConfig
}

func (api SimpleFINAPI) makeRequest(method string, url string, parameters map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.SetBasicAuth(api.Config.SFINAuth.Username, api.Config.SFINAuth.Password)
	if len(parameters) > 0 {
		q := req.URL.Query()
		for param, value := range parameters {
			q.Add(param, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code from SimpleFIN: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response from SimpleFIN: %v", err)
	}

	return body, nil

}

func (api SimpleFINAPI) GetAccounts(days int, sfinResp *SimpleFINResponse) error {
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -1*days)

	params := map[string]string{
		"start-date": strconv.FormatInt(startTime.Unix(), 10),
		"end-date":   strconv.FormatInt(endTime.Unix(), 10),
	}

	resp, err := api.makeRequest("GET", api.Config.SFINAuth.URL, params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp, &sfinResp); err != nil {
		return fmt.Errorf("error parsing SimpleFIN response: %v", err)
	}

	return nil
}

func (api SimpleFINAPI) ValidateAPICredentials() error {
	_, err := api.makeRequest("GET", api.Config.SFINAuth.URL, map[string]string{})
	if err != nil {
		return err
	}

	log.Println("API Credentials validated.")
	return nil
}

func (api SimpleFINAPI) GenerateAPICredentials() error {
	decodedToken, err := base64.StdEncoding.DecodeString(api.Config.SfinBridgeToken)
	if err != nil {
		return fmt.Errorf("unable to decode SimpleFIN token: %v", err)
	}
	genURL := string(decodedToken[:])

	resp, err := api.makeRequest("POST", genURL, map[string]string{})
	if err != nil {
		return err
	}
	accessUrl := string(resp[:])

	schemeSplit := strings.SplitN(accessUrl, "//", 2)
	authSplit := strings.SplitN(schemeSplit[1], "@", 2)
	url := schemeSplit[0] + "//" + authSplit[1] + "/accounts"
	credSplit := strings.SplitN(authSplit[0], ":", 2)

	sfinAuth := config.SimpleFINAuth{
		URL:      url,
		Username: credSplit[0],
		Password: credSplit[1],
	}

	api.Config.SFINAuth = sfinAuth
	if err := config.SaveAuth(sfinAuth); err != nil {
		return err
	}

	return nil
}
