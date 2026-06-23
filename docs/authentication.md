# Authentication and Session Flow

M-Pesa OpenAPI uses two token stages.

## 1. Generate a SessionID

Use your application API key from the portal. The SDK encrypts it with the configured platform public key and calls:

```text
GET /{sandbox|openapi}/ipg/v2/{market}/getSession/
```

```go
session, raw, err := client.GenerateSession(ctx)
```

## 2. Use the SessionID for API calls

Payment/query APIs use the generated SessionID as the encrypted bearer value:

```go
res, raw, err := client.C2BSingleStageWithSession(ctx, session, req)
```

## Activation delay

Official examples warn that a new SessionID can take up to 30 seconds to become active. Use:

```go
session, raw, err := client.GenerateSessionAndWait(ctx, 0)
```

Passing `0` uses the SDK default of `30 * time.Second`.

## Public key override

If Vodacom rotates the public key, set:

```bash
export MPESA_PUBLIC_KEY="base64-public-key-from-portal"
```

or:

```go
client, err := mpesa.NewClient(mpesa.Config{PublicKey: os.Getenv("MPESA_PUBLIC_KEY")})
```
