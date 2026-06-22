# go-vodacom-sdk

A Go SDK for Vodacom/Vodafone M-Pesa OpenAPI, focused first on **Vodacom DRC** and organized for adding the remaining OpenAPI endpoints cleanly.

Module:

```text
github.com/abubakar508/go-vodacom-sdk
```

Go directive:

```text
go 1.26.1
```

## Current support

- Generate SessionKey
- C2B Single Stage
- B2C Single Stage
- B2B Single Stage
- Sandbox and live/OpenAPI environments
- Vodacom DRC defaults: `vodacomDRC`, `DRC`, `USD`
- Other documented markets available as constants
- Standard config styles:
  - `mpesa.Config{...}`
  - `mpesa.DefaultConfig()`
  - `mpesa.NewClientWithOptions(...)`
  - `mpesa.ConfigFromEnv()`
- RSA PKCS#1 v1.5 bearer encryption compatible with the official PHP `openssl_public_encrypt` flow

## Folder structure

```text
.
├── go.mod
├── .env.example
├── README.md
├── internal/
│   └── crypto/
│       └── rsa.go              # internal RSA/public-key encryption helpers
├── mpesa/                      # public SDK package
│   ├── callbacks.go            # async callback helpers
│   ├── c2b.go                  # C2B Single Stage
│   ├── b2c.go                  # B2C Single Stage
│   ├── b2b.go                  # B2B Single Stage
│   ├── client.go               # HTTP client and request execution
│   ├── client_test.go
│   ├── config.go               # Config, defaults, validation
│   ├── config_env.go           # environment variable config
│   ├── environment.go          # sandbox/openapi environments
│   ├── errors.go               # typed API errors
│   ├── keys.go                 # sandbox/live public keys
│   ├── market.go               # supported markets
│   ├── options.go              # functional options
│   ├── response.go             # raw response wrapper
│   └── session.go              # Generate SessionKey
└── examples/
    ├── c2b/
    │   └── main.go
    ├── b2c/
    │   └── main.go
    ├── b2b/
    │   └── main.go
    └── session/
        └── main.go
```

## Install

```bash
go get github.com/abubakar508/go-vodacom-sdk
```

Import the public package:

```go
import "github.com/abubakar508/go-vodacom-sdk/mpesa"
```

## Configuration options

### 1. Direct config

```go
client, err := mpesa.NewClient(mpesa.Config{
    APIKey:      "your-application-api-key",
    Environment: mpesa.EnvironmentSandbox,
    Market:      mpesa.MarketDRC,
    Origin:      "*",
})
```

### 2. Functional options

```go
client, err := mpesa.NewClientWithOptions(
    mpesa.WithAPIKey("your-application-api-key"),
    mpesa.WithEnvironment(mpesa.EnvironmentSandbox),
    mpesa.WithMarket(mpesa.MarketDRC),
    mpesa.WithOrigin("*"),
)
```

### 3. Environment variables

Copy `.env.example` and set:

```bash
export MPESA_API_KEY="your-application-api-key"
export MPESA_ENVIRONMENT="sandbox" # or openapi
export MPESA_MARKET="vodacomDRC"
export MPESA_ORIGIN="*"
```

Then:

```go
cfg, err := mpesa.ConfigFromEnv()
if err != nil {
    log.Fatal(err)
}
client, err := mpesa.NewClient(cfg)
```

Supported environment variables:

| Variable | Description |
|---|---|
| `MPESA_API_KEY` | Application API key from the developer portal |
| `MPESA_PUBLIC_KEY` | Optional override for platform public key |
| `MPESA_ENVIRONMENT` | `sandbox` or `openapi` |
| `MPESA_MARKET` | `vodacomDRC`, `vodafoneGHA`, `vodacomTZN`, `vodacomLES`, `vodacomMOZ` |
| `MPESA_ORIGIN` | Origin header configured in your portal app |
| `MPESA_HOST` | Defaults to `openapi.m-pesa.com` |
| `MPESA_PORT` | Defaults to `443` |

## Generate SessionKey

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/abubakar508/go-vodacom-sdk/mpesa"
)

