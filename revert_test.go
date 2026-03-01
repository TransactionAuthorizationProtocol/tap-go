package tap

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestNewRevertMessage(t *testing.T) {
	body := &RevertBody{
		SettlementAddress: "payto://iban/DE75512108001245126199",
		Reason:            "Insufficient Originator Information",
	}
	msg, err := NewRevertMessage("did:web:beneficiary.vasp", []string{"did:web:originator.vasp"}, "thread-1", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeRevert {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewRevertMessage_MissingFields(t *testing.T) {
	tests := []struct {
		name string
		body *RevertBody
	}{
		{"missing settlementAddress", &RevertBody{Reason: "test"}},
		{"missing reason", &RevertBody{SettlementAddress: "eip155:1:0x1234"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRevertMessage("from", nil, "thid", tt.body)
			if !errors.Is(err, ErrInvalidBody) {
				t.Errorf("expected ErrInvalidBody, got %v", err)
			}
		})
	}
}

func TestRevertBody_JSONRoundTrip(t *testing.T) {
	body := RevertBody{
		Context:           TAPContext,
		Type:              TypeRevert,
		SettlementAddress: "payto://iban/DE75512108001245126199",
		Reason:            "Dispute",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got RevertBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.SettlementAddress != body.SettlementAddress {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestRevert_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/revert/valid-revert.json")
	if err != nil {
		t.Skipf("test vector not available: %v", err)
	}

	var tv struct {
		Body json.RawMessage `json:"body"`
	}
	if err := json.Unmarshal(data, &tv); err != nil {
		t.Fatalf("unmarshal test vector: %v", err)
	}

	var body RevertBody
	if err := json.Unmarshal(tv.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if body.SettlementAddress != "payto://iban/DE75512108001245126199" {
		t.Errorf("SettlementAddress: got %q", body.SettlementAddress)
	}
	if body.Reason != "Insufficient Originator Information" {
		t.Errorf("Reason: got %q", body.Reason)
	}
}

func TestRevertBody_ParseBody(t *testing.T) {
	body := &RevertBody{SettlementAddress: "eip155:1:0x1234", Reason: "dispute"}
	msg, err := NewRevertMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	rb, ok := parsed.(*RevertBody)
	if !ok {
		t.Fatalf("expected *RevertBody, got %T", parsed)
	}
	if rb.Reason != "dispute" {
		t.Errorf("Reason: got %q", rb.Reason)
	}
}
