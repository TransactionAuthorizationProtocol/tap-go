package tap

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestNewQuoteMessage(t *testing.T) {
	body := &QuoteBody{
		FromAsset:  "eip155:1/slip44:60",
		ToAsset:    "USD",
		FromAmount: "1.0",
		ToAmount:   "3000.00",
		Provider:   &Party{ID: "did:web:provider.example"},
		Agents:     []Agent{{ID: "did:web:provider.example"}},
		ExpiresAt:  "2025-03-22T15:00:00Z",
	}

	msg, err := NewQuoteMessage("did:web:provider.example", []string{"did:web:requester.example"}, "exchange-thread-1", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeQuote {
		t.Errorf("Type: got %q", msg.Type)
	}
	if msg.Thid != "exchange-thread-1" {
		t.Errorf("Thid: got %q", msg.Thid)
	}
}

func TestNewQuoteMessage_MissingFields(t *testing.T) {
	tests := []struct {
		name string
		body *QuoteBody
	}{
		{"missing fromAsset", &QuoteBody{ToAsset: "USD", FromAmount: "1", ToAmount: "3000", Provider: &Party{ID: "did:eg:p"}, Agents: []Agent{{ID: "a"}}, ExpiresAt: "2025-01-01T00:00:00Z"}},
		{"missing toAsset", &QuoteBody{FromAsset: "ETH", FromAmount: "1", ToAmount: "3000", Provider: &Party{ID: "did:eg:p"}, Agents: []Agent{{ID: "a"}}, ExpiresAt: "2025-01-01T00:00:00Z"}},
		{"missing fromAmount", &QuoteBody{FromAsset: "ETH", ToAsset: "USD", ToAmount: "3000", Provider: &Party{ID: "did:eg:p"}, Agents: []Agent{{ID: "a"}}, ExpiresAt: "2025-01-01T00:00:00Z"}},
		{"missing toAmount", &QuoteBody{FromAsset: "ETH", ToAsset: "USD", FromAmount: "1", Provider: &Party{ID: "did:eg:p"}, Agents: []Agent{{ID: "a"}}, ExpiresAt: "2025-01-01T00:00:00Z"}},
		{"missing provider", &QuoteBody{FromAsset: "ETH", ToAsset: "USD", FromAmount: "1", ToAmount: "3000", Agents: []Agent{{ID: "a"}}, ExpiresAt: "2025-01-01T00:00:00Z"}},
		{"missing agents", &QuoteBody{FromAsset: "ETH", ToAsset: "USD", FromAmount: "1", ToAmount: "3000", Provider: &Party{ID: "did:eg:p"}, ExpiresAt: "2025-01-01T00:00:00Z"}},
		{"missing expiresAt", &QuoteBody{FromAsset: "ETH", ToAsset: "USD", FromAmount: "1", ToAmount: "3000", Provider: &Party{ID: "did:eg:p"}, Agents: []Agent{{ID: "a"}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewQuoteMessage("from", nil, "thid", tt.body)
			if !errors.Is(err, ErrInvalidBody) {
				t.Errorf("expected ErrInvalidBody, got %v", err)
			}
		})
	}
}

func TestQuoteBody_JSONRoundTrip(t *testing.T) {
	body := QuoteBody{
		Context:    TAPContext,
		Type:       TypeQuote,
		FromAsset:  "eip155:1/slip44:60",
		ToAsset:    "USD",
		FromAmount: "1.0",
		ToAmount:   "3000.00",
		Provider:   &Party{ID: "did:web:provider.example"},
		Agents:     []Agent{{ID: "did:web:provider.example"}},
		ExpiresAt:  "2025-03-22T15:00:00Z",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got QuoteBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.FromAsset != body.FromAsset || got.ToAmount != body.ToAmount {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestQuoteBody_ParseBody(t *testing.T) {
	body := &QuoteBody{
		FromAsset:  "ETH",
		ToAsset:    "USD",
		FromAmount: "1",
		ToAmount:   "3000",
		Provider:   &Party{ID: "did:web:provider"},
		Agents:     []Agent{{ID: "did:web:provider"}},
		ExpiresAt:  "2025-01-01T00:00:00Z",
	}
	msg, err := NewQuoteMessage("from", nil, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	qb, ok := parsed.(*QuoteBody)
	if !ok {
		t.Fatalf("expected *QuoteBody, got %T", parsed)
	}
	if qb.FromAsset != "ETH" {
		t.Errorf("FromAsset: got %q", qb.FromAsset)
	}
}
