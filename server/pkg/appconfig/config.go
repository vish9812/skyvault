package appconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Auth   AuthConfig
	Media  MediaConfig
	Log    LogConfig
}

type ServerConfig struct {
	Path    string
	DataDir string
	Port    int
	Addr    string
}

type DBContainerConfig struct {
	Image string
	Name  string
	Port  int
}

type DBHostConfig struct {
	Name   string
	Port   int
	Volume string
}

type DBConfig struct {
	Container      DBContainerConfig
	Host           DBHostConfig
	Name           string
	User           string
	Pass           string
	ConnTimeoutSec int
	DSN            string
}

type JWTConfig struct {
	Key             []byte
	TokenTimeoutMin int
}

type AuthConfig struct {
	JWT JWTConfig
}

type MediaConfig struct {
	MaxUploadSizeMB       int64 // Max upload size even when including chunking strategy.
	MaxDirectUploadSizeMB int64 // Max size allowed for an upload, before chunking strategy is applied. This value must be less than MaxUploadSizeMB.
	MaxChunkSizeMB        int64 // Max size of a chunk. This value must be less than MaxDirectUploadSizeMB.
}

type LogConfig struct {
	Level string
}

// LoadConfig loads configuration from the specified environment file
func LoadConfig(path string, isDev bool) *Config {
	logger := initLogger(isDev)

	// Read env file
	envMap, err := godotenv.Read(path)
	if err != nil {
		logger.Fatal().Err(err).Msgf("error loading env file at %s", path)
	}

	config := &Config{}

	// Server config
	config.Server.Path = envMap["SERVER__PATH"]
	config.Server.DataDir = envMap["SERVER__DATA_DIR"]
	config.Server.Port = getIntOrZero(envMap["SERVER__PORT"])
	config.Server.Addr = envMap["SERVER__ADDR"]

	// DB config
	config.DB.Container.Image = envMap["DB__CONTAINER__IMAGE"]
	config.DB.Container.Name = envMap["DB__CONTAINER__NAME"]
	config.DB.Container.Port = getIntOrZero(envMap["DB__CONTAINER__PORT"])
	config.DB.Host.Name = envMap["DB__HOST__NAME"]
	config.DB.Host.Port = getIntOrZero(envMap["DB__HOST__PORT"])
	config.DB.Host.Volume = envMap["DB__HOST__VOLUME"]
	config.DB.Name = envMap["DB__NAME"]
	config.DB.User = envMap["DB__USER"]
	config.DB.Pass = envMap["DB__PASS"]
	config.DB.ConnTimeoutSec = getIntOrZero(envMap["DB__CONN_TIMEOUT_SEC"])
	config.DB.DSN = envMap["DB__DSN"]

	// Auth config
	config.Auth.JWT.Key = []byte(envMap["AUTH__JWT__KEY"])
	config.Auth.JWT.TokenTimeoutMin = getIntOrZero(envMap["AUTH__JWT__TOKEN_TIMEOUT_MIN"])

	// Media config
	config.Media.MaxUploadSizeMB = getInt64OrZero(envMap["MEDIA__MAX_UPLOAD_SIZE_MB"])
	config.Media.MaxDirectUploadSizeMB = getInt64OrZero(envMap["MEDIA__MAX_DIRECT_UPLOAD_SIZE_MB"])
	config.Media.MaxChunkSizeMB = getInt64OrZero(envMap["MEDIA__MAX_CHUNK_SIZE_MB"])

	// Log config
	config.Log.Level = envMap["LOG__LEVEL"]

	config.validate(logger, isDev)

	return config
}

