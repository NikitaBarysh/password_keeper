package server

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ServConfig struct {
	Endpoint             string
	DataBaseDSN          string
	PrivateCryptoKeyPath string
	Salt                 string
	SecretKey            string
}

type Option func(*ServConfig)

func withEndpoint(endpoint string) Option {
	return func(c *ServConfig) {
		c.Endpoint = endpoint
	}
}

func withDataBaseDSN(dsn string) Option {
	return func(c *ServConfig) {
		c.DataBaseDSN = dsn
	}
}

func withPrivateKey(path string) Option {
	return func(c *ServConfig) {
		c.PrivateCryptoKeyPath = path
	}
}

func withSalt(salt string) Option {
	return func(c *ServConfig) {
		c.Salt = salt
	}
}

func withSecretKey(secretKey string) Option {
	return func(c *ServConfig) {
		c.SecretKey = secretKey
	}
}

func NewServConfig(option ...Option) *ServConfig {
	cfg := &ServConfig{
		Endpoint:    "8080",
		DataBaseDSN: "postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable",
		Salt:        "cokdnvosavnsdfm3jr2034v=0wjv=4092h3jv",
		SecretKey:   "ewfwvfsvwgf42",
	}

	for _, opt := range option {
		opt(cfg)
	}

	return cfg
}

func NewServer() *ServConfig {
	var (
		endpoint       string
		database       string
		privateKeyPath string
		salt           string
		secretKey      string
	)

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	flag.StringVar(&endpoint, "a", "8080", "endpoint for server")
	flag.StringVar(&database, "d", "postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable", "url to connect to DB")

	flag.Parse()

	if envEndpoint, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		endpoint = envEndpoint
	}

	if envDataBase, ok := os.LookupEnv("DATABASE"); ok {
		database = envDataBase
	}

	if envPrivateKeyPath, ok := os.LookupEnv("PRIVATE_KEY"); ok {
		privateKeyPath = envPrivateKeyPath
	}

	if envSalt, ok := os.LookupEnv("SALT"); ok {
		salt = envSalt
	}

	if envSecretKey, ok := os.LookupEnv("SECRET_KEY"); ok {
		secretKey = envSecretKey
	}

	cfg := NewServConfig(withEndpoint(endpoint), withDataBaseDSN(database), withPrivateKey(privateKeyPath),
		withSalt(salt), withSecretKey(secretKey))

	return cfg
}
