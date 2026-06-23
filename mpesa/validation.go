package mpesa

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reAmount                   = regexp.MustCompile(`^\d*\.?\d+$`)
	reMSISDN                   = regexp.MustCompile(`^[0-9]{12,14}$`)
	reShortCode                = regexp.MustCompile(`^[0-9A-Za-z]{4,12}$`)
	reTransactionReference     = regexp.MustCompile(`^[0-9a-zA-Z \w+]{1,20}$`)
	reThirdPartyConversationID = regexp.MustCompile(`^[0-9a-zA-Z \w+]{1,40}$`)
	reThirdPartyReference      = regexp.MustCompile(`^[0-9a-zA-Z]{1,32}$`)
	reCurrency                 = regexp.MustCompile(`^[a-zA-Z]{1,3}$`)
	reYYYYMMDD                 = regexp.MustCompile(`^[0-9]{8}$`)
	reDayRange                 = regexp.MustCompile(`^[0-9]{2}$`)
	reTransactionID            = regexp.MustCompile(`^[0-9a-zA-Z]{1,13}$`)
	reMandateID                = regexp.MustCompile(`^[0-9]{1,12}$`)
	reMsisdnToken              = regexp.MustCompile(`^[0-9a-zA-Z \w=]{1,64}$`)
	reVoucherCode              = regexp.MustCompile(`^[0-9A-Za-z]{4,12}$`)
	reQueryReference           = regexp.MustCompile(`^[0-9a-zA-Z \w+.-]{1,64}$`)
)

func requireFields(fields map[string]string) error {
	for name, value := range fields {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required", name)
		}
	}
	return nil
}

func validatePattern(name, value string, re *regexp.Regexp) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	if !re.MatchString(value) {
		return fmt.Errorf("%s is invalid", name)
	}
	return nil
}

func validateAmount(name, value string) error    { return validatePattern(name, value, reAmount) }
func validateMSISDN(name, value string) error    { return validatePattern(name, value, reMSISDN) }
func validateShortCode(name, value string) error { return validatePattern(name, value, reShortCode) }
func validateTransactionReference(name, value string) error {
	return validatePattern(name, value, reTransactionReference)
}
func validateThirdPartyConversationID(name, value string) error {
	return validatePattern(name, value, reThirdPartyConversationID)
}
func validateThirdPartyReference(name, value string) error {
	return validatePattern(name, value, reThirdPartyReference)
}
func validateCurrency(name, value string) error     { return validatePattern(name, value, reCurrency) }
func validateDateYYYYMMDD(name, value string) error { return validatePattern(name, value, reYYYYMMDD) }
func validateDayRange(name, value string) error     { return validatePattern(name, value, reDayRange) }
func validateTransactionID(name, value string) error {
	return validatePattern(name, value, reTransactionID)
}
func validateMandateID(name, value string) error { return validatePattern(name, value, reMandateID) }
func validateMsisdnToken(name, value string) error {
	return validatePattern(name, value, reMsisdnToken)
}
func validateVoucherCode(name, value string) error {
	return validatePattern(name, value, reVoucherCode)
}
func validateQueryReference(name, value string) error {
	return validatePattern(name, value, reQueryReference)
}

func validateCommonTransactionFields(amount, country, currency, serviceProviderCode, transactionReference, thirdPartyConversationID string) error {
	if err := validateAmount("input_Amount", amount); err != nil {
		return err
	}
	if err := validateCurrency("input_Currency", currency); err != nil {
		return err
	}
	if err := validateShortCode("input_ServiceProviderCode", serviceProviderCode); err != nil {
		return err
	}
	if err := validateTransactionReference("input_TransactionReference", transactionReference); err != nil {
		return err
	}
	if err := validateThirdPartyConversationID("input_ThirdPartyConversationID", thirdPartyConversationID); err != nil {
		return err
	}
	return nil
}
