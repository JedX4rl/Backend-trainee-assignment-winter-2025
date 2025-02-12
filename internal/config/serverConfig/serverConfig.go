package serverConfig

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

const (
	ACCESS_TOKEN_EXPIRY_HOUR = "ACCESS_TOKEN_EXPIRY_HOUR"
	CONFIG_SERVER_PATH       = "CONFIG_SERVER_PATH"
)

type ServerConfig struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local"` //TODO remove?
	SecretKey  string `yaml:"secret_key"`
	HttpServer `yaml:"http_server" env-required:"true"`
}

type HttpServer struct {
	Address string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut time.Duration `yaml:"timeout" env-default:"4s"`
}

func MustLoadServerConfig() *ServerConfig {

	configPath := os.Getenv(CONFIG_SERVER_PATH)
	if configPath == "" {
		log.Fatal(CONFIG_SERVER_PATH + " environment variable not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf(CONFIG_SERVER_PATH+" does not exist %s", configPath)
	}

	var config ServerConfig

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Cannot load config file: %s", err)
	}

	return &config
}
