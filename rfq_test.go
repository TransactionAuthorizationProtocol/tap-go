package tap

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestNewRFQMessage(t *testing.T) {
	body := &RFQBody{
		FromAssets: []string{"eip155:1/slip44:60"},
		ToAssets:   []string{"USD"},
		FromAmount: "1.0",
		Requester:  &Party{ID: "did:eg:alice"},
		Agents:     []Agent{{ID: "did:web:exchange.example", For: NewForField("did:eg:alice")}},
	}

	msg, err := NewRFQMessage("did:web:exchange.example", []string{"did:web:provider.example"}, body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeRFQ {
		t.Errorf("Type: got %q", msg.Type)
	}

	var got RFQBody
	if err := json.Unmarshal(msg.Body, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.FromAmount != "1.0" {
		t.Errorf("FromAmount: got %q", got.FromAmount)
	}
}

func TestNewRFQMessage_MissingFromAssets(t *testing.T) {
	body := &RFQBody{
		ToAssets:   []string{"USD"},
		FromAmount: "1.0",
		Requester:  &Party{ID: "did:eg:alice"},
		Agents:     []Agent{{ID: "did:web:exchange.example"}},
	}
	_, err := NewRFQMessage("from", nil, body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestNewRFQMessage_MissingAmounts(t *testing.T) {
	body := &RFQBody{
		FromAssets: []string{"ETH"},
		ToAssets:   []string{"USD"},
		Requester:  &Party{ID: "did:eg:alice"},
		Agents:     []Agent{{ID: "did:web:exchange.example"}},
	}
	_, err := NewRFQMessage("from", nil, body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestRFQBody_JSONRoundTrip(t *testing.T) {
	body := RFQBody{
		Context:    TAPContext,
		Type:       TypeRFQ,
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

	var got RFQBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.FromAssets) != 1 || got.FromAssets[0] != "eip155:1/slip44:60" {
		t.Errorf("FromAssets: got %v", got.FromAssets)
	}
}

func TestRFQBody_ParseBody(t *testing.T) {
	body := &RFQBody{
		FromAssets: []string{"ETH"},
		ToAssets:   []string{"USD"},
		ToAmount:   "3000.00",
		Requester:  &Party{ID: "did:eg:alice"},
		Agents:     []Agent{{ID: "did:web:exchange.example"}},
	}
	msg, err := NewRFQMessage("did:web:exchange.example", nil, body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	rb, ok := parsed.(*RFQBody)
	if !ok {
		t.Fatalf("expected *RFQBody, got %T", parsed)
	}
	if rb.ToAmount != "3000.00" {
		t.Errorf("ToAmount: got %q", rb.ToAmount)
	}
}
