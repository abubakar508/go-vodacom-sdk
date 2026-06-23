package mpesa

import (
	"encoding/json"
	"net/http"
)

// DecodeCallbackJSON decodes an OpenAPI callback JSON request into dst.
func DecodeCallbackJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

// WriteJSON writes a JSON response with application/json content type.
func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

// CallbackHandler adapts a typed callback processor into an http.Handler.
// The ack function receives the decoded request and should return the response
// payload expected by OpenAPI, usually one of the Accept*Callback helpers.
func CallbackHandler[T any, A any](process func(*http.Request, T) error, ack func(T) A) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload T
		if err := DecodeCallbackJSON(r, &payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if process != nil {
			if err := process(r, payload); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		_ = WriteJSON(w, http.StatusOK, ack(payload))
	})
}

func C2BCallbackHandler(process func(*http.Request, C2BAsyncCallbackRequest) error) http.Handler {
	return CallbackHandler(process, AcceptC2BAsyncCallback)
}

func B2CCallbackHandler(process func(*http.Request, B2CAsyncCallbackRequest) error) http.Handler {
	return CallbackHandler(process, AcceptB2CAsyncCallback)
}

func B2BCallbackHandler(process func(*http.Request, B2BAsyncCallbackRequest) error) http.Handler {
	return CallbackHandler(process, AcceptB2BAsyncCallback)
}

func ReversalCallbackHandler(process func(*http.Request, ReversalAsyncCallbackRequest) error) http.Handler {
	return CallbackHandler(process, AcceptReversalAsyncCallback)
}

func DirectDebitCreateCallbackHandler(process func(*http.Request, DirectDebitCreateAsyncCallbackRequest) error) http.Handler {
	return CallbackHandler(process, AcceptDirectDebitCreateAsyncCallback)
}

func DirectDebitPaymentCallbackHandler(process func(*http.Request, DirectDebitPaymentAsyncCallbackRequest) error) http.Handler {
	return CallbackHandler(process, AcceptDirectDebitPaymentAsyncCallback)
}
