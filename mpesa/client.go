package mpesa

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	cryptoutil "github.com/abubakar508/go-vodacom-sdk/internal/crypto"
)

// Client is a Vodacom/Vodafone M-Pesa OpenAPI client.
type Client struct {
	cfg Config
}

// NewClient creates a client from Config. Missing fields are filled using
// DefaultConfig. Vodacom DRC sandbox is the default market/environment.
func NewClient(cfg Config) (*Client, error) {
	cfg = cfg.normalize()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &Client{cfg: cfg}, nil
}

// Config returns a copy of the resolved client configuration.
func (c *Client) Config() Config {
	return c.cfg
}

func (c *Client) country() string { return c.cfg.Market.Country }

func (c *Client) currency() string {
	if c.cfg.Currency != "" {
		return c.cfg.Currency
	}
	return c.cfg.Market.Currency
}

// EncryptBearerValue encrypts an API key or SessionID exactly as required by
// M-Pesa OpenAPI Authorization Bearer headers.
func EncryptBearerValue(value, publicKey string) (string, error) {
	return cryptoutil.EncryptPKCS1v15ToBase64(value, publicKey)
}

func (c *Client) endpoint(apiPath string) string {
	return c.endpointWithQuery(apiPath, nil)
}

func (c *Client) endpointWithQuery(apiPath string, query url.Values) string {
	host := c.cfg.Host
	if c.cfg.Port != 0 && c.cfg.Port != 443 {
		host = net.JoinHostPort(c.cfg.Host, fmt.Sprintf("%d", c.cfg.Port))
	}

	path := fmt.Sprintf("/%s/ipg/v2/%s/%s", c.cfg.Environment.BasePath(), c.cfg.Market.Context, strings.TrimLeft(apiPath, "/"))
	u := &url.URL{Scheme: "https", Host: host, Path: path}
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}
	return u.String()
}

func (c *Client) do(ctx context.Context, method, apiPath, bearerValue string, payload any, out any) (*RawResponse, error) {
	return c.doWithQuery(ctx, method, apiPath, nil, bearerValue, payload, out)
}

func (c *Client) doQuery(ctx context.Context, method, apiPath string, query url.Values, bearerValue string, out any) (*RawResponse, error) {
	return c.doWithQuery(ctx, method, apiPath, query, bearerValue, nil, out)
}

func (c *Client) doWithQuery(ctx context.Context, method, apiPath string, query url.Values, bearerValue string, payload any, out any) (*RawResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if strings.TrimSpace(bearerValue) == "" {
		return nil, errors.New("bearer value cannot be empty")
	}

	token, err := EncryptBearerValue(bearerValue, c.cfg.PublicKey)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.endpointWithQuery(apiPath, query), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", c.cfg.Origin)
	req.Host = c.cfg.Host

	resp, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, readErr := io.ReadAll(resp.Body)
	raw := &RawResponse{StatusCode: resp.StatusCode, Header: resp.Header.Clone(), Body: respBody}
	if readErr != nil {
		return raw, fmt.Errorf("read response body: %w", readErr)
	}

	var envelope responseEnvelope
	_ = json.Unmarshal(respBody, &envelope)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return raw, &APIError{
			StatusCode:   resp.StatusCode,
			ResponseCode: firstNonEmpty(envelope.OutputResponseCode, envelope.InputResultCode),
			Description:  firstNonEmpty(envelope.OutputResponseDesc, envelope.InputResultDesc),
			Body:         string(respBody),
		}
	}

	if out != nil && strings.TrimSpace(string(respBody)) != "" {
		if err := json.Unmarshal(respBody, out); err != nil {
			return raw, fmt.Errorf("decode response body: %w", err)
		}
	}

	code := firstNonEmpty(envelope.OutputResponseCode, envelope.InputResultCode)
	desc := firstNonEmpty(envelope.OutputResponseDesc, envelope.InputResultDesc)
	if code != "" && !isSuccessResponseCode(code) {
		return raw, &APIError{StatusCode: resp.StatusCode, ResponseCode: code, Description: desc, Body: string(respBody)}
	}

	return raw, nil
}

func isSuccessResponseCode(code string) bool {
	code = strings.TrimSpace(code)
	return code == "0" || code == "INS-0" || strings.HasSuffix(code, "-0")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
