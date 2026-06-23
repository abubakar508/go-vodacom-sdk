# Async Callbacks

The SDK provides callback request structs and acknowledgement helpers.

```go
mux.Handle("/callbacks/c2b", mpesa.C2BCallbackHandler(func(r *http.Request, cb mpesa.C2BAsyncCallbackRequest) error {
    // persist callback and update your application state
    return nil
}))
```

Generic handler:

```go
handler := mpesa.CallbackHandler(processFunc, mpesa.AcceptC2BAsyncCallback)
```

Manual helpers:

```go
var cb mpesa.C2BAsyncCallbackRequest
err := mpesa.DecodeCallbackJSON(r, &cb)
ack := mpesa.AcceptC2BAsyncCallback(cb)
_ = mpesa.WriteJSON(w, http.StatusOK, ack)
```

Most acknowledgements use response code `0`. Some APIs document specialized success codes; the SDK follows the endpoint documentation, for example:

- `AcceptUpdateTransactionStatusAsyncCallback` returns `INS-GAR-0`
- `AcceptQueryDirectDebitAsyncCallback` returns `INS-0`
