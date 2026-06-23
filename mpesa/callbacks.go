package mpesa

// AsyncTransactionCallbackRequest is the common OpenAPI -> ThirdParty callback
// payload shape used by async single-stage transaction APIs.
type AsyncTransactionCallbackRequest struct {
	InputOriginalConversationID   string `json:"input_OriginalConversationID"`
	InputTransactionID            string `json:"input_TransactionID"`
	InputResultCode               string `json:"input_ResultCode"`
	InputResultDesc               string `json:"input_ResultDesc"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
}

// AsyncTransactionCallbackResponse is the common acknowledgement your listener
// should return to OpenAPI.
type AsyncTransactionCallbackResponse struct {
	OutputOriginalConversationID   string `json:"output_OriginalConversationID"`
	OutputResponseCode             string `json:"output_ResponseCode"`
	OutputResponseDesc             string `json:"output_ResponseDesc"`
	OutputThirdPartyConversationID string `json:"output_ThirdPartyConversationID"`
}

// AcceptAsyncTransactionCallback builds the expected successful acknowledgement payload.
func AcceptAsyncTransactionCallback(req AsyncTransactionCallbackRequest) AsyncTransactionCallbackResponse {
	return AsyncTransactionCallbackResponse{
		OutputOriginalConversationID:   req.InputOriginalConversationID,
		OutputResponseCode:             "0",
		OutputResponseDesc:             "Successfully Accepted Result",
		OutputThirdPartyConversationID: req.InputThirdPartyConversationID,
	}
}

// C2BAsyncCallbackRequest is the OpenAPI -> ThirdParty callback payload for an async C2B flow.
type C2BAsyncCallbackRequest = AsyncTransactionCallbackRequest

// C2BAsyncCallbackResponse is the acknowledgement your listener should return to OpenAPI for C2B.
type C2BAsyncCallbackResponse = AsyncTransactionCallbackResponse

// AcceptC2BAsyncCallback builds the expected successful C2B acknowledgement payload.
func AcceptC2BAsyncCallback(req C2BAsyncCallbackRequest) C2BAsyncCallbackResponse {
	return AcceptAsyncTransactionCallback(req)
}

// B2CAsyncCallbackRequest is the OpenAPI -> ThirdParty callback payload for an async B2C flow.
type B2CAsyncCallbackRequest = AsyncTransactionCallbackRequest

// B2CAsyncCallbackResponse is the acknowledgement your listener should return to OpenAPI for B2C.
type B2CAsyncCallbackResponse = AsyncTransactionCallbackResponse

// AcceptB2CAsyncCallback builds the expected successful B2C acknowledgement payload.
func AcceptB2CAsyncCallback(req B2CAsyncCallbackRequest) B2CAsyncCallbackResponse {
	return AcceptAsyncTransactionCallback(req)
}

// B2BAsyncCallbackRequest is the OpenAPI -> ThirdParty callback payload for an async B2B flow.
type B2BAsyncCallbackRequest = AsyncTransactionCallbackRequest

// B2BAsyncCallbackResponse is the acknowledgement your listener should return to OpenAPI for B2B.
type B2BAsyncCallbackResponse = AsyncTransactionCallbackResponse

// AcceptB2BAsyncCallback builds the expected successful B2B acknowledgement payload.
func AcceptB2BAsyncCallback(req B2BAsyncCallbackRequest) B2BAsyncCallbackResponse {
	return AcceptAsyncTransactionCallback(req)
}

// ReversalAsyncCallbackRequest is the OpenAPI -> ThirdParty callback payload for an async Reversal flow.
type ReversalAsyncCallbackRequest = AsyncTransactionCallbackRequest

// ReversalAsyncCallbackResponse is the acknowledgement your listener should return to OpenAPI for Reversal.
type ReversalAsyncCallbackResponse = AsyncTransactionCallbackResponse

// AcceptReversalAsyncCallback builds the expected successful Reversal acknowledgement payload.
func AcceptReversalAsyncCallback(req ReversalAsyncCallbackRequest) ReversalAsyncCallbackResponse {
	return AcceptAsyncTransactionCallback(req)
}

// C2BMultiStageAsyncCallbackRequest is the OpenAPI -> ThirdParty callback payload for an async C2B Multi Stage flow.
type C2BMultiStageAsyncCallbackRequest = AsyncTransactionCallbackRequest

// C2BMultiStageAsyncCallbackResponse is the acknowledgement your listener should return to OpenAPI for C2B Multi Stage.
type C2BMultiStageAsyncCallbackResponse = AsyncTransactionCallbackResponse

// AcceptC2BMultiStageAsyncCallback builds the expected successful C2B Multi Stage acknowledgement payload.
func AcceptC2BMultiStageAsyncCallback(req C2BMultiStageAsyncCallbackRequest) C2BMultiStageAsyncCallbackResponse {
	return AcceptAsyncTransactionCallback(req)
}

// UpdateTransactionStatusAsyncCallbackRequest is the OpenAPI -> ThirdParty callback payload for an async Update Transaction Status flow.
type UpdateTransactionStatusAsyncCallbackRequest = AsyncTransactionCallbackRequest

// UpdateTransactionStatusAsyncCallbackResponse is the acknowledgement your listener should return to OpenAPI for Update Transaction Status.
type UpdateTransactionStatusAsyncCallbackResponse = AsyncTransactionCallbackResponse

// AcceptUpdateTransactionStatusAsyncCallback builds the documented successful
// Update Transaction Status acknowledgement payload. This API documents
// INS-GAR-0, not plain 0, as the acknowledgement response code.
func AcceptUpdateTransactionStatusAsyncCallback(req UpdateTransactionStatusAsyncCallbackRequest) UpdateTransactionStatusAsyncCallbackResponse {
	return UpdateTransactionStatusAsyncCallbackResponse{
		OutputOriginalConversationID:   req.InputOriginalConversationID,
		OutputResponseCode:             "INS-GAR-0",
		OutputResponseDesc:             "Successfully Accepted Result",
		OutputThirdPartyConversationID: req.InputThirdPartyConversationID,
	}
}

// DirectDebitCreateAsyncCallbackRequest is the OpenAPI -> ThirdParty callback
// payload for an async Direct Debit Create flow.
type DirectDebitCreateAsyncCallbackRequest struct {
	InputOriginalConversationID   string `json:"input_OriginalConversationID"`
	InputTransactionReference     string `json:"input_TransactionReference"`
	InputMsisdnToken              string `json:"input_MsisdnToken"`
	InputResultCode               string `json:"input_ResultCode"`
	InputResultDesc               string `json:"input_ResultDesc"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
	InputMandateID                string `json:"input_MandateID"`
}

