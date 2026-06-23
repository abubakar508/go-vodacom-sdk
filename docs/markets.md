# Markets

Built-in markets:

| Constant | Context | Country | Currency |
|---|---:|---:|---:|
| `mpesa.MarketGhana` | `vodafoneGHA` | `GHA` | `GHS` |
| `mpesa.MarketTanzania` | `vodacomTZN` | `TZN` | `TZS` |
| `mpesa.MarketLesotho` | `vodacomLES` | `LES` | `LSL` |
| `mpesa.MarketDRC` | `vodacomDRC` | `DRC` | `USD` |
| `mpesa.MarketMozambique` | `vodacomMOZ` | `MOZ` | `MZN` |

## Currency override

Some portal pages mention more than one currency for a market, for example DRC may support `USD` and `CDF` depending on the product/environment.

Use a currency override:

```go
client, err := mpesa.NewClientWithOptions(
    mpesa.WithMarket(mpesa.MarketDRC),
    mpesa.WithCurrency("CDF"),
)
```

or environment variable:

```bash
export MPESA_CURRENCY=CDF
```

## Custom market

```go
market := mpesa.CustomMarket("New Market", "vodacomXYZ", "XYZ", "XYZ")
```
