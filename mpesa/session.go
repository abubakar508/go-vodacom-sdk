package mpesa

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"
)

const getSessionPath = "getSession/"

// DefaultSessionActivationDelay is the delay recommended by the official
// combined examples before a fresh SessionID is used on transaction APIs.
const DefaultSessionActivationDelay = 30 * time.Second

// Session represents a generated M-Pesa OpenAPI session token.
type Session struct {
	ID          string
	Market      Market
	Environment Environment
	CreatedAt   time.Time
}

// Valid reports whether the session has a non-empty ID.
func (s Session) Valid() bool {
	return strings.TrimSpace(s.ID) != ""
}

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

// GenerateSession exchanges the application API key for a Session helper object.
// Transaction methods still accept a string SessionID, but this method makes it
// clear that the returned token is market/environment scoped and should be used
// for subsequent C2B/B2C/B2B calls.
func (c *Client) GenerateSession(ctx context.Context) (*Session, *RawResponse, error) {
	res, raw, err := c.GenerateSessionKey(ctx)
	if err != nil {
		return nil, raw, err
	}
	return &Session{
		ID:          res.OutputSessionID,
		Market:      c.cfg.Market,
		Environment: c.cfg.Environment,
		CreatedAt:   time.Now(),
	}, raw, nil
}

// GenerateSessionAndWait generates a Session and waits for the supplied delay
// before returning it. If delay is zero, DefaultSessionActivationDelay is used.
// This mirrors the official examples that warn a new SessionID can take up to
// 30 seconds to become active for transaction APIs.
func (c *Client) GenerateSessionAndWait(ctx context.Context, delay time.Duration) (*Session, *RawResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	session, raw, err := c.GenerateSession(ctx)
	if err != nil {
		return nil, raw, err
	}
	if delay == 0 {
		delay = DefaultSessionActivationDelay
	}
	if delay < 0 {
		return session, raw, nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return session, raw, ctx.Err()
	case <-timer.C:
		return session, raw, nil
	}
}
