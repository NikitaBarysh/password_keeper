// Package client - это пакет c конфигом для сервера
package server

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// ServConfig - структура, в которую передаем параметры для запуска сервера
type ServConfig struct {
	Endpoint             string `envconfig:"SERVER_ENDPOINT"`
	DataBaseDSN          string `envconfig:"DATABASE_DSN"`
	PrivateCryptoKeyPath string `envconfig:"PRIVATE_KEY"`
	PublicCryptoKeyPath  string `envconfig:"PUBLIC_KEY"`
	Salt                 string `envconfig:"SALT"`
	SecretKey            string `envconfig:"SECRET_KEY"`
	DBHost               string `envconfig:"DB_HOST"`
	DBPort               string `envconfig:"DB_PORT"`
	DBUsername           string `envconfig:"DB_USERNAME"`
	DBPassword           string `envconfig:"DB_PASSWORD"`
	DBDatabase           string `envconfig:"DB_DATABASE"`
}

type Option func(*ServConfig)

func WithEndpoint(endpoint string) Option {
	return func(c *ServConfig) {
		c.Endpoint = endpoint
	}
}

func WithDBAddress(host string) Option {
	return func(c *ServConfig) {
		c.DBHost = host
	}
}

func WithDBPort(port string) Option {
	return func(c *ServConfig) {
		c.DBPort = port
	}
}

func WithDBUsername(username string) Option {
	return func(c *ServConfig) {
		c.DBUsername = username
	}
}

func WithDBPassword(password string) Option {
	return func(c *ServConfig) {
		c.DBPassword = password
	}
}

func WithDBDatabase(database string) Option {
	return func(c *ServConfig) {
		c.DBDatabase = database
	}
}

func WithPublicCryptoKeyPath(path string) Option {
	return func(c *ServConfig) {
		c.PublicCryptoKeyPath = path
	}
}

func WithDataBaseDSN(dsn string) Option {
	return func(c *ServConfig) {
		c.DataBaseDSN = dsn
	}
}

func WithPrivateKey(path string) Option {
	return func(c *ServConfig) {
		c.PrivateCryptoKeyPath = path
	}
}

func WithSalt(salt string) Option {
	return func(c *ServConfig) {
		c.Salt = salt
	}
}

func WithSecretKey(secretKey string) Option {
	return func(c *ServConfig) {
		c.SecretKey = secretKey
	}
}

// NewServConfig - создает структуру ServConfig, (еще используеться для тестов)
func NewServConfig(option ...Option) *ServConfig {
	cfg := &ServConfig{
		Endpoint:  "8080",
		Salt:      "cokdnvosavnsdfm3jr2034v=0wjv=4092h3jv",
		SecretKey: "fad2osfj239vpsdlvmpKJV",
	}

	for _, opt := range option {
		opt(cfg)
	}

	return cfg
}

// NewServer - загружаем данные из переменных окружения или проставляем дефолтные
func NewServer() *ServConfig {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("NewServer: Error loading .env file: %s\n", err.Error())
	}

	cfg := &ServConfig{}
	if err := envconfig.Process("", cfg); err != nil {
		log.Printf("NewServer: Error process vars: %s\n", err.Error())
	}

	return cfg
}
