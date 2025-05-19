package config

import (
	"fmt"
	"github.com/andibalo/payd-test/backend/pkg/logger"
	"github.com/spf13/viper"
)

const (
	AppAddress         = ":8082"
	EnvDevEnvironment  = "DEV"
	EnvProdEnvironment = "PROD"
	ServiceName        = "roster-management-service"
)

type Config interface {
	Logger() logger.Logger

	AppVersion() string
	AppID() string
	AppName() string
	AppEnv() string
	AppAddress() string

	DBConnString() string

	GetAuthCfg() Auth
	GetFlags() Flag
}

type AppConfig struct {
	logger logger.Logger
	App    app
	Db     db
	Flag   Flag
	Auth   Auth
}

type app struct {
	AppEnv      string
	AppVersion  string
	Name        string
	Description string
	AppUrl      string
	AppID       string
}

type db struct {
	DSN      string
	User     string
	Password string
	Name     string
	Host     string
	Port     int
	MaxPool  int
}

type Flag struct {
	EnableSeedDB bool
}

type Auth struct {
	JWTSecret      string
	JWTStaticToken string
}

func InitConfig() *AppConfig {
	viper.SetConfigType("env")
	viper.SetConfigName(".env") // name of Config file (without extension)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	l := logger.GetLogger(logger.Options{
		DefaultFields: map[string]string{
			"service.name":    ServiceName,
			"service.version": viper.GetString("APP_VERSION"),
			"service.env":     viper.GetString("APP_ENV"),
		},
		ContextFields: map[string]string{
			"path":        "path",
			"method":      "method",
			"status_code": "status_code",
			"status":      "status",
			"error":       "error",
			"user_id":     "x-user-id",
			"user_email":  "x-user-email",
			"client_ip":   "x-forwarded-for",
			"payload":     "payload",
			"x-client-id": "x-client-id",
			"topic":       "topic",
			"broker":      "broker",
			"trace.id":    "trace.id",
			"span.id":     "span.id",
		},
		Level:     logger.LevelInfo,
		HookLevel: logger.LevelError,
	})

	if err := viper.ReadInConfig(); err != nil {
		l.Warn("Env config file not found")
	}

	return &AppConfig{
		logger: l,
		App: app{
			AppEnv:      viper.GetString("APP_ENV"),
			AppVersion:  viper.GetString("APP_VERSION"),
			Name:        ServiceName,
			Description: "core service",
			AppUrl:      viper.GetString("APP_URL"),
			AppID:       viper.GetString("APP_ID"),
		},
		Db: db{
			DSN:      getRequiredString("DB_DSN"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			Name:     viper.GetString("DB_NAME"),
			MaxPool:  viper.GetInt("DB_MAX_POOLING_CONNECTION"),
		},
		Flag: Flag{
			EnableSeedDB: viper.GetBool("ENABLE_SEED_DB"),
		},
		Auth: Auth{
			JWTSecret:      viper.GetString("JWT_SECRET"),
			JWTStaticToken: viper.GetString("JWT_STATIC_TOKEN"),
		},
	}
}

func getRequiredString(key string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}

	panic(fmt.Errorf("KEY %s IS MISSING", key))
}

func (c *AppConfig) Logger() logger.Logger {
	return c.logger
}

func (c *AppConfig) AppVersion() string {
	return c.App.AppVersion
}

func (c *AppConfig) AppID() string {
	return c.App.AppID
}

func (c *AppConfig) AppName() string {
	return c.App.Name
}

func (c *AppConfig) AppEnv() string {
	return c.App.AppEnv
}

func (c *AppConfig) AppAddress() string {
	return AppAddress
}

func (c *AppConfig) DBConnString() string {
	return c.Db.DSN
}

func (c *AppConfig) GetFlags() Flag {
	return c.Flag
}

func (c *AppConfig) GetAuthCfg() Auth {
	return c.Auth
}
