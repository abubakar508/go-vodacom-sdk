package mpesa

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const cancelDirectDebitPath = "directDebitCancel/"

// CancelDirectDebitRequest is the body for PUT /directDebitCancel/.
//
// Either InputMsisdnToken or InputCustomerMSISDN must be supplied. If both are
// supplied, M-Pesa requires them to match the same customer.
type CancelDirectDebitRequest struct {
	InputMsisdnToken              string `json:"input_MsisdnToken,omitempty"`
	InputCustomerMSISDN           string `json:"input_CustomerMSISDN,omitempty"`
	InputCountry                  string `json:"input_Country"`
	InputServiceProviderCode      string `json:"input_ServiceProviderCode"`
	InputThirdPartyReference      string `json:"input_ThirdPartyReference"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
	InputMandateID                string `json:"input_MandateID,omitempty"`
}

// NewCancelDirectDebitRequest creates a Direct Debit cancel request using the
// client's configured market country. Provide either customerMSISDN or
// msisdnToken before sending; mandateID is optional but recommended when
// available from Direct Debit Create/Query.
func (c *Client) NewCancelDirectDebitRequest(customerMSISDN, msisdnToken, serviceProviderCode, thirdPartyReference, thirdPartyConversationID, mandateID string) CancelDirectDebitRequest {
	return CancelDirectDebitRequest{
		InputMsisdnToken:              msisdnToken,
		InputCustomerMSISDN:           customerMSISDN,
		InputCountry:                  c.country(),
		InputServiceProviderCode:      serviceProviderCode,
		InputThirdPartyReference:      thirdPartyReference,
		InputThirdPartyConversationID: thirdPartyConversationID,
		InputMandateID:                mandateID,
	}
}

// NewCancelDirectDebitWithMSISDN creates a Direct Debit cancel request using a
// customer MSISDN as the customer identifier.
func (c *Client) NewCancelDirectDebitWithMSISDN(customerMSISDN, serviceProviderCode, thirdPartyReference, thirdPartyConversationID, mandateID string) CancelDirectDebitRequest {
	return c.NewCancelDirectDebitRequest(customerMSISDN, "", serviceProviderCode, thirdPartyReference, thirdPartyConversationID, mandateID)
}

// NewCancelDirectDebitWithToken creates a Direct Debit cancel request using the
// encrypted MSISDN token returned by Direct Debit APIs.
func (c *Client) NewCancelDirectDebitWithToken(msisdnToken, serviceProviderCode, thirdPartyReference, thirdPartyConversationID, mandateID string) CancelDirectDebitRequest {
	return c.NewCancelDirectDebitRequest("", msisdnToken, serviceProviderCode, thirdPartyReference, thirdPartyConversationID, mandateID)
}

// Validate checks mandatory Cancel Direct Debit fields. M-Pesa performs final
// regex and business validations server-side.
func (r CancelDirectDebitRequest) Validate() error {
	required := map[string]string{
		"input_Country":                  r.InputCountry,
		"input_ServiceProviderCode":      r.InputServiceProviderCode,
		"input_ThirdPartyReference":      r.InputThirdPartyReference,
		"input_ThirdPartyConversationID": r.InputThirdPartyConversationID,
	}
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			return errors.New(name + " is required")
		}
	}
	if strings.TrimSpace(r.InputMsisdnToken) == "" && strings.TrimSpace(r.InputCustomerMSISDN) == "" {
		return errors.New("either input_MsisdnToken or input_CustomerMSISDN is required")
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

// CancelDirectDebitResponse represents both sync and initial async responses.
type CancelDirectDebitResponse struct {
	OutputResponseCode             string `json:"output_ResponseCode"`
	OutputResponseDesc             string `json:"output_ResponseDesc"`
	OutputTransactionReference     string `json:"output_TransactionReference,omitempty"`
	OutputMsisdnToken              string `json:"output_MsisdnToken,omitempty"`
	OutputConversationID           string `json:"output_ConversationID"`
	OutputThirdPartyConversationID string `json:"output_ThirdPartyConversationID"`
}

// CancelDirectDebit cancels a Direct Debit mandate.
//
// Endpoint:
//   PUT /{sandbox|openapi}/ipg/v2/{market}/directDebitCancel/
//
// Bearer value:
//   RSA-encrypted SessionID returned by GenerateSessionKey/GenerateSession.
func (c *Client) CancelDirectDebit(ctx context.Context, sessionID string, request CancelDirectDebitRequest) (*CancelDirectDebitResponse, *RawResponse, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, nil, errors.New("sessionID is required; call GenerateSessionKey or GenerateSession first")
	}
	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	var decoded CancelDirectDebitResponse
	raw, err := c.do(ctx, http.MethodPut, cancelDirectDebitPath, sessionID, request, &decoded)
	if err != nil {
		return nil, raw, err
	}
	return &decoded, raw, nil
}

// CancelDirectDebitWithSession cancels a mandate using a Session returned by
// GenerateSession.
func (c *Client) CancelDirectDebitWithSession(ctx context.Context, session *Session, request CancelDirectDebitRequest) (*CancelDirectDebitResponse, *RawResponse, error) {
	if session == nil || !session.Valid() {
		return nil, nil, errors.New("valid session is required; call GenerateSession first")
	}
	return c.CancelDirectDebit(ctx, session.ID, request)
}
