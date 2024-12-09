package resource

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"rudder/config"
	"rudder/models"
	"strconv"
	"strings"
	"time"
)

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

func (api SimpleFINAPI) GetAccounts(days int, sfinResp *models.SimpleFINResponse) error {
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
