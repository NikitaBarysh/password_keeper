// Package server -пакет в котором создаем конфиг для клиента
package client

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// ClientConfig - структура, в которую передаем параметры для запуска клиента
type ClientConfig struct {
	Host          string `envconfig:"CLIENT_ENDPOINT"`
	Port          string `envconfig:"CLIENT_PORT"`
	PublicKeyPath string `envconfig:"PUBLIC_KEY"`
	HashKey       string `envconfig:"HASH_KEY"`
}

type Option func(*ClientConfig)

func WithUrl(url string) Option {
	return func(c *ClientConfig) {
		c.Host = url
	}
}

func WithPublicKey(publicKeyPath string) Option {
	return func(c *ClientConfig) {
		c.PublicKeyPath = publicKeyPath
	}
}

func WithHashKey(key string) Option {
	return func(c *ClientConfig) {
		c.HashKey = key
	}
}

// NewClientConfig - создает структуру ClientConfig
func NewClientConfig(option ...Option) *ClientConfig {
	cfg := &ClientConfig{
		Host:          "localhost:8000",
		PublicKeyPath: "public.rsa",
		HashKey:       "cm2984yf2v08ji23r0vhwssdkmvs",
	}

	for _, opt := range option {
		opt(cfg)
	}

	return cfg
}

// NewClient - загружаем данные из переменных окружения или проставляем дефолтные и возвращаем готовый конфиг
func NewClient() *ClientConfig {
	err := godotenv.Load()
	if err != nil {
		log.Printf("NewClient: Error loading .env file: %s\n", err.Error())
	}

	cfg := &ClientConfig{}

	if err := envconfig.Process("", cfg); err != nil {
		log.Printf("NewClient: Error process vars: %s\n", err.Error())
	}

	return cfg
}
