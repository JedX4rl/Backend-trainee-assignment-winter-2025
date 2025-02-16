package serverConfig

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"os"
	"strconv"
	"time"
)

const (
	CONFIG_SERVER_PATH       = "CONFIG_SERVER_PATH"
	ACCESS_TOKEN_SECRET      = "ACCESS_TOKEN_SECRET"
	ACCESS_TOKEN_EXPIRY_HOUR = "ACCESS_TOKEN_EXPIRY_HOUR"
)

type ServerConfig struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	HttpServer  `yaml:"http_server" env-required:"true"`
	SecretKey   string
	TokenExpiry int
}

type HttpServer struct {
	Address string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut time.Duration `yaml:"timeout" env-default:"4s"`
}

func MustLoadServerConfig() (*ServerConfig, error) {

	slog.Debug("Loading server config")

	configPath := os.Getenv(CONFIG_SERVER_PATH)
	if configPath == "" {
		return nil, fmt.Errorf("%s environment variable not set", CONFIG_SERVER_PATH)
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s does not exist %s", CONFIG_SERVER_PATH, configPath)
	}

	var config ServerConfig

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, fmt.Errorf("cannot load config file: %s", err)
	}

	if config.SecretKey = os.Getenv(ACCESS_TOKEN_SECRET); config.SecretKey == "" {
		return nil, fmt.Errorf("cannot load config file, secret key is missing")
	}
	tokenExpiry, err := strconv.Atoi(os.Getenv(ACCESS_TOKEN_EXPIRY_HOUR))
	if err != nil {
		return nil, fmt.Errorf("cannot load config file, expiry is missing")
	}
	config.TokenExpiry = tokenExpiry

	return &config, nil
}
