package storageConfig

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestLoadStorageConfig_Success(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	configData := `host: "localhost"
port: "5432"
username: "user"
password: "password"
db_name: "mydb"
ssl_mode: "disable"`
	_, err = tempFile.Write([]byte(configData))
	require.NoError(t, err)
	tempFile.Close()

	t.Setenv(CONFIG_STORAGE_PATH, tempFile.Name())

	config, err := MustLoadStorageConfig()

	assert.NoError(t, err)
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, "5432", config.Port)
	assert.Equal(t, "user", config.Username)
	assert.Equal(t, "password", config.Password)
	assert.Equal(t, "mydb", config.DBName)
	assert.Equal(t, "disable", config.SSLMode)
}

func TestLoadStorageConfig_MissingConfigPath(t *testing.T) {
	config, err := MustLoadStorageConfig()

	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestLoadStorageConfig_FileNotFound(t *testing.T) {
	t.Setenv(CONFIG_STORAGE_PATH, "nonexistent.yaml")
	config, err := MustLoadStorageConfig()

	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestLoadStorageConfig_InvalidConfig(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	configData := `host: "localhost"
port: "5432"`
	_, err = tempFile.Write([]byte(configData))
	require.NoError(t, err)
	tempFile.Close()

	t.Setenv(CONFIG_STORAGE_PATH, tempFile.Name())

	config, err := MustLoadStorageConfig()

	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestLoadStorageConfig_MissingRequiredField(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	configData := `host: "localhost"
port: "5432"
username: "user"`
	_, err = tempFile.Write([]byte(configData))
	require.NoError(t, err)
	tempFile.Close()

	t.Setenv(CONFIG_STORAGE_PATH, tempFile.Name())

	config, err := MustLoadStorageConfig()

	assert.Error(t, err)
	assert.Nil(t, config)
}
