package tap

import (
	"encoding/json"
	"testing"
)

func TestNewCaptureMessage(t *testing.T) {
	body := &CaptureBody{
		Amount:            "500.00",
		SettlementAddress: "eip155:1:0x1234",
	}
	msg, err := NewCaptureMessage("did:web:merchant", []string{"did:web:escrow"}, "escrow-thread-1", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeCapture {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewCaptureMessage_AllOptional(t *testing.T) {
	body := &CaptureBody{}
	msg, err := NewCaptureMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got CaptureBody
	if err := json.Unmarshal(msg.Body, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Context != TAPContext {
		t.Errorf("Context: got %q", got.Context)
	}
}

func TestCaptureBody_JSONRoundTrip(t *testing.T) {
	body := CaptureBody{
		Context:           TAPContext,
		Type:              TypeCapture,
		Amount:            "250.00",
		SettlementAddress: "eip155:1:0x5678",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got CaptureBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Amount != body.Amount || got.SettlementAddress != body.SettlementAddress {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestCaptureBody_ParseBody(t *testing.T) {
	body := &CaptureBody{Amount: "100.00"}
	msg, err := NewCaptureMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cb, ok := parsed.(*CaptureBody)
	if !ok {
		t.Fatalf("expected *CaptureBody, got %T", parsed)
	}
	if cb.Amount != "100.00" {
		t.Errorf("Amount: got %q", cb.Amount)
	}
}