// Helper functions to safely convert strings to numbers
func getIntOrZero(s string) int {
	if s == "" {
		return 0
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

func getInt64OrZero(s string) int64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

func getFloat64OrZero(s string) float64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func getBoolOrFalse(s string) bool {
	if s == "" {
		return false
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return v
}

// validate use default values if possible, otherwise log error
func (c *Config) validate(logger zerolog.Logger, isDev bool) {
	var foundErr bool

	// Server
	if c.Server.Path == "" {
		// Use current dir
		c.Server.Path = "."
		logger.Warn().Msg("app path not set, using current dir")
	}

	if c.Server.DataDir == "" {
		c.Server.DataDir = filepath.Join(c.Server.Path, ".data")
		logger.Warn().Msg("data dir not set, using default '.data' dir")
	}

	if c.Server.Port <= 0 {
		c.Server.Port = 8090
		logger.Warn().Msg("server port not set, using default port 8090")
	}

	if c.Server.Addr == "" {
		c.Server.Addr = fmt.Sprintf("localhost:%d", c.Server.Port)
		logger.Warn().Msg("server addr not set, using default addr localhost:8090")
	}

	// Database
	if c.DB.ConnTimeoutSec <= 0 {
		c.DB.ConnTimeoutSec = 30
		logger.Warn().Msg("database connection timeout not set, using default timeout 30 seconds")
	}

	if c.DB.Container.Image == "" {
		c.DB.Container.Image = "postgres:16"
		logger.Warn().Msg("database container image not set, using default image postgres:16")
	}

	if c.DB.Container.Name == "" {
		c.DB.Container.Name = "skyvault-db"
		logger.Warn().Msg("database container name not set, using default name skyvault-db")
	}

	if c.DB.Container.Port <= 0 {
		c.DB.Container.Port = 5432
		logger.Warn().Msg("database container port not set, using default port 5432")
	}

	if c.DB.Host.Name == "" {
		c.DB.Host.Name = "localhost"
		logger.Warn().Msg("database host name not set, using default localhost")
	}

	if c.DB.Host.Port <= 0 {
		c.DB.Host.Port = 5432
		logger.Warn().Msg("database host port not set, using default port 5432")
	}

	if c.DB.Host.Volume == "" {
		c.DB.Host.Volume = filepath.Join(c.Server.DataDir, "db")
		logger.Warn().Msgf("database host volume not set, using default volume: %s", c.DB.Host.Volume)
	}

	if c.DB.Name == "" {
		c.DB.Name = "skyvault"
		logger.Warn().Msg("database name not set, using default name skyvault")
	}

	if c.DB.User == "" {
		c.DB.User = "skyvault"
		logger.Warn().Msg("database user not set, using default user skyvault")
	}

	if c.DB.Pass == "" {
		if isDev {
			c.DB.Pass = "skyvault"
			logger.Info().Msg("database password not set, using default password skyvault")
		} else {
			foundErr = true
			logger.Error().Msg("database password must be set")
		}
	}

	if c.DB.DSN == "" {
		c.DB.DSN = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&connect_timeout=%d",
			c.DB.User, c.DB.Pass, c.DB.Host.Name, c.DB.Host.Port, c.DB.Name, c.DB.ConnTimeoutSec)
		logger.Info().Msgf("database DSN not set, constructing from other db settings")
	}

	// Auth
	if len(c.Auth.JWT.Key) < 32 {
		if isDev {
			c.Auth.JWT.Key = []byte("skyvault")
			logger.Info().Msg("JWT key length is less than 32, using default key skyvault")
		} else {
			foundErr = true
			logger.Error().Msg("JWT key must be at least 32 characters long")
		}
	}

	if c.Auth.JWT.TokenTimeoutMin <= 0 {
		c.Auth.JWT.TokenTimeoutMin = 10 // 10 minutes
		if isDev {
			c.Auth.JWT.TokenTimeoutMin = 60 * 24 * 30 // 30 days
		}
		logger.Warn().Msgf("JWT token timeout not set, using default timeout %d minutes", c.Auth.JWT.TokenTimeoutMin)
	}

	// Media
	if c.Media.MaxUploadSizeMB <= 0 {
		c.Media.MaxUploadSizeMB = 100
		if isDev {
			c.Media.MaxUploadSizeMB = 1024 // 1GB
		}
		logger.Warn().Msgf("media max size not set, using default size %dMB", c.Media.MaxUploadSizeMB)
	}

	if c.Media.MaxDirectUploadSizeMB > c.Media.MaxUploadSizeMB {
		c.Media.MaxDirectUploadSizeMB = 50
		logger.Warn().Msgf("media max direct upload size is greater than max upload size, using default size %dMB", c.Media.MaxDirectUploadSizeMB)
	}

	if c.Media.MaxDirectUploadSizeMB <= 0 {
		c.Media.MaxDirectUploadSizeMB = 50
		logger.Warn().Msgf("media max single upload size not set, using default size %dMB", c.Media.MaxDirectUploadSizeMB)
	}

	if c.Media.MaxChunkSizeMB <= 0 {
		c.Media.MaxChunkSizeMB = 10
		logger.Warn().Msgf("media max chunk size not set, using default size %dMB", c.Media.MaxChunkSizeMB)
	}

	if c.Media.MaxChunkSizeMB > c.Media.MaxDirectUploadSizeMB {
		c.Media.MaxChunkSizeMB = 10
		logger.Warn().Msgf("media max chunk size is greater than max direct upload size, using default size %dMB", c.Media.MaxChunkSizeMB)
	}

	// Logging
	if c.Log.Level == "" {
		c.Log.Level = "info"
		if isDev {
			c.Log.Level = "debug"
		}
		logger.Warn().Msgf("log level not set, using default level %s", c.Log.Level)
	}

	if foundErr {
		logger.Fatal().Msg("config validation failed: fix the errors and try again")
	}
}

// Just to log while loading config.
// Keep it same as in applog/logger.go and main.go to keep logs consistent
// and to avoid circular imports.
func initLogger(isDev bool) zerolog.Logger {
	// Configure zerolog
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	var level zerolog.Level
	if isDev {
		level = zerolog.DebugLevel
	} else {
		level = zerolog.InfoLevel
	}

	return zerolog.New(output).
		Level(level).
		With().
		Str("where", "appconfig").
		Timestamp().
		Logger()
}
