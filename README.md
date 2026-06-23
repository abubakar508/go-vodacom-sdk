# go-vodacom-sdk

A production-ready Go SDK for the Vodacom/Vodafone M-Pesa OpenAPI.

```go
import "github.com/abubakar508/go-vodacom-sdk/mpesa"
```

> Go directive: `go 1.26.1`

## What is included

- Session generation and bearer encryption
- C2B Single Stage and C2B Multi Stage
- B2C Single Stage
- B2B Single Stage
- Reversal
- Query Transaction Status
- Update Transaction Status
- Direct Debit Create
- Direct Debit Payment
- Query Beneficiary Name
- Query Direct Debit
- Cancel Direct Debit
- Async callback structs and acknowledgement helpers
- Request validation helpers
- Typed API errors and error-code helpers
- Configurable public key and currency override
- Examples for every endpoint
- GitHub Pages documentation site in `/docs`

## Install

```bash
go get github.com/abubakar508/go-vodacom-sdk
```

## Quick start

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/abubakar508/go-vodacom-sdk/mpesa"
)

func main() {
    ctx := context.Background()

    client, err := mpesa.NewClient(mpesa.Config{
        APIKey:      os.Getenv("MPESA_API_KEY"),
        PublicKey:   os.Getenv("MPESA_PUBLIC_KEY"), // optional override
        Environment: mpesa.EnvironmentSandbox,
        Market:      mpesa.MarketGhana,
        Origin:      "*",
    })
    if err != nil {
        log.Fatal(err)
    }

    session, _, err := client.GenerateSessionAndWait(ctx, 0)
    if err != nil {
        log.Fatal(err)
    }

    req := client.NewC2BSingleStageRequest(
        "10",
        "000000000001",
        "000000",
        "T1234C",
        "asv02e5958774f7ba228d83d0d689761",
        "Shoes",
    )

    res, raw, err := client.C2BSingleStageWithSession(ctx, session, req)
    _ = res
    _ = raw
    if err != nil {
        log.Fatal(err)
    }
}
```

## Environment configuration

```bash
export MPESA_API_KEY="your-application-api-key"
export MPESA_ENVIRONMENT="sandbox" # sandbox or openapi
export MPESA_MARKET="vodafoneGHA"
export MPESA_ORIGIN="*"

# Recommended when the portal publishes/rotates keys:
export MPESA_PUBLIC_KEY="base64-platform-public-key-from-portal"

# Optional currency override, e.g. DRC CDF:
export MPESA_CURRENCY="CDF"
```

Then:

```go
cfg, err := mpesa.ConfigFromEnv()
client, err := mpesa.NewClient(cfg)
```

## Supported markets

| Constant | Context | Country | Default currency |
|---|---:|---:|---:|
| `mpesa.MarketGhana` | `vodafoneGHA` | `GHA` | `GHS` |
| `mpesa.MarketTanzania` | `vodacomTZN` | `TZN` | `TZS` |
| `mpesa.MarketLesotho` | `vodacomLES` | `LES` | `LSL` |
| `mpesa.MarketDRC` | `vodacomDRC` | `DRC` | `USD` |
| `mpesa.MarketMozambique` | `vodacomMOZ` | `MOZ` | `MZN` |

Use `mpesa.WithCurrency("CDF")` or `MPESA_CURRENCY=CDF` if your portal product supports another currency.

## Session flow

All transaction/query APIs require a generated SessionID:

```go
session, _, err := client.GenerateSessionAndWait(ctx, 0)
res, raw, err := client.B2CSingleStageWithSession(ctx, session, req)
```

The default wait is 30 seconds because official examples warn that fresh SessionIDs can take time to become active.

## Error handling

```go
res, raw, err := client.C2BSingleStageWithSession(ctx, session, req)
if err != nil {
    if mpesa.IsDuplicate(err) {
        // handle duplicate transaction
    }
    if raw != nil {
        log.Println(raw.BodyString())
    }
    log.Fatal(err)
}
_ = res
```

## Callback handlers

```go
mux.Handle("/callbacks/c2b", mpesa.C2BCallbackHandler(func(r *http.Request, cb mpesa.C2BAsyncCallbackRequest) error {
    // persist callback and update business state
    return nil
}))
```

## Examples

```bash
go run ./examples/session
go run ./examples/c2b
go run ./examples/c2b_multi_stage
go run ./examples/b2c
go run ./examples/b2b
go run ./examples/reversal
go run ./examples/query_transaction_status
go run ./examples/update_transaction_status
go run ./examples/direct_debit_create
go run ./examples/direct_debit_payment
go run ./examples/query_beneficiary_name
go run ./examples/query_direct_debit
go run ./examples/cancel_direct_debit
go run ./examples/callback_listener
```

## Documentation

- [Authentication](docs/authentication.md)
- [Markets](docs/markets.md)
- [Endpoints](docs/endpoints.md)
- [Callbacks](docs/callbacks.md)
- [Validation and errors](docs/validation-and-errors.md)
- [GitHub Pages deployment](docs/github-pages-deployment.md)
- Static landing page: [`docs/index.html`](docs/index.html)

## Development

```bash
gofmt -w .
go vet ./...
go test -race ./...
```

A GitHub Actions workflow is included in `.github/workflows/ci.yml`.

## License

MIT. See [LICENSE](LICENSE).
