package mpesa

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const b2cSingleStagePath = "b2cPayment/"

// B2CSingleStageRequest is the body for POST /b2cPayment/.
//
// It is used for business-to-customer disbursements such as salary payments,
// business funds transfers, and charity payouts.
type B2CSingleStageRequest struct {
	InputAmount                   string `json:"input_Amount"`
	InputCustomerMSISDN           string `json:"input_CustomerMSISDN"`
	InputCountry                  string `json:"input_Country"`
	InputCurrency                 string `json:"input_Currency"`
	InputServiceProviderCode      string `json:"input_ServiceProviderCode"`
	InputTransactionReference     string `json:"input_TransactionReference"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
	InputPaymentItemsDesc         string `json:"input_PaymentItemsDesc"`
}

// NewB2CSingleStageRequest creates a B2C request using the client's configured
// market country/currency. You may still edit any field before sending.
func (c *Client) NewB2CSingleStageRequest(amount, customerMSISDN, serviceProviderCode, transactionReference, thirdPartyConversationID, paymentItemsDesc string) B2CSingleStageRequest {
	return B2CSingleStageRequest{
		InputAmount:                   amount,
		InputCustomerMSISDN:           customerMSISDN,
		InputCountry:                  c.country(),
		InputCurrency:                 c.currency(),
		InputServiceProviderCode:      serviceProviderCode,
		InputTransactionReference:     transactionReference,
		InputThirdPartyConversationID: thirdPartyConversationID,
		InputPaymentItemsDesc:         paymentItemsDesc,
	}
}

// Validate checks all mandatory B2C Single Stage fields. M-Pesa performs final
// regex and business validations server-side.
func (r B2CSingleStageRequest) Validate() error {
	required := map[string]string{
		"input_Amount":                   r.InputAmount,
		"input_CustomerMSISDN":           r.InputCustomerMSISDN,
		"input_Country":                  r.InputCountry,
		"input_Currency":                 r.InputCurrency,
		"input_ServiceProviderCode":      r.InputServiceProviderCode,
		"input_TransactionReference":     r.InputTransactionReference,
		"input_ThirdPartyConversationID": r.InputThirdPartyConversationID,
		"input_PaymentItemsDesc":         r.InputPaymentItemsDesc,
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

// B2CSingleStageResponse represents both sync and initial async responses.
// In async mode, OutputTransactionID may be empty until the callback arrives.
type B2CSingleStageResponse struct {
	OutputConversationID           string `json:"output_ConversationID"`
	OutputResponseCode             string `json:"output_ResponseCode"`
	OutputResponseDesc             string `json:"output_ResponseDesc"`
	OutputTransactionID            string `json:"output_TransactionID,omitempty"`
	OutputThirdPartyConversationID string `json:"output_ThirdPartyConversationID"`
}

// B2CSingleStage performs a Business-to-Customer single stage disbursement.
//
// Endpoint:
//   POST /{sandbox|openapi}/ipg/v2/{market}/b2cPayment/
//
// Bearer value:
//   RSA-encrypted SessionID returned by GenerateSessionKey.
func (c *Client) B2CSingleStage(ctx context.Context, sessionID string, request B2CSingleStageRequest) (*B2CSingleStageResponse, *RawResponse, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, nil, errors.New("sessionID is required; call GenerateSessionKey or GenerateSession first")
	}
	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	var decoded B2CSingleStageResponse
	raw, err := c.do(ctx, http.MethodPost, b2cSingleStagePath, sessionID, request, &decoded)
	if err != nil {
		return nil, raw, err
	}
	return &decoded, raw, nil
}

// B2CSingleStageWithSession performs B2C using a Session returned by GenerateSession.
func (c *Client) B2CSingleStageWithSession(ctx context.Context, session *Session, request B2CSingleStageRequest) (*B2CSingleStageResponse, *RawResponse, error) {
	if session == nil || !session.Valid() {
		return nil, nil, errors.New("valid session is required; call GenerateSession first")
	}
	return c.B2CSingleStage(ctx, session.ID, request)
}
