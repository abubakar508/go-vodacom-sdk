package mpesa

import "strings"

// Market describes an M-Pesa OpenAPI market.
//
// Context is the market URL segment, while Country and Currency are request
// body values used by transaction APIs such as C2B, B2C, and B2B.
type Market struct {
	Description string
	Context     string
	Country     string
	Currency    string
}

var (
	MarketGhana      = Market{Description: "Vodafone Ghana", Context: "vodafoneGHA", Country: "GHA", Currency: "GHS"}
	MarketTanzania   = Market{Description: "Vodacom Tanzania", Context: "vodacomTZN", Country: "TZN", Currency: "TZS"}
	MarketLesotho    = Market{Description: "Vodacom Lesotho", Context: "vodacomLES", Country: "LES", Currency: "LSL"}
	MarketDRC        = Market{Description: "Vodacom DR Congo", Context: "vodacomDRC", Country: "DRC", Currency: "USD"}
	MarketMozambique = Market{Description: "Vodacom Mozambique", Context: "vodacomMOZ", Country: "MOZ", Currency: "MZN"}
)

// Markets contains the markets currently documented by Vodacom/Vodafone M-Pesa OpenAPI.
//
// Keys are the official URL context values, e.g. "vodacomDRC".
var Markets = map[string]Market{
	MarketGhana.Context:      MarketGhana,
	MarketTanzania.Context:   MarketTanzania,
	MarketLesotho.Context:    MarketLesotho,
	MarketDRC.Context:        MarketDRC,
	MarketMozambique.Context: MarketMozambique,
}

// SupportedMarkets returns a copy of the documented markets supported by this SDK.
func SupportedMarkets() []Market {
	return []Market{MarketGhana, MarketTanzania, MarketLesotho, MarketDRC, MarketMozambique}
}

// MarketFromContext returns a documented market by URL context value.
// Matching is case-insensitive to make environment variable configuration easier.
func MarketFromContext(context string) (Market, bool) {
	context = strings.TrimSpace(context)
	if context == "" {
		return Market{}, false
	}
	if market, ok := Markets[context]; ok {
		return market, true
	}
	for _, market := range Markets {
		if strings.EqualFold(market.Context, context) {
			return market, true
		}
	}
	return Market{}, false
}

// CustomMarket creates a market not yet built into the SDK. This is useful if
// Vodacom/Vodafone adds a market before the SDK is updated.
func CustomMarket(description, context, country, currency string) Market {
	return Market{
		Description: strings.TrimSpace(description),
		Context:     strings.TrimSpace(context),
		Country:     strings.ToUpper(strings.TrimSpace(country)),
		Currency:    strings.ToUpper(strings.TrimSpace(currency)),
	}
}

func (m Market) valid() bool {
	return strings.TrimSpace(m.Context) != "" && strings.TrimSpace(m.Country) != "" && strings.TrimSpace(m.Currency) != ""
}
