package tap

import (
	"encoding/json"
	"os"
	"testing"
)

func TestNewAuthorizeMessage(t *testing.T) {
	body := &AuthorizeBody{
		SettlementAddress: "eip155:1:0x1234a96D359eC26a11e2C2b3d8f8B8942d5Bfcdb",
	}

	msg, err := NewAuthorizeMessage("did:web:beneficiary.vasp", []string{"did:web:originator.vasp"}, "1234567890", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeAuthorize {
		t.Errorf("Type: got %q", msg.Type)
	}
	if msg.Thid != "1234567890" {
		t.Errorf("Thid: got %q", msg.Thid)
	}
}

func TestNewAuthorizeMessage_AllFieldsOptional(t *testing.T) {
	body := &AuthorizeBody{}
	msg, err := NewAuthorizeMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got AuthorizeBody
	if err := json.Unmarshal(msg.Body, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Context != TAPContext {
		t.Errorf("Context: got %q", got.Context)
	}
}

func TestAuthorizeBody_JSONRoundTrip(t *testing.T) {
	body := AuthorizeBody{
		Context:           TAPContext,
		Type:              TypeAuthorize,
		SettlementAddress: "eip155:1:0x1234",
		Amount:            "100.00",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got AuthorizeBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.SettlementAddress != body.SettlementAddress {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestAuthorize_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/authorize/valid.json")
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

	var body AuthorizeBody
	if err := json.Unmarshal(tv.Message.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if body.SettlementAddress != "eip155:1:0x1234a96D359eC26a11e2C2b3d8f8B8942d5Bfcdb" {
		t.Errorf("SettlementAddress: got %q", body.SettlementAddress)
	}
}

func TestAuthorizeBody_ParseBody(t *testing.T) {
	body := &AuthorizeBody{SettlementAddress: "eip155:1:0x1234"}
	msg, err := NewAuthorizeMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	ab, ok := parsed.(*AuthorizeBody)
	if !ok {
		t.Fatalf("expected *AuthorizeBody, got %T", parsed)
	}
	if ab.SettlementAddress != "eip155:1:0x1234" {
		t.Errorf("SettlementAddress: got %q", ab.SettlementAddress)
	}
}
