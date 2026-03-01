package tap

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestNewAuthorizationRequiredMessage(t *testing.T) {
	body := &AuthorizationRequiredBody{
		AuthorizationURL: "https://vasp.example/authorize?request=abc123",
		Expires:          "2024-03-22T15:00:00Z",
	}

	msg, err := NewAuthorizationRequiredMessage("did:web:vasp.example", []string{"did:web:b2b-service.example"}, "thread-1", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeAuthorizationRequired {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewAuthorizationRequiredMessage_MissingURL(t *testing.T) {
	body := &AuthorizationRequiredBody{Expires: "2024-03-22T15:00:00Z"}
	_, err := NewAuthorizationRequiredMessage("from", nil, "thid", body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestNewAuthorizationRequiredMessage_MissingExpires(t *testing.T) {
	body := &AuthorizationRequiredBody{AuthorizationURL: "https://example.com"}
	_, err := NewAuthorizationRequiredMessage("from", nil, "thid", body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestAuthorizationRequiredBody_JSONRoundTrip(t *testing.T) {
	body := AuthorizationRequiredBody{
		Context:          TAPContext,
		Type:             TypeAuthorizationRequired,
		AuthorizationURL: "https://vasp.example/auth",
		Expires:          "2024-03-22T15:00:00Z",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got AuthorizationRequiredBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.AuthorizationURL != body.AuthorizationURL {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestAuthorizationRequired_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/authorization-required/valid-authorization-required.json")
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

	var body AuthorizationRequiredBody
	if err := json.Unmarshal(tv.Message.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if body.AuthorizationURL != "https://vasp.example/authorize?request=abc123" {
		t.Errorf("AuthorizationURL: got %q", body.AuthorizationURL)
	}
}

func TestAuthorizationRequiredBody_ParseBody(t *testing.T) {
	body := &AuthorizationRequiredBody{
		AuthorizationURL: "https://example.com/auth",
		Expires:          "2024-12-31T00:00:00Z",
	}
	msg, err := NewAuthorizationRequiredMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	arb, ok := parsed.(*AuthorizationRequiredBody)
	if !ok {
		t.Fatalf("expected *AuthorizationRequiredBody, got %T", parsed)
	}
	if arb.AuthorizationURL != "https://example.com/auth" {
		t.Errorf("AuthorizationURL: got %q", arb.AuthorizationURL)
	}
}
