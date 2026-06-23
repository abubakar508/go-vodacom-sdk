package mpesa

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const updateTransactionStatusPath = "updateTransactionStatus/"
const defaultUpdateTransactionStatusAPIVersion = "3.1"

const (
	// TransactionStatusUncommit indicates funds should be uncommitted.
	TransactionStatusUncommit = "0"
	// TransactionStatusCommit indicates funds should be committed.
	TransactionStatusCommit = "1"
)

// UpdateTransactionStatusRequest is the body for PUT /updateTransactionStatus/.
//
// The OpenAPI documentation names the operation field input_CustomerMSISDN even
// though its description says it is the operation type: "1" for commit and "0"
// for uncommit. The SDK keeps the JSON field name exactly as documented while
// the constructor parameter is named operationType for clarity.
type UpdateTransactionStatusRequest struct {
	InputCountry                  string `json:"input_Country"`
	InputVoucherCode              string `json:"input_VoucherCode"`
	InputCustomerMSISDN           string `json:"input_CustomerMSISDN"`
	InputServiceProviderCode      string `json:"input_ServiceProviderCode"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
	InputTransactionID            string `json:"input_TransactionID"`
	InputAPIVersion               string `json:"input_APIVersion"`
}

// NewUpdateTransactionStatusRequest creates an update-transaction-status request
// using the client's configured market country and default API version 3.1.
// operationType must be "1" to commit or "0" to uncommit.
func (c *Client) NewUpdateTransactionStatusRequest(operationType, voucherCode, serviceProviderCode, thirdPartyConversationID, transactionID string) UpdateTransactionStatusRequest {
	return UpdateTransactionStatusRequest{
		InputCountry:                  c.country(),
		InputVoucherCode:              voucherCode,
		InputCustomerMSISDN:           operationType,
		InputServiceProviderCode:      serviceProviderCode,
		InputThirdPartyConversationID: thirdPartyConversationID,
		InputTransactionID:            transactionID,
		InputAPIVersion:               defaultUpdateTransactionStatusAPIVersion,
	}
}

// Validate checks all mandatory Update Transaction Status fields. M-Pesa
// performs final regex and business validations server-side.
func (r UpdateTransactionStatusRequest) Validate() error {
	required := map[string]string{
		"input_Country":                  r.InputCountry,
		"input_VoucherCode":              r.InputVoucherCode,
		"input_CustomerMSISDN":           r.InputCustomerMSISDN,
		"input_ServiceProviderCode":      r.InputServiceProviderCode,
		"input_ThirdPartyConversationID": r.InputThirdPartyConversationID,
		"input_TransactionID":            r.InputTransactionID,
		"input_APIVersion":               r.InputAPIVersion,
	}
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			return errors.New(name + " is required")
		}
	}
	if r.InputCustomerMSISDN != TransactionStatusCommit && r.InputCustomerMSISDN != TransactionStatusUncommit {
		return errors.New("input_CustomerMSISDN must be operation type 1 (commit) or 0 (uncommit)")
	}
	if err := validateVoucherCode("input_VoucherCode", r.InputVoucherCode); err != nil {
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

// UpdateTransactionStatusResponse represents both sync and initial async responses.
// Sync success may return output_ResponseCode INS-GAR-0, while async accepted
// responses may return INS-A-0/INS-GAR-0 depending on the flow/market.
type UpdateTransactionStatusResponse struct {
	OutputConversationID           string `json:"output_ConversationID"`
	OutputResponseCode             string `json:"output_ResponseCode"`
	OutputResponseDesc             string `json:"output_ResponseDesc"`
	OutputTransactionID            string `json:"output_TransactionID,omitempty"`
	OutputThirdPartyConversationID string `json:"output_ThirdPartyConversationID"`
}

// UpdateTransactionStatus commits or uncommits a transaction.
//
// Endpoint:
//   PUT /{sandbox|openapi}/ipg/v2/{market}/updateTransactionStatus/
//
// Bearer value:
//   RSA-encrypted SessionID returned by GenerateSessionKey/GenerateSession.
func (c *Client) UpdateTransactionStatus(ctx context.Context, sessionID string, request UpdateTransactionStatusRequest) (*UpdateTransactionStatusResponse, *RawResponse, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, nil, errors.New("sessionID is required; call GenerateSessionKey or GenerateSession first")
	}
	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	var decoded UpdateTransactionStatusResponse
	raw, err := c.do(ctx, http.MethodPut, updateTransactionStatusPath, sessionID, request, &decoded)
	if err != nil {
		return nil, raw, err
	}
	return &decoded, raw, nil
}

// UpdateTransactionStatusWithSession performs update transaction status using a
// Session returned by GenerateSession.
func (c *Client) UpdateTransactionStatusWithSession(ctx context.Context, session *Session, request UpdateTransactionStatusRequest) (*UpdateTransactionStatusResponse, *RawResponse, error) {
	if session == nil || !session.Valid() {
		return nil, nil, errors.New("valid session is required; call GenerateSession first")
	}
	return c.UpdateTransactionStatus(ctx, session.ID, request)
}
