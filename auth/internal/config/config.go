package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string   `yaml:"env"         env:"ENV"                env-default:"local"`
	DB         database `yaml:"db"          env:"DATABASE_CONFIG"`
	GRPCServer gRPC     `yaml:"grpc_server" env:"GRPC_SERVER_CONFIG"`
}

type database struct {
	Port     int    `yaml:"port"     env:"DATABASE_PORT"     env-default:"5432"`
	Host     string `yaml:"host"     env:"DATABASE_HOST"     env-default:"localhost"`
	User     string `yaml:"user"     env:"DATABASE_USER"     env-default:"user"`
	Password string `yaml:"password" env:"DATABASE_PASSWORD"`
	Name     string `yaml:"name"     env:"DATABASE_NAME"     env-default:"postgres"`
}

type gRPC struct{}

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
