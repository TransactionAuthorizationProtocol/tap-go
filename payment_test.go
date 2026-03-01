package tap

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestNewPaymentMessage(t *testing.T) {
	body := &PaymentBody{
		Amount:   "250.00",
		Currency: "EUR",
		Merchant: &Party{ID: "did:example:merchant", Name: "Digital Goods Store"},
		Agents: []Agent{
			{ID: "did:example:merchant", For: NewForField("did:example:merchant")},
		},
	}

	msg, err := NewPaymentMessage("did:example:merchant", []string{"did:example:customer"}, body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypePayment {
		t.Errorf("Type: got %q", msg.Type)
	}

	var got PaymentBody
	if err := json.Unmarshal(msg.Body, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Amount != "250.00" || got.Currency != "EUR" {
		t.Errorf("mismatch: amount=%q currency=%q", got.Amount, got.Currency)
	}
}

func TestNewPaymentMessage_MissingAmount(t *testing.T) {
	body := &PaymentBody{
		Currency: "EUR",
		Merchant: &Party{ID: "did:example:merchant"},
		Agents:   []Agent{{ID: "did:example:merchant"}},
	}
	_, err := NewPaymentMessage("did:example:merchant", nil, body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestNewPaymentMessage_MissingMerchant(t *testing.T) {
	body := &PaymentBody{
		Amount:   "250.00",
		Currency: "EUR",
		Agents:   []Agent{{ID: "did:example:merchant"}},
	}
	_, err := NewPaymentMessage("did:example:merchant", nil, body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestNewPaymentMessage_MissingAssetAndCurrency(t *testing.T) {
	body := &PaymentBody{
		Amount:   "250.00",
		Merchant: &Party{ID: "did:example:merchant"},
		Agents:   []Agent{{ID: "did:example:merchant"}},
	}
	_, err := NewPaymentMessage("did:example:merchant", nil, body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestPaymentBody_JSONRoundTrip(t *testing.T) {
	body := PaymentBody{
		Context:  TAPContext,
		Type:     TypePayment,
		Amount:   "250.00",
		Currency: "EUR",
		Merchant: &Party{ID: "did:example:merchant", Name: "Shop"},
		Agents:   []Agent{{ID: "did:example:merchant"}},
		SupportedAssets: []any{
			"eip155:1/erc20:0xA0b86991c53D94fa4C0bCBf0C1C4DF2F15F1b7A8",
		},
		FallbackSettlementAddresses: []string{"eip155:137:0x8B5e7A2C3f4D1E6F9A0b3C5e7D9f1A2B4C6E8F0A"},
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got PaymentBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Amount != "250.00" || got.Currency != "EUR" {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestPayment_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/payment-request/valid-payment.json")
	if err != nil {
		t.Skipf("test vector not available: %v", err)
	}

	var tv struct {
		Body json.RawMessage `json:"body"`
	}
	if err := json.Unmarshal(data, &tv); err != nil {
		t.Fatalf("unmarshal test vector: %v", err)
	}

	var body PaymentBody
	if err := json.Unmarshal(tv.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if body.Amount != "250.00" {
		t.Errorf("Amount: got %q", body.Amount)
	}
	if body.Currency != "EUR" {
		t.Errorf("Currency: got %q", body.Currency)
	}
	if body.Merchant == nil || body.Merchant.ID != "did:example:merchant" {
		t.Errorf("Merchant: got %+v", body.Merchant)
	}
}

func TestPaymentBody_ParseBody(t *testing.T) {
	body := &PaymentBody{
		Amount:   "100.00",
		Currency: "USD",
		Merchant: &Party{ID: "did:example:merchant"},
		Agents:   []Agent{{ID: "did:example:merchant"}},
	}
	msg, err := NewPaymentMessage("did:example:merchant", []string{"did:example:customer"}, body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	pb, ok := parsed.(*PaymentBody)
	if !ok {
		t.Fatalf("expected *PaymentBody, got %T", parsed)
	}
	if pb.Amount != "100.00" {
		t.Errorf("Amount: got %q", pb.Amount)
	}
}
