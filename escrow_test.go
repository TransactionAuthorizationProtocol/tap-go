package tap

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestNewEscrowMessage(t *testing.T) {
	body := &EscrowBody{
		Asset:       "eip155:1/erc20:0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		Amount:      "1000.00",
		Originator:  &Party{ID: "did:web:buyer.example"},
		Beneficiary: &Party{ID: "did:web:seller.example"},
		Expiry:      "2024-03-15T00:00:00Z",
		Agents: []Agent{
			{ID: "did:web:escrow-service.example", Role: "EscrowAgent"},
		},
	}

	msg, err := NewEscrowMessage("did:web:buyer.example", []string{"did:web:escrow-service.example"}, body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeEscrow {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewEscrowMessage_MissingFields(t *testing.T) {
	base := EscrowBody{
		Asset:       "eip155:1/erc20:0xA0b86991",
		Amount:      "1000.00",
		Originator:  &Party{ID: "did:web:buyer"},
		Beneficiary: &Party{ID: "did:web:seller"},
		Expiry:      "2024-03-15T00:00:00Z",
		Agents:      []Agent{{ID: "did:web:escrow", Role: "EscrowAgent"}},
	}

	tests := []struct {
		name   string
		modify func(*EscrowBody)
	}{
		{"missing amount", func(b *EscrowBody) { b.Amount = "" }},
		{"missing originator", func(b *EscrowBody) { b.Originator = nil }},
		{"missing beneficiary", func(b *EscrowBody) { b.Beneficiary = nil }},
		{"missing expiry", func(b *EscrowBody) { b.Expiry = "" }},
		{"missing agents", func(b *EscrowBody) { b.Agents = nil }},
		{"missing asset and currency", func(b *EscrowBody) { b.Asset = "" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := base
			tt.modify(&b)
			_, err := NewEscrowMessage("from", nil, &b)
			if !errors.Is(err, ErrInvalidBody) {
				t.Errorf("expected ErrInvalidBody, got %v", err)
			}
		})
	}
}

func TestEscrowBody_JSONRoundTrip(t *testing.T) {
	body := EscrowBody{
		Context:     TAPContext,
		Type:        TypeEscrow,
		Asset:       "eip155:1/erc20:0xa0b86991",
		Amount:      "1000.00",
		Originator:  &Party{ID: "did:web:buyer.example"},
		Beneficiary: &Party{ID: "did:web:seller.example"},
		Expiry:      "2024-03-15T00:00:00Z",
		Agents:      []Agent{{ID: "did:web:escrow-service.example", Role: "EscrowAgent"}},
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got EscrowBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Amount != body.Amount || got.Asset != body.Asset {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestEscrow_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/escrow/valid-escrow.json")
	if err != nil {
		t.Skipf("test vector not available: %v", err)
	}

	var tv struct {
		Message struct {
			Body json.RawMessage `json:"body"`
		} `json:"message"`
	}
	if err := json.Unmarshal(data, &tv); err != nil {
		t.Fatalf("unmarshal test vector: %v", err)
	}

	var body EscrowBody
	if err := json.Unmarshal(tv.Message.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if body.Amount != "1000.00" {
		t.Errorf("Amount: got %q", body.Amount)
	}
	if body.Originator == nil || body.Originator.ID != "did:web:buyer.example" {
		t.Errorf("Originator: got %+v", body.Originator)
	}
	if len(body.Agents) != 3 {
		t.Errorf("Agents: got %d, want 3", len(body.Agents))
	}
}

func TestEscrowBody_ParseBody(t *testing.T) {
	body := &EscrowBody{
		Asset:       "eip155:1/erc20:0xa0b86991",
		Amount:      "500.00",
		Originator:  &Party{ID: "did:web:buyer"},
		Beneficiary: &Party{ID: "did:web:seller"},
		Expiry:      "2024-12-31T00:00:00Z",
		Agents:      []Agent{{ID: "did:web:escrow", Role: "EscrowAgent"}},
	}
	msg, err := NewEscrowMessage("did:web:buyer", []string{"did:web:escrow"}, body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	eb, ok := parsed.(*EscrowBody)
	if !ok {
		t.Fatalf("expected *EscrowBody, got %T", parsed)
	}
	if eb.Amount != "500.00" {
		t.Errorf("Amount: got %q", eb.Amount)
	}
}
