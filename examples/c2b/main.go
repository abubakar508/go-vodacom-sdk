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

	// Generates a session using MPESA_API_KEY and waits the default 30 seconds
	// recommended by the official examples before using it on transaction APIs.
	session, _, err := client.GenerateSessionAndWait(ctx, 0)
	if err != nil {
		log.Fatal(err)
	}

	req := client.NewC2BSingleStageRequest(
		"10",
		"000000000001",                     // customer MSISDN
		"000000",                           // service provider shortcode
		"T1234C",                           // transaction reference, max 20 chars
		"asv02e5958774f7ba228d83d0d689761", // third-party conversation ID
		"Shoes",
	)

	res, raw, err := client.C2BSingleStageWithSession(ctx, session, req)
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
	fmt.Println("TransactionID:", res.OutputTransactionID)
	fmt.Println("ThirdPartyConversationID:", res.OutputThirdPartyConversationID)
}
