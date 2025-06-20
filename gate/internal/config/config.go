package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server    ServerConfig    `env:"SERVER"   yaml:"server"`
	Services  ServicesConfig  `env:"SERVICES" yaml:"services"`
	Auth      AuthConfig      `env:"AUTH"     yaml:"auth"`
	Telemetry TelemertyConfig `               yaml:"telemetry"`
}

type ServerConfig struct {
	Port int    `env:"PORT" env-default:"8080"    yaml:"port"`
	Host string `env:"HOST" env-default:"0.0.0.0" yaml:"host"`
}

type ServicesConfig struct {
	Auth   ServiceConfig `env:"AUTH"   yaml:"auth"`
	User   ServiceConfig `env:"USER"   yaml:"user"`
	Menu   ServiceConfig `env:"MENU"   yaml:"menu"`
	Order  ServiceConfig `env:"ORDER"  yaml:"order"`
	Notify ServiceConfig `env:"NOTIFY" yaml:"notify"`
}

type ServiceConfig struct {
	Host string `env:"HOST" env-default:"localhost" yaml:"host"`
	Port int    `env:"PORT"                         yaml:"port"`
}

type AuthConfig struct {
	JWTSecret string `env:"JWT_SECRET" env-default:"your-secret-key" yaml:"jwt_secret"`
}

type TelemertyConfig struct {
	ServiceName    string `yaml:"serviceName"    env:"SERVICE_NAME"`
	ServiceVersion string `yaml:"serviceVersion" env:"SERVICE_VERSION"`
	Environment    string `yaml:"environment"    env:"ENVIRONMENT"`
	MetricsPort    int    `yaml:"metricsPort"    env:"METRICS_PORT"`
	TraceEndpoint  string `yaml:"traceEndpoint"  env:"TRACE_ENDPOINT"  env-default:"localhost:4317"`
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

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return cfg, nil
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

func (c *Config) GetServiceAddress(service string) string {
	var serviceConfig ServiceConfig

	switch service {
	case "auth":
		serviceConfig = c.Services.Auth
	case "user":
		serviceConfig = c.Services.User
	case "menu":
		serviceConfig = c.Services.Menu
	case "order":
		serviceConfig = c.Services.Order
	default:
		return ""
	}

	return fmt.Sprintf("%s:%d", serviceConfig.Host, serviceConfig.Port)
}
