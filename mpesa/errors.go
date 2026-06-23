package mpesa

import (
	"errors"
	"fmt"
)

const (
	CodeSuccess               = "INS-0"
	CodeAccepted              = "0"
	CodeInternalError         = "INS-1"
	CodeTransactionFailed     = "INS-6"
	CodeRequestTimeout        = "INS-9"
	CodeDuplicateTransaction  = "INS-10"
	CodeInvalidShortcode      = "INS-13"
	CodeInvalidAmount         = "INS-15"
	CodeInvalidTransactionRef = "INS-17"
	CodeMissingParameters     = "INS-20"
	CodeValidationFailed      = "INS-21"
	CodeInvalidCurrency       = "INS-26"
	CodeInvalidConversationID = "INS-28"
	CodeInvalidDescription    = "INS-30"
	CodeInvalidKYCQueryType   = "INS-32"
	CodeInvalidAgreedTC       = "INS-36"
	CodeInvalidMandateID      = "INS-51"
	CodeNoActiveMandate       = "INS-58"
	CodeInsufficientBalance   = "INS-2006"
	CodeAPIOutsideUsageTime   = "INS-996"
	CodeAPINotEnabled         = "INS-997"
	CodeInvalidMarket         = "INS-998"
	CodeMSISDNInvalid         = "INS-2051"
)

// APIError is returned for non-2xx HTTP responses and known non-success M-Pesa
// response codes.
type APIError struct {
	StatusCode   int
	ResponseCode string
	Description  string
	Body         string
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.ResponseCode != "" || e.Description != "" {
		return fmt.Sprintf("mpesa api error: http=%d code=%s desc=%s", e.StatusCode, e.ResponseCode, e.Description)
	}
	return fmt.Sprintf("mpesa api error: http=%d body=%s", e.StatusCode, e.Body)
}

// IsCode reports whether err is an *APIError with the supplied M-Pesa response code.
func IsCode(err error, code string) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.ResponseCode == code
}

func IsTimeout(err error) bool             { return IsCode(err, CodeRequestTimeout) }
func IsDuplicate(err error) bool           { return IsCode(err, CodeDuplicateTransaction) }
func IsInvalidMarket(err error) bool       { return IsCode(err, CodeInvalidMarket) }
func IsInsufficientBalance(err error) bool { return IsCode(err, CodeInsufficientBalance) }
func IsValidationFailed(err error) bool    { return IsCode(err, CodeValidationFailed) }
func IsAPINotEnabled(err error) bool       { return IsCode(err, CodeAPINotEnabled) }
