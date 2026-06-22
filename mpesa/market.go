package mpesa

// Market describes an M-Pesa OpenAPI market.
//
// Context is the market URL segment, while Country and Currency are request
// body values used by APIs such as C2B Single Stage.
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
var Markets = map[string]Market{
	MarketGhana.Context:      MarketGhana,
	MarketTanzania.Context:   MarketTanzania,
	MarketLesotho.Context:    MarketLesotho,
	MarketDRC.Context:        MarketDRC,
	MarketMozambique.Context: MarketMozambique,
}

func (m Market) valid() bool {
	return m.Context != "" && m.Country != "" && m.Currency != ""
}
