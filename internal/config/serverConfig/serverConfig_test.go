package serverConfig

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadServerConfig_Success(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	configData := `env: "test"
http_server:
  address: "127.0.0.1:8080"
  timeout: 5s`
	_, err = tempFile.Write([]byte(configData))
	require.NoError(t, err)
	tempFile.Close()

	// Устанавливаем переменные окружения
	t.Setenv(CONFIG_SERVER_PATH, tempFile.Name())
	t.Setenv(ACCESS_TOKEN_SECRET, "test_secret")
	t.Setenv(ACCESS_TOKEN_EXPIRY_HOUR, "24")

	config, err := MustLoadServerConfig()

	assert.NoError(t, err)
	assert.Equal(t, "test", config.Env)
	assert.Equal(t, "127.0.0.1:8080", config.HttpServer.Address)
	assert.Equal(t, 5*time.Second, config.HttpServer.TimeOut)
	assert.Equal(t, "test_secret", config.SecretKey)
	assert.Equal(t, 24, config.TokenExpiry)
}

func TestLoadServerConfig_MissingConfigPath(t *testing.T) {
	t.Setenv(ACCESS_TOKEN_SECRET, "test_secret")
	t.Setenv(ACCESS_TOKEN_EXPIRY_HOUR, "24")

	config, err := MustLoadServerConfig()

	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestLoadServerConfig_MissingSecretKey(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	configData := `env: "test"
http_server:
  address: "127.0.0.1:8080"
  timeout: 5s`
	_, err = tempFile.Write([]byte(configData))
	require.NoError(t, err)
	tempFile.Close()

	t.Setenv(CONFIG_SERVER_PATH, tempFile.Name())
	t.Setenv(ACCESS_TOKEN_EXPIRY_HOUR, "24")

	config, err := MustLoadServerConfig()

	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestLoadServerConfig_FileNotFound(t *testing.T) {
	t.Setenv(CONFIG_SERVER_PATH, "nonexistent.yaml")
	t.Setenv(ACCESS_TOKEN_SECRET, "test_secret")
	t.Setenv(ACCESS_TOKEN_EXPIRY_HOUR, "24")

	config, err := MustLoadServerConfig()

	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestLoadServerConfig_InvalidTokenExpiry(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	configData := `env: "test"
http_server:
  address: "127.0.0.1:8080"
  timeout: 5s`
	_, err = tempFile.Write([]byte(configData))
	require.NoError(t, err)
	tempFile.Close()

	t.Setenv(CONFIG_SERVER_PATH, tempFile.Name())
	t.Setenv(ACCESS_TOKEN_SECRET, "test_secret")
	t.Setenv(ACCESS_TOKEN_EXPIRY_HOUR, "invalid_value") // Некорректное значение для времени

	config, err := MustLoadServerConfig()

	assert.Error(t, err)
	assert.Nil(t, config)
}
