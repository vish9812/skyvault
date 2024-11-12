package common

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var Configs *config

type config struct {
	APP_SERVER_PATH    string
	DB_CONTAINER_IMAGE string
	DB_CONTAINER_NAME  string
	DB_HOST            string
	DB_HOST_PORT       int
	DB_CONTAINER_PORT  int
	DB_NAME            string
	DB_USER            string
	DB_PASS            string
	DB_CONN_STR        string
	DB_MIGRATION_PATH  string
}

func LoadConfig(path, nameWithoutExtension, extensionWithoutDot string) {
	logger := log.With().Str("path", path).Str("file_name", fmt.Sprintf("%s.%s", nameWithoutExtension, extensionWithoutDot)).Logger()

	viper.AddConfigPath(path)
	viper.SetConfigName(nameWithoutExtension)
	viper.SetConfigType(extensionWithoutDot)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to read config")
	}

	err = viper.Unmarshal(&Configs)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to unmarshal config")
	}

	logger.Info().Msg("loaded config")
}
