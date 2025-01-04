package common

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	APP_SERVER_PATH            string
	APP_DATA_FOLDER            string
	APP_PORT                   int
	APP_ADDR                   string
	DB_CONTAINER_IMAGE         string
	DB_CONTAINER_NAME          string
	DB_HOST                    string
	DB_HOST_PORT               int
	DB_CONTAINER_PORT          int
	DB_HOST_VOLUME             string
	DB_NAME                    string
	DB_USER                    string
	DB_PASS                    string
	DB_CONN_TIMEOUT_SEC        int
	DB_DSN                     string
	DB_MIGRATION_PATH          string
	DB_GEN_MODELS_PATH         string
	LOG_LEVEL                  string
	AUTH_JWT_KEY               string
	AUTH_JWT_TOKEN_TIMEOUT_MIN int
}

func LoadConfig(path, nameWithoutExt, extWithoutDot string) *Config {
	logger := log.With().Str("path", path).Str("file_name", fmt.Sprintf("%s.%s", nameWithoutExt, extWithoutDot)).Logger()

	viper.AddConfigPath(path)
	viper.SetConfigName(nameWithoutExt)
	viper.SetConfigType(extWithoutDot)
	viper.AutomaticEnv()

	// Allow nested environment variables using `__` as a separator
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to read config")
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to unmarshal config")
	}

	logger.Info().Msg("loaded config")

	return &config
}
