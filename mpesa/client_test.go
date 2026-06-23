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

func TestMarketFromContext(t *testing.T) {
	market, ok := MarketFromContext("VODACOMTZN")
	if !ok {
		t.Fatal("expected market")
	}
	if market.Context != MarketTanzania.Context || market.Country != "TZN" || market.Currency != "TZS" {
		t.Fatalf("unexpected market: %+v", market)
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

func TestReversalEndpointAndDecode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/sandbox/ipg/v2/vodacomDRC/reversal/" {
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

	req := client.NewReversalRequest("25", "000000", "third", "0000000000001")
	res, _, err := client.Reversal(context.Background(), "session-id", req)
	if err != nil {
		t.Fatal(err)
	}
	if res.OutputTransactionID != "tx" || res.OutputConversationID != "conv" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestQueryTransactionStatusEndpointAndDecode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/sandbox/ipg/v2/vodacomDRC/queryTransactionStatus/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("input_QueryReference"); got != "query-ref" {
			t.Fatalf("input_QueryReference = %s", got)
		}
		if got := r.URL.Query().Get("input_Country"); got != "DRC" {
			t.Fatalf("input_Country = %s", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"output_ConversationID":"conv","output_ResponseCode":"INS-0","output_ResponseDesc":"Request processed successfully","output_TransactionID":"tx","output_TransactionStatus":"Completed","output_ThirdPartyConversationID":"third"}`))
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

	req := client.NewQueryTransactionStatusRequest("query-ref", "000000", "third")
	res, _, err := client.QueryTransactionStatus(context.Background(), "session-id", req)
	if err != nil {
		t.Fatal(err)
	}
	if res.OutputTransactionStatus != "Completed" || res.OutputTransactionID != "tx" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestC2BMultiStageEndpointAndDecode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/sandbox/ipg/v2/vodacomDRC/c2bPayment/multiStage/" {
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

	req := client.NewC2BMultiStageRequest("10", "000000000001", "000000", "T1234C", "third", "Shoes")
	if req.InputAPIVersion != "3.1" {
		t.Fatalf("api version = %s", req.InputAPIVersion)
	}
	res, _, err := client.C2BMultiStage(context.Background(), "session-id", req)
	if err != nil {
		t.Fatal(err)
	}
	if res.OutputTransactionID != "tx" || res.OutputConversationID != "conv" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestUpdateTransactionStatusEndpointAndDecode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/sandbox/ipg/v2/vodacomDRC/updateTransactionStatus/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"output_ConversationID":"conv","output_ResponseCode":"INS-GAR-0","output_ResponseDesc":"Request processed successfully","output_TransactionID":"tx","output_ThirdPartyConversationID":"third"}`))
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

	req := client.NewUpdateTransactionStatusRequest(TransactionStatusCommit, "TGS813", "000000", "third", "0000000000001")
	if req.InputAPIVersion != "3.1" {
		t.Fatalf("api version = %s", req.InputAPIVersion)
	}
	res, _, err := client.UpdateTransactionStatus(context.Background(), "session-id", req)
	if err != nil {
		t.Fatal(err)
	}
	if res.OutputResponseCode != "INS-GAR-0" || res.OutputTransactionID != "tx" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestDirectDebitCreateEndpointAndDecode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/sandbox/ipg/v2/vodacomDRC/directDebitCreation/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"output_ResponseCode":"INS-0","output_ResponseDesc":"Request processed successfully","output_TransactionReference":"ref","output_MsisdnToken":"token","output_ConversationID":"conv","output_ThirdPartyConversationID":"third","output_MandateID":"15045"}`))
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

	req := client.NewDirectDebitCreateRequest("000000000001", "000000", "3333", "third", DirectDebitAgreedTCYes)
	req.InputFrequency = DirectDebitFrequencyHalfYearly
	req.InputFirstPaymentDate = "20160324"
	req.InputStartRangeOfDays = "01"
	req.InputEndRangeOfDays = "22"
	req.InputExpiryDate = "20161126"

	res, _, err := client.DirectDebitCreate(context.Background(), "session-id", req)
	if err != nil {
		t.Fatal(err)
	}
	if res.OutputMandateID != "15045" || res.OutputTransactionReference != "ref" || res.OutputMsisdnToken != "token" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestDirectDebitPaymentEndpointAndDecode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/sandbox/ipg/v2/vodacomDRC/directDebitPayment/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"output_ResponseCode":"INS-0","output_ResponseDesc":"Request processed successfully","output_TransactionReference":"ref","output_MsisdnToken":"token","output_ConversationID":"conv","output_ThirdPartyConversationID":"third"}`))
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

	req := client.NewDirectDebitPaymentWithMSISDN("10", "000000000001", "000000", "ref", "third", "15045")
	res, _, err := client.DirectDebitPayment(context.Background(), "session-id", req)
	if err != nil {
		t.Fatal(err)
	}
	if res.OutputTransactionReference != "ref" || res.OutputMsisdnToken != "token" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestQueryBeneficiaryNameEndpointAndDecode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/sandbox/ipg/v2/vodafoneGHA/queryBeneficiaryName/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("input_KycQueryType"); got != "Name" {
			t.Fatalf("input_KycQueryType = %s", got)
		}
		if got := r.URL.Query().Get("input_Country"); got != "GHA" {
			t.Fatalf("input_Country = %s", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"output_ConversationID":"conv","output_ResponseCode":"INS-0","output_ResponseDesc":"Request processed successfully","output_CustomerFirstName":"Jiazhen","output_CustomerLastName":"Wuu","output_ThirdPartyConversationID":"third"}`))
	}))
	defer server.Close()

	host, port, err := splitServerHostPort(server.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	client, err := NewClient(Config{Environment: EnvironmentSandbox, Market: MarketGhana, Host: host, Port: port, HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}

	req := client.NewQueryBeneficiaryNameRequest("254707161122", "000000", "third")
	res, _, err := client.QueryBeneficiaryName(context.Background(), "session-id", req)
	if err != nil {
		t.Fatal(err)
	}
	if res.OutputCustomerFirstName != "Jiazhen" || res.OutputCustomerLastName != "Wuu" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestQueryDirectDebitEndpointAndDecode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/sandbox/ipg/v2/vodacomDRC/queryDirectDebit/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("input_QueryBalanceAmount"); got != "True" {
			t.Fatalf("input_QueryBalanceAmount = %s", got)
		}
		if got := r.URL.Query().Get("input_Currency"); got != "USD" {
			t.Fatalf("input_Currency = %s", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"output_ResponseCode":"INS-0","output_ResponseDesc":"Request processed successfully","output_TransactionReference":"ref","output_ConversationID":"conv","output_ThirdPartyConversationID":"third","output_SufficientBalance":"True","output_MsisdnToken":"token","output_MandateID":"15132","output_MandateStatus":"Active","output_AccountStatus":"Active","output_FirstPaymentDate":"20231012","output_Frequency":"02","output_PaymentDayFrom":"01","output_PaymentDayTo":"25","output_ExpiryDate":"20230410"}`))
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

	req := client.NewQueryDirectDebitWithMSISDN(QueryBalanceAmountTrue, "100", "255744553111", "112244", "third", "Test123", "15132")
	res, _, err := client.QueryDirectDebit(context.Background(), "session-id", req)
	if err != nil {
		t.Fatal(err)
	}
	if res.OutputMandateStatus != "Active" || res.OutputSufficientBalance != "True" || res.OutputMandateID != "15132" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestCancelDirectDebitEndpointAndDecode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/sandbox/ipg/v2/vodacomDRC/directDebitCancel/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"output_ResponseCode":"INS-0","output_ResponseDesc":"Request processed successfully","output_TransactionReference":"ref","output_MsisdnToken":"token","output_ConversationID":"conv","output_ThirdPartyConversationID":"third"}`))
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

	req := client.NewCancelDirectDebitWithMSISDN("000000000001", "000000", "ref", "third", "15045")
	res, _, err := client.CancelDirectDebit(context.Background(), "session-id", req)
	if err != nil {
		t.Fatal(err)
	}
	if res.OutputTransactionReference != "ref" || res.OutputMsisdnToken != "token" {
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
