package main

import (
	"context"
	"fmt"
	"log"

	"github.com/abubakar508/go-vodacom-sdk/mpesa"
)

func main() {
	ctx := context.Background()

	cfg, err := mpesa.ConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	client, err := mpesa.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	session, _, err := client.GenerateSessionAndWait(ctx, 0)
	if err != nil {
		log.Fatal(err)
	}

	req := client.NewDirectDebitPaymentWithMSISDN(
		"10",                                     // amount
		"000000000001",                           // customer MSISDN; alternatively use token helper
		"000000",                                 // service provider shortcode
		"5db410b459bd433ca8e5",                   // mandate/third-party reference
		"AAA6d1f939c1005v2de053v4912jbasdj1j2kk", // third-party conversation ID
		"15045",                                  // mandate ID, optional but recommended when available
	)

	res, raw, err := client.DirectDebitPaymentWithSession(ctx, session, req)
	if err != nil {
		if raw != nil {
			log.Printf("raw response: %s", raw.BodyString())
		}
		log.Fatal(err)
	}

	fmt.Println("HTTP:", raw.StatusCode)
	fmt.Println("Code:", res.OutputResponseCode)
	fmt.Println("Desc:", res.OutputResponseDesc)
	fmt.Println("ConversationID:", res.OutputConversationID)
	fmt.Println("ThirdPartyConversationID:", res.OutputThirdPartyConversationID)
	fmt.Println("TransactionReference:", res.OutputTransactionReference)
	fmt.Println("TransactionID:", res.OutputTransactionID)
	fmt.Println("MsisdnToken:", res.OutputMsisdnToken)
}