// DirectDebitCreateAsyncCallbackResponse is the acknowledgement your listener
// should return to OpenAPI for Direct Debit Create.
type DirectDebitCreateAsyncCallbackResponse = AsyncTransactionCallbackResponse

// AcceptDirectDebitCreateAsyncCallback builds the expected successful Direct
// Debit Create acknowledgement payload.
func AcceptDirectDebitCreateAsyncCallback(req DirectDebitCreateAsyncCallbackRequest) DirectDebitCreateAsyncCallbackResponse {
	return DirectDebitCreateAsyncCallbackResponse{
		OutputOriginalConversationID:   req.InputOriginalConversationID,
		OutputResponseCode:             "0",
		OutputResponseDesc:             "Successfully Accepted Result",
		OutputThirdPartyConversationID: req.InputThirdPartyConversationID,
	}
}

// DirectDebitPaymentAsyncCallbackRequest is the OpenAPI -> ThirdParty callback
// payload for an async Direct Debit Payment flow.
type DirectDebitPaymentAsyncCallbackRequest struct {
	InputOriginalConversationID   string `json:"input_OriginalConversationID"`
	InputTransactionID            string `json:"input_TransactionID"`
	InputMsisdnToken              string `json:"input_MsisdnToken"`
	InputResultCode               string `json:"input_ResultCode"`
	InputResultDesc               string `json:"input_ResultDesc"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
}

// DirectDebitPaymentAsyncCallbackResponse is the acknowledgement your listener
// should return to OpenAPI for Direct Debit Payment.
type DirectDebitPaymentAsyncCallbackResponse = AsyncTransactionCallbackResponse

// AcceptDirectDebitPaymentAsyncCallback builds the expected successful Direct
// Debit Payment acknowledgement payload.
func AcceptDirectDebitPaymentAsyncCallback(req DirectDebitPaymentAsyncCallbackRequest) DirectDebitPaymentAsyncCallbackResponse {
	return DirectDebitPaymentAsyncCallbackResponse{
		OutputOriginalConversationID:   req.InputOriginalConversationID,
		OutputResponseCode:             "0",
		OutputResponseDesc:             "Successfully Accepted Result",
		OutputThirdPartyConversationID: req.InputThirdPartyConversationID,
	}
}

// QueryBeneficiaryNameAsyncCallbackRequest is the OpenAPI -> ThirdParty callback
// payload for an async Query Beneficiary Name flow.
type QueryBeneficiaryNameAsyncCallbackRequest struct {
	InputOriginalConversationID   string `json:"input_OriginalConversationID"`
	InputCustomerFirstName        string `json:"input_CustomerFirstName"`
	InputCustomerLastName         string `json:"input_CustomerLastName"`
	InputResultCode               string `json:"input_ResultCode"`
	InputResultDesc               string `json:"input_ResultDesc"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
}

