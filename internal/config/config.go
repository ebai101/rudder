package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type SimpleFINAuth struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AppConfig struct {
	Timezone        string `yaml:"timezone"`
	Location        *time.Location
	SfinBridgeToken string `yaml:"sfin_bridge_token"`
	DatabaseUrl     string `yaml:"db_url"`
	SetupPullDays   int    `yaml:"setup_pull_days"`
	WeeklyPullDays  int    `yaml:"weekly_pull_days"`
	DailyPullDays   int    `yaml:"daily_pull_days"`
	HourlyPullDays  int    `yaml:"hourly_pull_days"`
	SendRequests    bool   `yaml:"send_requests"`
	SaveCache       bool   `yaml:"save_cache"`
	SFINAuth        SimpleFINAuth
}

func loadAuth(sfinAuth *SimpleFINAuth) error {
	file, err := os.Open("sfin_auth.json")
	if err != nil {
		return fmt.Errorf("error opening auth file: %v", err)
	}
	defer file.Close()

	fileContents, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading auth file: %v", err)
	}

	if err := json.Unmarshal(fileContents, sfinAuth); err != nil {
		return err
	}

	return nil
}

func loadLocation(timezone string) (*time.Location, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("error loading config timezone: %v", err)
	}
	return loc, nil
}

func SaveAuth(sfinAuth SimpleFINAuth) error {
	sfinAuthJSON, err := json.Marshal(sfinAuth)
	if err != nil {
		return fmt.Errorf("error formatting auth json: %v", err)
	}

	if err := os.WriteFile("sfin_auth.json", sfinAuthJSON, 0644); err != nil {
		return fmt.Errorf("error writing auth file: %v", err)
	}

	return nil
}

func LoadConfig() (*AppConfig, error) {
	file, _ := os.Open("config.yml")
	defer file.Close()
	decoder := yaml.NewDecoder(file)

	var appConfig AppConfig
	if err := decoder.Decode(&appConfig); err != nil {
		return nil, fmt.Errorf("error loading main config: %v", err)
	}

	if err := loadAuth(&appConfig.SFINAuth); err != nil {
		return nil, err
	}

	loc, err := loadLocation(appConfig.Timezone)
	if err != nil {
		return nil, err
	}
	appConfig.Location = loc

	return &appConfig, nil
}
