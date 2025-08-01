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
	Env        string          `yaml:"env"        env:"ENV"                env-default:"local"`
	DB         DatabaseConfig  `yaml:"db"         env:"DATABASE_CONFIG"`
	GRPCServer gRPC            `yaml:"grpcServer" env:"GRPC_SERVER_CONFIG"`
	Telemetry  TelemertyConfig `yaml:"telemetry"`
	Kafka      Kafka           `yaml:"kafka"`
}

type Kafka struct {
	Brokers         []string `yaml:"brokers"         env:"KAFKA_BROKERS"          env-default:"localhost:9092"`
	Topic           string   `yaml:"topic"           env:"KAFKA_TOPIC"            env-default:"events"`
	GroupID         string   `yaml:"groupID"         env:"KAFKA_GROUP_ID"         env-default:"user_service_group"`
	RetryMax        int      `yaml:"retryMax"        env:"KAFKA_RETRY_MAX"        env-default:"5"`
	ReturnSuccesses bool     `yaml:"returnSuccesses" env:"KAFKA_RETURN_SUCCESSES" env-default:"true"`
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
	Address string `yaml:"address" env:"ADDRESS"`
	Port    int    `yaml:"port"    env:"PORT"`
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
