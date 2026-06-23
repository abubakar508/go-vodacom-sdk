package mpesa

import (
	"encoding/json"
	"net/http"
)

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

// JSONMap decodes the raw response body into a generic map. Use this when a
// market returns fields that are not yet represented by a typed SDK response.
func (r *RawResponse) JSONMap() (map[string]any, error) {
	if r == nil {
		return nil, nil
	}
	var out map[string]any
	if err := json.Unmarshal(r.Body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

type responseEnvelope struct {
	OutputResponseCode string `json:"output_ResponseCode"`
	OutputResponseDesc string `json:"output_ResponseDesc"`
	InputResultCode    string `json:"input_ResultCode"`
	InputResultDesc    string `json:"input_ResultDesc"`
}
