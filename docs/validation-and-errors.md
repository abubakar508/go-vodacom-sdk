# Validation and Errors

## Request validation

The SDK validates required fields and common regex constraints before sending requests:

- Amounts
- MSISDNs
- Shortcodes
- Currency codes
- Transaction references
- Third-party conversation IDs
- Mandate IDs
- Dates in `YYYYMMDD`

M-Pesa still performs final business validation server-side.

## API errors

Non-2xx responses and known non-success response codes return `*mpesa.APIError`.

```go
res, raw, err := client.C2BSingleStageWithSession(ctx, session, req)
if err != nil {
    if mpesa.IsDuplicate(err) {
        // handle duplicate transaction
    }
    if raw != nil {
        log.Println(raw.BodyString())
    }
}
```

Useful helpers:

```go
mpesa.IsTimeout(err)
mpesa.IsDuplicate(err)
mpesa.IsInvalidMarket(err)
mpesa.IsInsufficientBalance(err)
mpesa.IsValidationFailed(err)
mpesa.IsAPINotEnabled(err)
```
