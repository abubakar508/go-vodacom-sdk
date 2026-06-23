package mpesa

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const reversalPath = "reversal/"

// ReversalRequest is the body for POST /reversal/.
//
// It reverses a previously completed M-Pesa transaction. The transaction being
// reversed is identified by input_TransactionID, and the amount being reversed
// is sent as input_ReversalAmount.
type ReversalRequest struct {
	InputReversalAmount            string `json:"input_ReversalAmount"`
	InputCountry                   string `json:"input_Country"`
	InputServiceProviderCode       string `json:"input_ServiceProviderCode"`
	InputThirdPartyConversationID  string `json:"input_ThirdPartyConversationID"`
	InputTransactionID             string `json:"input_TransactionID"`
}

// NewReversalRequest creates a reversal request using the client's configured
// market country. You may still edit any field before sending.
func (c *Client) NewReversalRequest(reversalAmount, serviceProviderCode, thirdPartyConversationID, transactionID string) ReversalRequest {
	return ReversalRequest{
		InputReversalAmount:           reversalAmount,
		InputCountry:                  c.country(),
		InputServiceProviderCode:      serviceProviderCode,
		InputThirdPartyConversationID: thirdPartyConversationID,
		InputTransactionID:            transactionID,
	}
}

// Validate checks all mandatory Reversal fields. M-Pesa performs final regex
// and business validations server-side.
func (r ReversalRequest) Validate() error {
	required := map[string]string{
		"input_ReversalAmount":           r.InputReversalAmount,
		"input_Country":                  r.InputCountry,
		"input_ServiceProviderCode":      r.InputServiceProviderCode,
		"input_ThirdPartyConversationID": r.InputThirdPartyConversationID,
		"input_TransactionID":            r.InputTransactionID,
	}
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			return errors.New(name + " is required")
		}
	}
	if err := validateAmount("input_ReversalAmount", r.InputReversalAmount); err != nil {
		return err
	}
	if err := validateShortCode("input_ServiceProviderCode", r.InputServiceProviderCode); err != nil {
		return err
	}
	if err := validateThirdPartyConversationID("input_ThirdPartyConversationID", r.InputThirdPartyConversationID); err != nil {
		return err
	}
	if err := validateTransactionID("input_TransactionID", r.InputTransactionID); err != nil {
		return err
	}
	return nil
}

// ReversalResponse represents the Reversal API response.
//
// The documented synchronous/async response shape follows the common OpenAPI
// transaction response fields. OutputTransactionID is optional because some
// markets/flows may return only the conversation references for asynchronous
// processing.
type ReversalResponse struct {
	OutputConversationID           string `json:"output_ConversationID,omitempty"`
	OutputResponseCode             string `json:"output_ResponseCode"`
	OutputResponseDesc             string `json:"output_ResponseDesc"`
	OutputTransactionID            string `json:"output_TransactionID,omitempty"`
	OutputThirdPartyConversationID string `json:"output_ThirdPartyConversationID,omitempty"`
}

// Reversal performs a transaction reversal.
//
// Endpoint:
//   POST /{sandbox|openapi}/ipg/v2/{market}/reversal/
//
// Bearer value:
//   RSA-encrypted SessionID returned by GenerateSessionKey/GenerateSession.
func (c *Client) Reversal(ctx context.Context, sessionID string, request ReversalRequest) (*ReversalResponse, *RawResponse, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, nil, errors.New("sessionID is required; call GenerateSessionKey or GenerateSession first")
	}
	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	var decoded ReversalResponse
	raw, err := c.do(ctx, http.MethodPost, reversalPath, sessionID, request, &decoded)
	if err != nil {
		return nil, raw, err
	}
	return &decoded, raw, nil
}

// ReversalWithSession performs a reversal using a Session returned by GenerateSession.
func (c *Client) ReversalWithSession(ctx context.Context, session *Session, request ReversalRequest) (*ReversalResponse, *RawResponse, error) {
	if session == nil || !session.Valid() {
		return nil, nil, errors.New("valid session is required; call GenerateSession first")
	}
	return c.Reversal(ctx, session.ID, request)
}
