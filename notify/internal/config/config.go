package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string          `yaml:"env"`
	Bot           Bot             `yaml:"bot"`
	GRRPC         GRPC            `yaml:"grpc"`
	Shutdown      time.Duration   `yaml:"shutdown"       env:"SHUTDOWN_TIMEOUT" env-default:"5s"`
	StubRecipient int64           `yaml:"stub_recipient"`
	Telemetry     TelemertyConfig `yaml:"telemetry"`
}
type Bot struct {
	TelegramToken string        `yaml:"telegram_token" env:"TELEGRAM_TOKEN"`
	BotPoll       time.Duration `yaml:"bot_poll"       env:"BOT_POLL_TIMEOUT" env-default:"30s"`
}

type GRPC struct {
	Port int `yaml:"port" env:"port"`
}

type TelemertyConfig struct {
	ServiceName    string `yaml:"serviceName"    env:"SERVICE_NAME"`
	ServiceVersion string `yaml:"serviceVersion" env:"SERVICE_VERSION"`
	Environment    string `yaml:"environment"    env:"ENVIRONMENT"`
	MetricsPort    int    `yaml:"metricsPort"    env:"METRICS_PORT"`
	TraceEndpoint  string `yaml:"traceEndpoint"  env:"TRACE_ENDPOINT"  env-default:"localhost:4317"`
}

// Load загружает конфигурацию из файла или из переменных окружения
func Load() (*Config, error) {
	var cfg Config

	configPath := fetchConfigPath()

	// Загрузка конфигурации
	if configPath != "" {
		// Если путь к файлу указан, загружаем из YAML
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			return nil, fmt.Errorf("error readYaml config: %w", err)
		}
	} else {
		// Если путь не указан, загружаем из переменных окружения
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, fmt.Errorf("error readEnv config: %w", err)
		}
	}

	return &cfg, nil
}

// fetchConfigPath определяет путь к файлу конфигурации
// Приоритет: 1) аргумент командной строки, 2) переменная окружения, 3) значение по умолчанию
func fetchConfigPath() string {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Путь к файлу конфигурации")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	if configPath != "" {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return ""
		}
	}

	return configPath
}
