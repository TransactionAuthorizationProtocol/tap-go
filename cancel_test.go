package tap

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestNewCancelMessage(t *testing.T) {
	body := &CancelBody{By: "originator", Reason: "User requested cancellation"}
	msg, err := NewCancelMessage("did:web:originator.vasp", []string{"did:web:beneficiary.vasp"}, "1234567890", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeCancel {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewCancelMessage_MissingBy(t *testing.T) {
	body := &CancelBody{Reason: "test"}
	_, err := NewCancelMessage("from", nil, "thid", body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestCancelBody_JSONRoundTrip(t *testing.T) {
	body := CancelBody{
		Context: TAPContext,
		Type:    TypeCancel,
		By:      "originator",
		Reason:  "Changed mind",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got CancelBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.By != body.By || got.Reason != body.Reason {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestCancel_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/cancel/valid-transaction-cancel.json")
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

	var body CancelBody
	if err := json.Unmarshal(tv.Message.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if body.By != "originator" {
		t.Errorf("By: got %q", body.By)
	}
}

func TestCancelBody_ParseBody(t *testing.T) {
	body := &CancelBody{By: "originator"}
	msg, err := NewCancelMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cb, ok := parsed.(*CancelBody)
	if !ok {
		t.Fatalf("expected *CancelBody, got %T", parsed)
	}
	if cb.By != "originator" {
		t.Errorf("By: got %q", cb.By)
	}
}
