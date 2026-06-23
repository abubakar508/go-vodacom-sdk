package mpesa

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const c2bMultiStagePath = "c2bPayment/multiStage/"
const defaultC2BMultiStageAPIVersion = "3.1"

// C2BMultiStageRequest is the body for POST /c2bPayment/multiStage/.
//
// It is similar to C2B Single Stage, with input_APIVersion required for API
// version 3.1 and above.
type C2BMultiStageRequest struct {
	InputAmount                   string `json:"input_Amount"`
	InputCustomerMSISDN           string `json:"input_CustomerMSISDN"`
	InputCountry                  string `json:"input_Country"`
	InputCurrency                 string `json:"input_Currency"`
	InputServiceProviderCode      string `json:"input_ServiceProviderCode"`
	InputTransactionReference     string `json:"input_TransactionReference"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
	InputPurchasedItemsDesc       string `json:"input_PurchasedItemsDesc"`
	InputAPIVersion               string `json:"input_APIVersion"`
}

// NewC2BMultiStageRequest creates a C2B multi-stage request using the client's
// configured market country/currency and default API version 3.1. You may edit
// any field before sending.
func (c *Client) NewC2BMultiStageRequest(amount, customerMSISDN, serviceProviderCode, transactionReference, thirdPartyConversationID, purchasedItemsDesc string) C2BMultiStageRequest {
	return C2BMultiStageRequest{
		InputAmount:                   amount,
		InputCustomerMSISDN:           customerMSISDN,
		InputCountry:                  c.country(),
		InputCurrency:                 c.currency(),
		InputServiceProviderCode:      serviceProviderCode,
		InputTransactionReference:     transactionReference,
		InputThirdPartyConversationID: thirdPartyConversationID,
		InputPurchasedItemsDesc:       purchasedItemsDesc,
		InputAPIVersion:               defaultC2BMultiStageAPIVersion,
	}
}

// Validate checks all mandatory C2B Multi Stage fields. M-Pesa performs final
// regex and business validations server-side.
func (r C2BMultiStageRequest) Validate() error {
	required := map[string]string{
		"input_Amount":                   r.InputAmount,
		"input_CustomerMSISDN":           r.InputCustomerMSISDN,
		"input_Country":                  r.InputCountry,
		"input_Currency":                 r.InputCurrency,
		"input_ServiceProviderCode":      r.InputServiceProviderCode,
		"input_TransactionReference":     r.InputTransactionReference,
		"input_ThirdPartyConversationID": r.InputThirdPartyConversationID,
		"input_PurchasedItemsDesc":       r.InputPurchasedItemsDesc,
		"input_APIVersion":               r.InputAPIVersion,
	}
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			return errors.New(name + " is required")
		}
	}
	if err := validateCommonTransactionFields(r.InputAmount, r.InputCountry, r.InputCurrency, r.InputServiceProviderCode, r.InputTransactionReference, r.InputThirdPartyConversationID); err != nil {
		return err
	}
	if err := validateMSISDN("input_CustomerMSISDN", r.InputCustomerMSISDN); err != nil {
		return err
	}
	return nil
}

// C2BMultiStageResponse represents both sync and initial async responses.
// In async mode, OutputTransactionID may be empty until the callback arrives.
type C2BMultiStageResponse struct {
	OutputConversationID           string `json:"output_ConversationID"`
	OutputResponseCode             string `json:"output_ResponseCode"`
	OutputResponseDesc             string `json:"output_ResponseDesc"`
	OutputTransactionID            string `json:"output_TransactionID,omitempty"`
	OutputThirdPartyConversationID string `json:"output_ThirdPartyConversationID"`
}

// C2BMultiStage initiates a Customer-to-Business multi-stage payment.
//
// Endpoint:
//
//	POST /{sandbox|openapi}/ipg/v2/{market}/c2bPayment/multiStage/
//
// Bearer value:
//
//	RSA-encrypted SessionID returned by GenerateSessionKey/GenerateSession.
func (c *Client) C2BMultiStage(ctx context.Context, sessionID string, request C2BMultiStageRequest) (*C2BMultiStageResponse, *RawResponse, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, nil, errors.New("sessionID is required; call GenerateSessionKey or GenerateSession first")
	}
	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	var decoded C2BMultiStageResponse
	raw, err := c.do(ctx, http.MethodPost, c2bMultiStagePath, sessionID, request, &decoded)
	if err != nil {
		return nil, raw, err
	}
	return &decoded, raw, nil
}

// C2BMultiStageWithSession performs C2B Multi Stage using a Session returned by GenerateSession.
func (c *Client) C2BMultiStageWithSession(ctx context.Context, session *Session, request C2BMultiStageRequest) (*C2BMultiStageResponse, *RawResponse, error) {
	if session == nil || !session.Valid() {
		return nil, nil, errors.New("valid session is required; call GenerateSession first")
	}
	return c.C2BMultiStage(ctx, session.ID, request)
}