// QueryBeneficiaryNameAsyncCallbackResponse is the acknowledgement your listener
// should return to OpenAPI for Query Beneficiary Name.
type QueryBeneficiaryNameAsyncCallbackResponse = AsyncTransactionCallbackResponse

// AcceptQueryBeneficiaryNameAsyncCallback builds the expected successful Query
// Beneficiary Name acknowledgement payload.
func AcceptQueryBeneficiaryNameAsyncCallback(req QueryBeneficiaryNameAsyncCallbackRequest) QueryBeneficiaryNameAsyncCallbackResponse {
	return QueryBeneficiaryNameAsyncCallbackResponse{
		OutputOriginalConversationID:   req.InputOriginalConversationID,
		OutputResponseCode:             "0",
		OutputResponseDesc:             "Successfully Accepted Result",
		OutputThirdPartyConversationID: req.InputThirdPartyConversationID,
	}
}

// QueryDirectDebitAsyncCallbackRequest is the OpenAPI -> ThirdParty callback
// payload for an async Query Direct Debit flow.
type QueryDirectDebitAsyncCallbackRequest struct {
	InputResponseCode             string `json:"input_ResponseCode"`
	InputResponseDesc             string `json:"input_ResponseDesc"`
	InputTransactionReference     string `json:"input_TransactionReference"`
	InputOriginalConversationID   string `json:"input_OriginalConversationID"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
	InputSufficientBalance        string `json:"input_SufficientBalance"`
	InputMsisdnToken              string `json:"input_MsisdnToken"`
	InputMandateID                string `json:"input_MandateID"`
	InputMandateStatus            string `json:"input_MandateStatus"`
	InputAccountStatus            string `json:"input_AccountStatus"`
	InputFirstPaymentDate         string `json:"input_FirstPaymentDate"`
	InputFrequency                string `json:"input_Frequency"`
	InputPaymentDayFrom           string `json:"input_PaymentDayFrom"`
	InputPaymentDayTo             string `json:"input_PaymentDayTo"`
	InputExpiryDate               string `json:"input_ExpiryDate"`
}

// QueryDirectDebitAsyncCallbackResponse is the acknowledgement your listener
// should return to OpenAPI for Query Direct Debit.
type QueryDirectDebitAsyncCallbackResponse = AsyncTransactionCallbackResponse

// AcceptQueryDirectDebitAsyncCallback builds the expected successful Query
// Direct Debit acknowledgement payload. This API documents INS-0 as the
// acknowledgement response code.
func AcceptQueryDirectDebitAsyncCallback(req QueryDirectDebitAsyncCallbackRequest) QueryDirectDebitAsyncCallbackResponse {
	return QueryDirectDebitAsyncCallbackResponse{
		OutputOriginalConversationID:   req.InputOriginalConversationID,
		OutputResponseCode:             "INS-0",
		OutputResponseDesc:             "Successfully Accepted Result",
		OutputThirdPartyConversationID: req.InputThirdPartyConversationID,
	}
}

// CancelDirectDebitAsyncCallbackRequest is the OpenAPI -> ThirdParty callback
// payload for an async Cancel Direct Debit flow.
type CancelDirectDebitAsyncCallbackRequest struct {
	InputOriginalConversationID   string `json:"input_OriginalConversationID"`
	InputTransactionReference     string `json:"input_TransactionReference"`
	InputMsisdnToken              string `json:"input_MsisdnToken"`
	InputResultCode               string `json:"input_ResultCode"`
	InputResultDesc               string `json:"input_ResultDesc"`
	InputThirdPartyConversationID string `json:"input_ThirdPartyConversationID"`
}

// CancelDirectDebitAsyncCallbackResponse is the acknowledgement your listener
// should return to OpenAPI for Cancel Direct Debit.
type CancelDirectDebitAsyncCallbackResponse = AsyncTransactionCallbackResponse

// AcceptCancelDirectDebitAsyncCallback builds the expected successful Cancel
// Direct Debit acknowledgement payload.
func AcceptCancelDirectDebitAsyncCallback(req CancelDirectDebitAsyncCallbackRequest) CancelDirectDebitAsyncCallbackResponse {
	return CancelDirectDebitAsyncCallbackResponse{
		OutputOriginalConversationID:   req.InputOriginalConversationID,
		OutputResponseCode:             "0",
		OutputResponseDesc:             "Successfully Accepted Result",
		OutputThirdPartyConversationID: req.InputThirdPartyConversationID,
	}
}
