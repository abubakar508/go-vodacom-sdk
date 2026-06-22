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
