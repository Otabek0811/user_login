package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

const (
	// DebugMode indicates service mode is debug.
	DebugMode = "debug"
	// TestMode indicates service mode is test.
	TestMode = "test"
	// ReleaseMode indicates service mode is release.
	ReleaseMode = "release"

	TimeExpiredAt = time.Hour * 24
)

type Config struct {
	Environment string // debug, test, release

	ServerHost string
	ServerPort string

	PostgresHost           string
	PostgresUser           string
	PostgresDatabase       string
	PostgresPassword       string
	PostgresPort           string
	PostgresMaxConnections int32

	AuthSecretKey string

	DefaultOffset int
	DefaultLimit  int
}

func Load() Config {

	if err := godotenv.Load("./app.env"); err != nil {
		fmt.Println("No .env file found")
	}

	cfg := Config{}

	cfg.ServerHost = cast.ToString(getOrReturnDefaultValue("SERVICE_HOST", "localhost"))
	cfg.ServerPort = cast.ToString(getOrReturnDefaultValue("HTTP_PORT", ":8080"))

	cfg.PostgresHost = cast.ToString(getOrReturnDefaultValue("POSTGRES_HOST", "localhost"))
	cfg.PostgresPort = cast.ToString(getOrReturnDefaultValue("POSTGRES_PORT", 5432))
	cfg.PostgresUser = cast.ToString(getOrReturnDefaultValue("POSTGRES_USER", "postgres"))
	cfg.PostgresPassword = cast.ToString(getOrReturnDefaultValue("POSTGRES_PASSWORD", "1234"))
	cfg.PostgresDatabase = cast.ToString(getOrReturnDefaultValue("POSTGRES_DATABASE", "login"))
	cfg.PostgresMaxConnections = cast.ToInt32(getOrReturnDefaultValue("POSTGRES_MAXCONS", 20))

	cfg.DefaultOffset = cast.ToInt(getOrReturnDefaultValue("OFFSET", 0))
	cfg.DefaultLimit = cast.ToInt(getOrReturnDefaultValue("LIMIT", 10))

	cfg.AuthSecretKey = cast.ToString(getOrReturnDefaultValue("AUTH_SECRET_KEY", "secret"))

	return cfg
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
