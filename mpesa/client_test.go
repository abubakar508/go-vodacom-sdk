package mpesa

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEncryptBearerValueWithSandboxKey(t *testing.T) {
	token, err := EncryptBearerValue("test-api-key", DefaultSandboxPublicKey)
	if err != nil {
		t.Fatalf("EncryptBearerValue returned error: %v", err)
	}
	if token == "" {
		t.Fatal("expected token")
	}
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		t.Fatalf("token is not base64: %v", err)
	}
	if len(decoded) == 0 {
		t.Fatal("decoded token is empty")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Environment != EnvironmentSandbox {
		t.Fatalf("environment = %q", cfg.Environment)
	}
	if cfg.Market.Context != MarketDRC.Context {
		t.Fatalf("market = %+v", cfg.Market)
	}
	if cfg.PublicKey == "" || cfg.Host == "" || cfg.HTTPClient == nil {
		t.Fatalf("incomplete defaults: %+v", cfg)
	}
}

func TestGenerateSessionKeyEndpointAndDecode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/sandbox/ipg/v2/vodacomDRC/getSession/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			t.Fatalf("missing bearer authorization header")
		}
		if r.Header.Get("Origin") != "127.0.0.1" {
			t.Fatalf("origin = %s", r.Header.Get("Origin"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"output_ResponseCode":"INS-0","output_ResponseDesc":"Request processed successfully","output_SessionID":"abc123"}`))
	}))
	defer server.Close()

	host, port, err := splitServerHostPort(server.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	client, err := NewClient(Config{
		APIKey:      "app-key",
		Environment: EnvironmentSandbox,
		Market:      MarketDRC,
		Origin:      "127.0.0.1",
		Host:        host,
		Port:        port,
		HTTPClient:  server.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}

	res, raw, err := client.GenerateSessionKey(context.Background())
	if err != nil {
		t.Fatalf("GenerateSessionKey error: %v body=%s", err, raw.BodyString())
	}
	if raw.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", raw.StatusCode)
	}
	if res.OutputSessionID != "abc123" {
		t.Fatalf("session id = %s", res.OutputSessionID)
	}
}

func TestC2BSingleStageEndpointAndDecode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/openapi/ipg/v2/vodacomDRC/c2bPayment/singleStage/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"output_ConversationID":"conv","output_ResponseCode":"INS-0","output_ResponseDesc":"Request processed successfully","output_TransactionID":"tx","output_ThirdPartyConversationID":"third"}`))
	}))
	defer server.Close()

	host, port, err := splitServerHostPort(server.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	client, err := NewClient(Config{Environment: EnvironmentOpenAPI, Market: MarketDRC, Host: host, Port: port, HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}

	req := client.NewC2BSingleStageRequest("10", "000000000001", "000000", "T1234C", "third", "Shoes")
	res, _, err := client.C2BSingleStage(context.Background(), "session-id", req)
	if err != nil {
		t.Fatal(err)
	}
	if res.OutputTransactionID != "tx" || res.OutputConversationID != "conv" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestB2CSingleStageEndpointAndDecode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/sandbox/ipg/v2/vodacomDRC/b2cPayment/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"output_ConversationID":"conv","output_ResponseCode":"INS-0","output_ResponseDesc":"Request processed successfully","output_TransactionID":"tx","output_ThirdPartyConversationID":"third"}`))
	}))
	defer server.Close()

	host, port, err := splitServerHostPort(server.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	client, err := NewClient(Config{Environment: EnvironmentSandbox, Market: MarketDRC, Host: host, Port: port, HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}

	req := client.NewB2CSingleStageRequest("10", "000000000001", "000000", "T1234C", "third", "Salary payment")
	res, _, err := client.B2CSingleStage(context.Background(), "session-id", req)
	if err != nil {
		t.Fatal(err)
	}
	if res.OutputTransactionID != "tx" || res.OutputConversationID != "conv" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestB2BSingleStageEndpointAndDecode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/sandbox/ipg/v2/vodacomDRC/b2bPayment/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"output_ConversationID":"conv","output_ResponseCode":"INS-0","output_ResponseDesc":"Request processed successfully","output_TransactionID":"tx","output_ThirdPartyConversationID":"third"}`))
	}))
	defer server.Close()

	host, port, err := splitServerHostPort(server.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	client, err := NewClient(Config{Environment: EnvironmentSandbox, Market: MarketDRC, Host: host, Port: port, HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}

	req := client.NewB2BSingleStageRequest("10", "000000", "000001", "T1234C", "third", "Shoes")
	res, _, err := client.B2BSingleStage(context.Background(), "session-id", req)
	if err != nil {
		t.Fatal(err)
	}
	if res.OutputTransactionID != "tx" || res.OutputConversationID != "conv" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func splitServerHostPort(addr string) (string, int, error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}
	var port int
	_, err = fmt.Sscanf(portStr, "%d", &port)
	return host, port, err
}
