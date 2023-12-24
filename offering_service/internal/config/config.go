package config

import (
	"encoding/json"
	"offering_service/internal/models"
	"os"
)

const configPath = "./config/config.json"

func InitConfig() (*models.Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg models.Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
