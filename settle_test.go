package tap

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestNewSettleMessage(t *testing.T) {
	body := &SettleBody{
		SettlementAddress: "eip155:1:0x1234a96D359eC26a11e2C2b3d8f8B8942d5Bfcdb",
		SettlementID:      "eip155:1:tx/0x3edb98c24d46d148eb926c714f4fbaa117c47b0c0821f38bfce9763604457c33",
	}

	msg, err := NewSettleMessage("did:web:originator.vasp", []string{"did:web:beneficiary.vasp"}, "1234567890", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeSettle {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewSettleMessage_MissingAddress(t *testing.T) {
	body := &SettleBody{}
	_, err := NewSettleMessage("from", nil, "thid", body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestSettleBody_JSONRoundTrip(t *testing.T) {
	body := SettleBody{
		Context:           TAPContext,
		Type:              TypeSettle,
		SettlementAddress: "eip155:1:0x1234",
		SettlementID:      "eip155:1:tx/0xabc",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got SettleBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.SettlementAddress != body.SettlementAddress {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestSettle_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/settle/valid.json")
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

	var body SettleBody
	if err := json.Unmarshal(tv.Message.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if body.SettlementAddress != "eip155:1:0x1234a96D359eC26a11e2C2b3d8f8B8942d5Bfcdb" {
		t.Errorf("SettlementAddress: got %q", body.SettlementAddress)
	}
}

func TestSettleBody_ParseBody(t *testing.T) {
	body := &SettleBody{SettlementAddress: "eip155:1:0xabc"}
	msg, err := NewSettleMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	sb, ok := parsed.(*SettleBody)
	if !ok {
		t.Fatalf("expected *SettleBody, got %T", parsed)
	}
	if sb.SettlementAddress != "eip155:1:0xabc" {
		t.Errorf("SettlementAddress: got %q", sb.SettlementAddress)
	}
}
