package config

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/netscrawler/Restaurant_is/auth/internal/pkg"
)

type Config struct {
	Env          string          `yaml:"env"          env:"ENV"                env-default:"local"`
	DB           DatabaseConfig  `yaml:"db"           env:"DATABASE_CONFIG"    env-default:"db"`
	GRPCServer   gRPC            `yaml:"grpcServer"   env:"GRPC_SERVER_CONFIG" env-default:"grpc_server"`
	YandexOAuth  YandexOAuth     `yaml:"yandexOAuth"  env:"YANDEX_O_AUTH"      env-default:"yandex_o_auth"`
	JWT          JWTConfig       `yaml:"jwt"          env:"JWT"                env-default:"jwt"`
	JWTRaw       JWTConfigRaw    `yaml:"jwtRAW"       env:"JWT_RAW"            env-default:"jwt_raw"`
	CodeLife     time.Duration   `yaml:"codeLife"     env:"CODE_LIFE"          env-default:"5m"`
	NotifyClient NotifyClient    `yaml:"notifyClient" env:"NOTIFY_CLIENT"      env-default:"notify_client"`
	Telemetry    TelemertyConfig `yaml:"telemetry"    env:"TELEMETRY"          env-default:"telemetry"`
	Kafka        Kafka           `yaml:"kafka"        env:"KAFKA"              env-default:"kafka"`
}

type Kafka struct {
	Brokers         []string `yaml:"brokers"         env:"KAFKA_BROKERS"          env-default:"localhost:9092"`
	Topic           string   `yaml:"topic"           env:"KAFKA_TOPIC"            env-default:"events"`
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

type NotifyClient struct {
	Address           string        `yaml:"address"           env:"ADDRESS"`
	BaseDelay         time.Duration `yaml:"baseDelay"         env:"BASE_DELAY"`
	Multiplier        float64       `yaml:"multiplier"        env:"MULTIPLIER"`
	MaxDelay          time.Duration `yaml:"maxDelay"          env:"MAX_DELAY"`
	MinConnectTimeout time.Duration `yaml:"minConnectTimeout" env:"MIN_CONNECT_TIMEOUT"`
}

type gRPC struct {
	Address string `yaml:"address" env:"ADDRESS" env-default:"address"`
	Port    int    `yaml:"port"    env:"PORT"`
}
type YandexOAuth struct {
	ClientID     string `yaml:"yandexClientId"     env:"YANDEX_CLIENT_ID"`
	ClientSecret string `yaml:"yandexClientSecret" env:"YANDEX_CLIENT_SECRET"`
	RedirectURL  string `yaml:"yandexRedirectUrl"  env:"YANDEX_REDIRECT_URL"`
}
type JWTConfig struct {
	PrivateKey        *rsa.PrivateKey
	PublicKey         *rsa.PublicKey
	RefreshPrivateKey *rsa.PrivateKey
	RefreshPublicKey  *rsa.PublicKey
	AccessTTL         time.Duration
	RefreshTTL        time.Duration
	Issuer            string
}

type JWTConfigRaw struct {
	PrivateKey        string        `yaml:"privateKey"        env:"PRIVATE_KEY"`
	PublicKey         string        `yaml:"publicKey"         env:"PUBLIC_KEY"`
	RefreshPrivateKey string        `yaml:"refreshPrivateKey" env:"REFRESH_PRIVATE_KEY"`
	RefreshPublicKey  string        `yaml:"refreshPublicKey"  env:"REFRESH_PUBLIC_KEY"`
	AccessTTL         time.Duration `yaml:"accessTtl"         env:"ACCESS_TTL"`
	RefreshTTL        time.Duration `yaml:"refreshTtl"        env:"REFRESH_TTL"`
	Issuer            string        `yaml:"issuer"            env:"ISSUER"`
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

	err := NewJWTConfig(cfg.JWTRaw, &cfg.JWT)
	if err != nil {
		panic(err.Error())
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

func NewJWTConfig(raw JWTConfigRaw, out *JWTConfig) error {
	// Чтение приватного ключа из файла
	privData, err := os.ReadFile(raw.PrivateKey)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла с PrivateKey: %w", err)
	}

	priv, err := pkg.ParseRSAPrivateKey(privData)
	if err != nil {
		return fmt.Errorf("ошибка парсинга PrivateKey: %w", err)
	}

	// Чтение публичного ключа из файла
	pubData, err := os.ReadFile(raw.PublicKey)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла с PublicKey: %w", err)
	}

	pub, err := pkg.ParseRSAPublicKey(pubData)
	if err != nil {
		return fmt.Errorf("ошибка парсинга PublicKey: %w", err)
	}

	// Чтение приватного ключа для Refresh токенов
	refreshPrivData, err := os.ReadFile(raw.RefreshPrivateKey)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла с RefreshPrivateKey: %w", err)
	}

	refreshPriv, err := pkg.ParseRSAPrivateKey(refreshPrivData)
	if err != nil {
		return fmt.Errorf("ошибка парсинга RefreshPrivateKey: %w", err)
	}

	// Чтение публичного ключа для Refresh токенов
	refreshPubData, err := os.ReadFile(raw.RefreshPublicKey)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла с RefreshPublicKey: %w", err)
	}

	refreshPub, err := pkg.ParseRSAPublicKey(refreshPubData)
	if err != nil {
		return fmt.Errorf("ошибка парсинга RefreshPublicKey: %w", err)
	}

	// Присваиваем все значения в структуру JWTConfig
	*out = JWTConfig{
		PrivateKey:        priv,
		PublicKey:         pub,
		RefreshPrivateKey: refreshPriv,
		RefreshPublicKey:  refreshPub,
		AccessTTL:         raw.AccessTTL,
		RefreshTTL:        raw.RefreshTTL,
		Issuer:            raw.Issuer,
	}

	return nil
}
