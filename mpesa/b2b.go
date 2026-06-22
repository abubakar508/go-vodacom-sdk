package mpesa

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const b2bSingleStagePath = "b2bPayment/"

// B2BSingleStageRequest is the body for POST /b2bPayment/.
//
// It is used for business-to-business transfers such as stock purchases,
// bill payments, and ad-hoc business payments.
type B2BSingleStageRequest struct {
	InputAmount                   string `json:"input_Amount"`
	InputReceiverPartyCode        string `json:"input_ReceiverPartyCode"`
	InputCountry                  string `json:"input_Country"`
	InputCurrency                 string `json:"input_Currency"`
	InputPrimaryPartyCode         string `json:"input_PrimaryPartyCode"`
	InputTransactionReference     string `json:"input_TransactionReference"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
	InputPurchasedItemsDesc       string `json:"input_PurchasedItemsDesc"`
}

// NewB2BSingleStageRequest creates a B2B request using the client's configured
// market country/currency. You may still edit any field before sending.
func (c *Client) NewB2BSingleStageRequest(amount, primaryPartyCode, receiverPartyCode, transactionReference, thirdPartyConversationID, purchasedItemsDesc string) B2BSingleStageRequest {
	return B2BSingleStageRequest{
		InputAmount:                   amount,
		InputReceiverPartyCode:        receiverPartyCode,
		InputCountry:                  c.cfg.Market.Country,
		InputCurrency:                 c.cfg.Market.Currency,
		InputPrimaryPartyCode:         primaryPartyCode,
		InputTransactionReference:     transactionReference,
		InputThirdPartyConversationID: thirdPartyConversationID,
		InputPurchasedItemsDesc:       purchasedItemsDesc,
	}
}

// Validate checks all mandatory B2B Single Stage fields. M-Pesa performs final
// regex and business validations server-side.
func (r B2BSingleStageRequest) Validate() error {
	required := map[string]string{
		"input_Amount":                   r.InputAmount,
		"input_ReceiverPartyCode":        r.InputReceiverPartyCode,
		"input_Country":                  r.InputCountry,
		"input_Currency":                 r.InputCurrency,
		"input_PrimaryPartyCode":         r.InputPrimaryPartyCode,
		"input_TransactionReference":     r.InputTransactionReference,
		"input_ThirdPartyConversationID": r.InputThirdPartyConversationID,
		"input_PurchasedItemsDesc":       r.InputPurchasedItemsDesc,
	}
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			return errors.New(name + " is required")
		}
	}
	return nil
}

// B2BSingleStageResponse represents both sync and initial async responses.
// In async mode, OutputTransactionID may be empty until the callback arrives.
type B2BSingleStageResponse struct {
	OutputConversationID           string `json:"output_ConversationID"`
	OutputResponseCode             string `json:"output_ResponseCode"`
	OutputResponseDesc             string `json:"output_ResponseDesc"`
	OutputTransactionID            string `json:"output_TransactionID,omitempty"`
	OutputThirdPartyConversationID string `json:"output_ThirdPartyConversationID"`
}

// B2BSingleStage performs a Business-to-Business single stage transfer.
//
// Endpoint:
//   POST /{sandbox|openapi}/ipg/v2/{market}/b2bPayment/
//
// Bearer value:
//   RSA-encrypted SessionID returned by GenerateSessionKey.
func (c *Client) B2BSingleStage(ctx context.Context, sessionID string, request B2BSingleStageRequest) (*B2BSingleStageResponse, *RawResponse, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, nil, errors.New("sessionID is required")
	}
	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	var decoded B2BSingleStageResponse
	raw, err := c.do(ctx, http.MethodPost, b2bSingleStagePath, sessionID, request, &decoded)
	if err != nil {
		return nil, raw, err
	}
	return &decoded, raw, nil
}
