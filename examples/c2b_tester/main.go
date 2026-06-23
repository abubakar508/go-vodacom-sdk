package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/abubakar508/go-vodacom-sdk/mpesa"
)

const listenAddr = ":8088"

var page = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>M-Pesa C2B Tester</title>
  <script src="https://unpkg.com/htmx.org@1.9.12"></script>
  <style>
    :root { --red:#e60000; --dark:#211; --muted:#756b6b; --soft:#fff2f2; --card:#fff; }
    * { box-sizing: border-box; }
    body { margin:0; font-family:Inter,system-ui,-apple-system,Segoe UI,sans-serif; color:var(--dark); background:radial-gradient(circle at top left,#ffd6d6,transparent 32rem),linear-gradient(135deg,#fff,#fff6f6 55%,#f2eeee); min-height:100vh; }
    .wrap { max-width:1180px; margin:auto; padding:36px 20px; }
    .hero { display:grid; grid-template-columns:1fr .9fr; gap:24px; align-items:start; }
    h1 { font-size:clamp(34px,6vw,68px); line-height:.95; letter-spacing:-.06em; margin:16px 0; }
    h1 span { color:var(--red); }
    p { color:var(--muted); line-height:1.6; }
    .pill { display:inline-flex; background:#fff; color:#9b0000; border:1px solid #ffd0d0; border-radius:999px; padding:9px 13px; font-weight:800; }
    .card { background:rgba(255,255,255,.86); border:1px solid rgba(230,0,0,.12); box-shadow:0 24px 70px rgba(120,0,0,.14), inset 0 1px rgba(255,255,255,.8); border-radius:28px; padding:22px; }
    label { display:block; font-size:13px; font-weight:800; margin:14px 0 6px; }
    input, select, textarea { width:100%; border:1px solid #ead4d4; border-radius:14px; padding:12px 13px; font:inherit; background:#fff; }
    textarea { min-height:84px; font-family:ui-monospace,SFMono-Regular,Menlo,monospace; }
    .grid { display:grid; grid-template-columns:repeat(2,1fr); gap:12px; }
    button { margin-top:18px; border:0; border-radius:16px; background:var(--red); color:#fff; font-weight:900; padding:14px 18px; width:100%; cursor:pointer; box-shadow:0 16px 40px rgba(230,0,0,.26); }
    button:disabled { opacity:.65; }
    pre { white-space:pre-wrap; word-break:break-word; background:#241b1b; color:#ffecec; border-radius:20px; padding:18px; overflow:auto; max-height:640px; }
    .warn { background:#fff8e8; border:1px solid #ffe1a6; padding:12px 14px; border-radius:18px; color:#6d4a00; }
    .small { font-size:13px; }
    .htmx-indicator { display:none; }
    .htmx-request .htmx-indicator { display:inline; }
    @media(max-width:900px){ .hero,.grid{grid-template-columns:1fr;} }
  </style>
</head>
<body>
  <div class="wrap">
    <div class="hero">
      <section>
        <div class="pill">Local C2B tester · server-side Origin header</div>
        <h1>Test Vodacom M-Pesa <span>C2B</span> safely.</h1>
        <p>This page is served by a local Go backend. Your API key and SessionID stay on your machine and are not stored. Do not put real production credentials into a public/static frontend.</p>
        <div class="warn small"><b>Origin matters:</b> M-Pesa validates the exact <code>Origin</code> header against the value configured in your portal application. If you see <code>Origin not allowed</code>, replace the Origin field with the exact portal value, for example <code>127.0.0.1</code>, your server IP, or your configured domain. Do not assume <code>*</code> is allowed.</div>
      </section>
      <section class="card">
        <form hx-post="/api/c2b" hx-target="#result" hx-indicator="#loading">
          <h2>Credentials</h2>
          <label>Application API Key</label>
          <input name="api_key" type="password" required placeholder="Your M-Pesa application API key" />
          <label>Platform Public Key override</label>
          <textarea name="public_key" placeholder="Optional. Leave empty to use SDK default for selected environment."></textarea>
          <div class="grid">
            <div><label>Environment</label><select name="environment"><option value="sandbox">sandbox</option><option value="openapi">openapi/live</option></select></div>
            <div><label>Market</label><select name="market"><option value="vodafoneGHA">Ghana / GHS</option><option value="vodacomDRC">DRC / USD</option><option value="vodacomTZN">Tanzania / TZS</option><option value="vodacomLES">Lesotho / LSL</option><option value="vodacomMOZ">Mozambique / MZN</option></select></div>
          </div>
          <div class="grid">
            <div><label>Origin header to M-Pesa</label><input name="origin" value="127.0.0.1" placeholder="Exact origin configured in portal" /></div>
            <div><label>Session wait seconds</label><input name="wait_seconds" value="30" /></div>
          </div>
          <h2>C2B request</h2>
          <div class="grid">
            <div><label>Amount</label><input name="amount" value="10" required /></div>
            <div><label>Customer MSISDN</label><input name="customer_msisdn" value="000000000001" required /></div>
            <div><label>Service Provider Code</label><input name="service_provider_code" value="000000" required /></div>
            <div><label>Transaction Reference</label><input name="transaction_reference" value="T1234C" required /></div>
          </div>
          <label>Third Party Conversation ID</label>
          <input name="third_party_conversation_id" placeholder="Leave empty to auto-generate" />
          <label>Purchased Items Description</label>
          <input name="purchased_items_desc" value="Shoes" required />
          <button type="submit">Generate Session + Perform C2B <span id="loading" class="htmx-indicator">…</span></button>
        </form>
      </section>
    </div>
    <section class="card" style="margin-top:24px">
      <h2>Result</h2>
      <div id="result"><pre>Submit the form to see the session and C2B response.</pre></div>
    </section>
  </div>
</body>
</html>`))

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = page.Execute(w, nil)
	})
	mux.HandleFunc("/api/c2b", c2bHandler)

	log.Printf("C2B tester running at http://localhost%s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}

func c2bHandler(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		writePre(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	market, ok := mpesa.MarketFromContext(r.FormValue("market"))
	if !ok {
		writePre(w, http.StatusBadRequest, map[string]any{"error": "unsupported market"})
		return
	}

	cfg := mpesa.Config{
		APIKey:      strings.TrimSpace(r.FormValue("api_key")),
		PublicKey:   strings.TrimSpace(r.FormValue("public_key")),
		Environment: mpesa.Environment(strings.TrimSpace(r.FormValue("environment"))),
		Market:      market,
		Origin:      defaultString(r.FormValue("origin"), "127.0.0.1"),
	}
	client, err := mpesa.NewClient(cfg)
	if err != nil {
		writePre(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	resolved := client.Config()
	bearerToken, tokenErr := mpesa.EncryptBearerValue(cfg.APIKey, resolved.PublicKey)
	if tokenErr != nil {
		writePre(w, http.StatusBadRequest, map[string]any{"stage": "create_bearer_token", "error": tokenErr.Error()})
		return
	}
	log.Printf("[MPESA][TESTER][CONFIG] env=%s market=%s country=%s currency=%s origin=%q host=%s", resolved.Environment, resolved.Market.Context, resolved.Market.Country, resolved.Market.Currency, resolved.Origin, resolved.Host)
	log.Printf("[MPESA][TESTER][KEY] public_key_sha256=%s", fingerprint(resolved.PublicKey))
	log.Printf("[MPESA][TESTER][SESSION] bearer_token=%s", redact(bearerToken))

	waitSeconds, _ := strconv.Atoi(defaultString(r.FormValue("wait_seconds"), "30"))
	if waitSeconds < 0 {
		waitSeconds = 0
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(waitSeconds+90)*time.Second)
	defer cancel()

	log.Printf("[MPESA][SESSION][request] method=GET path=/%s/ipg/v2/%s/getSession/ origin=%q", resolved.Environment.BasePath(), resolved.Market.Context, resolved.Origin)
	session, sessionRaw, err := client.GenerateSessionAndWait(ctx, time.Duration(waitSeconds)*time.Second)
	log.Printf("[MPESA][SESSION][raw_body]=%s", rawString(sessionRaw))
	if err != nil {
		writePre(w, http.StatusBadGateway, map[string]any{
			"stage":       "generate_session",
			"error":       err.Error(),
			"hint":        "M-Pesa rejected the session request. If raw_response says Origin not allowed, set Origin to the exact value configured in your M-Pesa portal application, e.g. 127.0.0.1 or your domain/IP. The bearer token was generated from the API key and public key before this request.",
			"origin_used": resolved.Origin,
			"market":      resolved.Market.Context,
			"environment": resolved.Environment,
			"token_check": map[string]any{"created": true, "redacted_bearer_token": redact(bearerToken), "public_key_sha256": fingerprint(resolved.PublicKey)},
			"raw_response": rawString(sessionRaw),
		})
		return
	}
	log.Printf("[MPESA][SESSION][expected]=%s", mustJSON(map[string]any{"output_ResponseCode": "INS-0", "output_ResponseDesc": "Request processed successfully", "output_SessionID": redact(session.ID)}))

	conversationID := strings.TrimSpace(r.FormValue("third_party_conversation_id"))
	if conversationID == "" {
		conversationID = "go" + strconv.FormatInt(time.Now().UnixNano(), 36)
	}

	req := client.NewC2BSingleStageRequest(
		r.FormValue("amount"),
		r.FormValue("customer_msisdn"),
		r.FormValue("service_provider_code"),
		r.FormValue("transaction_reference"),
		conversationID,
		r.FormValue("purchased_items_desc"),
	)

	log.Printf("[MPESA][C2B][request]=%s", mustJSON(req))
	res, c2bRaw, err := client.C2BSingleStageWithSession(ctx, session, req)
	log.Printf("[MPESA][C2B][raw_body]=%s", rawString(c2bRaw))
	if err != nil {
		writePre(w, http.StatusBadGateway, map[string]any{
			"stage":        "c2b",
			"error":        err.Error(),
			"session_id":   redact(session.ID),
			"request":      req,
			"raw_response": rawString(c2bRaw),
		})
		return
	}

	writePre(w, http.StatusOK, map[string]any{
		"session": map[string]any{
			"id":          redact(session.ID),
			"market":      session.Market.Context,
			"country":     session.Market.Country,
			"currency":    client.Config().Market.Currency,
			"environment": session.Environment,
		},
		"request":      req,
		"response":     res,
		"raw_response": rawString(c2bRaw),
	})
}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func writePre(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	b, _ := json.MarshalIndent(payload, "", "  ")
	_, _ = fmt.Fprintf(w, "<pre>%s</pre>", template.HTMLEscapeString(string(b)))
}

func rawString(raw *mpesa.RawResponse) string {
	if raw == nil {
		return ""
	}
	return raw.BodyString()
}

func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}

func fingerprint(value string) string {
	clean := strings.Join(strings.Fields(value), "")
	sum := sha256.Sum256([]byte(clean))
	return hex.EncodeToString(sum[:])[:16]
}

func redact(value string) string {
	if len(value) <= 8 {
		return "***"
	}
	return value[:4] + "…" + value[len(value)-4:]
}

func defaultString(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
