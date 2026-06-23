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

// WithMarketContext configures a built-in market using its URL context value,
// e.g. vodacomDRC, vodafoneGHA, vodacomTZN, vodacomLES, or vodacomMOZ.
// Unknown context values are set as custom context only; for custom markets,
// prefer WithMarket(CustomMarket(...)) so country and currency are provided.
func WithMarketContext(context string) Option {
	return func(c *Config) {
		if market, ok := MarketFromContext(context); ok {
			c.Market = market
			return
		}
		c.Market = Market{Context: context}
	}
}

// WithCurrency overrides the configured market's default request currency.
// For example, use WithCurrency("CDF") for DRC CDF flows when enabled in the portal.
func WithCurrency(currency string) Option {
	return func(c *Config) { c.Currency = currency }
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
