package tap

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestNewExchangeMessage(t *testing.T) {
	body := &ExchangeBody{
		FromAssets: []string{"eip155:1/slip44:60"},
		ToAssets:   []string{"USD"},
		FromAmount: "1.0",
		Requester:  &Party{ID: "did:eg:alice"},
		Agents:     []Agent{{ID: "did:web:exchange.example", For: NewForField("did:eg:alice")}},
	}

	msg, err := NewExchangeMessage("did:web:exchange.example", []string{"did:web:provider.example"}, body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeExchange {
		t.Errorf("Type: got %q", msg.Type)
	}

	var got ExchangeBody
	if err := json.Unmarshal(msg.Body, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.FromAmount != "1.0" {
		t.Errorf("FromAmount: got %q", got.FromAmount)
	}
}

func TestNewExchangeMessage_MissingFromAssets(t *testing.T) {
	body := &ExchangeBody{
		ToAssets:   []string{"USD"},
		FromAmount: "1.0",
		Requester:  &Party{ID: "did:eg:alice"},
		Agents:     []Agent{{ID: "did:web:exchange.example"}},
	}
	_, err := NewExchangeMessage("from", nil, body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestNewExchangeMessage_MissingAmounts(t *testing.T) {
	body := &ExchangeBody{
		FromAssets: []string{"ETH"},
		ToAssets:   []string{"USD"},
		Requester:  &Party{ID: "did:eg:alice"},
		Agents:     []Agent{{ID: "did:web:exchange.example"}},
	}
	_, err := NewExchangeMessage("from", nil, body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestExchangeBody_JSONRoundTrip(t *testing.T) {
	body := ExchangeBody{
		Context:    TAPContext,
		Type:       TypeExchange,
		FromAssets: []string{"eip155:1/slip44:60"},
		ToAssets:   []string{"USD"},
		FromAmount: "1.0",
		Requester:  &Party{ID: "did:eg:alice"},
		Agents:     []Agent{{ID: "did:web:exchange.example"}},
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got ExchangeBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.FromAssets) != 1 || got.FromAssets[0] != "eip155:1/slip44:60" {
		t.Errorf("FromAssets: got %v", got.FromAssets)
	}
}

func TestExchangeBody_ParseBody(t *testing.T) {
	body := &ExchangeBody{
		FromAssets: []string{"ETH"},
		ToAssets:   []string{"USD"},
		ToAmount:   "3000.00",
		Requester:  &Party{ID: "did:eg:alice"},
		Agents:     []Agent{{ID: "did:web:exchange.example"}},
	}
	msg, err := NewExchangeMessage("did:web:exchange.example", nil, body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	eb, ok := parsed.(*ExchangeBody)
	if !ok {
		t.Fatalf("expected *ExchangeBody, got %T", parsed)
	}
	if eb.ToAmount != "3000.00" {
		t.Errorf("ToAmount: got %q", eb.ToAmount)
	}
}
