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

	req := client.NewQueryDirectDebitWithMSISDN(
		mpesa.QueryBalanceAmountTrue, // query balance sufficiency
		"100",                       // balance amount; required when QueryBalanceAmount is True
		"255744553111",              // customer MSISDN; alternatively use token helper
		"112244",                    // service provider shortcode
		"GPO3051656128",             // third-party conversation ID
		"Test123",                   // mandate/third-party reference
		"15045",                     // mandate ID, optional but recommended when available
	)

	res, raw, err := client.QueryDirectDebitWithSession(ctx, session, req)
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
	fmt.Println("SufficientBalance:", res.OutputSufficientBalance)
	fmt.Println("MsisdnToken:", res.OutputMsisdnToken)
	fmt.Println("MandateID:", res.OutputMandateID)
	fmt.Println("MandateStatus:", res.OutputMandateStatus)
	fmt.Println("AccountStatus:", res.OutputAccountStatus)
	fmt.Println("FirstPaymentDate:", res.OutputFirstPaymentDate)
	fmt.Println("Frequency:", res.OutputFrequency)
	fmt.Println("PaymentDayFrom:", res.OutputPaymentDayFrom)
	fmt.Println("PaymentDayTo:", res.OutputPaymentDayTo)
	fmt.Println("ExpiryDate:", res.OutputExpiryDate)
}
