package mpesa

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const c2bSingleStagePath = "c2bPayment/singleStage/"

// C2BSingleStageRequest is the body for POST /c2bPayment/singleStage/.
type C2BSingleStageRequest struct {
	InputAmount                   string `json:"input_Amount"`
	InputCustomerMSISDN           string `json:"input_CustomerMSISDN"`
	InputCountry                  string `json:"input_Country"`
	InputCurrency                 string `json:"input_Currency"`
	InputServiceProviderCode      string `json:"input_ServiceProviderCode"`
	InputTransactionReference     string `json:"input_TransactionReference"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
	InputPurchasedItemsDesc       string `json:"input_PurchasedItemsDesc"`
}

// NewC2BSingleStageRequest creates a C2B request using the client's configured
// market country/currency. You may still edit any field before sending.
func (c *Client) NewC2BSingleStageRequest(amount, customerMSISDN, serviceProviderCode, transactionReference, thirdPartyConversationID, purchasedItemsDesc string) C2BSingleStageRequest {
	return C2BSingleStageRequest{
		InputAmount:                   amount,
		InputCustomerMSISDN:           customerMSISDN,
		InputCountry:                  c.cfg.Market.Country,
		InputCurrency:                 c.cfg.Market.Currency,
		InputServiceProviderCode:      serviceProviderCode,
		InputTransactionReference:     transactionReference,
		InputThirdPartyConversationID: thirdPartyConversationID,
		InputPurchasedItemsDesc:       purchasedItemsDesc,
	}
}

// Validate checks all mandatory fields. M-Pesa performs final regex and
// business validations server-side.
func (r C2BSingleStageRequest) Validate() error {
	required := map[string]string{
		"input_Amount":                   r.InputAmount,
		"input_CustomerMSISDN":           r.InputCustomerMSISDN,
		"input_Country":                  r.InputCountry,
		"input_Currency":                 r.InputCurrency,
		"input_ServiceProviderCode":      r.InputServiceProviderCode,
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

// C2BSingleStageResponse represents both sync and initial async responses.
// In async mode, OutputTransactionID may be empty until the callback arrives.
type C2BSingleStageResponse struct {
	OutputConversationID           string `json:"output_ConversationID"`
	OutputResponseCode             string `json:"output_ResponseCode"`
	OutputResponseDesc             string `json:"output_ResponseDesc"`
	OutputTransactionID            string `json:"output_TransactionID,omitempty"`
	OutputThirdPartyConversationID string `json:"output_ThirdPartyConversationID"`
}

// C2BSingleStage performs a Customer-to-Business single stage payment.
//
// Endpoint:
//   POST /{sandbox|openapi}/ipg/v2/{market}/c2bPayment/singleStage/
//
// Bearer value:
//   RSA-encrypted SessionID returned by GenerateSessionKey.
func (c *Client) C2BSingleStage(ctx context.Context, sessionID string, request C2BSingleStageRequest) (*C2BSingleStageResponse, *RawResponse, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, nil, errors.New("sessionID is required")
	}
	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	var decoded C2BSingleStageResponse
	raw, err := c.do(ctx, http.MethodPost, c2bSingleStagePath, sessionID, request, &decoded)
	if err != nil {
		return nil, raw, err
	}
	return &decoded, raw, nil
}
