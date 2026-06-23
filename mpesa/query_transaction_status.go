package mpesa

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const queryTransactionStatusPath = "queryTransactionStatus/"

// QueryTransactionStatusRequest contains the URL query parameters for
// GET /queryTransactionStatus/.
type QueryTransactionStatusRequest struct {
	InputQueryReference            string `json:"input_QueryReference"`
	InputServiceProviderCode       string `json:"input_ServiceProviderCode"`
	InputThirdPartyConversationID  string `json:"input_ThirdPartyConversationID"`
	InputCountry                   string `json:"input_Country"`
}

// NewQueryTransactionStatusRequest creates a query request using the client's
// configured market country. queryReference is usually a transaction ID or a
// conversation reference returned by a previous OpenAPI transaction flow.
func (c *Client) NewQueryTransactionStatusRequest(queryReference, serviceProviderCode, thirdPartyConversationID string) QueryTransactionStatusRequest {
	return QueryTransactionStatusRequest{
		InputQueryReference:           queryReference,
		InputServiceProviderCode:      serviceProviderCode,
		InputThirdPartyConversationID: thirdPartyConversationID,
		InputCountry:                  c.country(),
	}
}

// Validate checks all mandatory Query Transaction Status fields.
func (r QueryTransactionStatusRequest) Validate() error {
	required := map[string]string{
		"input_QueryReference":           r.InputQueryReference,
		"input_ServiceProviderCode":      r.InputServiceProviderCode,
		"input_ThirdPartyConversationID": r.InputThirdPartyConversationID,
		"input_Country":                  r.InputCountry,
	}
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			return errors.New(name + " is required")
		}
	}
	if err := validateQueryReference("input_QueryReference", r.InputQueryReference); err != nil {
		return err
	}
	if err := validateShortCode("input_ServiceProviderCode", r.InputServiceProviderCode); err != nil {
		return err
	}
	if err := validateThirdPartyConversationID("input_ThirdPartyConversationID", r.InputThirdPartyConversationID); err != nil {
		return err
	}
	return nil
}

func (r QueryTransactionStatusRequest) values() url.Values {
	values := url.Values{}
	values.Set("input_QueryReference", r.InputQueryReference)
	values.Set("input_ServiceProviderCode", r.InputServiceProviderCode)
	values.Set("input_ThirdPartyConversationID", r.InputThirdPartyConversationID)
	values.Set("input_Country", r.InputCountry)
	return values
}

// QueryTransactionStatusResponse represents the Query Transaction Status API response.
//
// M-Pesa deployments may return additional fields depending on the queried flow.
// The common response fields are represented here; the raw response is returned
// alongside this struct if you need to inspect market-specific additions.
type QueryTransactionStatusResponse struct {
	OutputConversationID             string `json:"output_ConversationID,omitempty"`
	OutputOriginalConversationID     string `json:"output_OriginalConversationID,omitempty"`
	OutputResponseCode               string `json:"output_ResponseCode"`
	OutputResponseDesc               string `json:"output_ResponseDesc"`
	OutputTransactionID              string `json:"output_TransactionID,omitempty"`
	OutputTransactionStatus          string `json:"output_TransactionStatus,omitempty"`
	OutputResponseTransactionStatus  string `json:"output_ResponseTransactionStatus,omitempty"`
	OutputThirdPartyConversationID   string `json:"output_ThirdPartyConversationID,omitempty"`
}

// TransactionStatus returns the status regardless of which documented market
// field name the OpenAPI response uses.
func (r QueryTransactionStatusResponse) TransactionStatus() string {
	return firstNonEmpty(r.OutputTransactionStatus, r.OutputResponseTransactionStatus)
}

// QueryTransactionStatus queries the status of a previous M-Pesa transaction.
//
// Endpoint:
//   GET /{sandbox|openapi}/ipg/v2/{market}/queryTransactionStatus/
//
// Query parameters:
//   input_QueryReference, input_ServiceProviderCode,
//   input_ThirdPartyConversationID, input_Country
//
// Bearer value:
//   RSA-encrypted SessionID returned by GenerateSessionKey/GenerateSession.
func (c *Client) QueryTransactionStatus(ctx context.Context, sessionID string, request QueryTransactionStatusRequest) (*QueryTransactionStatusResponse, *RawResponse, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, nil, errors.New("sessionID is required; call GenerateSessionKey or GenerateSession first")
	}
	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	var decoded QueryTransactionStatusResponse
	raw, err := c.doQuery(ctx, http.MethodGet, queryTransactionStatusPath, request.values(), sessionID, &decoded)
	if err != nil {
		return nil, raw, err
	}
	return &decoded, raw, nil
}

// QueryTransactionStatusWithSession queries transaction status using a Session
// returned by GenerateSession.
func (c *Client) QueryTransactionStatusWithSession(ctx context.Context, session *Session, request QueryTransactionStatusRequest) (*QueryTransactionStatusResponse, *RawResponse, error) {
	if session == nil || !session.Valid() {
		return nil, nil, errors.New("valid session is required; call GenerateSession first")
	}
	return c.QueryTransactionStatus(ctx, session.ID, request)
}
