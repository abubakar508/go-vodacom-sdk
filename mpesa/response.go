package mpesa

import "net/http"

// RawResponse contains raw HTTP response data returned alongside decoded payloads.
type RawResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func (r *RawResponse) BodyString() string {
	if r == nil {
		return ""
	}
	return string(r.Body)
}

type responseEnvelope struct {
	OutputResponseCode string `json:"output_ResponseCode"`
	OutputResponseDesc string `json:"output_ResponseDesc"`
	InputResultCode    string `json:"input_ResultCode"`
	InputResultDesc    string `json:"input_ResultDesc"`
}
