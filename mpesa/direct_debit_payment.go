package mpesa

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const directDebitPaymentPath = "directDebitPayment/"

// DirectDebitPaymentRequest is the body for POST /directDebitPayment/.
//
// Either InputMsisdnToken or InputCustomerMSISDN must be supplied. If both are
// supplied, M-Pesa requires them to match the same customer.
type DirectDebitPaymentRequest struct {
	InputMsisdnToken              string `json:"input_MsisdnToken,omitempty"`
	InputCustomerMSISDN           string `json:"input_CustomerMSISDN,omitempty"`
	InputCountry                  string `json:"input_Country"`
	InputServiceProviderCode      string `json:"input_ServiceProviderCode"`
	InputThirdPartyReference      string `json:"input_ThirdPartyReference"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
	InputAmount                   string `json:"input_Amount"`
	InputCurrency                 string `json:"input_Currency"`
	InputMandateID                string `json:"input_MandateID,omitempty"`
}

// NewDirectDebitPaymentRequest creates a Direct Debit payment request using the
// client's configured market country/currency. Provide either customerMSISDN or
// msisdnToken before sending; mandateID is optional but recommended when you
// have it from Direct Debit Create.
func (c *Client) NewDirectDebitPaymentRequest(amount, customerMSISDN, msisdnToken, serviceProviderCode, thirdPartyReference, thirdPartyConversationID, mandateID string) DirectDebitPaymentRequest {
	return DirectDebitPaymentRequest{
		InputMsisdnToken:              msisdnToken,
		InputCustomerMSISDN:           customerMSISDN,
		InputCountry:                  c.country(),
		InputServiceProviderCode:      serviceProviderCode,
		InputThirdPartyReference:      thirdPartyReference,
		InputThirdPartyConversationID: thirdPartyConversationID,
		InputAmount:                   amount,
		InputCurrency:                 c.currency(),
		InputMandateID:                mandateID,
	}
}

// NewDirectDebitPaymentWithMSISDN creates a Direct Debit payment request using
// a customer MSISDN as the customer identifier.
func (c *Client) NewDirectDebitPaymentWithMSISDN(amount, customerMSISDN, serviceProviderCode, thirdPartyReference, thirdPartyConversationID, mandateID string) DirectDebitPaymentRequest {
	return c.NewDirectDebitPaymentRequest(amount, customerMSISDN, "", serviceProviderCode, thirdPartyReference, thirdPartyConversationID, mandateID)
}

// NewDirectDebitPaymentWithToken creates a Direct Debit payment request using
// the encrypted MSISDN token returned by Direct Debit Create/Payment.
func (c *Client) NewDirectDebitPaymentWithToken(amount, msisdnToken, serviceProviderCode, thirdPartyReference, thirdPartyConversationID, mandateID string) DirectDebitPaymentRequest {
	return c.NewDirectDebitPaymentRequest(amount, "", msisdnToken, serviceProviderCode, thirdPartyReference, thirdPartyConversationID, mandateID)
}

// Validate checks mandatory Direct Debit Payment fields. M-Pesa performs final
// regex and business validations server-side.
func (r DirectDebitPaymentRequest) Validate() error {
	required := map[string]string{
		"input_Country":                  r.InputCountry,
		"input_ServiceProviderCode":      r.InputServiceProviderCode,
		"input_ThirdPartyReference":      r.InputThirdPartyReference,
		"input_ThirdPartyConversationID": r.InputThirdPartyConversationID,
		"input_Amount":                   r.InputAmount,
		"input_Currency":                 r.InputCurrency,
	}
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			return errors.New(name + " is required")
		}
	}
	if strings.TrimSpace(r.InputMsisdnToken) == "" && strings.TrimSpace(r.InputCustomerMSISDN) == "" {
		return errors.New("either input_MsisdnToken or input_CustomerMSISDN is required")
	}
	if err := validateAmount("input_Amount", r.InputAmount); err != nil {
		return err
	}
	if err := validateCurrency("input_Currency", r.InputCurrency); err != nil {
		return err
	}
	if err := validateShortCode("input_ServiceProviderCode", r.InputServiceProviderCode); err != nil {
		return err
	}
	if err := validateThirdPartyReference("input_ThirdPartyReference", r.InputThirdPartyReference); err != nil {
		return err
	}
	if err := validateThirdPartyConversationID("input_ThirdPartyConversationID", r.InputThirdPartyConversationID); err != nil {
		return err
	}
	if err := validateMSISDN("input_CustomerMSISDN", r.InputCustomerMSISDN); err != nil {
		return err
	}
	if err := validateMsisdnToken("input_MsisdnToken", r.InputMsisdnToken); err != nil {
		return err
	}
	if err := validateMandateID("input_MandateID", r.InputMandateID); err != nil {
		return err
	}
	return nil
}

// DirectDebitPaymentResponse represents both sync and initial async responses.
// Some documentation samples return output_TransactionReference, while async
// callback samples use input_TransactionID, so both transaction fields are kept.
type DirectDebitPaymentResponse struct {
	OutputResponseCode             string `json:"output_ResponseCode"`
	OutputResponseDesc             string `json:"output_ResponseDesc"`
	OutputTransactionReference     string `json:"output_TransactionReference,omitempty"`
	OutputTransactionID            string `json:"output_TransactionID,omitempty"`
	OutputMsisdnToken              string `json:"output_MsisdnToken,omitempty"`
	OutputConversationID           string `json:"output_ConversationID"`
	OutputThirdPartyConversationID string `json:"output_ThirdPartyConversationID"`
}

// DirectDebitPayment pays/debits an existing Direct Debit mandate.
//
// Endpoint:
//   POST /{sandbox|openapi}/ipg/v2/{market}/directDebitPayment/
//
// Bearer value:
//   RSA-encrypted SessionID returned by GenerateSessionKey/GenerateSession.
func (c *Client) DirectDebitPayment(ctx context.Context, sessionID string, request DirectDebitPaymentRequest) (*DirectDebitPaymentResponse, *RawResponse, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, nil, errors.New("sessionID is required; call GenerateSessionKey or GenerateSession first")
	}
	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	var decoded DirectDebitPaymentResponse
	raw, err := c.do(ctx, http.MethodPost, directDebitPaymentPath, sessionID, request, &decoded)
	if err != nil {
		return nil, raw, err
	}
	return &decoded, raw, nil
}

// DirectDebitPaymentWithSession pays/debits a mandate using a Session returned
// by GenerateSession.
func (c *Client) DirectDebitPaymentWithSession(ctx context.Context, session *Session, request DirectDebitPaymentRequest) (*DirectDebitPaymentResponse, *RawResponse, error) {
	if session == nil || !session.Valid() {
		return nil, nil, errors.New("valid session is required; call GenerateSession first")
	}
	return c.DirectDebitPayment(ctx, session.ID, request)
}
