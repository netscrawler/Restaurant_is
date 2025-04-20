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
	Env          string         `yaml:"env"          env:"ENV"                env-default:"local"`
	DB           DatabaseConfig `yaml:"db"           env:"DATABASE_CONFIG"`
	GRPCServer   gRPC           `yaml:"grpcServer"   env:"GRPC_SERVER_CONFIG"`
	YandexOAuth  YandexOAuth    `yaml:"yandexOAuth"  env:"YANDEX_O_AUTH"`
	JWT          JWTConfig      `yaml:"jwt"`
	JWTRaw       JWTConfigRaw   `yaml:"jwtRAW"`
	CodeLife     time.Duration  `yaml:"codeLife"                              env-default:"5m"`
	NotifyClient NotifyClient   `yaml:"notifyClient"`
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
	Address           string        `yaml:"address"`
	BaseDelay         time.Duration `yaml:"baseDelay"`
	Multiplier        float64       `yaml:"multiplier"`
	MaxDelay          time.Duration `yaml:"maxDelay"`
	MinConnectTimeout time.Duration `yaml:"minConnectTimeout"`
}

type gRPC struct {
	Address string `yaml:"address" env:"address" env-default:"address"`
	Port    int    `yaml:"port"    env:"port"`
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
	PrivateKey        string        `yaml:"privateKey"`
	PublicKey         string        `yaml:"publicKey"`
	RefreshPrivateKey string        `yaml:"refreshPrivateKey"`
	RefreshPublicKey  string        `yaml:"refreshPublicKey"`
	AccessTTL         time.Duration `yaml:"accessTtl"`
	RefreshTTL        time.Duration `yaml:"refreshTtl"`
	Issuer            string        `yaml:"issuer"`
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