func main() {
    client, err := mpesa.NewClient(mpesa.Config{
        APIKey:      os.Getenv("MPESA_API_KEY"),
        Environment: mpesa.EnvironmentSandbox,
        Market:      mpesa.MarketDRC,
        Origin:      "*",
    })
    if err != nil {
        log.Fatal(err)
    }

    session, raw, err := client.GenerateSessionKey(context.Background())
    if err != nil {
        if raw != nil {
            log.Printf("raw response: %s", raw.BodyString())
        }
        log.Fatal(err)
    }

    fmt.Println(session.OutputResponseCode)
    fmt.Println(session.OutputResponseDesc)
    fmt.Println(session.OutputSessionID)
}
```

Endpoint generated for DRC sandbox:

```text
GET https://openapi.m-pesa.com/sandbox/ipg/v2/vodacomDRC/getSession/
```

Endpoint generated for DRC live/OpenAPI:

```text
GET https://openapi.m-pesa.com/openapi/ipg/v2/vodacomDRC/getSession/
```

Successful response shape:

```json
{
  "output_ResponseCode": "INS-0",
  "output_ResponseDesc": "Request processed successfully",
  "output_SessionID": "ed0ff3b37d5145f38885e3212aef9774"
}
```

## C2B Single Stage

```go
ctx := context.Background()

session, _, err := client.GenerateSessionKey(ctx)
if err != nil {
    log.Fatal(err)
}

// Official examples warn that a new SessionID can take up to 30 seconds
// to become active before transaction APIs accept it.
time.Sleep(30 * time.Second)

req := client.NewC2BSingleStageRequest(
    "10",
    "000000000001",
    "000000",
    "T1234C",
    "asv02e5958774f7ba228d83d0d689761",
    "Shoes",
)

res, raw, err := client.C2BSingleStage(ctx, session.OutputSessionID, req)
if err != nil {
    if raw != nil {
        log.Println(raw.BodyString())
    }
    log.Fatal(err)
}

fmt.Println(res.OutputResponseCode)
fmt.Println(res.OutputConversationID)
fmt.Println(res.OutputTransactionID)
```

Endpoint generated for DRC sandbox:

```text
POST https://openapi.m-pesa.com/sandbox/ipg/v2/vodacomDRC/c2bPayment/singleStage/
```

For Vodacom DRC, `NewC2BSingleStageRequest` automatically fills:

```json
{
  "input_Country": "DRC",
  "input_Currency": "USD"
}
```

Successful sync response shape:

```json
{
  "output_ConversationID": "d3502e5958774f7ba228d83d0d689761",
  "output_ResponseCode": "INS-0",
  "output_ResponseDesc": "Request processed successfully",
  "output_TransactionID": "49XCD123F6",
  "output_ThirdPartyConversationID": "asv02e5958774f7ba228d83d0d689761"
}
```

## B2C Single Stage

B2C Single Stage is used for business-to-customer disbursements such as salary payments, business payouts, and charity payouts. It uses the SessionID from `GenerateSessionKey` as the encrypted bearer value.

```go
ctx := context.Background()

session, _, err := client.GenerateSessionKey(ctx)
if err != nil {
    log.Fatal(err)
}

// Official examples warn that a new SessionID can take up to 30 seconds
// to become active before transaction APIs accept it.
time.Sleep(30 * time.Second)

req := client.NewB2CSingleStageRequest(
    "10",
    "000000000001",
    "000000",
    "T1234C",
    "asv02e5958774f7ba228d83d0d689761",
    "Salary payment",
)

res, raw, err := client.B2CSingleStage(ctx, session.OutputSessionID, req)
if err != nil {
    if raw != nil {
        log.Println(raw.BodyString())
    }
    log.Fatal(err)
}

fmt.Println(res.OutputResponseCode)
fmt.Println(res.OutputConversationID)
fmt.Println(res.OutputTransactionID)
```

Endpoint generated for DRC sandbox:

```text
POST https://openapi.m-pesa.com/sandbox/ipg/v2/vodacomDRC/b2cPayment/
```

Endpoint generated for DRC live/OpenAPI:

```text
POST https://openapi.m-pesa.com/openapi/ipg/v2/vodacomDRC/b2cPayment/
```

For Vodacom DRC, `NewB2CSingleStageRequest` automatically fills:

```json
{
  "input_Country": "DRC",
  "input_Currency": "USD"
}
```

Full B2C request shape:

```json
{
  "input_Amount": "10",
  "input_Country": "DRC",
  "input_Currency": "USD",
  "input_CustomerMSISDN": "000000000001",
  "input_ServiceProviderCode": "000000",
  "input_ThirdPartyConversationID": "asv02e5958774f7ba228d83d0d689761",
  "input_TransactionReference": "T1234C",
  "input_PaymentItemsDesc": "Salary payment"
}
```

Successful sync response shape:

```json
{
  "output_ConversationID": "d3502e5958774f7ba228d83d0d689761",
  "output_ResponseCode": "INS-0",
  "output_ResponseDesc": "Request processed successfully",
  "output_TransactionID": "49XCD123F6",
  "output_ThirdPartyConversationID": "asv02e5958774f7ba228d83d0d689761"
}
```

## B2B Single Stage

B2B Single Stage is used for business-to-business transfers such as stock purchases, bill payments, and ad-hoc payments. It uses the SessionID from `GenerateSessionKey` as the encrypted bearer value.

```go
ctx := context.Background()

