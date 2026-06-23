# Local C2B Tester

The SDK includes a local web UI for testing Generate SessionKey + C2B Single Stage.

## Why it is local, not GitHub Pages only

Do **not** put M-Pesa API keys into a public static page.

A browser-only GitHub Pages form is not suitable because:

1. It exposes your application API key to the browser.
2. Browsers cannot reliably set protected headers such as `Origin` manually.
3. M-Pesa OpenAPI may not allow cross-origin browser calls.

The tester runs a local Go backend. That backend uses the SDK server-side and sets the M-Pesa `Origin` header, defaulting to `*`.

## Run it

```bash
go run ./examples/c2b_tester
```

Open:

```text
http://localhost:8088
```

## Fields

- Application API key
- Optional platform public key override
- Environment: `sandbox` or `openapi`
- Market
- Origin header, default `*`
- C2B request fields

The tester:

1. Generates a SessionID
2. Waits the configured seconds, default `30`
3. Performs C2B Single Stage
4. Shows raw response and typed response

## Security

- The tester does not store credentials.
- It redacts the SessionID in output.
- Use sandbox credentials unless you intentionally want a live transaction.
- Do not deploy this tester as a public web service without authentication and secret handling.
