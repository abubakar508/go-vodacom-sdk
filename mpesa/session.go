package mpesa

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const getSessionPath = "getSession/"

// SessionKeyResponse is the synchronous response from GET /getSession/.
type SessionKeyResponse struct {
	OutputResponseCode string `json:"output_ResponseCode"`
	OutputResponseDesc string `json:"output_ResponseDesc"`
	OutputSessionID    string `json:"output_SessionID"`
}

// GenerateSessionKey exchanges the application API key for a SessionID.
//
// Endpoint:
//   GET /{sandbox|openapi}/ipg/v2/{market}/getSession/
//
// Bearer value:
//   RSA-encrypted application API key.
func (c *Client) GenerateSessionKey(ctx context.Context) (*SessionKeyResponse, *RawResponse, error) {
	if strings.TrimSpace(c.cfg.APIKey) == "" {
		return nil, nil, errors.New("api key is required to generate a session key")
	}

	var decoded SessionKeyResponse
	raw, err := c.do(ctx, http.MethodGet, getSessionPath, c.cfg.APIKey, nil, &decoded)
	if err != nil {
		return nil, raw, err
	}
	return &decoded, raw, nil
}
