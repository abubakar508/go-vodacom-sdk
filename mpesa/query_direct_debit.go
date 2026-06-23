package mpesa

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const queryDirectDebitPath = "queryDirectDebit/"

const (
	QueryBalanceAmountTrue  = "True"
	QueryBalanceAmountFalse = "False"
)

// QueryDirectDebitRequest contains the URL query parameters for
// GET /queryDirectDebit/.
//
// Either InputMsisdnToken or InputCustomerMSISDN must be supplied. If both are
// supplied, M-Pesa requires them to match the same customer. If
// InputQueryBalanceAmount is "True", InputBalanceAmount must be supplied.
type QueryDirectDebitRequest struct {
	InputQueryBalanceAmount      string `json:"input_QueryBalanceAmount"`
	InputBalanceAmount           string `json:"input_BalanceAmount,omitempty"`
	InputCountry                 string `json:"input_Country"`
	InputCustomerMSISDN          string `json:"input_CustomerMSISDN,omitempty"`
	InputMsisdnToken             string `json:"input_MsisdnToken,omitempty"`
	InputServiceProviderCode     string `json:"input_ServiceProviderCode"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
	InputThirdPartyReference     string `json:"input_ThirdPartyReference"`
	InputMandateID               string `json:"input_MandateID,omitempty"`
	InputCurrency                string `json:"input_Currency"`
}

// NewQueryDirectDebitRequest creates a Query Direct Debit request using the
// client's configured market country/currency. Provide either customerMSISDN or
// msisdnToken before sending. If queryBalanceAmount is "True", set
// balanceAmount.
func (c *Client) NewQueryDirectDebitRequest(queryBalanceAmount, balanceAmount, customerMSISDN, msisdnToken, serviceProviderCode, thirdPartyConversationID, thirdPartyReference, mandateID string) QueryDirectDebitRequest {
	return QueryDirectDebitRequest{
		InputQueryBalanceAmount:       queryBalanceAmount,
		InputBalanceAmount:            balanceAmount,
		InputCountry:                  c.country(),
		InputCustomerMSISDN:           customerMSISDN,
		InputMsisdnToken:              msisdnToken,
		InputServiceProviderCode:      serviceProviderCode,
		InputThirdPartyConversationID: thirdPartyConversationID,
		InputThirdPartyReference:      thirdPartyReference,
		InputMandateID:                mandateID,
		InputCurrency:                 c.currency(),
	}
}

// NewQueryDirectDebitWithMSISDN creates a query request using a customer MSISDN.
func (c *Client) NewQueryDirectDebitWithMSISDN(queryBalanceAmount, balanceAmount, customerMSISDN, serviceProviderCode, thirdPartyConversationID, thirdPartyReference, mandateID string) QueryDirectDebitRequest {
	return c.NewQueryDirectDebitRequest(queryBalanceAmount, balanceAmount, customerMSISDN, "", serviceProviderCode, thirdPartyConversationID, thirdPartyReference, mandateID)
}

// NewQueryDirectDebitWithToken creates a query request using an encrypted MSISDN token.
func (c *Client) NewQueryDirectDebitWithToken(queryBalanceAmount, balanceAmount, msisdnToken, serviceProviderCode, thirdPartyConversationID, thirdPartyReference, mandateID string) QueryDirectDebitRequest {
	return c.NewQueryDirectDebitRequest(queryBalanceAmount, balanceAmount, "", msisdnToken, serviceProviderCode, thirdPartyConversationID, thirdPartyReference, mandateID)
}

// Validate checks mandatory Query Direct Debit fields. M-Pesa performs final
// regex and business validations server-side.
func (r QueryDirectDebitRequest) Validate() error {
	required := map[string]string{
		"input_QueryBalanceAmount":      r.InputQueryBalanceAmount,
		"input_Country":                 r.InputCountry,
		"input_ServiceProviderCode":     r.InputServiceProviderCode,
		"input_ThirdPartyConversationID": r.InputThirdPartyConversationID,
		"input_ThirdPartyReference":     r.InputThirdPartyReference,
		"input_Currency":                r.InputCurrency,
	}
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			return errors.New(name + " is required")
		}
	}
	if strings.TrimSpace(r.InputCustomerMSISDN) == "" && strings.TrimSpace(r.InputMsisdnToken) == "" {
		return errors.New("either input_CustomerMSISDN or input_MsisdnToken is required")
	}
	if r.InputQueryBalanceAmount != QueryBalanceAmountTrue && r.InputQueryBalanceAmount != QueryBalanceAmountFalse {
		return errors.New("input_QueryBalanceAmount must be True or False")
	}
	if r.InputQueryBalanceAmount == QueryBalanceAmountTrue && strings.TrimSpace(r.InputBalanceAmount) == "" {
		return errors.New("input_BalanceAmount is required when input_QueryBalanceAmount is True")
	}
	if err := validateAmount("input_BalanceAmount", r.InputBalanceAmount); err != nil {
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

func (r QueryDirectDebitRequest) values() url.Values {
	values := url.Values{}
	values.Set("input_QueryBalanceAmount", r.InputQueryBalanceAmount)
	if r.InputBalanceAmount != "" {
		values.Set("input_BalanceAmount", r.InputBalanceAmount)
	}
	values.Set("input_Country", r.InputCountry)
	if r.InputCustomerMSISDN != "" {
		values.Set("input_CustomerMSISDN", r.InputCustomerMSISDN)
	}
	if r.InputMsisdnToken != "" {
		values.Set("input_MsisdnToken", r.InputMsisdnToken)
	}
	values.Set("input_ServiceProviderCode", r.InputServiceProviderCode)
	values.Set("input_ThirdPartyConversationID", r.InputThirdPartyConversationID)
	values.Set("input_ThirdPartyReference", r.InputThirdPartyReference)
	if r.InputMandateID != "" {
		values.Set("input_MandateID", r.InputMandateID)
	}
	values.Set("input_Currency", r.InputCurrency)
	return values
}

// QueryDirectDebitResponse represents both sync and initial async responses.
type QueryDirectDebitResponse struct {
	OutputResponseCode             string `json:"output_ResponseCode"`
	OutputResponseDesc             string `json:"output_ResponseDesc"`
	OutputTransactionReference     string `json:"output_TransactionReference,omitempty"`
	OutputConversationID           string `json:"output_ConversationID"`
	OutputThirdPartyConversationID string `json:"output_ThirdPartyConversationID"`
	OutputSufficientBalance        string `json:"output_SufficientBalance,omitempty"`
	OutputMsisdnToken              string `json:"output_MsisdnToken,omitempty"`
	OutputMandateID                string `json:"output_MandateID,omitempty"`
	OutputMandateStatus            string `json:"output_MandateStatus,omitempty"`
	OutputAccountStatus            string `json:"output_AccountStatus,omitempty"`
	OutputFirstPaymentDate         string `json:"output_FirstPaymentDate,omitempty"`
	OutputFrequency                string `json:"output_Frequency,omitempty"`
	OutputPaymentDayFrom           string `json:"output_PaymentDayFrom,omitempty"`
	OutputPaymentDayTo             string `json:"output_PaymentDayTo,omitempty"`
	OutputExpiryDate               string `json:"output_ExpiryDate,omitempty"`
}

// QueryDirectDebit queries the status of a Direct Debit mandate and optionally
// checks customer balance sufficiency against an amount.
//
// Endpoint:
//   GET /{sandbox|openapi}/ipg/v2/{market}/queryDirectDebit/
//
// Bearer value:
//   RSA-encrypted SessionID returned by GenerateSessionKey/GenerateSession.
func (c *Client) QueryDirectDebit(ctx context.Context, sessionID string, request QueryDirectDebitRequest) (*QueryDirectDebitResponse, *RawResponse, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, nil, errors.New("sessionID is required; call GenerateSessionKey or GenerateSession first")
	}
	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	var decoded QueryDirectDebitResponse
	raw, err := c.doQuery(ctx, http.MethodGet, queryDirectDebitPath, request.values(), sessionID, &decoded)
	if err != nil {
		return nil, raw, err
	}
	return &decoded, raw, nil
}

// QueryDirectDebitWithSession queries Direct Debit mandate status using a
// Session returned by GenerateSession.
func (c *Client) QueryDirectDebitWithSession(ctx context.Context, session *Session, request QueryDirectDebitRequest) (*QueryDirectDebitResponse, *RawResponse, error) {
	if session == nil || !session.Valid() {
		return nil, nil, errors.New("valid session is required; call GenerateSession first")
	}
	return c.QueryDirectDebit(ctx, session.ID, request)
}
