package main

import (
	"log"
	"net/http"

	"github.com/abubakar508/go-vodacom-sdk/mpesa"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/callbacks/c2b", mpesa.C2BCallbackHandler(func(r *http.Request, cb mpesa.C2BAsyncCallbackRequest) error {
		log.Printf("C2B callback: conversation=%s transaction=%s result=%s", cb.InputOriginalConversationID, cb.InputTransactionID, cb.InputResultCode)
		// Persist callback, update your order/payment state, then return nil.
		return nil
	}))

	mux.Handle("/callbacks/direct-debit-create", mpesa.DirectDebitCreateCallbackHandler(func(r *http.Request, cb mpesa.DirectDebitCreateAsyncCallbackRequest) error {
		log.Printf("Direct debit mandate callback: mandate=%s token=%s result=%s", cb.InputMandateID, cb.InputMsisdnToken, cb.InputResultCode)
		return nil
	}))

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
