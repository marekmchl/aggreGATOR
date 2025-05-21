package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed locating the home directory: %v", err)
	}
	return homeDir + "/" + configFileName, nil

}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("failed reading the config file: %v", err)
	}

	config := Config{}
	if err := json.Unmarshal(fileData, &config); err != nil {
		return Config{}, fmt.Errorf("failed unmarshalling: %v", err)
	}

	return config, nil
}

func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed opening the config file: %v", err)
	}
	defer file.Close()

	configJSON, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed marshalling: %v", err)
	}
	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("failed clearing the file: %v", err)
	}
	_, err = file.Write(configJSON)
	if err != nil {
		return fmt.Errorf("failed writing into the file: %v", err)
	}

	return nil
}

func (cfg Config) SetUser(newUsername string) error {
	cfg.CurrentUserName = newUsername
	err := write(cfg)
	if err != nil {
		return fmt.Errorf("failed setting new username: %v", err)
	}
	return nil
}
