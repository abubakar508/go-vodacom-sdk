package mpesa

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const directDebitCreatePath = "directDebitCreation/"

const (
	DirectDebitAgreedTCNo  = "0"
	DirectDebitAgreedTCYes = "1"

	DirectDebitFrequencyOnceOff    = "01"
	DirectDebitFrequencyDaily      = "02"
	DirectDebitFrequencyWeekly     = "03"
	DirectDebitFrequencyMonthly    = "04"
	DirectDebitFrequencyQuarterly  = "05"
	DirectDebitFrequencyHalfYearly = "06"
	DirectDebitFrequencyYearly     = "07"
	DirectDebitFrequencyOnDemand   = "08"
)

// DirectDebitCreateRequest is the body for POST /directDebitCreation/.
//
// Direct Debit Create requests customer consent to create a mandate that allows
// the organization to debit the customer's account at an agreed frequency and
// amount. Optional schedule fields use omitempty so they are not sent unless set.
type DirectDebitCreateRequest struct {
	InputCustomerMSISDN           string `json:"input_CustomerMSISDN"`
	InputCountry                  string `json:"input_Country"`
	InputServiceProviderCode      string `json:"input_ServiceProviderCode"`
	InputThirdPartyReference      string `json:"input_ThirdPartyReference"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
	InputAgreedTC                 string `json:"input_AgreedTC"`
	InputFirstPaymentDate         string `json:"input_FirstPaymentDate,omitempty"`
	InputFrequency                string `json:"input_Frequency,omitempty"`
	InputStartRangeOfDays         string `json:"input_StartRangeOfDays,omitempty"`
	InputEndRangeOfDays           string `json:"input_EndRangeOfDays,omitempty"`
	InputExpiryDate               string `json:"input_ExpiryDate,omitempty"`
}

// NewDirectDebitCreateRequest creates a Direct Debit mandate creation request
// using the client's configured market country. Optional mandate schedule fields
// can be set on the returned struct before sending.
func (c *Client) NewDirectDebitCreateRequest(customerMSISDN, serviceProviderCode, thirdPartyReference, thirdPartyConversationID, agreedTC string) DirectDebitCreateRequest {
	return DirectDebitCreateRequest{
		InputCustomerMSISDN:           customerMSISDN,
		InputCountry:                  c.country(),
		InputServiceProviderCode:      serviceProviderCode,
		InputThirdPartyReference:      thirdPartyReference,
		InputThirdPartyConversationID: thirdPartyConversationID,
		InputAgreedTC:                 agreedTC,
	}
}

// Validate checks mandatory Direct Debit Create fields and basic frequency rules.
// M-Pesa performs final regex and business validations server-side.
func (r DirectDebitCreateRequest) Validate() error {
	required := map[string]string{
		"input_CustomerMSISDN":           r.InputCustomerMSISDN,
		"input_Country":                  r.InputCountry,
		"input_ServiceProviderCode":      r.InputServiceProviderCode,
		"input_ThirdPartyReference":      r.InputThirdPartyReference,
		"input_ThirdPartyConversationID": r.InputThirdPartyConversationID,
		"input_AgreedTC":                 r.InputAgreedTC,
	}
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			return errors.New(name + " is required")
		}
	}
	if r.InputAgreedTC != DirectDebitAgreedTCYes && r.InputAgreedTC != DirectDebitAgreedTCNo {
		return errors.New("input_AgreedTC must be 1 or 0")
	}

	frequency := strings.TrimSpace(r.InputFrequency)
	firstPaymentDate := strings.TrimSpace(r.InputFirstPaymentDate)
	startRange := strings.TrimSpace(r.InputStartRangeOfDays)
	endRange := strings.TrimSpace(r.InputEndRangeOfDays)

	if frequency == "" {
		if firstPaymentDate != "" || startRange != "" || endRange != "" {
			return errors.New("input_FirstPaymentDate, input_StartRangeOfDays, and input_EndRangeOfDays must be empty when input_Frequency is empty")
		}
	} else {
		if !validDirectDebitFrequency(frequency) {
			return errors.New("input_Frequency must be one of 01, 02, 03, 04, 05, 06, 07, 08")
		}
		if firstPaymentDate == "" {
			return errors.New("input_FirstPaymentDate is required when input_Frequency is set")
		}

		switch frequency {
		case DirectDebitFrequencyOnceOff, DirectDebitFrequencyDaily, DirectDebitFrequencyWeekly, DirectDebitFrequencyOnDemand:
			if startRange != "" || endRange != "" {
				return errors.New("input_StartRangeOfDays and input_EndRangeOfDays must be empty for once-off, daily, weekly, or on-demand frequency")
			}
		}
	}

	if err := validateMSISDN("input_CustomerMSISDN", r.InputCustomerMSISDN); err != nil {
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
	if err := validateDateYYYYMMDD("input_FirstPaymentDate", r.InputFirstPaymentDate); err != nil {
		return err
	}
	if err := validateDateYYYYMMDD("input_ExpiryDate", r.InputExpiryDate); err != nil {
		return err
	}
	if err := validateDayRange("input_StartRangeOfDays", r.InputStartRangeOfDays); err != nil {
		return err
	}
	if err := validateDayRange("input_EndRangeOfDays", r.InputEndRangeOfDays); err != nil {
		return err
	}
	return nil
}

func validDirectDebitFrequency(frequency string) bool {
	switch frequency {
	case DirectDebitFrequencyOnceOff,
		DirectDebitFrequencyDaily,
		DirectDebitFrequencyWeekly,
		DirectDebitFrequencyMonthly,
		DirectDebitFrequencyQuarterly,
		DirectDebitFrequencyHalfYearly,
		DirectDebitFrequencyYearly,
		DirectDebitFrequencyOnDemand:
		return true
	default:
		return false
	}
}

// DirectDebitCreateResponse represents both sync and initial async responses.
type DirectDebitCreateResponse struct {
	OutputResponseCode             string `json:"output_ResponseCode"`
	OutputResponseDesc             string `json:"output_ResponseDesc"`
	OutputTransactionReference     string `json:"output_TransactionReference,omitempty"`
	OutputMsisdnToken              string `json:"output_MsisdnToken,omitempty"`
	OutputConversationID           string `json:"output_ConversationID"`
	OutputThirdPartyConversationID string `json:"output_ThirdPartyConversationID"`
	OutputMandateID                string `json:"output_MandateID,omitempty"`
}

// DirectDebitCreate creates a Direct Debit mandate.
//
// Endpoint:
//
//	POST /{sandbox|openapi}/ipg/v2/{market}/directDebitCreation/
//
// Bearer value:
//
//	RSA-encrypted SessionID returned by GenerateSessionKey/GenerateSession.
func (c *Client) DirectDebitCreate(ctx context.Context, sessionID string, request DirectDebitCreateRequest) (*DirectDebitCreateResponse, *RawResponse, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, nil, errors.New("sessionID is required; call GenerateSessionKey or GenerateSession first")
	}
	if err := request.Validate(); err != nil {
		return nil, nil, err
	}

	var decoded DirectDebitCreateResponse
	raw, err := c.do(ctx, http.MethodPost, directDebitCreatePath, sessionID, request, &decoded)
	if err != nil {
		return nil, raw, err
	}
	return &decoded, raw, nil
}

// DirectDebitCreateWithSession creates a Direct Debit mandate using a Session
// returned by GenerateSession.
func (c *Client) DirectDebitCreateWithSession(ctx context.Context, session *Session, request DirectDebitCreateRequest) (*DirectDebitCreateResponse, *RawResponse, error) {
	if session == nil || !session.Valid() {
		return nil, nil, errors.New("valid session is required; call GenerateSession first")
	}
	return c.DirectDebitCreate(ctx, session.ID, request)
}
