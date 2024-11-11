package common

import (
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

func LoadConfig(path, nameWithoutExtension, extension string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(nameWithoutExtension)
	viper.SetConfigType(extension)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return viper.Unmarshal(&Configs)
}
