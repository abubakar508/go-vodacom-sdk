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

	req := client.NewQueryBeneficiaryNameRequest(
		"254707161122",                      // customer MSISDN
		"000000",                            // service provider shortcode
		"asv02e5958774f7ba228d83d0d689761", // third-party conversation ID
	)
	// req.InputKycQueryType defaults to "Name".

	res, raw, err := client.QueryBeneficiaryNameWithSession(ctx, session, req)
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
	fmt.Println("FirstName:", res.OutputCustomerFirstName)
	fmt.Println("LastName:", res.OutputCustomerLastName)
	fmt.Println("ThirdPartyConversationID:", res.OutputThirdPartyConversationID)
}
