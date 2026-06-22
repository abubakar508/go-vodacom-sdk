package mpesa

import "net/http"

// Option configures the SDK client using the functional options pattern.
type Option func(*Config)

// NewClientWithOptions creates a Client from DefaultConfig plus the supplied options.
func NewClientWithOptions(options ...Option) (*Client, error) {
	cfg := DefaultConfig()
	for _, option := range options {
		if option != nil {
			option(&cfg)
		}
	}
	return NewClient(cfg)
}

func WithAPIKey(apiKey string) Option {
	return func(c *Config) { c.APIKey = apiKey }
}

func WithPublicKey(publicKey string) Option {
	return func(c *Config) { c.PublicKey = publicKey }
}

func WithEnvironment(environment Environment) Option {
	return func(c *Config) { c.Environment = environment }
}

func WithMarket(market Market) Option {
	return func(c *Config) { c.Market = market }
}

func WithOrigin(origin string) Option {
	return func(c *Config) { c.Origin = origin }
}

func WithHost(host string) Option {
	return func(c *Config) { c.Host = host }
}

func WithPort(port int) Option {
	return func(c *Config) { c.Port = port }
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Config) { c.HTTPClient = httpClient }
}
