// Package client -пакет в котором создаем конфиг для клиента
package client

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// ClientConfig - структура, в которую передаем параметры для запуска клиента
type ClientConfig struct {
	Url           string
	PublicKeyPath string
	HashKey       string
}

type Option func(*ClientConfig)

func withUrl(url string) Option {
	return func(c *ClientConfig) {
		c.Url = url
	}
}

func withPublicKey(publicKeyPath string) Option {
	return func(c *ClientConfig) {
		c.PublicKeyPath = publicKeyPath
	}
}

func withHashKey(key string) Option {
	return func(c *ClientConfig) {
		c.HashKey = key
	}
}

// newClientConfig - создает структуру ClientConfig
func newClientConfig(option ...Option) *ClientConfig {
	cfg := &ClientConfig{
		Url:           "localhost:8000",
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
	var (
		url       string
		publicKey string
		hashKey   string
	)

	err := godotenv.Load()
	if err != nil {
		log.Printf("NewClient: Error loading .env file: %s\n", err.Error())
	}

	if envPublicKeyPath, ok := os.LookupEnv("PUBLIC_KEY"); ok {
		publicKey = envPublicKeyPath
	}

	if envUrl, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		url = envUrl
	}

	if envHashKey, ok := os.LookupEnv("HASH_KEY"); ok {
		hashKey = envHashKey
	}

	return newClientConfig(withUrl(url), withPublicKey(publicKey), withHashKey(hashKey))
}