session, _, err := client.GenerateSessionKey(ctx)
if err != nil {
    log.Fatal(err)
}

// Official examples warn that a new SessionID can take up to 30 seconds
// to become active before transaction APIs accept it.
time.Sleep(30 * time.Second)

req := client.NewB2BSingleStageRequest(
    "10",
    "000000",
    "000001",
    "T1234C",
    "asv02e5958774f7ba228d83d0d689761",
    "Shoes",
)

res, raw, err := client.B2BSingleStage(ctx, session.OutputSessionID, req)
if err != nil {
    if raw != nil {
        log.Println(raw.BodyString())
    }
    log.Fatal(err)
}

fmt.Println(res.OutputResponseCode)
fmt.Println(res.OutputConversationID)
fmt.Println(res.OutputTransactionID)
```

Endpoint generated for DRC sandbox:

```text
POST https://openapi.m-pesa.com/sandbox/ipg/v2/vodacomDRC/b2bPayment/
```

Endpoint generated for DRC live/OpenAPI:

```text
POST https://openapi.m-pesa.com/openapi/ipg/v2/vodacomDRC/b2bPayment/
```

For Vodacom DRC, `NewB2BSingleStageRequest` automatically fills:

```json
{
  "input_Country": "DRC",
  "input_Currency": "USD"
}
```

Full B2B request shape:

```json
{
  "input_Amount": "10",
  "input_Country": "DRC",
  "input_Currency": "USD",
  "input_PrimaryPartyCode": "000000",
  "input_ReceiverPartyCode": "000001",
  "input_ThirdPartyConversationID": "asv02e5958774f7ba228d83d0d689761",
  "input_TransactionReference": "T1234C",
  "input_PurchasedItemsDesc": "Shoes"
}
```

Successful sync response shape:

```json
{
  "output_ConversationID": "d3502e5958774f7ba228d83d0d689761",
  "output_ResponseCode": "INS-0",
  "output_ResponseDesc": "Request processed successfully",
  "output_TransactionID": "49XCD123F6",
  "output_ThirdPartyConversationID": "asv02e5958774f7ba228d83d0d689761"
}
```

Initial async accepted response shape:

```json
{
  "output_ResponseCode": "INS-0",
  "output_ResponseDesc": "Successfully Accepted Request",
  "output_ConversationID": "d3502e5958774f7ba228d83d0d689761",
  "output_ThirdPartyConversationID": "asv02e5958774f7ba228d83d0d689761"
}
```

## Markets

```go
mpesa.MarketDRC        // vodacomDRC / DRC / USD
mpesa.MarketGhana      // vodafoneGHA / GHA / GHS
mpesa.MarketTanzania   // vodacomTZN / TZN / TZS
mpesa.MarketLesotho    // vodacomLES / LES / LSL
mpesa.MarketMozambique // vodacomMOZ / MOZ / MZN
```

## Async callback acknowledgement

C2B and B2C async callbacks use the same payload shape, so the SDK provides generic and API-specific helpers:

```go
c2bAck := mpesa.AcceptC2BAsyncCallback(c2bCallback)
b2cAck := mpesa.AcceptB2CAsyncCallback(b2cCallback)
b2bAck := mpesa.AcceptB2BAsyncCallback(b2bCallback)
// or:
ack := mpesa.AcceptAsyncTransactionCallback(callback)
```

This returns the expected success acknowledgement:

```json
{
  "output_ResponseCode": "0",
  "output_ResponseDesc": "Successfully Accepted Result"
}
```

## Run examples

```bash
export MPESA_API_KEY="your-api-key"
export MPESA_ENVIRONMENT="sandbox"
export MPESA_MARKET="vodacomDRC"
export MPESA_ORIGIN="*"

go run ./examples/session
# go run ./examples/c2b
# go run ./examples/b2c
# go run ./examples/b2b
```

## Test

```bash
go test ./...
```

Note: I could not run tests inside the current workspace because the sandbox does not have the `go` command installed.

## Next endpoint folders/files

The `mpesa` package is ready for additional endpoint files:

```text
mpesa/reversal.go
mpesa/query_transaction_status.go
mpesa/c2b_multi_stage.go
mpesa/update_transaction_status.go
mpesa/direct_debit_create.go
mpesa/direct_debit_payment.go
mpesa/query_beneficiary_name.go
mpesa/query_direct_debit.go
mpesa/cancel_direct_debit.go
```
