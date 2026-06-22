package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	session, _, err := client.GenerateSessionKey(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Official examples warn that a new SessionID can take up to 30 seconds
	// to become usable by transaction APIs.
	time.Sleep(30 * time.Second)

	req := client.NewB2CSingleStageRequest(
		"10",
		"000000000001",                      // customer MSISDN to credit
		"000000",                            // service provider shortcode
		"T1234C",                            // transaction reference, max 20 chars
		"asv02e5958774f7ba228d83d0d689761", // third-party conversation ID
		"Salary payment",
	)

	res, raw, err := client.B2CSingleStage(ctx, session.OutputSessionID, req)
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
