package tap

import (
	"encoding/json"
	"os"
	"testing"
)

func TestNewRejectMessage(t *testing.T) {
	body := &RejectBody{Reason: "Beneficiary name mismatch"}
	msg, err := NewRejectMessage("did:web:beneficiary.vasp", []string{"did:web:originator.vasp"}, "1234567890", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeReject {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewRejectMessage_NoReasonOptional(t *testing.T) {
	body := &RejectBody{}
	_, err := NewRejectMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRejectBody_JSONRoundTrip(t *testing.T) {
	body := RejectBody{
		Context: TAPContext,
		Type:    TypeReject,
		Reason:  "Compliance check failed",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got RejectBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Reason != body.Reason {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestReject_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/reject/valid.json")
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

	var body RejectBody
	if err := json.Unmarshal(tv.Message.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if body.Reason != "Beneficiary name mismatch" {
		t.Errorf("Reason: got %q", body.Reason)
	}
}

func TestRejectBody_ParseBody(t *testing.T) {
	body := &RejectBody{Reason: "test"}
	msg, err := NewRejectMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	rb, ok := parsed.(*RejectBody)
	if !ok {
		t.Fatalf("expected *RejectBody, got %T", parsed)
	}
	if rb.Reason != "test" {
		t.Errorf("Reason: got %q", rb.Reason)
	}
}
