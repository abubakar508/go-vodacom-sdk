package mpesa

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const queryBeneficiaryNamePath = "queryBeneficiaryName/"

const (
	// KYCQueryTypeName is currently the only documented KYC query type.
	KYCQueryTypeName = "Name"
)

// QueryBeneficiaryNameRequest contains the URL query parameters for
// GET /queryBeneficiaryName/.
//
// The documentation currently lists this API only for Vodafone Ghana, but the
// SDK uses the configured market so future/portal-enabled markets can still be
// used if M-Pesa enables them.
type QueryBeneficiaryNameRequest struct {
	InputCustomerMSISDN           string `json:"input_CustomerMSISDN"`
	InputCountry                  string `json:"input_Country"`
	InputServiceProviderCode      string `json:"input_ServiceProviderCode"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
	InputKycQueryType             string `json:"input_KycQueryType"`
}

// NewQueryBeneficiaryNameRequest creates a request using the client's configured
// market country. kycQueryType defaults to "Name" because it is currently the
// only documented/supported value.
func (c *Client) NewQueryBeneficiaryNameRequest(customerMSISDN, serviceProviderCode, thirdPartyConversationID string) QueryBeneficiaryNameRequest {
	return QueryBeneficiaryNameRequest{
		InputCustomerMSISDN:           customerMSISDN,
		InputCountry:                  c.country(),
		InputServiceProviderCode:      serviceProviderCode,
		InputThirdPartyConversationID: thirdPartyConversationID,
		InputKycQueryType:             KYCQueryTypeName,
	}
}

// Validate checks all mandatory Query Beneficiary Name fields.
func (r QueryBeneficiaryNameRequest) Validate() error {
	required := map[string]string{
		"input_CustomerMSISDN":           r.InputCustomerMSISDN,
		"input_Country":                  r.InputCountry,
		"input_ServiceProviderCode":      r.InputServiceProviderCode,
		"input_ThirdPartyConversationID": r.InputThirdPartyConversationID,
		"input_KycQueryType":             r.InputKycQueryType,
	}
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			return errors.New(name + " is required")
		}
	}
	if r.InputKycQueryType != KYCQueryTypeName {
		return errors.New("input_KycQueryType must be Name")
	}
	if err := validateMSISDN("input_CustomerMSISDN", r.InputCustomerMSISDN); err != nil {
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

func (r QueryBeneficiaryNameRequest) values() url.Values {
	values := url.Values{}
	values.Set("input_CustomerMSISDN", r.InputCustomerMSISDN)
	values.Set("input_Country", r.InputCountry)
	values.Set("input_ServiceProviderCode", r.InputServiceProviderCode)
	values.Set("input_ThirdPartyConversationID", r.InputThirdPartyConversationID)
	values.Set("input_KycQueryType", r.InputKycQueryType)
	return values
}

// QueryBeneficiaryNameResponse represents both sync and initial async responses.
type QueryBeneficiaryNameResponse struct {
	OutputConversationID           string `json:"output_ConversationID"`
	OutputResponseCode             string `json:"output_ResponseCode"`
	OutputResponseDesc             string `json:"output_ResponseDesc"`
	OutputCustomerFirstName        string `json:"output_CustomerFirstName,omitempty"`
	OutputCustomerLastName         string `json:"output_CustomerLastName,omitempty"`
	OutputThirdPartyConversationID string `json:"output_ThirdPartyConversationID"`
}

// QueryBeneficiaryName queries KYC name details for a customer.
//
// Endpoint:
//   GET /{sandbox|openapi}/ipg/v2/{market}/queryBeneficiaryName/
//
// Query parameters:
//   input_CustomerMSISDN, input_Country, input_ServiceProviderCode,
//   input_ThirdPartyConversationID, input_KycQueryType
//
// Bearer value:
//   RSA-encrypted SessionID returned by GenerateSessionKey/GenerateSession.
func (c *Client) QueryBeneficiaryName(ctx context.Context, sessionID string, request QueryBeneficiaryNameRequest) (*QueryBeneficiaryNameResponse, *RawResponse, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, nil, errors.New("sessionID is required; call GenerateSessionKey or GenerateSession first")
	}
	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	var decoded QueryBeneficiaryNameResponse
	raw, err := c.doQuery(ctx, http.MethodGet, queryBeneficiaryNamePath, request.values(), sessionID, &decoded)
	if err != nil {
		return nil, raw, err
	}
	return &decoded, raw, nil
}

// QueryBeneficiaryNameWithSession queries beneficiary name using a Session
// returned by GenerateSession.
func (c *Client) QueryBeneficiaryNameWithSession(ctx context.Context, session *Session, request QueryBeneficiaryNameRequest) (*QueryBeneficiaryNameResponse, *RawResponse, error) {
	if session == nil || !session.Valid() {
		return nil, nil, errors.New("valid session is required; call GenerateSession first")
	}
	return c.QueryBeneficiaryName(ctx, session.ID, request)
}
