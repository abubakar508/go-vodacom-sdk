package mpesa

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	cryptoutil "github.com/abubakar508/go-vodacom-sdk/internal/crypto"
)

// Config contains all configuration needed by the SDK.
//
// Standard usage is either:
//   cfg := mpesa.Config{APIKey: "..."}
//   client, err := mpesa.NewClient(cfg)
//
// or the functional option style:
//   client, err := mpesa.NewClientWithOptions(mpesa.WithAPIKey("..."))
type Config struct {
	// APIKey is your application API key from the M-Pesa developer portal.
	// It is used only by GenerateSessionKey. Transaction APIs use SessionID.
	APIKey string

	// PublicKey is the platform public key used to encrypt APIKey/SessionID.
	// If empty, the SDK uses the default key for Environment.
	PublicKey string

	// Environment controls /sandbox/ vs /openapi/. Defaults to sandbox.
	Environment Environment

	// Market controls URL market context and default country/currency.
	// Defaults to Vodacom DRC.
	Market Market

	// Currency optionally overrides the market default currency for requests.
	// This is useful for markets such as DRC where the portal may allow multiple
	// currencies, e.g. USD and CDF.
	Currency string

	// Origin should match the origin configured on the M-Pesa application.
	// Defaults to "*" to match the official examples.
	Origin string

	// Host defaults to openapi.m-pesa.com.
	Host string

	// Port defaults to 443.
	Port int

	// HTTPClient defaults to an http.Client with a 60 second timeout.
	HTTPClient *http.Client
}

// DefaultConfig returns a Config with safe SDK defaults for Vodacom DRC sandbox.
func DefaultConfig() Config {
	return Config{
		Environment: EnvironmentSandbox,
		Market:      MarketDRC,
		PublicKey:   DefaultSandboxPublicKey,
		Origin:      "*",
		Host:        defaultHost,
		Port:        443,
		HTTPClient:  &http.Client{Timeout: 60 * time.Second},
	}
}

func (c Config) normalize() Config {
	defaults := DefaultConfig()

	if c.Environment == "" {
		c.Environment = defaults.Environment
	}
	if c.Market.Context == "" {
		c.Market = defaults.Market
	}
	if strings.TrimSpace(c.PublicKey) == "" {
		c.PublicKey = defaultPublicKeyForEnvironment(c.Environment)
	}
	c.Currency = strings.ToUpper(strings.TrimSpace(c.Currency))
	if strings.TrimSpace(c.Origin) == "" {
		c.Origin = defaults.Origin
	}
	if strings.TrimSpace(c.Host) == "" {
		c.Host = defaults.Host
	}
	if c.Port == 0 {
		c.Port = defaults.Port
	}
	if c.HTTPClient == nil {
		c.HTTPClient = defaults.HTTPClient
	}

	return c
}

func (c Config) validate() error {
	if !c.Environment.Valid() {
		return fmt.Errorf("unsupported environment %q", c.Environment)
	}
	if !c.Market.valid() {
		return errors.New("market context, country, and currency are required")
	}
	if c.Currency != "" {
		if err := validateCurrency("currency override", c.Currency); err != nil {
			return err
		}
	}
	if c.Port < 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port %d", c.Port)
	}
	if _, err := cryptoutil.ParseRSAPublicKey(c.PublicKey); err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}
	return nil
}
