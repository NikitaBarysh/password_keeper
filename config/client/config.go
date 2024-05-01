package client

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

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

func newClientConfig(option ...Option) *ClientConfig {
	cfg := &ClientConfig{Url: "localhost:8080"}

	for _, opt := range option {
		opt(cfg)
	}

	return cfg
}

func NewClient() *ClientConfig {
	var (
		url       string
		publicKey string
		hashKey   string
	)

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
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
