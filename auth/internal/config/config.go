package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string         `yaml:"env"         env:"ENV"                env-default:"local"`
	DB          DatabaseConfig `yaml:"db"          env:"DATABASE_CONFIG"`
	GRPCServer  gRPC           `yaml:"grpcServer"  env:"GRPC_SERVER_CONFIG"`
	YandexOAuth YandexOAuth    `yaml:"yandexOAuth" env:"YANDEX_O_AUTH"`
	JWT         JWT            `yaml:"jwt"         env:"JWT"`
}

type DatabaseConfig struct {
	Type                string        `yaml:"type"                env:"DATABASE_TYPE"          env-default:"postgres"`
	Port                int           `yaml:"port"                env:"DATABASE_PORT"          env-default:"5432"`
	Host                string        `yaml:"host"                env:"DATABASE_HOST"          env-default:"localhost"`
	User                string        `yaml:"user"                env:"DATABASE_USER"          env-default:"user"`
	Password            string        `yaml:"password"            env:"DATABASE_PASSWORD"      env-default:"password"`
	Name                string        `yaml:"name"                env:"DATABASE_NAME"          env-default:"postgres"`
	SSLMode             string        `yaml:"sslMode"             env:"SSL_MODE"               env-default:"false"`
	PoolMaxConn         int           `yaml:"poolMaxConn"         env:"POOL_MAX_CONN"          env-default:"10"`
	PoolMaxConnLifetime time.Duration `yaml:"poolMaxConnLifetime" env:"POOL_MAX_CONN_LIFETIME" env-default:"1h30m"`
}
type gRPC struct {
	Address string `yaml:"address" env:"address" env-default:"address"`
	Port    int    `yaml:port`
}
type YandexOAuth struct {
	ClientID     string `yaml:"yandexClientId"     env:"YANDEX_CLIENT_ID"`
	ClientSecret string `yaml:"yandexClientSecret" env:"YANDEX_CLIENT_SECRET"`
	RedirectURL  string `yaml:"yandexRedirectUrl"  env:"YANDEX_REDIRECT_URL"`
}

type JWT struct {
	Secret          string        `yaml:"jwtSecret"       env:"JWT_SECRET"`
	AccessDuration  time.Duration `yaml:"accessDuration"  env:"JWT_ACCESS_DURATION"  env-default:"15m"`
	RefreshDuration time.Duration `yaml:"refreshDuration" env:"JWT_REFRESH_DURATION" env-default:"720h"` // 30 days
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func (db *DatabaseConfig) GetURL() string {
	encodedPassword := url.QueryEscape(db.Password)

	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s&pool_max_conns=%d&pool_max_conn_lifetime=%s",
		db.Type,
		db.User,
		encodedPassword,
		db.Host,
		db.Port,
		db.Name,
		db.SSLMode,
		db.PoolMaxConn,
		db.PoolMaxConnLifetime.String(),
	)
}
