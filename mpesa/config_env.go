package mpesa

import (
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	EnvAPIKey      = "MPESA_API_KEY"
	EnvPublicKey   = "MPESA_PUBLIC_KEY"
	EnvEnvironment = "MPESA_ENVIRONMENT"
	EnvMarket      = "MPESA_MARKET"
	EnvOrigin      = "MPESA_ORIGIN"
	EnvHost        = "MPESA_HOST"
	EnvPort        = "MPESA_PORT"
)

// ConfigFromEnv builds a Config from environment variables, falling back to
// DefaultConfig for values that are not provided.
//
// Supported variables:
//   MPESA_API_KEY
//   MPESA_PUBLIC_KEY
//   MPESA_ENVIRONMENT=sandbox|openapi
//   MPESA_MARKET=vodacomDRC|vodafoneGHA|vodacomTZN|vodacomLES|vodacomMOZ
//   MPESA_ORIGIN
//   MPESA_HOST
//   MPESA_PORT
func ConfigFromEnv() (Config, error) {
	return ConfigFromEnvWithClient(nil)
}

// ConfigFromEnvWithClient is ConfigFromEnv with an explicit HTTP client.
func ConfigFromEnvWithClient(httpClient *http.Client) (Config, error) {
	cfg := DefaultConfig()

	if value := os.Getenv(EnvAPIKey); value != "" {
		cfg.APIKey = value
	}
	if value := os.Getenv(EnvEnvironment); value != "" {
		cfg.Environment = Environment(strings.ToLower(value))
		// Re-select default key for that environment unless MPESA_PUBLIC_KEY overrides below.
		cfg.PublicKey = defaultPublicKeyForEnvironment(cfg.Environment)
	}
	if value := os.Getenv(EnvPublicKey); value != "" {
		cfg.PublicKey = value
	}
	if value := os.Getenv(EnvMarket); value != "" {
		if market, ok := Markets[value]; ok {
			cfg.Market = market
		} else {
			cfg.Market = Market{Context: value}
		}
	}
	if value := os.Getenv(EnvOrigin); value != "" {
		cfg.Origin = value
	}
	if value := os.Getenv(EnvHost); value != "" {
		cfg.Host = value
	}
	if value := os.Getenv(EnvPort); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil {
			return Config{}, err
		}
		cfg.Port = port
	}
	if httpClient != nil {
		cfg.HTTPClient = httpClient
	}

	return cfg.normalize(), nil
}
