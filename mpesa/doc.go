// Package mpesa provides a Go client for the Vodacom/Vodafone M-Pesa OpenAPI.
//
// Basic usage:
//
//	client, err := mpesa.NewClient(mpesa.Config{
//		APIKey:      os.Getenv("MPESA_API_KEY"),
//		Environment: mpesa.EnvironmentSandbox,
//		Market:      mpesa.MarketGhana,
//		PublicKey:   os.Getenv("MPESA_PUBLIC_KEY"), // optional override
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	session, _, err := client.GenerateSessionAndWait(ctx, 0)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	req := client.NewC2BSingleStageRequest("10", "000000000001", "000000", "T1234C", "conversation-id", "Shoes")
//	res, raw, err := client.C2BSingleStageWithSession(ctx, session, req)
//	_ = res
//	_ = raw
//
// The SDK supports sandbox/openapi environments, all documented Vodacom/Vodafone
// markets, configurable platform public keys, and both synchronous and
// asynchronous OpenAPI flows.
package mpesa
