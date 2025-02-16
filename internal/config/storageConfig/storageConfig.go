package storageConfig

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

const (
	CONFIG_STORAGE_PATH = "CONFIG_STORAGE_PATH"
)

type StorageConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"db_name" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env-required:"true"`
}

func MustLoadStorageConfig() (*StorageConfig, error) {

	configPath := os.Getenv(CONFIG_STORAGE_PATH)
	if configPath == "" {
		return nil, fmt.Errorf("%s environment variable not set", CONFIG_STORAGE_PATH)
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s does not exist %s", CONFIG_STORAGE_PATH, configPath)
	}
	var config StorageConfig

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, fmt.Errorf("cannot load database config file: %s", err)
	}

	return &config, nil
}
