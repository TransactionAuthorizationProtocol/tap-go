package tap

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestNewLockMessage(t *testing.T) {
	body := &LockBody{
		Asset:       "eip155:1/erc20:0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		Amount:      "1000.00",
		Originator:  &Party{ID: "did:web:buyer.example"},
		Beneficiary: &Party{ID: "did:web:seller.example"},
		Expiry:      "2024-03-15T00:00:00Z",
		Agents: []Agent{
			{ID: "did:web:escrow-service.example", Role: "EscrowAgent"},
		},
	}

	msg, err := NewLockMessage("did:web:buyer.example", []string{"did:web:escrow-service.example"}, body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeLock {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewLockMessage_MissingFields(t *testing.T) {
	base := LockBody{
		Asset:       "eip155:1/erc20:0xA0b86991",
		Amount:      "1000.00",
		Originator:  &Party{ID: "did:web:buyer"},
		Beneficiary: &Party{ID: "did:web:seller"},
		Expiry:      "2024-03-15T00:00:00Z",
		Agents:      []Agent{{ID: "did:web:escrow", Role: "EscrowAgent"}},
	}

	tests := []struct {
		name   string
		modify func(*LockBody)
	}{
		{"missing amount", func(b *LockBody) { b.Amount = "" }},
		{"missing originator", func(b *LockBody) { b.Originator = nil }},
		{"missing beneficiary", func(b *LockBody) { b.Beneficiary = nil }},
		{"missing expiry", func(b *LockBody) { b.Expiry = "" }},
		{"missing agents", func(b *LockBody) { b.Agents = nil }},
		{"missing asset and currency", func(b *LockBody) { b.Asset = "" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := base
			tt.modify(&b)
			_, err := NewLockMessage("from", nil, &b)
			if !errors.Is(err, ErrInvalidBody) {
				t.Errorf("expected ErrInvalidBody, got %v", err)
			}
		})
	}
}

func TestLockBody_JSONRoundTrip(t *testing.T) {
	body := LockBody{
		Context:     TAPContext,
		Type:        TypeLock,
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

	var got LockBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Amount != body.Amount || got.Asset != body.Asset {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestLock_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/lock/valid-lock.json")
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

	var body LockBody
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

func TestLockBody_ParseBody(t *testing.T) {
	body := &LockBody{
		Asset:       "eip155:1/erc20:0xa0b86991",
		Amount:      "500.00",
		Originator:  &Party{ID: "did:web:buyer"},
		Beneficiary: &Party{ID: "did:web:seller"},
		Expiry:      "2024-12-31T00:00:00Z",
		Agents:      []Agent{{ID: "did:web:escrow", Role: "EscrowAgent"}},
	}
	msg, err := NewLockMessage("did:web:buyer", []string{"did:web:escrow"}, body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	lb, ok := parsed.(*LockBody)
	if !ok {
		t.Fatalf("expected *LockBody, got %T", parsed)
	}
	if lb.Amount != "500.00" {
		t.Errorf("Amount: got %q", lb.Amount)
	}
}
