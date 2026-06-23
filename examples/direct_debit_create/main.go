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

	req := client.NewDirectDebitCreateRequest(
		"000000000001",                         // customer MSISDN
		"000000",                               // service provider shortcode
		"3333",                                 // mandate reference
		"AAA6d1f9391a0052de0b5334a912jbsj1j2kk", // third-party conversation ID
		mpesa.DirectDebitAgreedTCYes,            // customer agreed to T&C: "1"
	)

	// Optional schedule fields. Leave empty if frequency is not used.
	req.InputFrequency = mpesa.DirectDebitFrequencyHalfYearly
	req.InputFirstPaymentDate = "20160324"
	req.InputStartRangeOfDays = "01"
	req.InputEndRangeOfDays = "22"
	req.InputExpiryDate = "20161126"

	res, raw, err := client.DirectDebitCreateWithSession(ctx, session, req)
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
	fmt.Println("MsisdnToken:", res.OutputMsisdnToken)
	fmt.Println("MandateID:", res.OutputMandateID)
}
